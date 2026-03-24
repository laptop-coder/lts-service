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

type StaffPositionService interface {
	CreatePosition(ctx context.Context, dto CreateStaffPositionDTO) (*StaffPositionResponseDTO, error)
	GetPositions(ctx context.Context, filter repository.StaffPositionFilter) ([]StaffPositionResponseDTO, error)
	UpdatePosition(ctx context.Context, id uint8, dto UpdateStaffPositionDTO) (*StaffPositionResponseDTO, error)
	DeletePosition(ctx context.Context, id uint8) error
}

type CreateStaffPositionDTO struct {
	Name string `form:"name" validate:"required,min=4,max=100"`
}

type UpdateStaffPositionDTO struct {
	Name *string `form:"name,omitempty" validate:"max=100"`
}

type StaffPositionResponseDTO struct {
	ID        uint8  `json:"id"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	Name      string `json:"name"`
}

type staffPositionService struct {
	staffPositionRepo repository.StaffPositionRepository
	db          *gorm.DB
	log         logger.Logger
}

func NewStaffPositionService(
	staffPositionRepo repository.StaffPositionRepository,
	db *gorm.DB,
	log logger.Logger,
) StaffPositionService {
	return &staffPositionService{
		staffPositionRepo: staffPositionRepo,
		db:          db,
		log:         log,
	}
}

func (s *staffPositionService) CreatePosition(ctx context.Context, dto CreateStaffPositionDTO) (*StaffPositionResponseDTO, error) {
	// Input data validation
	if err := s.validateCreatePositionDTO(&dto); err != nil {
		return nil, fmt.Errorf("validation error during staffPosition creation: %w", err)
	}
	// Check name uniqueness
	exists, err := s.staffPositionRepo.ExistsByName(ctx, &dto.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to check name uniqueness: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("staffPosition with name '%s' already exists", dto.Name)
	}
	// Creating model object
	staffPosition := &model.StaffPosition{
		Name: dto.Name,
	}
	// Create staffPosition
	if err := s.staffPositionRepo.Create(ctx, staffPosition); err != nil {
		return nil, fmt.Errorf("failed to create staffPosition: %w", err)
	}
	// Get created staffPosition for response
	createdStaffPosition, err := s.staffPositionRepo.FindByID(ctx, &staffPosition.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch created staffPosition: %w", err)
	}
	return StaffPositionToDTO(createdStaffPosition), nil
}

func (s *staffPositionService) GetPositions(ctx context.Context, filter repository.StaffPositionFilter) ([]StaffPositionResponseDTO, error) {
	staffPositions, err := s.staffPositionRepo.FindAll(ctx, &filter)
	if err != nil {
		s.log.Error(
			"failed to get staffPositions from repository",
			"limit",
			filter.Limit,
			"offset",
			filter.Offset,
			"error",
			err,
		)
		return nil, fmt.Errorf(
			"failed to get staffPositions from repository (limit: %d, offset: %d): %w",
			filter.Limit,
			filter.Offset,
			err,
		)
	}
	staffPositionDTOs := make([]StaffPositionResponseDTO, len(staffPositions))
	for i, staffPosition := range staffPositions {
		staffPositionDTOs[i] = *StaffPositionToDTO(&staffPosition)
	}
	s.log.Info("successfully received the list of staffPositions")
	return staffPositionDTOs, nil
}

func (s *staffPositionService) UpdatePosition(ctx context.Context, id uint8, dto UpdateStaffPositionDTO) (*StaffPositionResponseDTO, error) {
	// Input data validation
	if err := s.validateUpdatePositionDTO(&dto); err != nil {
		return nil, fmt.Errorf("validation error during staffPosition updating: %w", err)
	}
	// Getting existing staffPosition
	staffPosition, err := s.staffPositionRepo.FindByID(ctx, &id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.log.Error("StaffPosition for update was not found by id", "staffPosition id", id, "error", err)
			return nil, fmt.Errorf("staffPosition with id %s was not found: %w", id, err)
		}
		s.log.Error("Failed to get staffPosition for update", "staffPosition id", id, "error", err)
		return nil, fmt.Errorf("failed to get staffPosition for update: %w", err)
	}
	// Update field if provided and was changed
	if dto.Name != nil && *dto.Name != staffPosition.Name {
		// Check name uniqueness
		exists, err := s.staffPositionRepo.ExistsByName(ctx, dto.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to check name uniqueness: %w", err)
		}
		if exists {
			return nil, fmt.Errorf("staffPosition with name '%s' already exists", dto.Name)
		}
		staffPosition.Name = *dto.Name
	} else {
		// No changes to update, return existing staffPosition
		s.log.Info("No changes to update staffPosition", "staffPosition ID", id)
		return StaffPositionToDTO(staffPosition), nil
	}
	// Update staffPosition in DB
	if err := s.staffPositionRepo.Update(ctx, staffPosition); err != nil {
		s.log.Error("Failed to update the staffPosition")
		return nil, fmt.Errorf("failed to update the staffPosition: %w", err)
	}
	// Get updated staffPosition for response
	updatedStaffPosition, err := s.staffPositionRepo.FindByID(ctx, &staffPosition.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated staffPosition: %w", err)
	}
	return StaffPositionToDTO(updatedStaffPosition), nil
}

func (s *staffPositionService) DeletePosition(ctx context.Context, id uint8) error {
	s.log.Info("Starting staffPosition deletion...")
	if err := s.staffPositionRepo.Delete(ctx, &id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.log.Error("StaffPosition does not exist", "staffPosition id", id, "error", err)
			return fmt.Errorf("staffPosition with id %d does not exist: %w", id, err)
		}
		s.log.Error("Failed to delete the staffPosition")
		return fmt.Errorf("failed to delete the staffPosition: %w", err)
	}
	s.log.Info("StaffPosition deleted successfully")
	return nil
}

func (s *staffPositionService) validateCreatePositionDTO(dto *CreateStaffPositionDTO) error {
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

func (s *staffPositionService) validateUpdatePositionDTO(dto *UpdateStaffPositionDTO) error {
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

func StaffPositionToDTO(staffPosition *model.StaffPosition) *StaffPositionResponseDTO {
	return &StaffPositionResponseDTO{
		ID:        staffPosition.ID,
		CreatedAt: staffPosition.CreatedAt.Format(time.RFC3339),
		UpdatedAt: staffPosition.UpdatedAt.Format(time.RFC3339),
		Name:      staffPosition.Name,
	}
}
