package service

import (
	"backend/internal/model"
	"backend/internal/repository"
	"backend/pkg/logger"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type InstitutionAdministratorService interface {
	GetInstitutionAdministratorByID(ctx context.Context, id uuid.UUID) (*InstitutionAdministratorResponseDTO, error)
	GetInstitutionAdministrator(ctx context.Context, filter repository.InstitutionAdministratorFilter) ([]InstitutionAdministratorResponseDTO, error)
}

type InstitutionAdministratorResponseDTO struct {
	User     UserResponseDTO
	Position InstitutionAdministratorPositionResponseDTO
}

type InstitutionAdministratorPositionResponseDTO struct {
	ID        uint8  `json:"id"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	Name      string `json:"name"`
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("institutionAdministrator with id %s was not found: %w", id, err)
		}
		return nil, fmt.Errorf("failed to get institutionAdministrator: %w", err)
	}
	return InstitutionAdministratorToDTO(institutionAdministrator), nil
}

func (s *institutionAdministratorService) GetInstitutionAdministrator(ctx context.Context, filter repository.InstitutionAdministratorFilter) ([]InstitutionAdministratorResponseDTO, error) {
	institutionAdministrators, err := s.institutionAdministratorRepo.FindAll(ctx, &filter)
	if err != nil {
		s.log.Error(
			"failed to get institutionAdministrator from repository",
			"limit",
			filter.Limit,
			"offset",
			filter.Offset,
			"error",
			err,
		)
		return nil, fmt.Errorf(
			"failed to get institutionAdministrator from repository (limit: %d, offset: %d): %w",
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
		User:     *UserToDTO(&institutionAdministrator.User),
		Position: *InstitutionAdministratorPositionToDTO(&institutionAdministrator.Position),
	}
}

func InstitutionAdministratorPositionToDTO(position *model.InstitutionAdministratorPosition) *InstitutionAdministratorPositionResponseDTO {
	return &InstitutionAdministratorPositionResponseDTO{
		ID:        position.ID,
		CreatedAt: position.CreatedAt.Format(time.RFC3339),
		UpdatedAt: position.UpdatedAt.Format(time.RFC3339),
		Name:      position.Name,
	}
}
