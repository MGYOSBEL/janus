package service

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/MGYOSBEL/janus/backend/mqtt"
	"github.com/MGYOSBEL/janus/pkg/log"
)

func Run() {
	logger := log.New("ZapLogger", "debug")

	logger.Debugf("Hola Service")

	mqttCfg := mqtt.Config{
		Brokers: []string{
			"localhost:1883",
		},
		Username:  "",
		Password:  "",
		Topic:     "#",
		BatchSize: 1,
	}
	client, err := mqtt.New(mqttCfg, logger)
	if err != nil {
		logger.Errorf("error creating MQTT client")
		return
	}
	defer client.Disconnect()
	// MQTT connect
	err = client.Connect()
	if err != nil {
		logger.Fatalf("mqtt backend connect failed %s", err.Error())
	}

	mqttChannel, err := client.SubscribeTopic()
	if err != nil {
		logger.Errorf("error subscribing", err)
		return
	}

	logger.Infof("receiving data from MQTT")
	go func() {
		for data := range mqttChannel {
			logger.Debugf(string(data))
		}
	}()

	WaitSignal()

}

// WaitSignal catching exit signal
func WaitSignal() os.Signal {
	ch := make(chan os.Signal, 2)
	signal.Notify(
		ch,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)
	for {
		sig := <-ch
		switch sig {
		case syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM:
			return sig
		}
	}
}
