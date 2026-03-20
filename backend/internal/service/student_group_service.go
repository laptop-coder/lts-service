package service

import (
	"backend/internal/repository"
	"backend/pkg/logger"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StudentGroupService interface {
	GetStudentGroups(ctx context.Context, filter repository.StudentGroupFilter) ([]StudentGroupResponseDTO, error)
	GetStudentGroupByID(ctx context.Context, id uint16) (*StudentGroupResponseDTO, error)
	GetAdvisorByGroupID(ctx context.Context, id uint16) (*UserResponseDTO, error)
	DeleteStudentGroup(ctx context.Context, id uint16) error
	AssignAdvisor(ctx context.Context, groupID uint16, userID uuid.UUID) error
	UnassignAdvisor(ctx context.Context, groupID uint16) error
}

type studentGroupService struct {
	userRepo         repository.UserRepository
	studentGroupRepo repository.StudentGroupRepository
	db               *gorm.DB
	log              logger.Logger
}

func NewStudentGroupService(
	userRepo         repository.UserRepository,
	studentGroupRepo repository.StudentGroupRepository,
	db *gorm.DB,
	log logger.Logger,
) StudentGroupService {
	return &studentGroupService{
		userRepo:         userRepo,
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

func (s *studentGroupService) AssignAdvisor(ctx context.Context, groupID uint16, userID uuid.UUID) error {
	// Get user
	user, err := s.userRepo.FindByID(ctx, &userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user with ID %s was not found", userID)
		}
		return fmt.Errorf("failed to find user: %w", err)
	}
	// Get group
	// TODO: fix like here in other services/handlers, change direct access
	// to the DB to repository FindByID method (maybe search in the code by
	// /backend/internal/model import in services/handlers)
	group, err := s.studentGroupRepo.FindByID(ctx, &groupID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("student group with ID %d was not found", groupID)
		}
		return fmt.Errorf("failed to find student group: %w", err)
	}
	// Check if user has a teacher role
	isTeacher := false
	for _, role := range user.Roles {
		if role.Name == "teacher" {
			isTeacher = true
			break
		}
	}
	if !isTeacher {
		return fmt.Errorf("forbidden: to be student group advisor the teacher role is required")
	}
	// Update group advisor
	group.GroupAdvisorID = &userID
	if err := s.studentGroupRepo.Update(ctx, group); err != nil {
		return fmt.Errorf("failed to update student group advisor: %w", err)
	}
	s.log.Info("Advisor was successfully assigned to student group", "group ID", groupID, "user ID", userID)
	return nil
}

func (s *studentGroupService) UnassignAdvisor(ctx context.Context, groupID uint16) error {
	// Get group
	group, err := s.studentGroupRepo.FindByID(ctx, &groupID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("student group with ID %d was not found", groupID)
		}
		return fmt.Errorf("failed to find student group: %w", err)
	}
	// Unassign group advisor
	group.GroupAdvisorID = nil
	if err := s.studentGroupRepo.Update(ctx, group); err != nil {
		return fmt.Errorf("failed to unassign student group advisor: %w", err)
	}
	s.log.Info("Advisor was successfully unassigned from student group", "group ID", groupID)
	return nil
}
