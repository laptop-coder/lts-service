package service

import (
	"backend/internal/model"
	"backend/pkg/logger"
	"context"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type RoleService interface {
	GetRolePermissions(ctx context.Context, roleID uint8) ([]PermissionResponseDTO, error)
	// replace old permissions with new ones
	AssignPermissionsToRole(ctx context.Context, roleID uint8, permissionIDs []uint8) error
	AddPermissionsToRole(ctx context.Context, roleID uint8, permissionIDs []uint8) error
	RemovePermissionsFromRole(ctx context.Context, roleID uint8, permissionIDs []uint8) error
}

type PermissionResponseDTO struct {
	ID        uint8  `json:"id"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	Name      string `json:"name"`
}

type roleService struct {
	db  *gorm.DB
	log logger.Logger
}

func NewRoleService(db *gorm.DB, log logger.Logger) RoleService {
	return &roleService{
		db:  db,
		log: log,
	}
}

func (s *roleService) GetRolePermissions(ctx context.Context, roleID uint8) ([]PermissionResponseDTO, error) {
	// Get role by ID
	var role model.Role
	if err := s.db.WithContext(ctx).
		Preload("Permissions").
		First(&role, roleID).Error; err != nil {
		return nil, fmt.Errorf("role not found: %w", err)
	}
	// Get role's permissions, convert them to DTO
	dtos := make([]PermissionResponseDTO, len(role.Permissions))
	for i, permission := range role.Permissions {
		dtos[i] = *PermissionToDTO(&permission)
	}
	// Return response
	return dtos, nil
}

func (s *roleService) AssignPermissionsToRole(ctx context.Context, roleID uint8, permissionIDs []uint8) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Get role by ID
		var role model.Role
		if err := tx.First(&role, roleID).Error; err != nil {
			return fmt.Errorf("role not found: %w", err)
		}
		// Get permissions by IDs
		var permissions []model.Permission
		if len(permissionIDs) > 0 {
			if err := tx.Where("id IN (?)", permissionIDs).Find(&permissions).Error; err != nil {
				return fmt.Errorf("failed to fetch permissions: %w", err)
			}
			if len(permissions) != len(permissionIDs) {
				return fmt.Errorf("some permissions not found")
			}
		}
		// Return response
		return tx.Model(&role).Association("Permissions").Replace(&permissions)
	})
}

func (s *roleService) AddPermissionsToRole(ctx context.Context, roleID uint8, permissionIDs []uint8) error {
	// Get permissions by IDs
	var permissions []model.Permission
	if err := s.db.Where("id IN (?)", permissionIDs).Find(&permissions).Error; err != nil {
		return fmt.Errorf("failed to fetch permissions: %w", err)
	}
	// Add permission to role, return response
	var role model.Role
	role.ID = roleID
	return s.db.Model(&role).Association("Permissions").Append(&permissions)
}

func (s *roleService) RemovePermissionsFromRole(ctx context.Context, roleID uint8, permissionIDs []uint8) error {
	// Get permissions by ID
	var permissions []model.Permission
	if err := s.db.Where("id IN (?)", permissionIDs).Find(&permissions).Error; err != nil {
		return fmt.Errorf("failed to fetch permissions: %w", err)
	}
	// Delete permissions from role, return response
	var role model.Role
	role.ID = roleID
	return s.db.Model(&role).Association("Permissions").Delete(&permissions)
}

func PermissionToDTO(permission *model.Permission) *PermissionResponseDTO {
	return &PermissionResponseDTO{
		ID:        permission.ID,
		CreatedAt: permission.CreatedAt.Format(time.RFC3339),
		UpdatedAt: permission.UpdatedAt.Format(time.RFC3339),
		Name:      permission.Name,
	}
}
