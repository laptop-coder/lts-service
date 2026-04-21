package repository

import (
	"backend/internal/model"
	"backend/pkg/apperrors"
	"backend/pkg/logger"
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type SubjectRepository interface {
	Create(ctx context.Context, subject *model.Subject) error
	FindAll(ctx context.Context, filter *SubjectFilter) ([]model.Subject, error)
	FindByID(ctx context.Context, id *uint8) (*model.Subject, error)
	Update(ctx context.Context, subject *model.Subject) error
	Delete(ctx context.Context, id *uint8) error
	ExistsByName(ctx context.Context, name *string) (bool, error)
}

type subjectRepository struct {
	db  *gorm.DB
	log logger.Logger
}

type SubjectFilter struct {
	Limit  int
	Offset int
}

func NewSubjectRepository(db *gorm.DB, log logger.Logger) SubjectRepository {
	if db == nil {
		log.Error("DB is nil")
		panic("DB is nil")
	}
	return &subjectRepository{db: db, log: log}
}

func (r *subjectRepository) Create(ctx context.Context, subject *model.Subject) error {
	if subject == nil {
		return fmt.Errorf("subject cannot be nil: %w", apperrors.ErrRequiredField)
	}
	result := r.db.WithContext(ctx).Create(subject)
	if result.Error != nil {
		return fmt.Errorf("failed to create new subject: %w", result.Error)
	}
	return nil
}

func (r *subjectRepository) FindAll(ctx context.Context, filter *SubjectFilter) ([]model.Subject, error) {
	if filter == nil {
		return nil, fmt.Errorf("subjects list filter cannot be nil: %w", apperrors.ErrRequiredField)
	}
	var subjects []model.Subject
	query := r.db.WithContext(ctx).Model(&model.Subject{})
	// Filters
	// offset (for pagination):
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}
	// limit (for pagination):
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	// Sort subjects in the alphabetical order
	query = query.Order("name")
	// Find subjects
	result := query.Find(&subjects)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch subjects list: %w", result.Error)
	}
	// Return response
	return subjects, nil
}

func (r *subjectRepository) FindByID(ctx context.Context, id *uint8) (*model.Subject, error) {
	if id == nil {
		return nil, fmt.Errorf("subject id cannot be nil: %w", apperrors.ErrRequiredField)
	}
	var subject model.Subject
	result := r.db.WithContext(ctx).First(&subject, *id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("subject with id %d was not found: %s: %w", *id, result.Error.Error(), apperrors.ErrNotFound)
		}
		return nil, fmt.Errorf("failed to fetch subject by id (%d): %w", *id, result.Error)
	}
	return &subject, nil
}

func (r *subjectRepository) Update(ctx context.Context, subject *model.Subject) error {
	if subject == nil {
		return fmt.Errorf("subject cannot be nil: %w", apperrors.ErrRequiredField)
	}
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.Subject{}).
		Where("id = ?", subject.ID).
		Count(&count).Error
	if err != nil {
		return fmt.Errorf("failed to check subject existence: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("subject with id %d was not found: %w", subject.ID, apperrors.ErrNotFound)
	}
	result := r.db.WithContext(ctx).Save(subject)
	if result.Error != nil {
		return fmt.Errorf("failed to update subject: %w", result.Error)
	}
	return nil
}

func (r *subjectRepository) Delete(ctx context.Context, id *uint8) error {
	result := r.db.WithContext(ctx).Unscoped().Delete(&model.Subject{}, *id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete subject with id %d: %w", *id, result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("subject to delete was not found by id: %w", apperrors.ErrNotFound)
	}
	return nil
}

func (r *subjectRepository) ExistsByName(ctx context.Context, name *string) (bool, error) {
	if name == nil {
		return false, fmt.Errorf("name cannot be nil: %w", apperrors.ErrRequiredField)
	}
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.Subject{}).
		Where("name = ?", name).
		Count(&count).Error
	return count > 0, err
}
