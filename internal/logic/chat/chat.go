package chat

import (
	"context"
	"web-chat/api/http_model/chat"
)

type Logic interface {
	ResponseStream(ctx context.Context, req *chat.Completion, userID string) (chat.MessageSteam, string, error)
	CreateConversation(ctx context.Context, req *chat.CreateConversationReq, userID string) (*chat.CreateConversationResp, error)
	ListConversations(ctx context.Context, req *chat.ListConversationsReq, userID string) (*chat.ListConversationsResp, error)
	ListMessages(ctx context.Context, req *chat.ListMessagesReq, userID string) (*chat.ListMessagesResp, error)
	GetConversation(ctx context.Context, req *chat.GetConversationReq, userID string) (*chat.ConversationItem, error)
	UpdateConversationTitle(ctx context.Context, req *chat.UpdateConversationTitleReq, userID string) error
	DeleteConversation(ctx context.Context, req *chat.DeleteConversationReq, userID string) error
	ClearMessages(ctx context.Context, req *chat.ClearMessagesReq, userID string) error
	PullModules(ctx context.Context) (*chat.ModelListResp, error)
}
