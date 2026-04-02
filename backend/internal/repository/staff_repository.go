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

type StaffRepository interface {
	FindByID(ctx context.Context, userID *uuid.UUID) (*model.Staff, error)
	FindAll(ctx context.Context, filter *StaffFilter) ([]model.Staff, error)
}

type StaffFilter struct {
	Limit  int
	Offset int
}

type staffRepository struct {
	db  *gorm.DB
	log logger.Logger
}

func NewStaffRepository(db *gorm.DB, log logger.Logger) StaffRepository {
	if db == nil {
		log.Error("DB is nil")
		panic("DB is nil")
	}
	return &staffRepository{db: db, log: log}
}

func (r *staffRepository) FindByID(ctx context.Context, userID *uuid.UUID) (*model.Staff, error) {
	if userID == nil {
		return nil, fmt.Errorf("user id cannot be nil")
	}
	var staff model.Staff
	result := r.db.WithContext(ctx).
		Preload("Position").
		First(&staff, *userID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("staff with user id %s was not found: %w", *userID, result.Error)
		}
		return nil, fmt.Errorf("failed to fetch staff by user id (%s): %w", *userID, result.Error)
	}
	return &staff, nil
}

func (r *staffRepository) FindAll(ctx context.Context, filter *StaffFilter) ([]model.Staff, error) {
	if filter == nil {
		return nil, fmt.Errorf("staff list filter cannot be nil")
	}
	var staffList []model.Staff
	query := r.db.WithContext(ctx).Model(&model.Staff{})
	// Filters
	// offset (for pagination):
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}
	// limit (for pagination):
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	// Sort staff in the alphabetical order
	query = query.Order("last_name")
	// Find staff
	result := query.
		Preload("Position").
		Find(&staffList)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch staff list: %w", result.Error)
	}
	// Return response
	return staffList, nil
}
