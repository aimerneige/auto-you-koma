package queue

import "context"

type TaskMessage struct {
	ID      string
	Payload []byte
}

type TaskQueue interface {
	Publish(ctx context.Context, topic string, message *TaskMessage) error
	Subscribe(ctx context.Context, topic string) (<-chan *TaskMessage, error)
	Ack(ctx context.Context, messageID string) error
}
