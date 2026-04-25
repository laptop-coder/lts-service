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
)

type TeacherService interface {
	GetTeacherByID(ctx context.Context, id uuid.UUID) (*TeacherResponseDTO, error)
	GetTeachers(ctx context.Context, filter repository.TeacherFilter) ([]TeacherResponseDTO, error)
	GetTeacherClassroom(ctx context.Context, userID uuid.UUID) (*RoomResponseDTO, error)
	GetTeacherSubjects(ctx context.Context, userID uuid.UUID) ([]SubjectResponseDTO, error)
	AssignClassroom(ctx context.Context, userID uuid.UUID, classroomID uint8) error
	UnassignClassroom(ctx context.Context, userID uuid.UUID) error
	AddSubjects(ctx context.Context, userID uuid.UUID, subjectIDs []uint8) error
	AssignSubjects(ctx context.Context, userID uuid.UUID, subjectIDs []uint8) error
	UnassignSubject(ctx context.Context, userID uuid.UUID, subjectID uint8) error
}

func TeacherToDTO(teacher *model.Teacher) *TeacherResponseDTO {
	// Get subjects list
	var subjects []SubjectResponseDTO
	for _, subject := range teacher.Subjects {
		subjects = append(subjects, *SubjectToDTO(&subject))
	}
	// Get teacher classroom (if exists)
	var classroom *RoomResponseDTO
	if teacher.Classroom != nil {
		classroom = RoomToDTO(teacher.Classroom)
	}
	// Get student group where the teacher is advisor (if exists)
	var studentGroups []StudentGroupResponseDTO
	if teacher.StudentGroups != nil && len(*teacher.StudentGroups) > 0 {
		for _, group := range *teacher.StudentGroups {
			studentGroups = append(studentGroups, *StudentGroupToDTO(&group))
		}
	}
	// Return response
	return &TeacherResponseDTO{
		UserID:        teacher.UserID,
		Subjects:      subjects,
		Classroom:     classroom,
		StudentGroups: studentGroups,
	}
}

type TeacherResponseDTO struct {
	UserID        uuid.UUID                    `json:"userId"`
	Subjects      []SubjectResponseDTO      `json:"subjects"`
	Classroom     *RoomResponseDTO          `json:"classroom,omitempty"`
	StudentGroups []StudentGroupResponseDTO `json:"studentGroups,omitempty"`
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

func (s *teacherService) GetTeacherClassroom(ctx context.Context, userID uuid.UUID) (*RoomResponseDTO, error) {
	// Find teacher by ID
	teacher, err := s.teacherRepo.FindByID(ctx, &userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get teacher: %w", err)
	}
	// Return response
	if teacher.Classroom == nil {
		return nil, nil
	}
	return RoomToDTO(teacher.Classroom), nil
}

func (s *teacherService) GetTeacherSubjects(ctx context.Context, userID uuid.UUID) ([]SubjectResponseDTO, error) {
	// Find teacher by ID
	teacher, err := s.teacherRepo.FindByID(ctx, &userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get teacher: %w", err)
	}
	// Check if there are subjects
	if len(teacher.Subjects) == 0 {
		return nil, fmt.Errorf("teacher must have at least one subject")
	}
	// Collect subjects, convert to DTO
	var subjects []SubjectResponseDTO
	for _, subject := range teacher.Subjects {
		subjects = append(subjects, *SubjectToDTO(&subject))
	}
	// Return response
	return subjects, nil
}

func (s *teacherService) AssignClassroom(ctx context.Context, userID uuid.UUID, classroomID uint8) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Check teacher existence
		var teacher model.Teacher
		if err := tx.WithContext(ctx).
			Where("user_id = ?", userID).
			First(&teacher).Error; err != nil {
			return fmt.Errorf("failed to find teacher by ID (%s)", userID)
		}
		// Check room existence
		var classroom model.Room
		if err := tx.WithContext(ctx).
			Where("id = ?", classroomID).
			First(&classroom).Error; err != nil {
			return fmt.Errorf("failed to find room by ID (%d)", classroomID)
		}
		// Error if the room already has another teacher
		if classroom.TeacherID != nil && *classroom.TeacherID != userID {
			return fmt.Errorf("the room with id %d already has another teacher with id %s: %w", classroomID, userID, apperrors.ErrRoomAlreadyHasTeacherAssignedToIt)
		}
		// Unassign if the teacher already has another room
		if teacher.Classroom != nil && teacher.Classroom.ID != classroomID {
			oldClassroomID := teacher.Classroom.ID
			if err := tx.Model(&model.Room{}).
				Where("id = ?", oldClassroomID).
				Update("teacher_id", nil).Error; err != nil {
				return fmt.Errorf("failed to unassign teacher (ID %s) from the old classroom (ID %d): %w", userID, oldClassroomID, err)
			}
		}
		// Assign teacher to new room
		if err := tx.Model(&model.Room{}).
			Where("id = ?", classroomID).
			Update("teacher_id", userID).Error; err != nil {
			return fmt.Errorf("failed to assign teacher (ID %s) to the new classroom (ID %d): %w", userID, classroomID, err)
		}
		s.log.Info("classroom assigned to teacher successfully", "teacher ID", userID, "classroom ID", classroomID)
		return nil
	})
}

