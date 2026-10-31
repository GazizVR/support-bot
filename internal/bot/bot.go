package bot

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/mymmrac/telego"
)

type Bot struct {
	token        string
	targetChatID *int64
	client       *telego.Bot
}

func NewBot(token, targetChat string) (*Bot, error) {
	logger, err := NewFileLogger("bot.log")
	if err != nil {
		return nil, err
	}
	bot, err := telego.NewBot(token, telego.WithLogger(logger))
	if err != nil {
		return nil, err
	}
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
		return &Bot{
			token:        token,
			targetChatID: &chatID,
			client:       bot,
		}, nil
	}
	return &Bot{token: token, client: bot}, nil
}

const startText string = "🗣 Ask any questions and we'll try to help you!"

func (b *Bot) handleUpdate(ctx context.Context, upd telego.Update) error {
	if upd.Message != nil {
		chatID := upd.Message.Chat.ID
		switch upd.Message.Text {
		case "/id", "/chatid":
			if _, err := b.client.SendMessage(
				ctx,
				&telego.SendMessageParams{
					ChatID:    telego.ChatID{ID: chatID},
					Text:      fmt.Sprintf("ChatID: `%d`", chatID),
					ParseMode: "Markdown",
				},
			); err != nil {
				return err
			}
			return nil
		}
		if b.targetChatID == nil {
			return ErrTargetChatIDEmpty
		}
		targetChatID := *b.targetChatID
		if chatID == targetChatID {
			replyMsg := upd.Message.ReplyToMessage
			if replyMsg != nil {
				var originChatID int64
				switch v := replyMsg.ForwardOrigin.(type) {
				case *telego.MessageOriginUser:
					originChatID = v.SenderUser.ID
				case *telego.MessageOriginChat:
					originChatID = v.SenderChat.ID
				case *telego.MessageOriginChannel:
					originChatID = v.Chat.ID
				default:
					b.client.SendMessage(
						ctx,
						&telego.SendMessageParams{
							ChatID: telego.ChatID{ID: chatID},
							Text:   "Message sender is not accessible",
							ReplyParameters: &telego.ReplyParameters{
								MessageID: upd.Message.MessageID,
							},
						},
					)
					return errors.New("unexpected MessageOrigin")
				}
				if _, err := b.client.CopyMessage(
					ctx,
					&telego.CopyMessageParams{
						ChatID:     telego.ChatID{ID: originChatID},
						FromChatID: telego.ChatID{ID: chatID},
						MessageID:  upd.Message.GetMessageID(),
					},
				); err != nil {
					return err
				}
			}
		} else {
			if upd.Message.Text == "/start" {
				if _, err := b.client.SendMessage(
					ctx,
					&telego.SendMessageParams{
						ChatID: telego.ChatID{ID: chatID},
						Text:   startText,
					},
				); err != nil {
					return err
				}
				return nil
			}
			if _, err := b.client.ForwardMessage(
				ctx,
				&telego.ForwardMessageParams{
					ChatID:     telego.ChatID{ID: targetChatID},
					FromChatID: telego.ChatID{ID: chatID},
					MessageID:  upd.Message.GetMessageID(),
				},
			); err != nil {
				return err
			}
		}
	}
	return nil
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
