package repository

import (
	"backend/internal/model"
	"backend/pkg/logger"
	"context"
	"fmt"
	"gorm.io/gorm"
)

type RoleRepository interface {
	FindAll(ctx context.Context) ([]model.Role, error)
}

type roleRepository struct {
	db  *gorm.DB
	log logger.Logger
}

func NewRoleRepository(db *gorm.DB, log logger.Logger) RoleRepository {
	if db == nil {
		log.Error("DB is nil")
		panic("DB is nil")
	}
	return &roleRepository{db: db, log: log}
}

func (r *roleRepository) FindAll(ctx context.Context) ([]model.Role, error) {
	var roles []model.Role

	err := r.db.WithContext(ctx).
		Model(&model.Role{}).
		Order("name").
		Find(&roles).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch roles list: %w", err)
	}

	return roles, nil
}
