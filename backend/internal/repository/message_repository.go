package repository

import (
	"backend/pkg/logger"
    "context"
    "github.com/google/uuid"
    "gorm.io/gorm"
    "backend/internal/model"
)

type MessageRepository interface {
    Create(ctx context.Context, message *model.Message) error
    FindByConversationID(ctx context.Context, conversationID uuid.UUID, limit, offset int) ([]model.Message, error)
    MarkAsRead(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID) error
    CountUnread(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID) (int64, error)
    FindLastMessage(ctx context.Context, conversationID uuid.UUID) (*model.Message, error)
}

type messageRepository struct {
    db *gorm.DB
	log logger.Logger
}

func NewMessageRepository(db *gorm.DB, log logger.Logger) MessageRepository {
	if db == nil {
		log.Error("DB is nil")
		panic("DB is nil")
	}
	return &messageRepository{db: db, log: log}
}

func (r *messageRepository) Create(ctx context.Context, message *model.Message) error {
    return r.db.WithContext(ctx).Create(message).Error
}

func (r *messageRepository) FindByConversationID(ctx context.Context, conversationID uuid.UUID, limit, offset int) ([]model.Message, error) {
    var messages []model.Message
    
    query := r.db.WithContext(ctx).
        Preload("Sender").
        Where("conversation_id = ?", conversationID).
        Order("created_at ASC")
    
    if limit > 0 {
        query = query.Limit(limit)
    }
    if offset > 0 {
        query = query.Offset(offset)
    }
    
    err := query.Find(&messages).Error
    return messages, err
}

func (r *messageRepository) MarkAsRead(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID) error {
    return r.db.WithContext(ctx).
        Model(&model.Message{}).
        Where("conversation_id = ? AND sender_id != ? AND is_read = false", conversationID, userID).
        Update("is_read", true).Error
}

func (r *messageRepository) CountUnread(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID) (int64, error) {
    var count int64
    err := r.db.WithContext(ctx).
        Model(&model.Message{}).
        Where("conversation_id = ? AND sender_id != ? AND is_read = false", conversationID, userID).
        Count(&count).Error
    return count, err
}

func (r *messageRepository) FindLastMessage(ctx context.Context, conversationID uuid.UUID) (*model.Message, error) {
    var message model.Message
    err := r.db.WithContext(ctx).
        Preload("Sender").
        Where("conversation_id = ?", conversationID).
        Order("created_at DESC").
        First(&message).Error
    
    if err == gorm.ErrRecordNotFound {
        return nil, nil
    }
    return &message, err
}
