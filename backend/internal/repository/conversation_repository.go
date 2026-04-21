package repository

import (
	"fmt"
	"backend/internal/model"
	"backend/pkg/logger"
	"backend/pkg/apperrors"
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
		return nil, fmt.Errorf("conversation not found by id: %s: %w", err.Error(), apperrors.ErrConversationNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find conversation by id: %w", err)
	}
	return &conversation, nil
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
		return nil, fmt.Errorf("conversation not found by post and users: %s: %w", err.Error(), apperrors.ErrConversationNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find conversation by post and users: %w", err)
	}
	return &conversation, nil
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
	if err != nil {
		return nil, fmt.Errorf("failed to find conversations by user id: %w", err)
	}
	return conversations, nil
}

func (r *conversationRepository) Update(ctx context.Context, conversation *model.Conversation) error {
	return r.db.WithContext(ctx).Save(conversation).Error
}
