package queue

import (
	"context"
	"sync"
)

type channelQueue struct {
	mu     sync.RWMutex
	topics map[string]chan *TaskMessage
}

func NewChannelQueue() TaskQueue {
	return &channelQueue{
		topics: make(map[string]chan *TaskMessage),
	}
}

func (q *channelQueue) getTopic(topic string) chan *TaskMessage {
	q.mu.Lock()
	defer q.mu.Unlock()
	ch, exists := q.topics[topic]
	if !exists {
		ch = make(chan *TaskMessage, 100)
		q.topics[topic] = ch
	}
	return ch
}

func (q *channelQueue) Publish(ctx context.Context, topic string, message *TaskMessage) error {
	ch := q.getTopic(topic)
	select {
	case ch <- message:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (q *channelQueue) Subscribe(ctx context.Context, topic string) (<-chan *TaskMessage, error) {
	return q.getTopic(topic), nil
}

func (q *channelQueue) Ack(ctx context.Context, messageID string) error {
	// Simple channel queue does not implement manual ACK logic yet
	return nil
}
