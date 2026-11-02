package bot

import (
	"context"
	"errors"
	"fmt"

	"github.com/gazizvr/support-bot/internal/domain"
	"github.com/mymmrac/telego"
)

type TopicStore interface {
	Get(ctx context.Context, id int64) (int, error)
	Set(ctx context.Context, id int64, topicID int) error
}

type UserStore interface {
	Get(ctx context.Context, id int) (int64, error)
	Set(ctx context.Context, id int, userID int64) error
}

type Handler struct {
	bot          *telego.Bot
	targetChatID telego.ChatID
	topicStore   TopicStore
	userStore    UserStore
}

func NewHandler(
	bot *telego.Bot,
	targetChatID telego.ChatID,
	topicStore TopicStore,
	userStore UserStore,
) *Handler {
	return &Handler{
		bot:          bot,
		targetChatID: targetChatID,
		topicStore:   topicStore,
		userStore:    userStore,
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
	topicID, err := h.topicStore.Get(ctx, chatID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
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
			topicID = topic.MessageThreadID
			if err := h.topicStore.Set(ctx, chatID, topicID); err != nil {
				return err
			}
			if err := h.userStore.Set(ctx, topicID, chatID); err != nil {
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
			MessageThreadID: topicID,
			MessageID:       u.Message.MessageID,
			FromChatID:      u.Message.Chat.ChatID(),
		},
	); err != nil {
		return err
	}
	return nil
}

func (h *Handler) OnSupportMessage(ctx context.Context, u telego.Update) error {
	userID, err := h.userStore.Get(ctx, u.Message.MessageThreadID)
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
