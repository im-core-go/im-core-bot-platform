package chat

import (
	"github.com/im-core-go/im-core-bot-platform/internal/model"

	"gorm.io/gorm"
)

type chatDaoImpl struct {
	db *gorm.DB
}

func NewDao(db *gorm.DB) Dao {
	return &chatDaoImpl{db: db}
}

func (c *chatDaoImpl) CreateConversation(conversation model.Conversation) error {
	return c.db.Create(&conversation).Error
}

func (c *chatDaoImpl) UpdateConversation(conversationID string, updateMap map[string]interface{}) error {
	return c.db.Model(&model.Conversation{}).Where("uuid = ?", conversationID).Updates(updateMap).Error
}

func (c *chatDaoImpl) GetConversationByID(conversationID string) (*model.Conversation, error) {
	var entity model.Conversation
	if err := c.db.Where("uuid = ?", conversationID).First(&entity).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (c *chatDaoImpl) DeleteConversation(conversationID string) error {
	return c.db.Where("uuid = ?", conversationID).Delete(&model.Conversation{}).Error
}

func (c *chatDaoImpl) ListConversationsByUser(userID string, offset, limit int) ([]model.Conversation, int64, error) {
	var (
		items []model.Conversation
		total int64
	)
	query := c.db.Model(&model.Conversation{}).Where("user_id = ?", userID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	if err := query.Order("updated_at desc").Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (c *chatDaoImpl) CreateMessage(message model.Message) error {
	return c.db.Create(&message).Error
}

func (c *chatDaoImpl) ListNonSummaryMessagesAfterSequence(conversationID string, afterSequence int64) ([]model.Message, error) {
	var messages []model.Message
	query := c.db.Where("conversation_id = ? AND is_summary = ?", conversationID, false)
	if afterSequence > 0 {
		query = query.Where("sequence > ?", afterSequence)
	}
	if err := query.Order("sequence asc").Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

func (c *chatDaoImpl) ListRecentNonSummaryMessages(conversationID string, limit int) ([]model.Message, error) {
	var messages []model.Message
	query := c.db.Where("conversation_id = ? AND is_summary = ?", conversationID, false)
	if limit > 0 {
		query = query.Limit(limit)
	}
	if err := query.Order("sequence desc").Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

func (c *chatDaoImpl) ListSummaryMessages(conversationID string, limit int) ([]model.Message, error) {
	var messages []model.Message
	query := c.db.Where("conversation_id = ? AND is_summary = ?", conversationID, true).Order("sequence desc")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if err := query.Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

func (c *chatDaoImpl) GetLastSummary(conversationID string) (*model.Message, error) {
	var message model.Message
	err := c.db.Where("conversation_id = ? AND is_summary = ?", conversationID, true).
		Order("sequence desc").
		Limit(1).
		First(&message).Error
	if err != nil {
		return nil, err
	}
	return &message, nil
}

func (c *chatDaoImpl) ListMessagesByConversation(conversationID string, offset, limit int) ([]model.Message, int64, error) {
	var (
		items []model.Message
		total int64
	)
	query := c.db.Model(&model.Message{}).Where("conversation_id = ?", conversationID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	if err := query.Order("sequence desc").Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (c *chatDaoImpl) DeleteMessagesByConversation(conversationID string) error {
	return c.db.Where("conversation_id = ?", conversationID).Delete(&model.Message{}).Error
}
