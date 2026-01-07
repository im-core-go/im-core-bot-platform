package memory

import (
	"context"
	"web-chat/internal/model"
)

type PromptMessage struct {
	Role    string
	Content string
}

type Summarizer func(ctx context.Context, modelName string, messages []PromptMessage) (string, error)

type MessageInput struct {
	Role        string
	ContentType string
	Content     string
	Meta        string
}

type Manager interface {
	EnsureConversation(ctx context.Context, userID, conversationID string) (string, error)
	SaveUserMessage(ctx context.Context, conversationID string, msg MessageInput) (model.Message, error)
	SaveAssistantMessage(ctx context.Context, conversationID, content string) error
	SaveSummaryMessage(ctx context.Context, conversationID, content string, fromID, toID int64) error
	BuildPrompt(ctx context.Context, conversationID string, latest model.Message, modelName string, summarize Summarizer) ([]PromptMessage, error)
	BuildTitleMessages(ctx context.Context, conversationID string, limit int) ([]PromptMessage, error)
	GetConversation(ctx context.Context, conversationID string) (*model.Conversation, error)
	UpdateConversationTitle(ctx context.Context, conversationID, title string) error
	ListConversations(ctx context.Context, userID string, offset, limit int) ([]model.Conversation, int64, error)
	ListMessages(ctx context.Context, conversationID string, offset, limit int) ([]model.Message, int64, error)
	DeleteConversation(ctx context.Context, conversationID string) error
	ClearMessages(ctx context.Context, conversationID string) error
}
