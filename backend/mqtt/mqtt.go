// Copyright (c) 2021, CELSA Group All rights reserved.
// Author: Albert Espín <
// Contributors:
package mqtt

import (
	"crypto/tls"
	"fmt"
	"sync"
	"time"

	"github.com/MGYOSBEL/janus/pkg/log"
	paho "github.com/eclipse/paho.mqtt.golang"
)

var (
	// PublishQoS Quality of Service Level
	PublishQoS byte = 0x00
	// SubscribeQoS Quality of Service Level
	SubscribeQoS byte = 0x00
	// PublishTimeout is the timeout before returning from publish without checking error
	PublishTimeout = 50 * time.Millisecond
	// BufferSize indicates the maximum number of MQTT messages that should be buffered
	BufferSize = 10
	// ConnectRetries says how many times the client should retry a failed connection
	ConnectRetries = 10
	// ConnectRetryDelay says how long the client should wait between retries
	ConnectRetryDelay = time.Second
)

// Config contains configuration for MQTT
type Config struct {
	Brokers   []string
	Username  string
	Password  string
	Topic     string
	TLSConfig *tls.Config
	BatchSize int
}

type subscription struct {
	handler paho.MessageHandler
	cancel  func()
}

// MQTT side of the bridge
type MQTT struct {
	logger        log.Logger
	client        paho.Client
	subscriptions map[string]subscription
	mu            sync.Mutex
	cfg           Config
}

// New returns a new MQTT
func New(config Config, logger log.Logger) (*MQTT, error) {
	mqtt := new(MQTT)
	mqtt.logger = logger
	mqtt.cfg = config
	mqttOpts := paho.NewClientOptions()
	for _, broker := range config.Brokers {
		mqttOpts.AddBroker(broker)
	}
	if config.TLSConfig != nil {
		mqttOpts.SetTLSConfig(config.TLSConfig)
	}
	mqttOpts.SetClientID(fmt.Sprintf("hermes-mqtt-tpa-%s", time.Now().Format("20060102150405")))
	mqttOpts.SetUsername(config.Username)
	mqttOpts.SetPassword(config.Password)
	mqttOpts.SetKeepAlive(30 * time.Second)
	mqttOpts.SetPingTimeout(10 * time.Second)
	mqttOpts.SetCleanSession(true)
	mqttOpts.SetDefaultPublishHandler(func(_ paho.Client, msg paho.Message) {
		mqtt.logger.Errorf("received message on unknown topic %s", msg.Topic())
	})

	mqtt.subscriptions = make(map[string]subscription)
	var reconnecting bool
	mqttOpts.SetConnectionLostHandler(func(_ paho.Client, err error) {
		mqtt.logger.Warnf("failed mqtt disconnected %s. Reconnecting...", err.Error())
		reconnecting = true
	})
	mqttOpts.SetOnConnectHandler(func(_ paho.Client) {
		mqtt.logger.Infof("connected to MQTT")
		if reconnecting {
			mqtt.resubscribe()
			reconnecting = false
		}
	})
	mqtt.client = paho.NewClient(mqttOpts)
	return mqtt, nil
}

// Connect to MQTT
func (c *MQTT) Connect() error {
	var err error
	for retries := 0; retries < ConnectRetries; retries++ {
		token := c.client.Connect()
		finished := token.WaitTimeout(1 * time.Second)
		if !finished {
			c.logger.Warnf("connection to MQTT broker timed out")
			token.Wait()
		}
		err = token.Error()
		if err == nil {
			break
		}
		c.logger.Warnf("connection to MQTT broker failed: %s retry: %s", err.Error(), ConnectRetryDelay.String())
		<-time.After(ConnectRetryDelay)
	}
	if err != nil {
		return fmt.Errorf("could not connect to MQTT failed: %s", err.Error())
	}
	return err
}

// Disconnect from MQTT
func (c *MQTT) Disconnect() error {
	c.client.Disconnect(100)
	return nil
}

func (c *MQTT) resubscribe() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for topic, subscription := range c.subscriptions {
		c.client.Subscribe(topic, SubscribeQoS, subscription.handler)
	}
}

func (c *MQTT) unsubscribe(topic string) paho.Token {
	c.mu.Lock()
	defer c.mu.Unlock()
	if subscription, ok := c.subscriptions[topic]; ok && subscription.cancel != nil {
		subscription.cancel()
	}
	delete(c.subscriptions, topic)
	return c.client.Unsubscribe(topic)
}

func (c *MQTT) subscribe(topic string, handler paho.MessageHandler, cancel func()) paho.Token {
	c.mu.Lock()
	defer c.mu.Unlock()
	wrappedHandler := func(client paho.Client, msg paho.Message) {
		if msg.Retained() {
			c.logger.Debugf("received retained message topic: %s", msg.Topic())
			return
		}
		handler(client, msg)
	}
	c.subscriptions[topic] = subscription{wrappedHandler, cancel}
	return c.client.Subscribe(topic, SubscribeQoS, wrappedHandler)
}

// SubscribeTopic subscribes to telemetry messages
func (c *MQTT) SubscribeTopic() (<-chan []byte, error) {
	messages := make(chan []byte, BufferSize)
	token := c.subscribe(c.cfg.Topic, func(_ paho.Client, msg paho.Message) {
		messages <- msg.Payload()
	}, func() {
		close(messages)
	})
	token.Wait()
	return messages, token.Error()
}

// UnsubscribeTopic unsubscribes from topic messages
func (c *MQTT) UnsubscribeTopic() error {
	token := c.unsubscribe(c.cfg.Topic)
	token.Wait()
	return token.Error()
}
