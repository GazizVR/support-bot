package bot

import (
	"context"
	"fmt"

	"github.com/mymmrac/telego"
)

type TicketStore interface {
	Get(ctx context.Context, id int64) (int64, error)
	Set(ctx context.Context, id int64) error
	Del(ctx context.Context, id int64) error
}

type Handler struct {
	bot          *telego.Bot
	targetChatID telego.ChatID
	store        TicketStore
}

func NewHandler(
	bot *telego.Bot,
	targetChatID telego.ChatID,
	store TicketStore,
) *Handler {
	return &Handler{
		bot:          bot,
		targetChatID: targetChatID,
		store:        store,
	}
}

const startText string = "🗣 Ask any questions and we'll try to help you!"

func (h *Handler) OnStart(ctx context.Context, u telego.Update) error {
	chatID := u.Message.Chat.ChatID()
	if _, err := h.bot.SendMessage(
		ctx,
		&telego.SendMessageParams{
			ChatID: chatID,
			Text:   startText,
		},
	); err != nil {
		return err
	}
	return nil
}

func (h *Handler) OnGetChatID(ctx context.Context, u telego.Update) error {
	chatID := u.Message.Chat.ChatID()
	if _, err := h.bot.SendMessage(
		ctx,
		&telego.SendMessageParams{
			ChatID:    chatID,
			Text:      fmt.Sprintf("ChatID: `%d`", chatID.ID),
			ParseMode: "Markdown",
		},
	); err != nil {
		return err
	}
	return nil
}

func (h *Handler) OnUserMessage(ctx context.Context, u telego.Update) error {
	topic, err := h.bot.CreateForumTopic(
		ctx,
		&telego.CreateForumTopicParams{
			ChatID: h.targetChatID,
			Name:   fmt.Sprintf("%d", u.Message.Chat.ID),
		},
	)
	if err != nil {
		return err
	}
	if _, err := h.bot.CopyMessage(
		ctx,
		&telego.CopyMessageParams{
			ChatID:          h.targetChatID,
			MessageThreadID: topic.MessageThreadID,
			MessageID:       u.Message.MessageID,
			FromChatID:      u.Message.Chat.ChatID(),
		},
	); err != nil {
		return err
	}
	return nil
}

func (h *Handler) OnSupportMessage(ctx context.Context, u telego.Update) error {
	// if _, err := h.bot.CopyMessage(
	// 	ctx,
	// 	&telego.CopyMessageParams{
	// 		ChatID:          ,
	// 		MessageThreadID: u.Message.MessageThreadID,
	// 		MessageID:       u.Message.MessageID,
	// 		FromChatID:      u.Message.Chat.ChatID(),
	// 	},
	// ); err != nil {
	// 	return err
	// }
	return nil
}
