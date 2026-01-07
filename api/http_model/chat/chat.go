package chat

type ModelListResp struct {
	Data []struct {
		ID        string `json:"id"`
		CreatedAt int64  `json:"created"`
	} `json:"data"`
}

type Message struct {
	Role        string `json:"role"`
	ContentType string `json:"content_type,omitempty"`
	Content     string `json:"content"`
	Meta        string `json:"meta,omitempty"`
}

type Completion struct {
	ConversationID string `json:"conversation_id"`
	Model          string `json:"model"`
	Messages       []Message `json:"messages"`
	Stream bool `json:"stream"`
}

type CreateConversationReq struct {
	Model   string  `json:"model"`
	Message Message `json:"message"`
}

type CreateConversationResp struct {
	ConversationID string  `json:"conversation_id"`
	Title          string  `json:"title"`
	Reply          Message `json:"reply"`
}

type ListConversationsReq struct {
	Page     int `form:"page" binding:"omitempty,min=1"`
	PageSize int `form:"page_size" binding:"omitempty,min=1,max=100"`
}

type ConversationItem struct {
	ConversationID string `json:"conversation_id"`
	Title          string `json:"title"`
	CreatedAt      int64  `json:"created_at"`
	UpdatedAt      int64  `json:"updated_at"`
}

type ListConversationsResp struct {
	Total    int64              `json:"total"`
	Page     int                `json:"page"`
	PageSize int                `json:"page_size"`
	Items    []ConversationItem `json:"items"`
}

type ListMessagesReq struct {
	ConversationID string `form:"conversation_id" binding:"required"`
	Page           int    `form:"page" binding:"omitempty,min=1"`
	PageSize       int    `form:"page_size" binding:"omitempty,min=1,max=100"`
}

type MessageItem struct {
	ID          int64  `json:"id"`
	Sequence    int64  `json:"sequence"`
	Role        string `json:"role"`
	ContentType string `json:"content_type"`
	Content     string `json:"content"`
	Meta        string `json:"meta,omitempty"`
	IsSummary   bool   `json:"is_summary"`
	CreatedAt   int64  `json:"created_at"`
}

type ListMessagesResp struct {
	Total    int64         `json:"total"`
	Page     int           `json:"page"`
	PageSize int           `json:"page_size"`
	Items    []MessageItem `json:"items"`
}

type GetConversationReq struct {
	ConversationID string `form:"conversation_id" binding:"required"`
}

type UpdateConversationTitleReq struct {
	ConversationID string `json:"conversation_id" binding:"required"`
	Title          string `json:"title" binding:"required"`
}

type DeleteConversationReq struct {
	ConversationID string `form:"conversation_id" binding:"required"`
}

type ClearMessagesReq struct {
	ConversationID string `form:"conversation_id" binding:"required"`
}
