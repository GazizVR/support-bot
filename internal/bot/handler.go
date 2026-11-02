package bot

import (
	"context"
	"errors"
	"fmt"

	"github.com/gazizvr/support-bot/internal/domain"
	"github.com/mymmrac/telego"
)

type Store interface {
	Get(ctx context.Context, id int64) (*domain.Ticket, error)
	UserIDByThreadID(ctx context.Context, threadID int) (int64, error)
	Create(ctx context.Context, userID int64, threadID int) error
}

type Handler struct {
	bot          *telego.Bot
	targetChatID telego.ChatID
	store        Store
}

func NewHandler(
	bot *telego.Bot,
	targetChatID telego.ChatID,
	store Store,
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
	chatID := u.Message.Chat.ID
	ticket, err := h.store.Get(ctx, chatID)
	if err != nil {
		if errors.Is(err, domain.ErrTicketNotFound) {
			topic, err := h.bot.CreateForumTopic(
				ctx,
				&telego.CreateForumTopicParams{
					ChatID: h.targetChatID,
					Name:   u.Message.From.FirstName,
				},
			)
			if err != nil {
				return err
			}
			ticket = &domain.Ticket{
				UserID:   chatID,
				ThreadID: topic.MessageThreadID,
			}
			if err := h.store.Create(ctx, chatID, topic.MessageThreadID); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	if _, err := h.bot.CopyMessage(
		ctx,
		&telego.CopyMessageParams{
			ChatID:          h.targetChatID,
			MessageThreadID: ticket.ThreadID,
			MessageID:       u.Message.MessageID,
			FromChatID:      telego.ChatID{ID: ticket.UserID},
		},
	); err != nil {
		return err
	}
	return nil
}

func (h *Handler) OnSupportMessage(ctx context.Context, u telego.Update) error {
	userID, err := h.store.UserIDByThreadID(ctx, u.Message.MessageThreadID)
	if err != nil {
		return err
	}
	if _, err := h.bot.CopyMessage(
		ctx,
		&telego.CopyMessageParams{
			ChatID:     telego.ChatID{ID: userID},
			MessageID:  u.Message.MessageID,
			FromChatID: u.Message.Chat.ChatID(),
		},
	); err != nil {
		return err
	}
	return nil
}
