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

type InstitutionAdministratorRepository interface {
	FindByID(ctx context.Context, userID *uuid.UUID) (*model.InstitutionAdministrator, error)
	FindAll(ctx context.Context, filter *InstitutionAdministratorFilter) ([]model.InstitutionAdministrator, error)
}

type InstitutionAdministratorFilter struct {
	Limit  int
	Offset int
}

type institutionAdministratorRepository struct {
	db  *gorm.DB
	log logger.Logger
}

func NewInstitutionAdministratorRepository(db *gorm.DB, log logger.Logger) InstitutionAdministratorRepository {
	if db == nil {
		log.Error("DB is nil")
		panic("DB is nil")
	}
	return &institutionAdministratorRepository{db: db, log: log}
}

func (r *institutionAdministratorRepository) FindByID(ctx context.Context, userID *uuid.UUID) (*model.InstitutionAdministrator, error) {
	if userID == nil {
		return nil, fmt.Errorf("user id cannot be nil: %w", apperrors.ErrRequiredField)
	}
	var institutionAdministrator model.InstitutionAdministrator
	result := r.db.WithContext(ctx).
		Preload("Position").
		First(&institutionAdministrator, *userID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("institutionAdministrator with user id %s was not found: %s: %w", *userID, result.Error.Error(), apperrors.ErrNotFound)
		}
		return nil, fmt.Errorf("failed to fetch institutionAdministrator by user id (%s): %w", *userID, result.Error)
	}
	return &institutionAdministrator, nil
}

func (r *institutionAdministratorRepository) FindAll(ctx context.Context, filter *InstitutionAdministratorFilter) ([]model.InstitutionAdministrator, error) {
	if filter == nil {
		return nil, fmt.Errorf("institutionAdministrators list filter cannot be nil: %w", apperrors.ErrRequiredField)
	}
	var institutionAdministrators []model.InstitutionAdministrator
	query := r.db.WithContext(ctx).Model(&model.InstitutionAdministrator{})
	// Filters
	// offset (for pagination):
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}
	// limit (for pagination):
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	// Sort institutionAdministrators in the alphabetical order
	query = query.Order("last_name")
	// Find institutionAdministrators
	result := query.
		Preload("Position").
		Find(&institutionAdministrators)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch institutionAdministrators list: %w", result.Error)
	}
	// Return response
	return institutionAdministrators, nil
}
