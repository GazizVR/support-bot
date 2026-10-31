package config

import (
	"os"
)

type Config struct {
	Token      string
	TargetChat string
}

func Load() (*Config, error) {
	token := os.Getenv("TOKEN")
	cfg := &Config{}
	if token == "" {
		return nil, ErrTokenEmpty
	}
	cfg.Token = token
	targetChat := os.Getenv("TARGET_CHAT")
	if targetChat == "" {
		return cfg, ErrTargetChatEmpty
	}
	cfg.TargetChat = targetChat
	return cfg, nil
}
