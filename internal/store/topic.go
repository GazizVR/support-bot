package store

import (
	"context"
	"sync"

	"github.com/gazizvr/support-bot/internal/domain"
)

type TopicStore struct {
	m      sync.RWMutex
	topics map[int64]int
}

func NewTopicStore() *TopicStore {
	return &TopicStore{
		m:      sync.RWMutex{},
		topics: make(map[int64]int),
	}
}

func (s *TopicStore) Get(ctx context.Context, id int64) (int, error) {
	s.m.RLock()
	defer s.m.RUnlock()
	userID, ok := s.topics[id]
	if !ok {
		return 0, domain.ErrNotFound
	}
	return userID, nil
}

func (s *TopicStore) Set(ctx context.Context, id int64, topicID int) error {
	s.m.Lock()
	s.topics[id] = topicID
	s.m.Unlock()
	return nil
}
