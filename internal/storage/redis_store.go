package storage

import (
	"context"
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

const userIDKeyPrefix = "user:id"
const threadIDKeyPrefix = "thread:id"

var getScript *redis.Script = redis.NewScript(`
	local id = ARGV[1]

	local userIDKeyPrefix = KEYS[1]
	local userIDKey = userIDKeyPrefix .. id
	local userID = tonumber(redis.call('get',userIDKey))

	local threadIDKeyPrefix = KEYS[2]
	local threadIDKey = threadIDKeyPrefix .. userID
	local threadID = tonumber(redis.call('get',threadIDKey))

	return {id, userID, threadID}	
`)

func (s *Store) Get(ctx context.Context, id int64) (*domain.Ticket, error) {
	intr, err := getScript.Run(
		ctx,
		s.client,
		[]string{userIDKeyPrefix, threadIDKeyPrefix},
		id,
	).Result()
	if err != nil {
		return nil, err
	}
	ticket, ok := intr.(domain.Ticket)
	if !ok {
		return nil, domain.ErrTicketNotFound
	}
	return &ticket, nil
}

const seqKey string = "ticket:id:seq"

var setScript *redis.Script = redis.NewScript(`
	local seqKey = KEYS[1]
	local id = redis.call('incr',seqKey)

	local userIDKeyPrefix = KEYS[2]
	local userIDKey = userIDKeyPrefix .. id
	local userID = ARGV[1]
	redis.call('set',userIDKey,userID)

	local threadIDKeyPrefix = KEYS[3]
	local threadIDKey = threadIDKeyPrefix .. userID
	local threadID = ARGV[2]
	redis.call('set',threadIDKey,threadID)
`)

func (s *Store) Create(ctx context.Context, userID int64, threadID int) error {
	return setScript.Run(ctx, s.client, []string{seqKey, userIDKeyPrefix, threadIDKeyPrefix}, userID, threadID).Err()
}

func (s *Store) UserIDByThreadID(ctx context.Context, threadID int) (int64, error) {
	return s.client.Get(ctx, fmt.Sprintf("%s:%d", threadIDKeyPrefix, threadID)).Int64()
}
