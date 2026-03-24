package repository

import (
	"backend/internal/model"
	"backend/pkg/logger"
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type StaffPositionRepository interface {
	Create(ctx context.Context, staffPosition *model.StaffPosition) error
	FindAll(ctx context.Context, filter *StaffPositionFilter) ([]model.StaffPosition, error)
	FindByID(ctx context.Context, id *uint8) (*model.StaffPosition, error)
	Update(ctx context.Context, staffPosition *model.StaffPosition) error
	Delete(ctx context.Context, id *uint8) error
}

type staffPositionRepository struct {
	db  *gorm.DB
	log logger.Logger
}

type StaffPositionFilter struct {
	Limit  int
	Offset int
}

func NewStaffPositionRepository(db *gorm.DB, log logger.Logger) StaffPositionRepository {
	if db == nil {
		log.Error("DB is nil")
		panic("DB is nil")
	}
	return &staffPositionRepository{db: db, log: log}
}

func (r *staffPositionRepository) Create(ctx context.Context, staffPosition *model.StaffPosition) error {
	if staffPosition == nil {
		return fmt.Errorf("staffPosition cannot be nil")
	}
	result := r.db.WithContext(ctx).Create(staffPosition)
	if result.Error != nil {
		return fmt.Errorf("failed to create new staffPosition: %w", result.Error)
	}
	return nil
}

func (r *staffPositionRepository) FindAll(ctx context.Context, filter *StaffPositionFilter) ([]model.StaffPosition, error) {
	if filter == nil {
		return nil, fmt.Errorf("staffPositions list filter cannot be nil")
	}
	var staffPositions []model.StaffPosition
	query := r.db.WithContext(ctx).Model(&model.StaffPosition{})
	// Filters
	// offset (for pagination):
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}
	// limit (for pagination):
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	// Sort staffPositions in the alphabetical order
	query = query.Order("name")
	// Find staffPositions
	result := query.Find(&staffPositions)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch staffPositions list: %w", result.Error)
	}
	return staffPositions, nil
}

func (r *staffPositionRepository) FindByID(ctx context.Context, id *uint8) (*model.StaffPosition, error) {
	if id == nil {
		return nil, fmt.Errorf("staffPosition id cannot be nil")
	}
	var staffPosition model.StaffPosition
	result := r.db.WithContext(ctx).First(&staffPosition, *id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("staffPosition with id %d was not found: %w", *id, result.Error)
		}
		return nil, fmt.Errorf("failed to fetch staffPosition by id (%d): %w", *id, result.Error)
	}
	return &staffPosition, nil
}

func (r *staffPositionRepository) Update(ctx context.Context, staffPosition *model.StaffPosition) error {
	if staffPosition == nil {
		return fmt.Errorf("staffPosition cannot be nil")
	}
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.StaffPosition{}).
		Where("id = ?", staffPosition.ID).
		Count(&count).Error
	if err != nil {
		return fmt.Errorf("failed to check staffPosition existence: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("staffPosition with id %d was not found", staffPosition.ID)
	}
	result := r.db.WithContext(ctx).Save(staffPosition)
	if result.Error != nil {
		return fmt.Errorf("failed to update staffPosition: %w", result.Error)
	}
	return nil
}

func (r *staffPositionRepository) Delete(ctx context.Context, id *uint8) error {
	result := r.db.WithContext(ctx).Unscoped().Delete(&model.StaffPosition{}, *id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete staffPosition with id %d: %w", *id, result.Error)
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