func (s *teacherService) UnassignClassroom(ctx context.Context, userID uuid.UUID) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Check teacher existence
		var teacher model.Teacher
		if err := tx.WithContext(ctx).
			Where("user_id = ?", userID).
			First(&teacher).Error; err != nil {
			return fmt.Errorf("failed to find teacher by ID (%s)", userID)
		}
		// Unassign if the teacher has room
		if teacher.Classroom == nil {
			return nil // idempotence
		}
		classroomID := teacher.Classroom.ID
		if err := tx.Model(&model.Room{}).
			Where("id = ?", classroomID).
			Update("teacher_id", nil).Error; err != nil {
			return fmt.Errorf("failed to unassign teacher (ID %s) from the classroom (ID %d): %w", userID, classroomID, err)
		}
		s.log.Info("classroom unassigned from teacher successfully", "teacher ID", userID)
		return nil
	})
}

func (s *teacherService) AddSubjects(ctx context.Context, userID uuid.UUID, subjectIDs []uint8) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Check teacher existence
		var teacher model.Teacher
		if err := tx.WithContext(ctx).
			Preload("Subjects").
			Where("user_id = ?", userID).
			First(&teacher).Error; err != nil {
			return fmt.Errorf("failed to find teacher: %w", err)
		}
		// Get subjects by IDs
		var subjects []model.Subject
		if err := tx.WithContext(ctx).
			Where("id IN (?)", subjectIDs).
			Find(&subjects).Error; err != nil {
			return fmt.Errorf("failed to fetch subjects: %w", err)
		}
		// Check if all subjects were found
		if len(subjects) != len(subjectIDs) {
			return fmt.Errorf("some subjects not found: %w", apperrors.ErrNotFound)
		}
		// Add subjects
		if err := tx.Model(&teacher).Association("Subjects").Append(&subjects); err != nil {
			return fmt.Errorf("failed to add subjects: %w", err)
		}
		// Return response
		s.log.Info("subjects was successfully added to teacher", "teacher ID", userID, "subject IDs", subjectIDs)
		return nil
	})
}

func (s *teacherService) UnassignSubject(ctx context.Context, userID uuid.UUID, subjectID uint8) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Check teacher existence
		var teacher model.Teacher
		if err := tx.WithContext(ctx).
			Where("user_id = ?", userID).
			First(&teacher).Error; err != nil {
			return fmt.Errorf("failed to find teacher: %w", err)
		}
		// Get subject by ID
		var subject model.Subject
		if err := tx.WithContext(ctx).
			First(&subject, subjectID).Error; err != nil {
			return fmt.Errorf("subject not found: %s: %w", err.Error(), apperrors.ErrNotFound)
		}
		// Remove subject from teacher
		if err := tx.Model(&teacher).Association("Subjects").Delete(&subject); err != nil {
			return fmt.Errorf("failed to remove subject: %w", err)
		}
		// Return response
		s.log.Info("subject was successfully removed from teacher", "teacher ID", userID, "subject ID", subjectID)
		return nil
	})
}

func (s *teacherService) AssignSubjects(ctx context.Context, userID uuid.UUID, subjectIDs []uint8) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Check teacher existence
		var teacher model.Teacher
		if err := tx.WithContext(ctx).
			Where("user_id = ?", userID).
			First(&teacher).Error; err != nil {
			return fmt.Errorf("failed to find teacher: %w", err)
		}
		// Get subjects by IDs
		var subjects []model.Subject
		if err := tx.WithContext(ctx).
			Where("id IN (?)", subjectIDs).
			Find(&subjects).Error; err != nil {
			return fmt.Errorf("failed to fetch subjects: %w", err)
		}
		// Check if all subjects were found
		if len(subjects) != len(subjectIDs) {
			return fmt.Errorf("some subjects not found: %w", apperrors.ErrNotFound)
		}
		// Replace subjects
		if err := tx.Model(&teacher).Association("Subjects").Replace(&subjects); err != nil {
			return fmt.Errorf("failed to replace subjects: %w", err)
		}
		// Return response
		s.log.Info("subjects was successfully assigned to teacher (old subjects were replaced by new ones)", "teacher ID", userID, "subject IDs", subjectIDs)
		return nil
	})
}
