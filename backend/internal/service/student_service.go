// Package service provides business logic and use cases.
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

type StudentService interface {
	GetStudentByID(ctx context.Context, id uuid.UUID) (*StudentResponseDTO, error)
	GetStudents(ctx context.Context, filter repository.StudentFilter) ([]StudentResponseDTO, error)
}

type StudentResponseDTO struct {
	User UserResponseDTO
	StudentGroupID uint16
}

type studentService struct {
	studentRepo repository.StudentRepository
	userRepo repository.UserRepository
	db       *gorm.DB
	log      logger.Logger
}

func NewStudentService(
	studentRepo repository.StudentRepository,
	userRepo repository.UserRepository,
	db *gorm.DB,
	log logger.Logger,
) StudentService {
	return &studentService{
		studentRepo: studentRepo,
		userRepo: userRepo,
		db:       db,
		log:      log,
	}
}

func (s *studentService) GetStudentByID(ctx context.Context, id uuid.UUID) (*StudentResponseDTO, error) {
	student, err := s.studentRepo.FindByID(ctx, &id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("student with id %s was not found: %w", id, err)
		}
		return nil, fmt.Errorf("failed to get student: %w", err)
	}
	return StudentToDTO(student), nil
}

func (s *studentService) GetStudents(ctx context.Context, filter repository.StudentFilter) ([]StudentResponseDTO, error) {
	students, err := s.studentRepo.FindAll(ctx, &filter)
	if err != nil {
		s.log.Error(
			"failed to get students from repository",
			"limit",
			filter.Limit,
			"offset",
			filter.Offset,
			"error",
			err,
		)
		return nil, fmt.Errorf(
			"failed to get students from repository (limit: %d, offset: %d): %w",
			filter.Limit,
			filter.Offset,
			err,
		)
	}
	studentDTOs := make([]StudentResponseDTO, len(students))
	for i, student := range students {
		studentDTOs[i] = *StudentToDTO(&student)
	}
	s.log.Info("successfully received the list of students")
	return studentDTOs, nil
}

func StudentToDTO(student *model.Student) *StudentResponseDTO {
	return &StudentResponseDTO{
		User: *UserToDTO(&student.User),
		StudentGroupID: student.StudentGroupID,
	}
}

