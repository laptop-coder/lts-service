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

type TeacherService interface {
	GetTeacherByID(ctx context.Context, id uuid.UUID) (*TeacherResponseDTO, error)
	GetTeachers(ctx context.Context, filter repository.TeacherFilter) ([]TeacherResponseDTO, error)
}

func TeacherToDTO(teacher *model.Teacher) *TeacherResponseDTO {
	// Get subjects list
	var subjects []SubjectResponseDTO
	for _, subject := range teacher.Subjects {
		subjects = append(subjects, SubjectResponseDTO{
			ID:        subject.ID,
			CreatedAt: subject.CreatedAt.Format(time.RFC3339),
			UpdatedAt: subject.UpdatedAt.Format(time.RFC3339),
			Name:      subject.Name,
		})
	}
	// Get teacher classroom (if exists)
	var classroom *RoomResponseDTO
	if teacher.Classroom != nil {
		classroom = RoomToDTO(teacher.Classroom)
	}
	// Get student group where the teacher is advisor (if exists)
	var studentGroups []StudentGroupResponseDTO
	if teacher.StudentGroups != nil {
		for _, group := range *teacher.StudentGroups {
			studentGroups = append(studentGroups, *StudentGroupToDTO(&group))
		}
	}
	// Return response
	return &TeacherResponseDTO{
		CreatedAt:     teacher.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     teacher.UpdatedAt.Format(time.RFC3339),
		User:          *UserToDTO(&teacher.User),
		Subjects:      subjects,
		Classroom:     classroom,
		StudentGroups: studentGroups,
	}
}

type TeacherResponseDTO struct {
	CreatedAt     string                    `json:"createdAt"`
	UpdatedAt     string                    `json:"updatedAt"`
	User          UserResponseDTO           `json:"user"`
	Subjects      []SubjectResponseDTO      `json:"subjects"`
	Classroom     *RoomResponseDTO          `json:"classroom,omitempty"`
	StudentGroups []StudentGroupResponseDTO `json:"studentGroups,omitempty"`
}

type StudentGroupResponseDTO struct {
	ID        uint16 `json:"id"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	Name      string `json:"name"`
}

func StudentGroupToDTO(studentGroup *model.StudentGroup) *StudentGroupResponseDTO {
	return &StudentGroupResponseDTO{
		CreatedAt: studentGroup.CreatedAt.Format(time.RFC3339),
		UpdatedAt: studentGroup.UpdatedAt.Format(time.RFC3339),
		ID:        studentGroup.ID,
		Name:      studentGroup.Name,
	}
}

type teacherService struct {
	teacherRepo repository.TeacherRepository
	userRepo    repository.UserRepository
	db          *gorm.DB
	log         logger.Logger
}

func NewTeacherService(
	teacherRepo repository.TeacherRepository,
	userRepo repository.UserRepository,
	db *gorm.DB,
	log logger.Logger,
) TeacherService {
	return &teacherService{
		teacherRepo: teacherRepo,
		userRepo:    userRepo,
		db:          db,
		log:         log,
	}
}

func (s *teacherService) GetTeacherByID(ctx context.Context, id uuid.UUID) (*TeacherResponseDTO, error) {
	teacher, err := s.teacherRepo.FindByID(ctx, &id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("teacher with id %s was not found: %w", id, err)
		}
		return nil, fmt.Errorf("failed to get teacher: %w", err)
	}
	return TeacherToDTO(teacher), nil
}

func (s *teacherService) GetTeachers(ctx context.Context, filter repository.TeacherFilter) ([]TeacherResponseDTO, error) {
	teachers, err := s.teacherRepo.FindAll(ctx, &filter)
	if err != nil {
		s.log.Error(
			"failed to get teachers from repository",
			"limit",
			filter.Limit,
			"offset",
			filter.Offset,
			"error",
			err,
		)
		return nil, fmt.Errorf(
			"failed to get teachers from repository (limit: %d, offset: %d): %w",
			filter.Limit,
			filter.Offset,
			err,
		)
	}
	teacherDTOs := make([]TeacherResponseDTO, len(teachers))
	for i, teacher := range teachers {
		teacherDTOs[i] = *TeacherToDTO(&teacher)
	}
	s.log.Info("successfully received the list of teachers")
	return teacherDTOs, nil
}
