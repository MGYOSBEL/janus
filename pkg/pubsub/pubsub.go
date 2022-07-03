package pubsub

import "context"

type Publisher interface {
	Publish(topic string, message []byte) error
	Close() error
}

type Subscriber interface {
	Subscribe(ctx context.Context, topic string) (chan []byte, error)
	Close() error
}
