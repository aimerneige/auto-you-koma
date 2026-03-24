package queue

import "context"

type MemoryQueue struct {
	topics map[string]chan *TaskMessage
}

func NewMemoryQueue() TaskQueue {
	return &MemoryQueue{
		topics: make(map[string]chan *TaskMessage),
	}
}

func (q *MemoryQueue) Publish(ctx context.Context, topic string, message *TaskMessage) error {
	if _, ok := q.topics[topic]; !ok {
		q.topics[topic] = make(chan *TaskMessage, 100)
	}
	q.topics[topic] <- message
	return nil
}

func (q *MemoryQueue) Subscribe(ctx context.Context, topic string) (<-chan *TaskMessage, error) {
	if _, ok := q.topics[topic]; !ok {
		q.topics[topic] = make(chan *TaskMessage, 100)
	}
	return q.topics[topic], nil
}

func (q *MemoryQueue) Ack(ctx context.Context, messageID string) error {
	// Memory queue doesn't require explicit ack mechanism
	return nil
}
