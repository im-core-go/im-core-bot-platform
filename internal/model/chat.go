package model

type Conversation struct {
	UUID   string `gorm:"primaryKey;type:varchar(36)"`
	UserID string `gorm:"column:user_id;index;type:varchar(36)"`
	Title  string `gorm:"column:title;type:varchar(255);default:'New'"`
	CommonPartNoUnique
}
type Message struct {
	ID             int64  `gorm:"primaryKey"`
	Sequence       int64  `gorm:"column:sequence;index"`
	ConversationID string `gorm:"column:conversation_id;index;type:varchar(36)"`
	Role           string `gorm:"column:role;type:varchar(32)"`
	ContentType    string `gorm:"column:content_type;type:varchar(32)"`
	Content        string `gorm:"column:content;type:text"`
	Meta           *string `gorm:"column:meta;type:json"`
	IsSummary      bool   `gorm:"column:is_summary;index"`
	SummaryFromID  int64  `gorm:"column:summary_from_id;index"`
	SummaryToID    int64  `gorm:"column:summary_to_id;index"`
	CommonPartNoUnique
}

func (Message) TableName() string      { return "message" }
func (Conversation) TableName() string { return "conversation" }
