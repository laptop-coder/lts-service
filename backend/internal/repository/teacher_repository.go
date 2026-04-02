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

type TeacherRepository interface {
	FindByID(ctx context.Context, userID *uuid.UUID) (*model.Teacher, error)
	FindAll(ctx context.Context, filter *TeacherFilter) ([]model.Teacher, error)
}

type TeacherFilter struct {
	Limit  int
	Offset int
}

type teacherRepository struct {
	db  *gorm.DB
	log logger.Logger
}

func NewTeacherRepository(db *gorm.DB, log logger.Logger) TeacherRepository {
	if db == nil {
		log.Error("DB is nil")
		panic("DB is nil")
	}
	return &teacherRepository{db: db, log: log}
}

func (r *teacherRepository) FindByID(ctx context.Context, userID *uuid.UUID) (*model.Teacher, error) {
	if userID == nil {
		return nil, fmt.Errorf("user id cannot be nil")
	}
	var teacher model.Teacher
	result := r.db.WithContext(ctx).
		Preload("Classroom").
		Preload("Subjects").
		Preload("StudentGroups").
		First(&teacher, *userID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("teacher with user id %s was not found: %w", *userID, result.Error)
		}
		return nil, fmt.Errorf("failed to fetch teacher by user id (%s): %w", *userID, result.Error)
	}
	return &teacher, nil
}

func (r *teacherRepository) FindAll(ctx context.Context, filter *TeacherFilter) ([]model.Teacher, error) {
	if filter == nil {
		return nil, fmt.Errorf("teachers list filter cannot be nil")
	}
	var teachers []model.Teacher
	query := r.db.WithContext(ctx).Model(&model.Teacher{})
	// Filters
	// offset (for pagination):
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}
	// limit (for pagination):
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	// Sort teachers in the alphabetical order
	query = query.Order("last_name")
	// Find teachers
	result := query.
		Preload("Classroom").
		Preload("Subjects").
		Preload("StudentGroups").
		Find(&teachers)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch teachers list: %w", result.Error)
	}
	// Return response
	return teachers, nil
}
