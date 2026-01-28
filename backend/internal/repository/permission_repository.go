package repository

import (
	"backend/internal/model"
	"backend/pkg/logger"
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type PermissionRepository interface {
	Create(ctx context.Context, permission *model.Permission) error
	FindAll(ctx context.Context) ([]model.Permission, error)
	FindByID(ctx context.Context, id *uint8) (*model.Permission, error)
	Update(ctx context.Context, permission *model.Permission) error
	Delete(ctx context.Context, id *uint8) error
}

type permissionRepository struct {
	db  *gorm.DB
	log logger.Logger
}

func NewPermissionRepository(db *gorm.DB, log logger.Logger) PermissionRepository {
	if db == nil {
		log.Error("DB is nil")
		panic("DB is nil")
	}
	return &permissionRepository{db: db, log: log}
}

func (r *permissionRepository) Create(ctx context.Context, permission *model.Permission) error {
	if permission == nil {
		return fmt.Errorf("permission cannot be nil")
	}

	result := r.db.WithContext(ctx).Create(permission)
	if result.Error != nil {
		return fmt.Errorf("failed to create new permission: %w", result.Error)
	}

	return nil
}

func (r *permissionRepository) FindAll(ctx context.Context) ([]model.Permission, error) {
	var permissions []model.Permission

	err := r.db.WithContext(ctx).
		Model(&model.Permission{}).
		Order("name").
		Find(&permissions).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch permissions list: %w", err)
	}

	return permissions, nil
}

func (r *permissionRepository) FindByID(ctx context.Context, id *uint8) (*model.Permission, error) {
	if id == nil {
		return nil, fmt.Errorf("permission id cannot be nil")
	}

	var permission model.Permission
	result := r.db.WithContext(ctx).First(&permission, *id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("permission with id %d was not found: %w", *id, result.Error)
		}
		return nil, fmt.Errorf("failed to fetch permission by id (%d): %w", *id, result.Error)
	}

	return &permission, nil
}

func (r *permissionRepository) Update(ctx context.Context, permission *model.Permission) error {
	if permission == nil {
		return fmt.Errorf("permission cannot be nil")
	}

	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.Permission{}).
		Where("id = ?", permission.ID).
		Count(&count).Error

	if err != nil {
		return fmt.Errorf("failed to check permission existence: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("permission with id %d was not found", permission.ID)
	}

	result := r.db.WithContext(ctx).Save(permission)
	if result.Error != nil {
		return fmt.Errorf("failed to update permission: %w", result.Error)
	}

	return nil
}

func (r *permissionRepository) Delete(ctx context.Context, id *uint8) error {
	result := r.db.WithContext(ctx).Unscoped().Delete(&model.Permission{}, *id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete permission with id %d: %w", *id, result.Error)
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
