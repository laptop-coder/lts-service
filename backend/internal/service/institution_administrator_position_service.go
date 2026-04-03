package service

import (
	"backend/internal/model"
	"backend/internal/repository"
	"backend/pkg/logger"
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"strings"
	"time"
)

type InstitutionAdministratorPositionService interface {
	CreatePosition(ctx context.Context, dto CreateInstitutionAdministratorPositionDTO) (*InstitutionAdministratorPositionResponseDTO, error)
	GetPositions(ctx context.Context, filter repository.InstitutionAdministratorPositionFilter) ([]InstitutionAdministratorPositionResponseDTO, error)
	UpdatePosition(ctx context.Context, id uint8, dto UpdateInstitutionAdministratorPositionDTO) (*InstitutionAdministratorPositionResponseDTO, error)
	DeletePosition(ctx context.Context, id uint8) error
}

type CreateInstitutionAdministratorPositionDTO struct {
	Name string `form:"name" validate:"required,min=4,max=100"`
}

type UpdateInstitutionAdministratorPositionDTO struct {
	Name *string `form:"name,omitempty" validate:"max=100"`
}

type InstitutionAdministratorPositionResponseDTO struct {
	ID        uint8  `json:"id"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	Name      string `json:"name"`
}

type institutionAdministratorPositionService struct {
	institutionAdministratorPositionRepo repository.InstitutionAdministratorPositionRepository
	db                                   *gorm.DB
	log                                  logger.Logger
}

func NewInstitutionAdministratorPositionService(
	institutionAdministratorPositionRepo repository.InstitutionAdministratorPositionRepository,
	db *gorm.DB,
	log logger.Logger,
) InstitutionAdministratorPositionService {
	return &institutionAdministratorPositionService{
		institutionAdministratorPositionRepo: institutionAdministratorPositionRepo,
		db:                                   db,
		log:                                  log,
	}
}

func (s *institutionAdministratorPositionService) CreatePosition(ctx context.Context, dto CreateInstitutionAdministratorPositionDTO) (*InstitutionAdministratorPositionResponseDTO, error) {
	// Input data validation
	if err := s.validateCreatePositionDTO(&dto); err != nil {
		return nil, fmt.Errorf("validation error during institutionAdministratorPosition creation: %w", err)
	}
	// Check name uniqueness
	exists, err := s.institutionAdministratorPositionRepo.ExistsByName(ctx, &dto.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to check name uniqueness: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("institutionAdministratorPosition with name '%s' already exists", dto.Name)
	}
	// Creating model object
	institutionAdministratorPosition := &model.InstitutionAdministratorPosition{
		Name: dto.Name,
	}
	// Create institutionAdministratorPosition
	if err := s.institutionAdministratorPositionRepo.Create(ctx, institutionAdministratorPosition); err != nil {
		return nil, fmt.Errorf("failed to create institutionAdministratorPosition: %w", err)
	}
	// Get created institutionAdministratorPosition for response
	createdInstitutionAdministratorPosition, err := s.institutionAdministratorPositionRepo.FindByID(ctx, &institutionAdministratorPosition.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch created institutionAdministratorPosition: %w", err)
	}
	return InstitutionAdministratorPositionToDTO(createdInstitutionAdministratorPosition), nil
}

func (s *institutionAdministratorPositionService) GetPositions(ctx context.Context, filter repository.InstitutionAdministratorPositionFilter) ([]InstitutionAdministratorPositionResponseDTO, error) {
	institutionAdministratorPositions, err := s.institutionAdministratorPositionRepo.FindAll(ctx, &filter)
	if err != nil {
		s.log.Error(
			"failed to get institutionAdministratorPositions from repository",
			"limit",
			filter.Limit,
			"offset",
			filter.Offset,
			"error",
			err,
		)
		return nil, fmt.Errorf(
			"failed to get institutionAdministratorPositions from repository (limit: %d, offset: %d): %w",
			filter.Limit,
			filter.Offset,
			err,
		)
	}
	institutionAdministratorPositionDTOs := make([]InstitutionAdministratorPositionResponseDTO, len(institutionAdministratorPositions))
	for i, institutionAdministratorPosition := range institutionAdministratorPositions {
		institutionAdministratorPositionDTOs[i] = *InstitutionAdministratorPositionToDTO(&institutionAdministratorPosition)
	}
	s.log.Info("successfully received the list of institutionAdministratorPositions")
	return institutionAdministratorPositionDTOs, nil
}

func (s *institutionAdministratorPositionService) UpdatePosition(ctx context.Context, id uint8, dto UpdateInstitutionAdministratorPositionDTO) (*InstitutionAdministratorPositionResponseDTO, error) {
	// Input data validation
	if err := s.validateUpdatePositionDTO(&dto); err != nil {
		return nil, fmt.Errorf("validation error during institutionAdministratorPosition updating: %w", err)
	}
	// Getting existing institutionAdministratorPosition
	institutionAdministratorPosition, err := s.institutionAdministratorPositionRepo.FindByID(ctx, &id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.log.Error("InstitutionAdministratorPosition for update was not found by id", "institutionAdministratorPosition id", id, "error", err)
			return nil, fmt.Errorf("institutionAdministratorPosition with id %s was not found: %w", id, err)
		}
		s.log.Error("Failed to get institutionAdministratorPosition for update", "institutionAdministratorPosition id", id, "error", err)
		return nil, fmt.Errorf("failed to get institutionAdministratorPosition for update: %w", err)
	}
	// Update field if provided and was changed
	if dto.Name != nil && *dto.Name != institutionAdministratorPosition.Name {
		// Check name uniqueness
		exists, err := s.institutionAdministratorPositionRepo.ExistsByName(ctx, dto.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to check name uniqueness: %w", err)
		}
		if exists {
			return nil, fmt.Errorf("institutionAdministratorPosition with name '%s' already exists", dto.Name)
		}
		institutionAdministratorPosition.Name = *dto.Name
	} else {
		// No changes to update, return existing institutionAdministratorPosition
		s.log.Info("No changes to update institutionAdministratorPosition", "institutionAdministratorPosition ID", id)
		return InstitutionAdministratorPositionToDTO(institutionAdministratorPosition), nil
	}
	// Update institutionAdministratorPosition in DB
	if err := s.institutionAdministratorPositionRepo.Update(ctx, institutionAdministratorPosition); err != nil {
		s.log.Error("Failed to update the institutionAdministratorPosition")
		return nil, fmt.Errorf("failed to update the institutionAdministratorPosition: %w", err)
	}
	// Get updated institutionAdministratorPosition for response
	updatedInstitutionAdministratorPosition, err := s.institutionAdministratorPositionRepo.FindByID(ctx, &institutionAdministratorPosition.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated institutionAdministratorPosition: %w", err)
	}
	return InstitutionAdministratorPositionToDTO(updatedInstitutionAdministratorPosition), nil
}

func (s *institutionAdministratorPositionService) DeletePosition(ctx context.Context, id uint8) error {
	s.log.Info("Starting institutionAdministratorPosition deletion...")
	if err := s.institutionAdministratorPositionRepo.Delete(ctx, &id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.log.Error("InstitutionAdministratorPosition does not exist", "institutionAdministratorPosition id", id, "error", err)
			return fmt.Errorf("institutionAdministratorPosition with id %d does not exist: %w", id, err)
		}
		s.log.Error("Failed to delete the institutionAdministratorPosition")
		return fmt.Errorf("failed to delete the institutionAdministratorPosition: %w", err)
	}
	s.log.Info("InstitutionAdministratorPosition deleted successfully")
	return nil
}

func (s *institutionAdministratorPositionService) validateCreatePositionDTO(dto *CreateInstitutionAdministratorPositionDTO) error {
	if strings.TrimSpace(dto.Name) == "" {
		return fmt.Errorf("name cannot be empty or only whitespace")
	}
	if len(dto.Name) < 4 {
		return fmt.Errorf("name must be at least 4 characters")
	}
	if len(dto.Name) > 100 {
		return fmt.Errorf("name must be at most 100 characters")
	}
	return nil
}

func (s *institutionAdministratorPositionService) validateUpdatePositionDTO(dto *UpdateInstitutionAdministratorPositionDTO) error {
	if dto.Name != nil {
		if strings.TrimSpace(*dto.Name) == "" {
			return fmt.Errorf("name cannot be only whitespace")
		}
		if len(*dto.Name) < 4 {
			return fmt.Errorf("name must be at least 4 characters")
		}
		if len(*dto.Name) > 100 {
			return fmt.Errorf("name must be at most 100 characters")
		}
	}
	return nil
}

func InstitutionAdministratorPositionToDTO(institutionAdministratorPosition *model.InstitutionAdministratorPosition) *InstitutionAdministratorPositionResponseDTO {
	return &InstitutionAdministratorPositionResponseDTO{
		ID:        institutionAdministratorPosition.ID,
		CreatedAt: institutionAdministratorPosition.CreatedAt.Format(time.RFC3339),
		UpdatedAt: institutionAdministratorPosition.UpdatedAt.Format(time.RFC3339),
		Name:      institutionAdministratorPosition.Name,
	}
}
