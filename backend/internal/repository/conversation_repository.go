package repository

import (
	"backend/internal/model"
	"backend/pkg/logger"
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ConversationRepository interface {
	Create(ctx context.Context, conversation *model.Conversation) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.Conversation, error)
	FindByPostAndUsers(ctx context.Context, postID, authorID, requesterID uuid.UUID) (*model.Conversation, error)
	FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.Conversation, error)
	Update(ctx context.Context, conversation *model.Conversation) error
}

type conversationRepository struct {
	db  *gorm.DB
	log logger.Logger
}

func NewConversationRepository(db *gorm.DB, log logger.Logger) ConversationRepository {
	return &conversationRepository{db: db, log: log}
}

func (r *conversationRepository) Create(ctx context.Context, conversation *model.Conversation) error {
	return r.db.WithContext(ctx).Create(conversation).Error
}

func (r *conversationRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Conversation, error) {
	var conversation model.Conversation
	err := r.db.WithContext(ctx).
		Preload("Post").
		Preload("Post.Author").
		Preload("Requester").
		Preload("Messages", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC")
		}).
		Where("id = ?", id).
		First(&conversation).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &conversation, err
}

func (r *conversationRepository) FindByPostAndUsers(ctx context.Context, postID, authorID, requesterID uuid.UUID) (*model.Conversation, error) {
	var conversation model.Conversation
	err := r.db.WithContext(ctx).
		Preload("Post").
		Preload("Post.Author").
		Preload("Requester").
		Preload("Messages", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC")
		}).
		Where("post_id = ? AND author_id = ? AND requester_id = ?", postID, authorID, requesterID).
		First(&conversation).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &conversation, err
}

func (r *conversationRepository) FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.Conversation, error) {
	var conversations []model.Conversation

	query := r.db.WithContext(ctx).
		Preload("Post").
		Preload("Post.Author").
		Preload("Requester").
		Preload("Messages", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC")
		}).
		Where("author_id = ? OR requester_id = ?", userID, userID).
		Order("updated_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&conversations).Error
	return conversations, err
}

func (r *conversationRepository) Update(ctx context.Context, conversation *model.Conversation) error {
	return r.db.WithContext(ctx).Save(conversation).Error
}
