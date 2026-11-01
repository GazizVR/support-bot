package bot

import (
	"context"
	"errors"
	"log/slog"
	"strconv"

	"github.com/mymmrac/telego"
)

type Bot struct {
	token        string
	targetChatID *int64
	client       *telego.Bot
}

func NewBot(token, targetChat string) (*Bot, error) {
	bot, err := telego.NewBot(token)
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
	handler := NewHandler(b.client, telego.ChatID{ID: *b.targetChatID}, nil)
	routes := []Route{
		{
			Match: func(u telego.Update) bool {
				if u.Message != nil {
					if b.targetChatID != nil {
						if u.Message.Chat.ID != *b.targetChatID {
							if u.Message.Text == "/start" {
								return true
							}
						}
					}
				}
				return false
			},
			Handle: handler.OnStart,
		},
		{
			Match: func(u telego.Update) bool {
				if u.Message != nil {
					if b.targetChatID == nil {
						text := u.Message.Text
						if text == "/id" || text == "/chatid" {
							return true
						}
					}
				}
				return false
			},
			Handle: handler.OnGetChatID,
		},
		{
			Match: func(u telego.Update) bool {
				if u.Message != nil {
					if b.targetChatID != nil {
						if u.Message.Chat.ID != *b.targetChatID {
							return true
						}
					}
				}
				return false
			},
			Handle: handler.OnUserMessage,
		},
		{
			Match: func(u telego.Update) bool {
				if u.Message != nil {
					if b.targetChatID != nil {
						if u.Message.Chat.ID == *b.targetChatID {
							return true
						}
					}
				}
				return false
			},
			Handle: handler.OnSupportMessage,
		},
	}
	router := NewRouter(routes)
	for upd := range updChan {
		go func() {
			if err := router.Route(ctx, upd); err != nil {
				slog.ErrorContext(ctx, err.Error())
			}
		}()
	}
	return nil
}
