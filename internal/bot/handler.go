package bot

import (
	"context"
	"errors"
	"fmt"

	"github.com/mymmrac/telego"
)

const startText string = "🗣 Ask any questions and we'll try to help you!"

func (b *Bot) handleUpdate(ctx context.Context, upd telego.Update) error {
	if upd.Message != nil {
		chatID := upd.Message.Chat.ID
		if b.targetChatID == nil {
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
