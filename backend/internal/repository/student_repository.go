package repository

import (
	"backend/internal/model"
	log "backend/pkg/logger"
	"context"
	"fmt"
	"gorm.io/gorm"
)

type StudentRepository interface {
	FindByGroupID(ctx context.Context, id *uint16) ([]model.Student, error)
}

type studentRepository struct {
	db *gorm.DB
}

func NewStudentRepository(db *gorm.DB) StudentRepository {
	if db == nil {
		log.Error("DB is nil")
		panic("DB is nil")
	}
	return &studentRepository{db: db}
}

func (r *studentRepository) FindByGroupID(ctx context.Context, id *uint16) ([]model.Student, error) {
	if id == nil {
		return nil, fmt.Errorf("student group id cannot be nil")
	}

	var students []model.Student
	result := r.db.WithContext(ctx).
		Where("student_group_id = ?", *id).
		Find(&students)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch students list by group id (%d): %w", *id, result.Error)
	}

	return students, nil
}
