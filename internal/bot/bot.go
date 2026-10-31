package bot

import (
	"context"
	"errors"
	"strconv"

	"github.com/mymmrac/telego"
)

type Bot struct {
	token        string
	targetChatID *int64
	client       *telego.Bot
}

func NewBot(token, targetChat string) (*Bot, error) {
	bot, err := telego.NewBot(token, telego.WithDefaultDebugLogger())
	if err != nil {
		return nil, err
	}
	var targetChatID *int64
	if targetChat != "" {
		chatID, err := strconv.ParseInt(targetChat, 10, 64)
		if err != nil {
			switch {
			case errors.Is(err, strconv.ErrSyntax):
				chatInfo, err := bot.GetChat(
					context.Background(),
					&telego.GetChatParams{ChatID: telego.ChatID{Username: targetChat}},
				)
				if err != nil {
					return nil, err
				}
				chatID = chatInfo.ID
			default:
				return nil, err
			}
		}
		targetChatID = &chatID
	}
	return &Bot{
		token:        token,
		targetChatID: targetChatID,
		client:       bot,
	}, nil
}

func (b *Bot) Run() error {
	ctx := context.Background()
	updChan, err := b.client.UpdatesViaLongPolling(ctx, &telego.GetUpdatesParams{Timeout: 45})
	if err != nil {
		return err
	}
	for upd := range updChan {
		go b.handleUpdate(ctx, upd)
	}
	return nil
}
