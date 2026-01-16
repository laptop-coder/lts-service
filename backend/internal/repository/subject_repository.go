package repository

import (
	"errors"
	"backend/internal/model"
	log "backend/pkg/logger"
	"context"
	"fmt"
	"gorm.io/gorm"
)

type SubjectRepository interface {
	Create(ctx context.Context, subject *model.Subject) error
	FindAll(ctx context.Context) ([]model.Subject, error)
	FindByID(ctx context.Context, id *uint8) (*model.Subject, error)
	Update(ctx context.Context, subject *model.Subject) error
	Delete(ctx context.Context, id *uint8) error
}

type subjectRepository struct {
	db *gorm.DB
}

func NewSubjectRepository(db *gorm.DB) SubjectRepository {
	if db == nil {
		log.Error("DB is nil")
		panic("DB is nil")
	}
	return &subjectRepository{db: db}
}

func (r *subjectRepository) Create(ctx context.Context, subject *model.Subject) error {
	if subject == nil {
		return fmt.Errorf("subject cannot be nil")
	}

	result := r.db.WithContext(ctx).Create(subject)
	if result.Error != nil {
		return fmt.Errorf("failed to create new subject: %w", result.Error)
	}

	return nil
}

func (r *subjectRepository) FindAll(ctx context.Context) ([]model.Subject, error) {
	var subjects []model.Subject

	err := r.db.WithContext(ctx).
	    Model(&model.Subject{}).
		Order("name").
		Find(&subjects).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch subjects list: %w", err)
	}

	return subjects, nil
}

func (r *subjectRepository) FindByID(ctx context.Context, id *uint8) (*model.Subject, error) {
	if id == nil {
		return nil, fmt.Errorf("subject id cannot be nil")
	}

	var subject model.Subject
	result := r.db.WithContext(ctx).First(&subject, *id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("subject with id %d was not found: %w", *id, result.Error)
		}
		return nil, fmt.Errorf("failed to fetch subject by id (%d): %w", *id, result.Error)
	}

	return &subject, nil
}

func (r *subjectRepository) Update(ctx context.Context, subject *model.Subject) error {
	if subject == nil {
		return fmt.Errorf("subject cannot be nil")
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
		return fmt.Errorf("subject with id %d was not found", subject.ID)
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
		return gorm.ErrRecordNotFound
	}
	return nil
}
