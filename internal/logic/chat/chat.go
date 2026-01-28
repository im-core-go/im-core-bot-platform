package chat

import "context"

type Logic interface {
	ResponseStream(ctx context.Context, req *Completion, userID string) (MessageStream, string, error)
	CreateConversation(ctx context.Context, req *CreateConversationReq, userID string) (*CreateConversationResp, error)
	ListConversations(ctx context.Context, req *ListConversationsReq, userID string) (*ListConversationsResp, error)
	ListMessages(ctx context.Context, req *ListMessagesReq, userID string) (*ListMessagesResp, error)
	GetConversation(ctx context.Context, req *GetConversationReq, userID string) (*ConversationItem, error)
	UpdateConversationTitle(ctx context.Context, req *UpdateConversationTitleReq, userID string) error
	DeleteConversation(ctx context.Context, req *DeleteConversationReq, userID string) error
	ClearMessages(ctx context.Context, req *ClearMessagesReq, userID string) error
	PullModules(ctx context.Context) (*ModelListResp, error)
	BuildUserSystemPrompt(ctx context.Context, userID string) (string, error)
}
