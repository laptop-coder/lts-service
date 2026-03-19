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
	GetStudentGroups(ctx context.Context, filter repository.StudentGroupFilter) ([]StudentGroupResponseDTO, error)
	GetStudentGroupByID(ctx context.Context, id uint16) (*StudentGroupResponseDTO, error)
	GetAdvisorByGroupID(ctx context.Context, id uint16) (*UserResponseDTO, error)
	DeleteStudentGroup(ctx context.Context, id uint16) error
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

func (s *studentGroupService) GetStudentGroups(ctx context.Context, filter repository.StudentGroupFilter) ([]StudentGroupResponseDTO, error) {
	studentGroups, err := s.studentGroupRepo.FindAll(ctx, &filter)
	if err != nil {
		groupAdvisorID := ""
		if filter.GroupAdvisorID != nil {
			groupAdvisorID = (*filter.GroupAdvisorID).String()
		}
		s.log.Error(
			"failed to get student groups from repository",
			"group advisor ID",
			groupAdvisorID,
			"limit",
			filter.Limit,
			"offset",
			filter.Offset,
			"error",
			err,
		)
		return nil, fmt.Errorf(
			"failed to get student groups from repository (groupAdvisorID: %s, limit: %d, offset: %d): %w",
			groupAdvisorID,
			filter.Limit,
			filter.Offset,
			err,
		)
	}
	studentGroupDTOs := make([]StudentGroupResponseDTO, len(studentGroups))
	for i, studentGroup := range studentGroups {
		studentGroupDTOs[i] = *StudentGroupToDTO(&studentGroup)
	}
	s.log.Info("successfully received the list of student groups")
	return studentGroupDTOs, nil
}

func (s *studentGroupService) DeleteStudentGroup(ctx context.Context, id uint16) error {
	s.log.Info("Starting student group deletion...")
	if err := s.studentGroupRepo.Delete(ctx, &id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.log.Error("Student group does not exist", "group id", id, "error", err)
			return fmt.Errorf("student group with id %d does not exist: %w", id, err)
		}
		s.log.Error("Failed to delete the student group")
		return fmt.Errorf("failed to delete the student group: %w", err)
	}
	s.log.Info("Student group deleted successfully")
	return nil
}
