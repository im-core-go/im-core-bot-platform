package chat

import "web-chat/internal/model"

type Dao interface {
	CreateConversation(conversation model.Conversation) error
	UpdateConversation(conversationID string, updateMap map[string]interface{}) error
	GetConversationByID(conversationID string) (*model.Conversation, error)
	ListConversationsByUser(userID string, offset, limit int) ([]model.Conversation, int64, error)
	DeleteConversation(conversationID string) error

	CreateMessage(message model.Message) error
	ListNonSummaryMessagesAfterSequence(conversationID string, afterSequence int64) ([]model.Message, error)
	ListRecentNonSummaryMessages(conversationID string, limit int) ([]model.Message, error)
	ListSummaryMessages(conversationID string, limit int) ([]model.Message, error)
	GetLastSummary(conversationID string) (*model.Message, error)
	ListMessagesByConversation(conversationID string, offset, limit int) ([]model.Message, int64, error)
	DeleteMessagesByConversation(conversationID string) error
}
