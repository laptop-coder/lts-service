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
)

type ParentService interface {
	GetParentByID(ctx context.Context, id uuid.UUID) (*ParentResponseDTO, error)
	GetParents(ctx context.Context, filter repository.ParentFilter) ([]ParentResponseDTO, error)
}

type ParentResponseDTO struct {
	User          UserResponseDTO
	ParentGroupID uint16
}

type parentService struct {
	parentRepo repository.ParentRepository
	userRepo   repository.UserRepository
	db         *gorm.DB
	log        logger.Logger
}

func NewParentService(
	parentRepo repository.ParentRepository,
	userRepo repository.UserRepository,
	db *gorm.DB,
	log logger.Logger,
) ParentService {
	return &parentService{
		parentRepo: parentRepo,
		userRepo:   userRepo,
		db:         db,
		log:        log,
	}
}

func (s *parentService) GetParentByID(ctx context.Context, id uuid.UUID) (*ParentResponseDTO, error) {
	parent, err := s.parentRepo.FindByID(ctx, &id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("parent with id %s was not found: %w", id, err)
		}
		return nil, fmt.Errorf("failed to get parent: %w", err)
	}
	return ParentToDTO(parent), nil
}

func (s *parentService) GetParents(ctx context.Context, filter repository.ParentFilter) ([]ParentResponseDTO, error) {
	parents, err := s.parentRepo.FindAll(ctx, &filter)
	if err != nil {
		s.log.Error(
			"failed to get parents from repository",
			"limit",
			filter.Limit,
			"offset",
			filter.Offset,
			"error",
			err,
		)
		return nil, fmt.Errorf(
			"failed to get parents from repository (limit: %d, offset: %d): %w",
			filter.Limit,
			filter.Offset,
			err,
		)
	}
	parentDTOs := make([]ParentResponseDTO, len(parents))
	for i, parent := range parents {
		parentDTOs[i] = *ParentToDTO(&parent)
	}
	s.log.Info("successfully received the list of parents")
	return parentDTOs, nil
}

func ParentToDTO(parent *model.Parent) *ParentResponseDTO {
	return &ParentResponseDTO{
		User: *UserToDTO(&parent.User),
	}
}
