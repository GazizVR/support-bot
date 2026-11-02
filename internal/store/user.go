package store

import (
	"context"
	"sync"

	"github.com/gazizvr/support-bot/internal/domain"
)

type UserStore struct {
	m     sync.RWMutex
	users map[int]int64
}

func NewUserStore() *UserStore {
	return &UserStore{
		m:     sync.RWMutex{},
		users: make(map[int]int64),
	}
}

func (s *UserStore) Get(ctx context.Context, id int) (int64, error) {
	s.m.RLock()
	defer s.m.RUnlock()
	userID, ok := s.users[id]
	if !ok {
		return 0, domain.ErrNotFound
	}
	return userID, nil
}

func (s *UserStore) Set(ctx context.Context, id int, userID int64) error {
	s.m.Lock()
	s.users[id] = userID
	s.m.Unlock()
	return nil
}
