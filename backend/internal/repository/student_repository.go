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

type StudentRepository interface {
	FindByGroupID(ctx context.Context, id *uint16) ([]model.Student, error)
	FindByID(ctx context.Context, userID *uuid.UUID) (*model.Student, error)
	FindAll(ctx context.Context, filter *StudentFilter) ([]model.Student, error)
}

type StudentFilter struct {
	Limit  int
	Offset int
}

type studentRepository struct {
	db  *gorm.DB
	log logger.Logger
}

func NewStudentRepository(db *gorm.DB, log logger.Logger) StudentRepository {
	if db == nil {
		log.Error("DB is nil")
		panic("DB is nil")
	}
	return &studentRepository{db: db, log: log}
}

func (r *studentRepository) FindByGroupID(ctx context.Context, id *uint16) ([]model.Student, error) {
	if id == nil {
		return nil, fmt.Errorf("student group id cannot be nil")
	}
	var students []model.Student
	result := r.db.WithContext(ctx).
		Preload("User").
		Preload("StudentGroup").
		Preload("Parents").
		Where("student_group_id = ?", *id).
		Find(&students)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch students list by group id (%d): %w", *id, result.Error)
	}
	return students, nil
}

func (r *studentRepository) FindByID(ctx context.Context, userID *uuid.UUID) (*model.Student, error) {
	if userID == nil {
		return nil, fmt.Errorf("user id cannot be nil")
	}
	var student model.Student
	result := r.db.WithContext(ctx).
		Preload("User").
		Preload("StudentGroup").
		Preload("Parents").
		First(&student, *userID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("student with user id %s was not found: %w", *userID, result.Error)
		}
		return nil, fmt.Errorf("failed to fetch student by user id (%s): %w", *userID, result.Error)
	}
	return &student, nil
}

func (r *studentRepository) FindAll(ctx context.Context, filter *StudentFilter) ([]model.Student, error) {
	if filter == nil {
		return nil, fmt.Errorf("students list filter cannot be nil")
	}
	var students []model.Student
	query := r.db.WithContext(ctx).Model(&model.Student{})
	// Filters
	// offset (for pagination):
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}
	// limit (for pagination):
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	// Sort students in the alphabetical order
	query = query.Order("name")
	// Find students
	result := query.
		Preload("User").
		Preload("StudentGroup").
		Preload("Parents").
		Find(&students)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch students list: %w", result.Error)
	}
	// Return response
	return students, nil
}
