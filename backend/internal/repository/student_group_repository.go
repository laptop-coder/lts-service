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

type StudentGroupRepository interface {
	Create(ctx context.Context, studentGroup *model.StudentGroup) error
	FindAll(ctx context.Context, filter *StudentGroupFilter) ([]model.StudentGroup, error)
	FindByID(ctx context.Context, id *uint16) (*model.StudentGroup, error)
	FindAdvisorByGroupID(ctx context.Context, id *uint16) (*model.User, error)
	Update(ctx context.Context, studentGroup *model.StudentGroup) error
	Delete(ctx context.Context, id *uint16) error
	ExistsByName(ctx context.Context, name *string) (bool, error)
}

type studentGroupRepository struct {
	db  *gorm.DB
	log logger.Logger
}

type StudentGroupFilter struct {
	GroupAdvisorID *uuid.UUID
	Limit          int
	Offset         int
}

func NewStudentGroupRepository(db *gorm.DB, log logger.Logger) StudentGroupRepository {
	if db == nil {
		log.Error("DB is nil")
		panic("DB is nil")
	}
	return &studentGroupRepository{db: db, log: log}
}

func (r *studentGroupRepository) Create(ctx context.Context, studentGroup *model.StudentGroup) error {
	if studentGroup == nil {
		return fmt.Errorf("student group cannot be nil: %w", apperrors.ErrRequiredField)
	}
	result := r.db.WithContext(ctx).Create(studentGroup)
	if result.Error != nil {
		return fmt.Errorf("failed to create new student group: %w", result.Error)
	}
	return nil
}

func (r *studentGroupRepository) FindAll(ctx context.Context, filter *StudentGroupFilter) ([]model.StudentGroup, error) {
	if filter == nil {
		return nil, fmt.Errorf("student groups list filter cannot be nil: %w", apperrors.ErrRequiredField)
	}
	var studentGroups []model.StudentGroup
	query := r.db.WithContext(ctx).Model(&model.StudentGroup{})
	// Filters
	// by group advisor:
	if filter.GroupAdvisorID != nil {
		query = query.Where("group_advisor_id = ?", *filter.GroupAdvisorID)
	}
	// offset (for pagination):
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}
	// limit (for pagination):
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	// Sort student groups in the alphabetical order
	query = query.Order("name")
	// Find student groups
	result := query.Find(&studentGroups)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch student groups list: %w", result.Error)
	}
	// Return response
	return studentGroups, nil
}

func (r *studentGroupRepository) FindByID(ctx context.Context, id *uint16) (*model.StudentGroup, error) {
	if id == nil {
		return nil, fmt.Errorf("student group id cannot be nil: %w", apperrors.ErrRequiredField)
	}
	var studentGroup model.StudentGroup
	result := r.db.WithContext(ctx).First(&studentGroup, *id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("student group with id %d was not found: %s: %w", *id, result.Error.Error(), apperrors.ErrNotFound)
		}
		return nil, fmt.Errorf("failed to fetch student group by id (%d): %w", *id, result.Error)
	}
	return &studentGroup, nil
}

func (r *studentGroupRepository) FindAdvisorByGroupID(ctx context.Context, id *uint16) (*model.User, error) {
	if id == nil {
		return nil, fmt.Errorf("student group id cannot be nil: %w", apperrors.ErrRequiredField)
	}
	var user model.User
	result := r.db.WithContext(ctx).
		Joins("JOIN student_groups ON student_groups.group_advisor_id = users.id").
		Where("student_groups.id = ?", *id).
		First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user (student group advisor) was not found by group id (%d): %s: %w", *id, result.Error.Error(), apperrors.ErrNotFound)
		}
		return nil, fmt.Errorf("failed to fetch user (student group advisor) by group id (%d): %w", *id, result.Error)
	}
	return &user, nil
}

func (r *studentGroupRepository) Update(ctx context.Context, studentGroup *model.StudentGroup) error {
	if studentGroup == nil {
		return fmt.Errorf("student group cannot be nil: %w", apperrors.ErrRequiredField)
	}
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.StudentGroup{}).
		Where("id = ?", studentGroup.ID).
		Count(&count).Error
	if err != nil {
		return fmt.Errorf("failed to check student group existence (id %d): %w", studentGroup.ID, err)
	}
	if count == 0 {
		return fmt.Errorf("student group with id %d was not found: %w", studentGroup.ID, apperrors.ErrNotFound)
	}
	result := r.db.WithContext(ctx).Save(studentGroup)
	if result.Error != nil {
		return fmt.Errorf("failed to update studentGroup: %w", result.Error)
	}
	return nil
}

func (r *studentGroupRepository) Delete(ctx context.Context, id *uint16) error {
	result := r.db.WithContext(ctx).Unscoped().Delete(&model.StudentGroup{}, *id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete student group with id %d: %w", *id, result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("student group to delete was not found by id: %w", apperrors.ErrNotFound)
	}
	return nil
}

func (r *studentGroupRepository) ExistsByName(ctx context.Context, name *string) (bool, error) {
	if name == nil {
		return false, fmt.Errorf("name cannot be nil: %w", apperrors.ErrRequiredField)
	}
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.StudentGroup{}).
		Where("name = ?", name).
		Count(&count).Error
	return count > 0, err
}
