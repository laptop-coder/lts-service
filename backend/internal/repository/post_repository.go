package repository

import (
	"backend/internal/model"
	"backend/pkg/logger"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostRepository interface {
	Create(ctx context.Context, post *model.Post) error
	FindAll(ctx context.Context, filter *PostFilter) ([]model.Post, error)
	FindByID(ctx context.Context, id *uuid.UUID) (*model.Post, error)
	Update(ctx context.Context, post *model.Post) error
	Delete(ctx context.Context, id *uuid.UUID) error
}

type postRepository struct {
	db  *gorm.DB
	log logger.Logger
}

type PostFilter struct {
	AuthorID             *uuid.UUID
	Verified             *bool
	ThingReturnedToOwner *bool
	Limit                int
	Offset               int
}

func NewPostRepository(db *gorm.DB, log logger.Logger) PostRepository {
	if db == nil {
		log.Error("DB is nil")
		panic("DB is nil")
	}
	return &postRepository{db: db, log: log}
}

func (r *postRepository) FindAll(ctx context.Context, filter *PostFilter) ([]model.Post, error) {
	if filter == nil {
		return nil, fmt.Errorf("posts list filter cannot be nil")
	}
	var posts []model.Post
	query := r.db.WithContext(ctx).Model(&model.Post{})
	// Filters
	// by post's author:
	if filter.AuthorID != nil {
		query = query.
			Where("posts.author_id = ?", *filter.AuthorID)
	}
	// by verification status:
	if filter.Verified != nil {
		query = query.
			Where("posts.verified = ?", *filter.Verified)
	}
	// by thing return status:
	if filter.ThingReturnedToOwner != nil {
		query = query.
			Where("posts.thing_returned_to_owner = ?", *filter.ThingReturnedToOwner)
	}
	// offset (for pagination):
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}
	// limit (for pagination):
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	// Sort posts by name in the alphabetical order
	query = query.Order("name")
	// Find posts
	result := query.Find(&posts)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch posts list: %w", result.Error)
	}
	// Return response
	return posts, nil
}

func (r *postRepository) FindByID(ctx context.Context, id *uuid.UUID) (*model.Post, error) {
	if id == nil {
		return nil, fmt.Errorf("post id cannot be nil")
	}
	var post model.Post
	result := r.db.WithContext(ctx).First(&post, *id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("post with id %s was not found: %w", *id, result.Error)
		}
		return nil, fmt.Errorf("failed to fetch post by id (%s): %w", *id, result.Error)
	}
	return &post, nil
}

func (r *postRepository) Create(ctx context.Context, post *model.Post) error {
	if post == nil {
		return fmt.Errorf("post cannot be nil")
	}
	result := r.db.WithContext(ctx).Create(post)
	if result.Error != nil {
		return fmt.Errorf("failed to create new post: %w", result.Error)
	}
	return nil
}

func (r *postRepository) Update(ctx context.Context, post *model.Post) error {
	if post == nil {
		return fmt.Errorf("post cannot be nil")
	}
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.Post{}).
		Where("id = ?", post.ID).
		Count(&count).Error
	if err != nil {
		return fmt.Errorf("failed to check post existence: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("post with id %d was not found", post.ID)
	}
	result := r.db.WithContext(ctx).Save(post)
	if result.Error != nil {
		return fmt.Errorf("failed to update post: %w", result.Error)
	}
	return nil
}

func (r *postRepository) Delete(ctx context.Context, id *uuid.UUID) error {
	result := r.db.WithContext(ctx).Unscoped().Delete(&model.Post{}, *id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete post with id %s: %w", *id, result.Error)
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
