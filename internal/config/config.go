package config

import (
	"os"
)

type Config struct {
	BotToken   string
	TargetChat string
	RedisAddr  string
	RedisPswd  string
}

func Load() (*Config, error) {
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		return nil, ErrBotTokenEmpty
	}
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		return nil, ErrRedisAddrEmpty
	}
	return &Config{
		BotToken:   botToken,
		TargetChat: os.Getenv("TARGET_CHAT"),
		RedisAddr:  redisAddr,
		RedisPswd:  os.Getenv("REDIS_PASSWORD"),
	}, nil
}
