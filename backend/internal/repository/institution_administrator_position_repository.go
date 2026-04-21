package repository

import (
	"backend/internal/model"
	"backend/pkg/apperrors"
	"backend/pkg/logger"
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type InstitutionAdministratorPositionRepository interface {
	Create(ctx context.Context, institutionAdministratorPosition *model.InstitutionAdministratorPosition) error
	FindAll(ctx context.Context, filter *InstitutionAdministratorPositionFilter) ([]model.InstitutionAdministratorPosition, error)
	FindByID(ctx context.Context, id *uint8) (*model.InstitutionAdministratorPosition, error)
	Update(ctx context.Context, institutionAdministratorPosition *model.InstitutionAdministratorPosition) error
	Delete(ctx context.Context, id *uint8) error
	ExistsByName(ctx context.Context, name *string) (bool, error)
}

type institutionAdministratorPositionRepository struct {
	db  *gorm.DB
	log logger.Logger
}

type InstitutionAdministratorPositionFilter struct {
	Limit  int
	Offset int
}

func NewInstitutionAdministratorPositionRepository(db *gorm.DB, log logger.Logger) InstitutionAdministratorPositionRepository {
	if db == nil {
		log.Error("DB is nil")
		panic("DB is nil")
	}
	return &institutionAdministratorPositionRepository{db: db, log: log}
}

func (r *institutionAdministratorPositionRepository) Create(ctx context.Context, institutionAdministratorPosition *model.InstitutionAdministratorPosition) error {
	if institutionAdministratorPosition == nil {
		return fmt.Errorf("institutionAdministratorPosition cannot be nil: %w", apperrors.ErrRequiredField)
	}
	result := r.db.WithContext(ctx).Create(institutionAdministratorPosition)
	if result.Error != nil {
		return fmt.Errorf("failed to create new institutionAdministratorPosition: %w", result.Error)
	}
	return nil
}

func (r *institutionAdministratorPositionRepository) FindAll(ctx context.Context, filter *InstitutionAdministratorPositionFilter) ([]model.InstitutionAdministratorPosition, error) {
	if filter == nil {
		return nil, fmt.Errorf("institutionAdministratorPositions list filter cannot be nil: %w", apperrors.ErrRequiredField)
	}
	var institutionAdministratorPositions []model.InstitutionAdministratorPosition
	query := r.db.WithContext(ctx).Model(&model.InstitutionAdministratorPosition{})
	// Filters
	// offset (for pagination):
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}
	// limit (for pagination):
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	// Sort institutionAdministratorPositions in the alphabetical order
	query = query.Order("name")
	// Find institutionAdministratorPositions
	result := query.Find(&institutionAdministratorPositions)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch institutionAdministratorPositions list: %w", result.Error)
	}
	return institutionAdministratorPositions, nil
}

func (r *institutionAdministratorPositionRepository) FindByID(ctx context.Context, id *uint8) (*model.InstitutionAdministratorPosition, error) {
	if id == nil {
		return nil, fmt.Errorf("institutionAdministratorPosition id cannot be nil: %w", apperrors.ErrRequiredField)
	}
	var institutionAdministratorPosition model.InstitutionAdministratorPosition
	result := r.db.WithContext(ctx).First(&institutionAdministratorPosition, *id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("institutionAdministratorPosition with id %d was not found: %s: %w", *id, result.Error.Error(), apperrors.ErrNotFound)
		}
		return nil, fmt.Errorf("failed to fetch institutionAdministratorPosition by id (%d): %w", *id, result.Error)
	}
	return &institutionAdministratorPosition, nil
}

func (r *institutionAdministratorPositionRepository) Update(ctx context.Context, institutionAdministratorPosition *model.InstitutionAdministratorPosition) error {
	if institutionAdministratorPosition == nil {
		return fmt.Errorf("institutionAdministratorPosition cannot be nil: %w", apperrors.ErrRequiredField)
	}
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.InstitutionAdministratorPosition{}).
		Where("id = ?", institutionAdministratorPosition.ID).
		Count(&count).Error
	if err != nil {
		return fmt.Errorf("failed to check institutionAdministratorPosition existence: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("institutionAdministratorPosition with id %d was not found: %w", institutionAdministratorPosition.ID, apperrors.ErrNotFound)
	}
	result := r.db.WithContext(ctx).Save(institutionAdministratorPosition)
	if result.Error != nil {
		return fmt.Errorf("failed to update institutionAdministratorPosition: %w", result.Error)
	}
	return nil
}

func (r *institutionAdministratorPositionRepository) Delete(ctx context.Context, id *uint8) error {
	result := r.db.WithContext(ctx).Unscoped().Delete(&model.InstitutionAdministratorPosition{}, *id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete institutionAdministratorPosition with id %d: %w", *id, result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("institution administrator position not found by id: %w", apperrors.ErrNotFound)
	}
	return nil
}

func (r *institutionAdministratorPositionRepository) ExistsByName(ctx context.Context, name *string) (bool, error) {
	if name == nil {
		return false, fmt.Errorf("name cannot be nil: %w", apperrors.ErrRequiredField)
	}
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.InstitutionAdministratorPosition{}).
		Where("name = ?", name).
		Count(&count).Error
	return count > 0, err
}
