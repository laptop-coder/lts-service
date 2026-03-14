package service

import (
	"backend/internal/repository"
	"backend/pkg/logger"
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type StudentGroupService interface {
	GetStudentGroupByID(ctx context.Context, id uint16) (*StudentGroupResponseDTO, error)
	GetAdvisorByGroupID(ctx context.Context, id uint16) (*UserResponseDTO, error)
}

type studentGroupService struct {
	studentGroupRepo repository.StudentGroupRepository
	db               *gorm.DB
	log              logger.Logger
}

func NewStudentGroupService(
	studentGroupRepo repository.StudentGroupRepository,
	db *gorm.DB,
	log logger.Logger,
) StudentGroupService {
	return &studentGroupService{
		studentGroupRepo: studentGroupRepo,
		db:               db,
		log:              log,
	}
}

func (s *studentGroupService) GetStudentGroupByID(ctx context.Context, id uint16) (*StudentGroupResponseDTO, error) {
	group, err := s.studentGroupRepo.FindByID(ctx, &id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("student group with id %d was not found: %w", id, err)
		}
		return nil, fmt.Errorf("failed to get student group: %w", err)
	}
	return StudentGroupToDTO(group), nil
}

func (s *studentGroupService) GetAdvisorByGroupID(ctx context.Context, id uint16) (*UserResponseDTO, error) {
	user, err := s.studentGroupRepo.FindAdvisorByGroupID(ctx, &id)
	if err != nil {
		s.log.Error(
			"failed to get student group advisor by group id",
			"group id", id,
			"error", err,
		)
		return nil, fmt.Errorf("failed to get student group advisor by group id (%d): %w", id, err)
	}
	s.log.Info("successfully received student group advisor")
	return UserToDTO(user), nil
}
