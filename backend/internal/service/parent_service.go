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
	GetParentStudents(ctx context.Context, userID uuid.UUID) ([]StudentResponseDTO, error)
	GetStudentGroupsOwn(ctx context.Context, userID uuid.UUID) ([]StudentGroupResponseDTO, error)
	AddStudents(ctx context.Context, userID uuid.UUID, studentIDs []uuid.UUID) error
	UnassignStudent(ctx context.Context, userID uuid.UUID, studentID uuid.UUID) error
}

type ParentResponseDTO struct {
	Students []StudentResponseDTO `json:"students"`
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

func (s *parentService) GetParentStudents(ctx context.Context, userID uuid.UUID) ([]StudentResponseDTO, error) {
	// Find parent by ID
	parent, err := s.parentRepo.FindByID(ctx, &userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("parent with id %s was not found: %w", userID, err)
		}
		return nil, fmt.Errorf("failed to get parent: %w", err)
	}
	// Check if there are connected students
	if parent.Students != nil && len(*parent.Students) > 0 {
		// Get students
		var students []StudentResponseDTO
		for _, student := range *parent.Students {
			students = append(students, *StudentToDTO(&student))
		}
		// Return response
		return students, nil
	}
	return []StudentResponseDTO{}, nil
}

func (s *parentService) GetStudentGroupsOwn(ctx context.Context, userID uuid.UUID) ([]StudentGroupResponseDTO, error) {
	// Get parent
	parent, err := s.GetParentByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get parent by id: %w", err)
	}
	// Check if students assigned
	if len(parent.Students) == 0 {
		return []StudentGroupResponseDTO{}, nil
	}
	// Get student groups of students assigned to parent
	studentGroups := make([]StudentGroupResponseDTO, len(parent.Students))
	for i, student := range parent.Students {
		studentGroups[i] = student.StudentGroup
	}
	return studentGroups, nil
}

func ParentToDTO(parent *model.Parent) *ParentResponseDTO {
	var students []StudentResponseDTO
	if parent.Students != nil && len(*parent.Students) > 0 {
		for _, student := range *parent.Students {
			students = append(students, *StudentToDTO(&student))
		}
	}
	return &ParentResponseDTO{
		Students: students,
	}
}

func (s *parentService) AddStudents(ctx context.Context, userID uuid.UUID, studentIDs []uuid.UUID) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Check parent existence
		var parent model.Parent
		if err := tx.WithContext(ctx).
			Preload("Students").
			Where("user_id = ?", userID).
			First(&parent).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("parent with ID %s was not found", userID)
			}
			return fmt.Errorf("failed to found parent: %w", err)
		}
		// Get students by IDs
		var students []model.Student
		if err := tx.WithContext(ctx).
			Where("user_id IN (?)", studentIDs).
			Find(&students).Error; err != nil {
			return fmt.Errorf("failed to fetch students: %w", err)
		}
		// Check if all students were found
		if len(students) != len(studentIDs) {
			return fmt.Errorf("some students not found")
		}
		// Init "Students" field if empty
		if parent.Students == nil {
			parent.Students = &[]model.Student{}
		}
		// Add students
		if err := tx.Model(&parent).Association("Students").Append(&students); err != nil {
			return fmt.Errorf("failed to add students: %w", err)
		}
		// Return response
		s.log.Info("Students was successfully added to parent", "parent ID", userID, "student IDs", studentIDs)
		return nil
	})
}

func (s *parentService) UnassignStudent(ctx context.Context, userID uuid.UUID, studentID uuid.UUID) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Check parent existence
		var parent model.Parent
		if err := tx.WithContext(ctx).
			Preload("Students").
			Where("user_id = ?", userID).
			First(&parent).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("parent with ID %s was not found", userID)
			}
			return fmt.Errorf("failed to found parent: %w", err)
		}
		// Get student by ID
		var student model.Student
		if err := tx.WithContext(ctx).
			Where("user_id = ?", studentID).
			First(&student).Error; err != nil {
			return fmt.Errorf("student not found: %w", err)
		}
		// Remove student from parent
		if err := tx.Model(&parent).Association("Students").Delete(&student); err != nil {
			return fmt.Errorf("failed to remove student: %w", err)
		}
		// Return response
		s.log.Info("Student was successfully unassigned from parent", "parent ID", userID, "student ID", studentID)
		return nil
	})
}
