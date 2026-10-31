package main

import (
	"errors"
	"log"
	"log/slog"

	"github.com/gazizvr/support-bot/internal/bot"
	"github.com/gazizvr/support-bot/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		if errors.Is(err, config.ErrTargetChatEmpty) {
			slog.Warn("TARGET_CHAT_EMPTY")
		} else {
			log.Fatal(err)
		}
	}
	bot, err := bot.NewBot(cfg.Token, cfg.TargetChat)
	if err != nil {
		log.Fatal(err)
	}
	if err := bot.Run(); err != nil {
		log.Fatal(err)
	}
}
