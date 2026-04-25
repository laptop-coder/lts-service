package service

import (
	"backend/internal/model"
	"backend/internal/repository"
	"backend/pkg/apperrors"
	"backend/pkg/logger"
	"context"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
	"time"
)

type StudentGroupService interface {
	CreateStudentGroup(ctx context.Context, dto CreateStudentGroupDTO) (*StudentGroupResponseDTO, error)
	GetStudentGroups(ctx context.Context, filter repository.StudentGroupFilter) ([]StudentGroupResponseDTO, error)
	GetStudentGroupByID(ctx context.Context, id uint16) (*StudentGroupResponseDTO, error)
	GetAdvisorByGroupID(ctx context.Context, id uint16) (*UserResponseDTO, error)
	UpdateStudentGroup(ctx context.Context, id uint16, dto UpdateStudentGroupDTO) (*StudentGroupResponseDTO, error)
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
	userRepo repository.UserRepository,
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

type CreateStudentGroupDTO struct {
	Name           string     `form:"name" validate:"required,min=1,max=20"`
	GroupAdvisorID *uuid.UUID `form:"advisorId,omitempty"`
}

type UpdateStudentGroupDTO struct {
	Name           *string    `form:"name,omitempty" validate:"omitempty,min=1,max=20"`
	GroupAdvisorID *uuid.UUID `form:"advisorId,omitempty"`
}

type StudentGroupResponseDTO struct {
	ID             uint16                           `json:"id"`
	CreatedAt      string                           `json:"createdAt"`
	UpdatedAt      string                           `json:"updatedAt"`
	Name           string                           `json:"name"`
	GroupAdvisorID *uuid.UUID                       `json:"advisorId,omitempty"`
	Students       []StudentGroupStudentResponseDTO `json:"students"`
}

func StudentGroupToDTO(studentGroup *model.StudentGroup) *StudentGroupResponseDTO {
	var students []StudentGroupStudentResponseDTO
	for _, student := range studentGroup.Students {
		students = append(students, *StudentGroupStudentToDTO(&student))
	}
	return &StudentGroupResponseDTO{
		CreatedAt: studentGroup.CreatedAt.Format(time.RFC3339),
		UpdatedAt: studentGroup.UpdatedAt.Format(time.RFC3339),
		ID:        studentGroup.ID,
		Name:      studentGroup.Name,
		Students:  students,
	}
}

func (s *studentGroupService) CreateStudentGroup(ctx context.Context, dto CreateStudentGroupDTO) (*StudentGroupResponseDTO, error) {
	// Input data validation
	if err := s.validateCreateStudentGroupDTO(&dto); err != nil {
		return nil, fmt.Errorf("validation error during student group creation: %w", err)
	}
	// Check name uniqueness
	exists, err := s.studentGroupRepo.ExistsByName(ctx, &dto.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to check name uniqueness: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("student group with name %s already exists: %w", dto.Name, apperrors.ErrStudentGroupAlreadyExists)
	}
	// Check if advisor exists and has the teacher role (if provided)
	if dto.GroupAdvisorID != nil {
		user, err := s.userRepo.FindByID(ctx, dto.GroupAdvisorID)
		if err != nil {
			return nil, fmt.Errorf("user with ID %s was not found: %w", dto.GroupAdvisorID, err)
		}
		isTeacher := false
		for _, role := range user.Roles {
			if role.Name == "teacher" {
				isTeacher = true
				break
			}
		}
		if !isTeacher {
			return nil, fmt.Errorf("user with ID %s is not a teacher and cannot be student group advisor: %w", dto.GroupAdvisorID, apperrors.ErrForbidden)
		}
	}
	// Create student group
	group := &model.StudentGroup{
		Name:           dto.Name,
		GroupAdvisorID: dto.GroupAdvisorID,
	}
	if err := s.studentGroupRepo.Create(ctx, group); err != nil {
		return nil, fmt.Errorf("failed to create student group: %w", err)
	}
	// Get created group for response
	createdGroup, err := s.studentGroupRepo.FindByID(ctx, &group.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch create student group: %w", err)
	}
	s.log.Info("student group was created successfully", "group ID", createdGroup.ID, "group name", createdGroup.Name)
	return StudentGroupToDTO(createdGroup), nil
}

func (s *studentGroupService) UpdateStudentGroup(ctx context.Context, id uint16, dto UpdateStudentGroupDTO) (*StudentGroupResponseDTO, error) {
	// Input data validation
	if err := s.validateUpdateStudentGroupDTO(&dto); err != nil {
		return nil, fmt.Errorf("validation error during student group update: %w", err)
	}
	// Get existing group
	group, err := s.studentGroupRepo.FindByID(ctx, &id)
	if err != nil {
		return nil, fmt.Errorf("failed to get student group with ID %s: %w", id, err)
	}
	// Track updated fields
	updatedFields := make([]string, 0)
	// Update name if provided and changed
	if dto.Name != nil && *dto.Name != group.Name {
		// Check name uniqueness
		exists, err := s.studentGroupRepo.ExistsByName(ctx, dto.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to check name uniqueness: %w", err)
		}
		if exists {
			return nil, fmt.Errorf("student group with name %s already exists: %w", dto.Name, apperrors.ErrStudentGroupAlreadyExists)
		}
		group.Name = *dto.Name
		updatedFields = append(updatedFields, "name")
	}
	// Update advisor ID if provided and changed
	if dto.GroupAdvisorID != nil &&
		(group.GroupAdvisorID == nil || (*dto.GroupAdvisorID != *group.GroupAdvisorID)) {
		// Check if advisor exists and has the teacher role
		user, err := s.userRepo.FindByID(ctx, dto.GroupAdvisorID)
		if err != nil {
			return nil, err
		}
		isTeacher := false
		for _, role := range user.Roles {
			if role.Name == "teacher" {
				isTeacher = true
				break
			}
		}
		if !isTeacher {
			return nil, fmt.Errorf("user with ID %s is not a teacher and cannot be student group advisor: %w", dto.GroupAdvisorID, apperrors.ErrForbidden)
		}
		group.GroupAdvisorID = dto.GroupAdvisorID
		updatedFields = append(updatedFields, "advisor ID")
	}
	// No changes to update, return existing group
	if len(updatedFields) == 0 {
		s.log.Info("No changes to update student group", "group ID", id)
		return StudentGroupToDTO(group), nil
	}
	// Update student group in DB
	if err := s.studentGroupRepo.Update(ctx, group); err != nil {
		return nil, fmt.Errorf("failed to update student group: %w", err)
	}
	// Get updated group for response
	updatedGroup, err := s.studentGroupRepo.FindByID(ctx, &group.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated student group: %w", err)
	}
	s.log.Info("student group was updated successfully", "group ID", id, "updated fields", updatedFields)
	return StudentGroupToDTO(updatedGroup), nil
}

func (s *studentGroupService) validateCreateStudentGroupDTO(dto *CreateStudentGroupDTO) error {
	if strings.TrimSpace(dto.Name) == "" {
		return fmt.Errorf("name cannot be empty or only whitespace: %w", apperrors.ErrRequiredField)
	}
	if len(dto.Name) < 1 {
		return fmt.Errorf("name must be at least 1 character: %w", apperrors.ErrValueTooShort)
	}
	if len(dto.Name) > 20 {
		return fmt.Errorf("name must be at most 20 characters: %w", apperrors.ErrValueTooLong)
	}
	return nil
}

func (s *studentGroupService) validateUpdateStudentGroupDTO(dto *UpdateStudentGroupDTO) error {
	if dto.Name != nil {
		if strings.TrimSpace(*dto.Name) == "" {
			return fmt.Errorf("name cannot be empty or only whitespace: %w", apperrors.ErrRequiredField)
		}
		if len(*dto.Name) < 1 {
			return fmt.Errorf("name must be at least 1 character: %w", apperrors.ErrValueTooShort)
		}
		if len(*dto.Name) > 20 {
			return fmt.Errorf("name must be at most 20 characters: %w", apperrors.ErrValueTooLong)
		}
	}
	return nil
}

func (s *studentGroupService) GetStudentGroupByID(ctx context.Context, id uint16) (*StudentGroupResponseDTO, error) {
	group, err := s.studentGroupRepo.FindByID(ctx, &id)
	if err != nil {
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
		return fmt.Errorf("failed to delete the student group: %w", err)
	}
	s.log.Info("student group deleted successfully")
	return nil
}

func (s *studentGroupService) AssignAdvisor(ctx context.Context, groupID uint16, userID uuid.UUID) error {
	// Get user
	user, err := s.userRepo.FindByID(ctx, &userID)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}
	// Get group
	// TODO: fix like here in other services/handlers, change direct access
	// to the DB to repository FindByID method (maybe search in the code by
	// /backend/internal/model import in services/handlers)
	group, err := s.studentGroupRepo.FindByID(ctx, &groupID)
	if err != nil {
		return fmt.Errorf("failed to find student group: %w", err)
	}
	// Check if group already has an advisor
	if group.GroupAdvisorID != nil {
		return fmt.Errorf("student group already has an advisor: %w", apperrors.ErrStudentGroupAlreadyHasAdvisor)
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
		return fmt.Errorf("to be student group advisor the teacher role is required: %w", apperrors.ErrForbidden)
	}
	// Update group advisor
	group.GroupAdvisorID = &userID
	if err := s.studentGroupRepo.Update(ctx, group); err != nil {
		return fmt.Errorf("failed to update student group advisor: %w", err)
	}
	s.log.Info("advisor was successfully assigned to student group", "group ID", groupID, "user ID", userID)
	return nil
}

func (s *studentGroupService) UnassignAdvisor(ctx context.Context, groupID uint16) error {
	// Get group
	group, err := s.studentGroupRepo.FindByID(ctx, &groupID)
	if err != nil {
		return fmt.Errorf("failed to find student group: %w", err)
	}
	// Unassign group advisor
	group.GroupAdvisorID = nil
	if err := s.studentGroupRepo.Update(ctx, group); err != nil {
		return fmt.Errorf("failed to unassign student group advisor: %w", err)
	}
	s.log.Info("advisor was successfully unassigned from student group", "group ID", groupID)
	return nil
}

type StudentGroupStudentResponseDTO struct {
	UserID uuid.UUID `json:"userId"`
}

func StudentGroupStudentToDTO(student *model.Student) *StudentGroupStudentResponseDTO {
	return &StudentGroupStudentResponseDTO{
		UserID: student.UserID,
	}
}
