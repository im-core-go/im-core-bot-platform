package chat

type ModelInfo struct {
	ID        string
	CreatedAt int64
}

type ModelListResp struct {
	Data []ModelInfo
}

type Message struct {
	Role        string
	ContentType string
	Content     string
	Meta        string
}

type Completion struct {
	ConversationID string
	Model          string
	Messages       []Message
	Stream         bool
}

type CreateConversationReq struct {
	Model   string
	Message Message
}

type CreateConversationResp struct {
	ConversationID string
	Title          string
	Reply          Message
}

type ListConversationsReq struct {
	Page     int
	PageSize int
}

type ConversationItem struct {
	ConversationID string
	Title          string
	CreatedAt      int64
	UpdatedAt      int64
}

type ListConversationsResp struct {
	Total    int64
	Page     int
	PageSize int
	Items    []ConversationItem
}

type ListMessagesReq struct {
	ConversationID string
	Page           int
	PageSize       int
}

type MessageItem struct {
	ID          int64
	Sequence    int64
	Role        string
	ContentType string
	Content     string
	Meta        string
	IsSummary   bool
	CreatedAt   int64
}

type ListMessagesResp struct {
	Total    int64
	Page     int
	PageSize int
	Items    []MessageItem
}

type GetConversationReq struct {
	ConversationID string
}

type UpdateConversationTitleReq struct {
	ConversationID string
	Title          string
}

type DeleteConversationReq struct {
	ConversationID string
}

type ClearMessagesReq struct {
	ConversationID string
}

type StreamEventType string

const (
	EventTextDelta StreamEventType = "text.delta"
	EventImage     StreamEventType = "image"
	EventDone      StreamEventType = "done"
	EventError     StreamEventType = "error"
)

type StreamEvent struct {
	Type           StreamEventType
	Delta          string
	ConversationID string
	Title          string
}

type MessageStream interface {
	Next() (StreamEvent, bool, error)
	Close() error
}
