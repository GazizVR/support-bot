package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/gazizvr/support-bot/internal/domain"
	"github.com/redis/go-redis/v9"
)

type Store struct {
	client *redis.Client
}

func NewStore(client *redis.Client) *Store {
	return &Store{
		client: client,
	}
}

const threadIDKeyPrefix string = "thread:id:"

func (s *Store) GetThreadID(ctx context.Context, userID int64) (int, error) {
	threadID, err := s.client.Get(ctx, fmt.Sprintf("%s%d", threadIDKeyPrefix, userID)).Int()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, domain.ErrThreadIDNotFound
		} else {
			return 0, err
		}
	}
	return threadID, nil
}

const userIDKeyPrefix string = "user:id:"

func (s *Store) GetUserID(ctx context.Context, threadID int) (int64, error) {
	userID, err := s.client.Get(ctx, fmt.Sprintf("%s%d", userIDKeyPrefix, threadID)).Int64()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, domain.ErrUserIDNotFound
		} else {
			return 0, err
		}
	}
	return userID, nil
}

var setScript *redis.Script = redis.NewScript(`
	local userIDKeyPrefix = KEYS[1]
    local threadIDKeyPrefix = KEYS[2]
    
	local userID = ARGV[1]
    local threadID = ARGV[2]
    
    local userIDKey = userIDKeyPrefix .. threadID
	redis.call('set', userIDKey, userID)
	
	local threadIDKey = threadIDKeyPrefix .. userID
	redis.call('set', threadIDKey, threadID)
`)

func (s *Store) Create(ctx context.Context, userID int64, threadID int) error {
	err := setScript.Run(ctx, s.client, []string{userIDKeyPrefix, threadIDKeyPrefix}, userID, threadID).Err()
	if err != nil && !errors.Is(err, redis.Nil) {
		return err
	}
	return nil
}
