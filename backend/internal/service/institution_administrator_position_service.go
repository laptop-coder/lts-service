package service

import (
	"backend/internal/model"
	"backend/internal/repository"
	"backend/pkg/apperrors"
	"backend/pkg/logger"
	"context"
	"fmt"
	"gorm.io/gorm"
	"strings"
	"time"
)

type InstitutionAdministratorPositionService interface {
	CreatePosition(ctx context.Context, dto CreateInstitutionAdministratorPositionDTO) (*InstitutionAdministratorPositionResponseDTO, error)
	GetPositions(ctx context.Context, filter repository.InstitutionAdministratorPositionFilter) ([]InstitutionAdministratorPositionResponseDTO, error)
	UpdatePosition(ctx context.Context, id uint16, dto UpdateInstitutionAdministratorPositionDTO) (*InstitutionAdministratorPositionResponseDTO, error)
	DeletePosition(ctx context.Context, id uint16) error
}

type CreateInstitutionAdministratorPositionDTO struct {
	Name string `form:"name" validate:"required,min=4,max=100"`
}

type UpdateInstitutionAdministratorPositionDTO struct {
	Name *string `form:"name,omitempty" validate:"max=100"`
}

type InstitutionAdministratorPositionResponseDTO struct {
	ID        uint16 `json:"id"`
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
		return nil, fmt.Errorf("institutionAdministratorPosition with name '%s' already exists: %w", dto.Name, apperrors.ErrInstitutionAdministratorPositionAlreadyExists)
	}
	// Creating model object
	institutionAdministratorPosition := &model.InstitutionAdministratorPosition{
		Name: dto.Name,
	}
	// Create institutionAdministratorPosition
	if err := s.institutionAdministratorPositionRepo.Create(ctx, institutionAdministratorPosition); err != nil {
		return nil, fmt.Errorf("failed to create institution administrator position: %w", err)
	}
	// Get created institutionAdministratorPosition for response
	createdInstitutionAdministratorPosition, err := s.institutionAdministratorPositionRepo.FindByID(ctx, &institutionAdministratorPosition.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch created institution administrator position: %w", err)
	}
	return InstitutionAdministratorPositionToDTO(createdInstitutionAdministratorPosition), nil
}

func (s *institutionAdministratorPositionService) GetPositions(ctx context.Context, filter repository.InstitutionAdministratorPositionFilter) ([]InstitutionAdministratorPositionResponseDTO, error) {
	institutionAdministratorPositions, err := s.institutionAdministratorPositionRepo.FindAll(ctx, &filter)
	if err != nil {
		s.log.Error(
			"failed to get institution administrator positions from repository",
			"limit",
			filter.Limit,
			"offset",
			filter.Offset,
			"error",
			err,
		)
		return nil, fmt.Errorf(
			"failed to get institution administrator positions from repository (limit: %d, offset: %d): %w",
			filter.Limit,
			filter.Offset,
			err,
		)
	}
	institutionAdministratorPositionDTOs := make([]InstitutionAdministratorPositionResponseDTO, len(institutionAdministratorPositions))
	for i, institutionAdministratorPosition := range institutionAdministratorPositions {
		institutionAdministratorPositionDTOs[i] = *InstitutionAdministratorPositionToDTO(&institutionAdministratorPosition)
	}
	s.log.Info("successfully received the list of institution administrator positions")
	return institutionAdministratorPositionDTOs, nil
}

func (s *institutionAdministratorPositionService) UpdatePosition(ctx context.Context, id uint16, dto UpdateInstitutionAdministratorPositionDTO) (*InstitutionAdministratorPositionResponseDTO, error) {
	// Input data validation
	if err := s.validateUpdatePositionDTO(&dto); err != nil {
		return nil, fmt.Errorf("validation error during institution administrator position updating: %w", err)
	}
	// Getting existing institutionAdministratorPosition
	institutionAdministratorPosition, err := s.institutionAdministratorPositionRepo.FindByID(ctx, &id)
	if err != nil {
		return nil, fmt.Errorf("failed to get institution administrator position for update: %w", err)
	}
	// Update field if provided and was changed
	if dto.Name != nil && *dto.Name != institutionAdministratorPosition.Name {
		// Check name uniqueness
		exists, err := s.institutionAdministratorPositionRepo.ExistsByName(ctx, dto.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to check name uniqueness: %w", err)
		}
		if exists {
			return nil, fmt.Errorf("institution administrator position with name '%s' already exists: %w", dto.Name, apperrors.ErrInstitutionAdministratorPositionAlreadyExists)
		}
		institutionAdministratorPosition.Name = *dto.Name
	} else {
		// No changes to update, return existing institutionAdministratorPosition
		s.log.Info("no changes to update institution administrator position", "position_id", id)
		return InstitutionAdministratorPositionToDTO(institutionAdministratorPosition), nil
	}
	// Update institutionAdministratorPosition in DB
	if err := s.institutionAdministratorPositionRepo.Update(ctx, institutionAdministratorPosition); err != nil {
		s.log.Error("failed to update the institution administrator position")
		return nil, fmt.Errorf("failed to update the institution administrator position: %w", err)
	}
	// Get updated institutionAdministratorPosition for response
	updatedInstitutionAdministratorPosition, err := s.institutionAdministratorPositionRepo.FindByID(ctx, &institutionAdministratorPosition.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated institution administrator position: %w", err)
	}
	return InstitutionAdministratorPositionToDTO(updatedInstitutionAdministratorPosition), nil
}

func (s *institutionAdministratorPositionService) DeletePosition(ctx context.Context, id uint16) error {
	s.log.Info("starting institution administrator position deletion...")
	if err := s.institutionAdministratorPositionRepo.Delete(ctx, &id); err != nil {
		s.log.Error("failed to delete the institution administrator position")
		return fmt.Errorf("failed to delete the institution administrator position: %w", err)
	}
	s.log.Info("institution administrator position deleted successfully")
	return nil
}

func (s *institutionAdministratorPositionService) validateCreatePositionDTO(dto *CreateInstitutionAdministratorPositionDTO) error {
	if strings.TrimSpace(dto.Name) == "" {
		return fmt.Errorf("name cannot be empty or only whitespace: %w", apperrors.ErrRequiredField)
	}
	if len(dto.Name) < 4 {
		return fmt.Errorf("name must be at least 4 characters: %w", apperrors.ErrValueTooShort)
	}
	if len(dto.Name) > 100 {
		return fmt.Errorf("name must be at most 100 characters: %w", apperrors.ErrValueTooLong)
	}
	return nil
}

func (s *institutionAdministratorPositionService) validateUpdatePositionDTO(dto *UpdateInstitutionAdministratorPositionDTO) error {
	if dto.Name != nil {
		if strings.TrimSpace(*dto.Name) == "" {
			return fmt.Errorf("name cannot be only whitespace: %w", apperrors.ErrRequiredField)
		}
		if len(*dto.Name) < 4 {
			return fmt.Errorf("name must be at least 4 characters: %w", apperrors.ErrValueTooShort)
		}
		if len(*dto.Name) > 100 {
			return fmt.Errorf("name must be at most 100 characters: %w", apperrors.ErrValueTooLong)
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
