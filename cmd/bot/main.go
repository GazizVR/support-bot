package main

import (
	"log"
	"log/slog"

	"github.com/gazizvr/support-bot/internal/bot"
	"github.com/gazizvr/support-bot/internal/config"
	"github.com/gazizvr/support-bot/internal/storage"
	"github.com/redis/go-redis/v9"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	if cfg.TargetChat == "" {
		slog.Warn("TARGET_CHAT_EMPTY", "message", "please pass the TARGET_CHAT env varible")
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPswd,
	})
	store := storage.NewStore(rdb)
	bot, err := bot.NewBot(cfg.BotToken, cfg.TargetChat, store)
	if err != nil {
		log.Fatal(err)
	}
	if err := bot.Run(); err != nil {
		log.Fatal(err)
	}
}
