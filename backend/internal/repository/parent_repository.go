package repository

import (
	"backend/internal/model"
	"backend/pkg/apperrors"
	"backend/pkg/logger"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ParentRepository interface {
	FindByID(ctx context.Context, userID *uuid.UUID) (*model.Parent, error)
	FindAll(ctx context.Context, filter *ParentFilter) ([]model.Parent, error)
}

type ParentFilter struct {
	Limit  int
	Offset int
}

type parentRepository struct {
	db  *gorm.DB
	log logger.Logger
}

func NewParentRepository(db *gorm.DB, log logger.Logger) ParentRepository {
	if db == nil {
		log.Error("DB is nil")
		panic("DB is nil")
	}
	return &parentRepository{db: db, log: log}
}

func (r *parentRepository) FindByID(ctx context.Context, userID *uuid.UUID) (*model.Parent, error) {
	if userID == nil {
		return nil, fmt.Errorf("user id cannot be nil: %w", apperrors.ErrRequiredField)
	}
	var parent model.Parent
	result := r.db.WithContext(ctx).
		Preload("Students").
		Preload("Students.User").
		Preload("Students.StudentGroup").
		First(&parent, *userID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("parent with user id %s was not found: %s: %w", *userID, result.Error.Error(), apperrors.ErrNotFound)
		}
		return nil, fmt.Errorf("failed to fetch parent by user id (%s): %w", *userID, result.Error)
	}
	return &parent, nil
}

func (r *parentRepository) FindAll(ctx context.Context, filter *ParentFilter) ([]model.Parent, error) {
	if filter == nil {
		return nil, fmt.Errorf("parents list filter cannot be nil: %w", apperrors.ErrRequiredField)
	}
	var parents []model.Parent
	query := r.db.WithContext(ctx).Model(&model.Parent{})
	// Filters
	// offset (for pagination):
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}
	// limit (for pagination):
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	// Sort parents in the alphabetical order
	query = query.Order("last_name")
	// Find parents
	result := query.
		Preload("Students").
		Preload("Students.User").
		Preload("Students.StudentGroup").
		Find(&parents)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch parents list: %w", result.Error)
	}
	// Return response
	return parents, nil
}
