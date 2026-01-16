// Package repository provides data access abstraction.
package repository

import (
	"errors"
	"backend/internal/model"
	log "backend/pkg/logger"
	"context"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	FindAll(ctx context.Context, filter *UserFilter) ([]model.User, error)
	FindByID(ctx context.Context, id *uuid.UUID) (*model.User, error)
	FindByEmail(ctx context.Context, email *string) (*model.User, error)
	FindStudentGroupAdvisorByGroupID(ctx context.Context, id *uint16) (*model.User, error)
	// ID must be set to update
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id *uuid.UUID) error
}

type userRepository struct {
	db *gorm.DB
}

type UserFilter struct {
	RoleID *uint8
	Limit  int
	Offset int
}

func NewUserRepository(db *gorm.DB) UserRepository {
	if db == nil {
		log.Error("DB is nil")
		panic("DB is nil")
	}
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	if user == nil {
		return fmt.Errorf("user cannot be nil")
	}
	result := r.db.WithContext(ctx).Create(user)
	if result.Error != nil {
		return fmt.Errorf("failed to create new user: %w", result.Error)
	}
	return nil
}

func (r *userRepository) FindAll(ctx context.Context, filter *UserFilter) ([]model.User, error) {
	if filter == nil {
		return nil, fmt.Errorf("users list filter cannot be nil")
	}

	var users []model.User
	query := r.db.WithContext(ctx).Model(&model.User{})

	// Filters
	// By user's role:
	if filter.RoleID != nil {
		query = query.
			Joins("JOIN user_roles ON user_roles.user_id = users.id").
			Where("user_roles.role_id = ?", *filter.RoleID)
	}
	// Offset (for pagination):
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}
	// Limit (for pagination):
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}

	// Sort users in the alphabetical order
	query = query.Order("name")

	result := query.Find(&users)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch users list: %w", result.Error)
	}
	return users, nil
}

func (r *userRepository) FindByID(ctx context.Context, id *uuid.UUID) (*model.User, error) {
	if id == nil {
		return nil, fmt.Errorf("user id cannot be nil")
	}

	var user model.User
	result := r.db.WithContext(ctx).First(&user, *id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user with id %s was not found: %w", *id, result.Error)
		}
		return nil, fmt.Errorf("failed to fetch user by id (%s): %w", *id, result.Error)
	}

	return &user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email *string) (*model.User, error) {
	if email == nil {
		return nil, fmt.Errorf("user email cannot be nil")
	}

	var user model.User
	result := r.db.WithContext(ctx).Where("email = ?", *email).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user was not found by email: %w", result.Error)
		}
		return nil, fmt.Errorf("failed to fetch user by email: %w", result.Error)
	}

	return &user, nil
}

func (r *userRepository) FindStudentGroupAdvisorByGroupID(ctx context.Context, id *uint16) (*model.User, error) {
	if id == nil {
		return nil, fmt.Errorf("student group id cannot be nil")
	}

	var user model.User
	result := r.db.WithContext(ctx).
		Joins("JOIN student_groups ON student_groups.group_advisor_id = users.id").
		Where("student_groups.id = ?", *id).
		First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user (student group advisor) was not found by group id (%d): %w", *id, result.Error)
		}
		return nil, fmt.Errorf("failed to fetch user (student group advisor) by group id (%d): %w", *id, result.Error)
	}

	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	if user == nil {
		return fmt.Errorf("user cannot be nil")
	}

	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", user.ID).
		Count(&count).Error
	
	if err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("user with id %s was not found", user.ID)
	}

	result := r.db.WithContext(ctx).Save(user)
	if result.Error != nil {
		return fmt.Errorf("failed to update user: %w", result.Error)
	}

	return nil
}

func (r *userRepository) Delete(ctx context.Context, id *uuid.UUID) error {
	result := r.db.WithContext(ctx).Unscoped().Delete(&model.User{}, *id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete user with id %s: %w", *id, result.Error)
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
