package service

import (
	"backend/internal/repository"
	"backend/pkg/logger"
	"context"
	"fmt"
	"gorm.io/gorm"
)

type StudentGroupService interface {
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
