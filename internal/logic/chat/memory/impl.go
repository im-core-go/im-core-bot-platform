package memory

import (
	"context"
	"errors"
	"github.com/im-core-go/im-core-bot-platform/internal/dao/chat"
	"github.com/im-core-go/im-core-bot-platform/internal/model"
	"strings"
	"time"

	"gorm.io/gorm"
)

const (
	defaultSummaryThreshold    = 5
	defaultSummaryIncludeLimit = 3
)

type manager struct {
	dao                 chat.Dao
	newID               func() int64
	newUUID             func() string
	summaryThreshold    int
	summaryIncludeLimit int
}

func NewManager(dao chat.Dao, newID func() int64, newUUID func() string) Manager {
	return &manager{
		dao:                 dao,
		newID:               newID,
		newUUID:             newUUID,
		summaryThreshold:    defaultSummaryThreshold,
		summaryIncludeLimit: defaultSummaryIncludeLimit,
	}
}

func (m *manager) EnsureConversation(ctx context.Context, userID, conversationID string) (string, error) {
	if conversationID == "" {
		if userID == "" {
			return "", errors.New("missing user")
		}
		conversationID = m.newUUID()
		conversation := model.Conversation{
			UUID:   conversationID,
			UserID: userID,
			Title:  "New",
		}
		if err := m.dao.CreateConversation(conversation); err != nil {
			return "", err
		}
		return conversationID, nil
	}
	conversation, err := m.dao.GetConversationByID(conversationID)
	if err != nil {
		return "", err
	}
	if userID != "" && conversation.UserID != userID {
		return "", errors.New("forbidden")
	}
	return conversationID, nil
}

func (m *manager) SaveUserMessage(ctx context.Context, conversationID string, msg MessageInput) (model.Message, error) {
	role := msg.Role
	if role == "" {
		role = "user"
	}
	contentType := msg.ContentType
	if contentType == "" {
		contentType = "text"
	}
	content := strings.TrimSpace(msg.Content)
	if content == "" {
		return model.Message{}, errors.New("empty message")
	}
	var meta *string
	if strings.TrimSpace(msg.Meta) != "" {
		meta = &msg.Meta
	}
	id := m.newID()
	entity := model.Message{
		ID:             id,
		Sequence:       id,
		ConversationID: conversationID,
		Role:           role,
		ContentType:    contentType,
		Content:        content,
		Meta:           meta,
		IsSummary:      false,
	}
	if err := m.dao.CreateMessage(entity); err != nil {
		return model.Message{}, err
	}
	m.touchConversation(conversationID)
	return entity, nil
}

func (m *manager) SaveAssistantMessage(ctx context.Context, conversationID, content string) error {
	trimmed := strings.TrimSpace(content)
	if trimmed == "" {
		return nil
	}
	id := m.newID()
	entity := model.Message{
		ID:             id,
		Sequence:       id,
		ConversationID: conversationID,
		Role:           "assistant",
		ContentType:    "text",
		Content:        trimmed,
	}
	if err := m.dao.CreateMessage(entity); err != nil {
		return err
	}
	m.touchConversation(conversationID)
	return nil
}

func (m *manager) SaveSummaryMessage(ctx context.Context, conversationID, content string, fromID, toID int64) error {
	trimmed := strings.TrimSpace(content)
	if trimmed == "" {
		return nil
	}
	id := m.newID()
	entity := model.Message{
		ID:             id,
		Sequence:       id,
		ConversationID: conversationID,
		Role:           "system",
		ContentType:    "text",
		Content:        trimmed,
		IsSummary:      true,
		SummaryFromID:  fromID,
		SummaryToID:    toID,
	}
	if err := m.dao.CreateMessage(entity); err != nil {
		return err
	}
	m.touchConversation(conversationID)
	return nil
}

