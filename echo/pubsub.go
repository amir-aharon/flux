package echo

import "sync"

type Topic string

type Broker[T any] struct {
	topics       map[Topic]map[*Subscription[T]]chan T
	mu           sync.RWMutex
	chanCapacity int
}

func NewBroker[T any](capacity int) *Broker[T] {
	return &Broker[T]{
		topics:       make(map[Topic]map[*Subscription[T]]chan T),
		chanCapacity: capacity,
	}
}

func (b *Broker[T]) Subscribe(topic Topic) *Subscription[T] {
	b.mu.Lock()
	defer b.mu.Unlock()

	if _, topicExists := b.topics[topic]; !topicExists {
		b.topics[topic] = make(map[*Subscription[T]]chan T)
	}

	subChan := make(chan T, b.chanCapacity)
	sub := &Subscription[T]{
		C:      subChan,
		topic:  topic,
		broker: b,
	}

	b.topics[topic][sub] = subChan

	return sub
}

func (b *Broker[T]) Publish(topic Topic, msg T) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	topicMap := b.topics[topic]
	for _, ch := range topicMap {
		ch <- msg
	}
}

type Subscription[T any] struct {
	C      <-chan T
	topic  Topic
	broker *Broker[T]
}

func (s *Subscription[T]) Unsubscribe() {
	s.broker.mu.Lock()
	defer s.broker.mu.Unlock()

	if _, exists := s.broker.topics[s.topic][s]; !exists {
		return
	}
	close(s.broker.topics[s.topic][s])
	delete(s.broker.topics[s.topic], s)
	if len(s.broker.topics[s.topic]) == 0 {
		delete(s.broker.topics, s.topic)
	}
}
