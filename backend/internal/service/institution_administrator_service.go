package service

import (
	"backend/pkg/apperrors"
	"backend/internal/model"
	"backend/internal/repository"
	"backend/pkg/logger"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InstitutionAdministratorService interface {
	GetInstitutionAdministratorByID(ctx context.Context, id uuid.UUID) (*InstitutionAdministratorResponseDTO, error)
	GetInstitutionAdministrator(ctx context.Context, filter repository.InstitutionAdministratorFilter) ([]InstitutionAdministratorResponseDTO, error)
	// Position
	AssignPosition(ctx context.Context, userID uuid.UUID, positionID uint8) error
	GetPosition(ctx context.Context, userID uuid.UUID) (*InstitutionAdministratorPositionResponseDTO, error)
}

type InstitutionAdministratorResponseDTO struct {
	UserID   string                                      `json:"userId"`
	Position InstitutionAdministratorPositionResponseDTO `json:"position"`
}

type institutionAdministratorService struct {
	institutionAdministratorRepo repository.InstitutionAdministratorRepository
	userRepo                     repository.UserRepository
	db                           *gorm.DB
	log                          logger.Logger
}

func NewInstitutionAdministratorService(
	institutionAdministratorRepo repository.InstitutionAdministratorRepository,
	userRepo repository.UserRepository,
	db *gorm.DB,
	log logger.Logger,
) InstitutionAdministratorService {
	return &institutionAdministratorService{
		institutionAdministratorRepo: institutionAdministratorRepo,
		userRepo:                     userRepo,
		db:                           db,
		log:                          log,
	}
}

func (s *institutionAdministratorService) GetInstitutionAdministratorByID(ctx context.Context, id uuid.UUID) (*InstitutionAdministratorResponseDTO, error) {
	institutionAdministrator, err := s.institutionAdministratorRepo.FindByID(ctx, &id)
	if err != nil {
		return nil, fmt.Errorf("failed to get institution administrator: %w", err)
	}
	return InstitutionAdministratorToDTO(institutionAdministrator), nil
}

func (s *institutionAdministratorService) GetInstitutionAdministrator(ctx context.Context, filter repository.InstitutionAdministratorFilter) ([]InstitutionAdministratorResponseDTO, error) {
	institutionAdministrators, err := s.institutionAdministratorRepo.FindAll(ctx, &filter)
	if err != nil {
		s.log.Error(
			"failed to get institution administrator from repository",
			"limit",
			filter.Limit,
			"offset",
			filter.Offset,
			"error",
			err,
		)
		return nil, fmt.Errorf(
			"failed to get institution administrator from repository (limit: %d, offset: %d): %w",
			filter.Limit,
			filter.Offset,
			err,
		)
	}
	institutionAdministratorDTOs := make([]InstitutionAdministratorResponseDTO, len(institutionAdministrators))
	for i, institutionAdministrator := range institutionAdministrators {
		institutionAdministratorDTOs[i] = *InstitutionAdministratorToDTO(&institutionAdministrator)
	}
	s.log.Info("successfully received the list of institutionAdministrator")
	return institutionAdministratorDTOs, nil
}

func InstitutionAdministratorToDTO(institutionAdministrator *model.InstitutionAdministrator) *InstitutionAdministratorResponseDTO {
	return &InstitutionAdministratorResponseDTO{
		UserID:   institutionAdministrator.UserID.String(),
		Position: *InstitutionAdministratorPositionToDTO(&institutionAdministrator.Position),
	}
}

func (s *institutionAdministratorService) AssignPosition(ctx context.Context, userID uuid.UUID, positionID uint8) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Get institution administrator
		var institutionAdministrator model.InstitutionAdministrator
		if err := tx.WithContext(ctx).
			First(&institutionAdministrator, "user_id = ?", userID).Error; err != nil {
				return fmt.Errorf("institution administrator with user ID %s was not found: %s: %w", userID, err.Error(), apperrors.ErrNotFound)
		}
		// Check position existence
		var count int64
		if err := tx.WithContext(ctx).
			Model(&model.InstitutionAdministratorPosition{}).
			Where("id = ?", positionID).
			Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			return fmt.Errorf("position with ID %d was not found: %w", positionID, apperrors.ErrNotFound)
		}
		institutionAdministrator.PositionID = positionID
		if err := tx.WithContext(ctx).Save(&institutionAdministrator).Error; err != nil {
			return fmt.Errorf("failed to assign position to institution administrator: %w", err)
		}
		s.log.Info("position was successfully assigned to the institution administrator")
		return nil
	})
}

func (s *institutionAdministratorService) GetPosition(ctx context.Context, userID uuid.UUID) (*InstitutionAdministratorPositionResponseDTO, error) {
	institutionAdministrator, err := s.institutionAdministratorRepo.FindByID(ctx, &userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get institution administrator: %w", err)
	}
	if institutionAdministrator.Position.ID == 0 {
		return nil, fmt.Errorf("institution administrator has no position assigned")
	}
	return InstitutionAdministratorPositionToDTO(&institutionAdministrator.Position), nil
}