func (m *manager) BuildPrompt(ctx context.Context, conversationID string, latest model.Message, modelName string, summarize Summarizer) ([]PromptMessage, error) {
	lastSummary, err := m.dao.GetLastSummary(conversationID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	var afterSequence int64
	if lastSummary != nil {
		afterSequence = lastSummary.SummaryToID
	}

	messages, err := m.dao.ListNonSummaryMessagesAfterSequence(conversationID, afterSequence)
	if err != nil {
		return nil, err
	}
	if len(messages) == 0 {
		return []PromptMessage{{Role: latest.Role, Content: latest.Content}}, nil
	}

	if len(messages) > m.summaryThreshold {
		summarySource := messages[:len(messages)-1]
		lastMsg := messages[len(messages)-1]
		summaryPrompt := buildSummaryPrompt(summarySource)
		if len(summaryPrompt) > 0 {
			summaryText, sumErr := summarize(ctx, modelName, summaryPrompt)
			if sumErr == nil && strings.TrimSpace(summaryText) != "" {
				if err := m.SaveSummaryMessage(ctx, conversationID, summaryText, summarySource[0].ID, summarySource[len(summarySource)-1].ID); err != nil {
					return nil, err
				}
				return []PromptMessage{
					{Role: "system", Content: summaryText},
					{Role: lastMsg.Role, Content: lastMsg.Content},
				}, nil
			}
		}
	}

	summaries, err := m.dao.ListSummaryMessages(conversationID, m.summaryIncludeLimit)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	prompt := make([]PromptMessage, 0, len(summaries)+len(messages))
	for i := len(summaries) - 1; i >= 0; i-- {
		prompt = append(prompt, PromptMessage{Role: "system", Content: summaries[i].Content})
	}
	for _, msg := range messages {
		if msg.ContentType != "text" {
			continue
		}
		prompt = append(prompt, PromptMessage{Role: msg.Role, Content: msg.Content})
	}
	return prompt, nil
}

func (m *manager) BuildTitleMessages(ctx context.Context, conversationID string, limit int) ([]PromptMessage, error) {
	if limit <= 0 {
		limit = 4
	}
	messages, err := m.dao.ListRecentNonSummaryMessages(conversationID, limit)
	if err != nil {
		return nil, err
	}
	prompt := make([]PromptMessage, 0, len(messages))
	for i := len(messages) - 1; i >= 0; i-- {
		if messages[i].ContentType != "text" {
			continue
		}
		prompt = append(prompt, PromptMessage{
			Role:    messages[i].Role,
			Content: messages[i].Content,
		})
	}
	return prompt, nil
}

func (m *manager) GetConversation(ctx context.Context, conversationID string) (*model.Conversation, error) {
	return m.dao.GetConversationByID(conversationID)
}

func (m *manager) UpdateConversationTitle(ctx context.Context, conversationID, title string) error {
	title = strings.TrimSpace(title)
	if title == "" {
		return errors.New("empty title")
	}
	return m.dao.UpdateConversation(conversationID, map[string]interface{}{
		"title": title,
	})
}

func (m *manager) ListConversations(ctx context.Context, userID string, offset, limit int) ([]model.Conversation, int64, error) {
	return m.dao.ListConversationsByUser(userID, offset, limit)
}

func (m *manager) ListMessages(ctx context.Context, conversationID string, offset, limit int) ([]model.Message, int64, error) {
	return m.dao.ListMessagesByConversation(conversationID, offset, limit)
}

func (m *manager) DeleteConversation(ctx context.Context, conversationID string) error {
	if conversationID == "" {
		return errors.New("missing conversation_id")
	}
	return m.dao.DeleteConversation(conversationID)
}

func (m *manager) ClearMessages(ctx context.Context, conversationID string) error {
	if conversationID == "" {
		return errors.New("missing conversation_id")
	}
	return m.dao.DeleteMessagesByConversation(conversationID)
}

func buildSummaryPrompt(messages []model.Message) []PromptMessage {
	prompt := []PromptMessage{
		{
			Role:    "system",
			Content: "Summarize the conversation briefly, focusing on key facts and decisions.",
		},
	}
	for _, msg := range messages {
		if msg.ContentType != "text" {
			continue
		}
		prompt = append(prompt, PromptMessage{Role: msg.Role, Content: msg.Content})
	}
	if len(prompt) <= 1 {
		return nil
	}
	return prompt
}

func (m *manager) touchConversation(conversationID string) {
	_ = m.dao.UpdateConversation(conversationID, map[string]interface{}{
		"updated_at": time.Now().Unix(),
	})
}
