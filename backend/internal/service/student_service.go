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
	GetStudentClassroom(ctx context.Context, userID uuid.UUID) (*RoomResponseDTO, error)
	GetStudentAdvisor(ctx context.Context, userID uuid.UUID) (*TeacherResponseDTO, error)
	GetStudentParents(ctx context.Context, userID uuid.UUID) ([]ParentResponseDTO, error)
}

type StudentResponseDTO struct {
	StudentGroup StudentGroupResponseDTO
	Parents      []ParentResponseDTO
}

type studentService struct {
	studentRepo      repository.StudentRepository
	studentGroupRepo repository.StudentGroupRepository
	userRepo         repository.UserRepository
	teacherRepo      repository.TeacherRepository
	db               *gorm.DB
	log              logger.Logger
}

func NewStudentService(
	studentRepo repository.StudentRepository,
	studentGroupRepo repository.StudentGroupRepository,
	userRepo repository.UserRepository,
	teacherRepo repository.TeacherRepository,
	db *gorm.DB,
	log logger.Logger,
) StudentService {
	return &studentService{
		studentRepo:      studentRepo,
		studentGroupRepo: studentGroupRepo,
		userRepo:         userRepo,
		teacherRepo:      teacherRepo,
		db:               db,
		log:              log,
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

func (s *studentService) GetStudentClassroom(ctx context.Context, userID uuid.UUID) (*RoomResponseDTO, error) {
	// Find student by ID
	student, err := s.studentRepo.FindByID(ctx, &userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("student with id %s was not found: %w", userID, err)
		}
		return nil, fmt.Errorf("failed to get student: %w", err)
	}
	// Find group by ID
	group, err := s.studentGroupRepo.FindByID(ctx, &student.StudentGroupID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student group: %w", err)
	}
	// Find group advisor by ID
	if group.GroupAdvisorID == nil {
		return nil, nil
	}
	teacher, err := s.teacherRepo.FindByID(ctx, group.GroupAdvisorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student group advisor: %w", err)
	}
	// Return response
	if teacher.Classroom == nil {
		return nil, nil
	}
	return RoomToDTO(teacher.Classroom), nil
}

func (s *studentService) GetStudentAdvisor(ctx context.Context, userID uuid.UUID) (*TeacherResponseDTO, error) {
	// Find student by ID
	student, err := s.studentRepo.FindByID(ctx, &userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("student with id %s was not found: %w", userID, err)
		}
		return nil, fmt.Errorf("failed to get student: %w", err)
	}
	// Find group by ID
	group, err := s.studentGroupRepo.FindByID(ctx, &student.StudentGroupID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student group: %w", err)
	}
	// Find group advisor by ID
	if group.GroupAdvisorID == nil {
		return nil, nil
	}
	teacher, err := s.teacherRepo.FindByID(ctx, group.GroupAdvisorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student group advisor: %w", err)
	}
	// Return response
	return TeacherToDTO(teacher), nil
}

func (s *studentService) GetStudentParents(ctx context.Context, userID uuid.UUID) ([]ParentResponseDTO, error) {
	// Find student by ID
	student, err := s.studentRepo.FindByID(ctx, &userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("student with id %s was not found: %w", userID, err)
		}
		return nil, fmt.Errorf("failed to get student: %w", err)
	}
	// Check if there are connected parents
	if student.Parents != nil && len(*student.Parents) > 0 {
		// Get parents
		var parents []ParentResponseDTO
		for _, parent := range *student.Parents {
			parents = append(parents, *ParentToDTO(&parent))
		}
		// Return response
		return parents, nil
	}
	return nil, nil
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
	var parents []ParentResponseDTO
	if student.Parents != nil && len(*student.Parents) > 0 {
		for _, parent := range *student.Parents {
			parents = append(parents, *ParentToDTO(&parent))
		}
	}
	return &StudentResponseDTO{
		StudentGroup: *StudentGroupToDTO(&student.StudentGroup),
		Parents:      parents,
	}
}
