package service

import (
	"backend/internal/model"
	"backend/internal/repository"
	"backend/pkg/apperrors"
	"backend/pkg/logger"
	"context"
	"fmt"
	"gorm.io/gorm"
	"strings"
	"time"
)

type SubjectService interface {
	CreateSubject(ctx context.Context, dto CreateSubjectDTO) (*SubjectResponseDTO, error)
	GetSubjects(ctx context.Context, filter repository.SubjectFilter) ([]SubjectResponseDTO, error)
	UpdateSubject(ctx context.Context, id uint16, dto UpdateSubjectDTO) (*SubjectResponseDTO, error)
	DeleteSubject(ctx context.Context, id uint16) error
}

type CreateSubjectDTO struct {
	Name string `form:"name" validate:"required,min=3,max=100"`
}

type UpdateSubjectDTO struct {
	Name *string `form:"name,omitempty" validate:"max=100"` // TODO: check other (and this too) validation rules like here. Maybe there are mistakes
}

type SubjectResponseDTO struct {
	ID        uint16  `json:"id"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	Name      string `json:"name"`
}

type subjectService struct {
	subjectRepo repository.SubjectRepository
	db          *gorm.DB
	log         logger.Logger
}

func NewSubjectService(
	subjectRepo repository.SubjectRepository,
	db *gorm.DB,
	log logger.Logger,
) SubjectService {
	return &subjectService{
		subjectRepo: subjectRepo,
		db:          db,
		log:         log,
	}
}

func (s *subjectService) CreateSubject(ctx context.Context, dto CreateSubjectDTO) (*SubjectResponseDTO, error) {
	// Input data validation
	if err := s.validateCreateSubjectDTO(&dto); err != nil {
		return nil, fmt.Errorf("validation error during subject creation: %w", err)
	}
	// Check name uniqueness
	exists, err := s.subjectRepo.ExistsByName(ctx, &dto.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to check name uniqueness: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("subject with name '%s' already exists: %w", dto.Name, apperrors.ErrSubjectAlreadyExists) // TODO: maybe remove quotation marks from errors (are they safety?)
	}
	// Creating model object
	subject := &model.Subject{
		Name: dto.Name,
	}
	// Create subject
	if err := s.subjectRepo.Create(ctx, subject); err != nil {
		return nil, fmt.Errorf("failed to create subject: %w", err)
	}
	// Get created subject for response
	createdSubject, err := s.subjectRepo.FindByID(ctx, &subject.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch created subject: %w", err)
	}
	return SubjectToDTO(createdSubject), nil
}

func (s *subjectService) GetSubjects(ctx context.Context, filter repository.SubjectFilter) ([]SubjectResponseDTO, error) {
	subjects, err := s.subjectRepo.FindAll(ctx, &filter)
	if err != nil {
		s.log.Error(
			"failed to get subjects from repository",
			"limit",
			filter.Limit,
			"offset",
			filter.Offset,
			"error",
			err,
		)
		return nil, fmt.Errorf(
			"failed to get subjects from repository (limit: %d, offset: %d): %w",
			filter.Limit,
			filter.Offset,
			err,
		)
	}
	subjectDTOs := make([]SubjectResponseDTO, len(subjects))
	for i, subject := range subjects {
		subjectDTOs[i] = *SubjectToDTO(&subject)
	}
	s.log.Info("successfully received the list of subjects")
	return subjectDTOs, nil
}

func (s *subjectService) UpdateSubject(ctx context.Context, id uint16, dto UpdateSubjectDTO) (*SubjectResponseDTO, error) {
	// Input data validation
	if err := s.validateUpdateSubjectDTO(&dto); err != nil {
		return nil, fmt.Errorf("validation error during subject updating: %w", err)
	}
	// Getting existing subject
	subject, err := s.subjectRepo.FindByID(ctx, &id)
	if err != nil {
		return nil, fmt.Errorf("failed to get subject for update: %w", err)
	}
	// Update field if provided and was changed
	if dto.Name != nil && *dto.Name != subject.Name {
		// Check name uniqueness
		exists, err := s.subjectRepo.ExistsByName(ctx, dto.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to check name uniqueness: %w", err)
		}
		if exists {
			return nil, fmt.Errorf("subject with name '%s' already exists: %w", dto.Name, apperrors.ErrSubjectAlreadyExists)
		}
		subject.Name = *dto.Name
	} else {
		// No changes to update, return existing subject
		s.log.Info("no changes to update subject", "subject_id", id)
		return SubjectToDTO(subject), nil
	}
	// Update subject in DB
	if err := s.subjectRepo.Update(ctx, subject); err != nil {
		s.log.Error("failed to update the subject")
		return nil, fmt.Errorf("failed to update the subject: %w", err)
	}
	// Get updated subject for response
	updatedSubject, err := s.subjectRepo.FindByID(ctx, &subject.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated subject: %w", err)
	}
	return SubjectToDTO(updatedSubject), nil
}

func (s *subjectService) DeleteSubject(ctx context.Context, id uint16) error {
	s.log.Info("starting subject deletion...")
	if err := s.subjectRepo.Delete(ctx, &id); err != nil {
		return fmt.Errorf("failed to delete the subject: %w", err)
	}
	s.log.Info("subject deleted successfully")
	return nil
}

func (s *subjectService) validateCreateSubjectDTO(dto *CreateSubjectDTO) error {
	// TODO: check if teacher with this ID exists in DB
	if strings.TrimSpace(dto.Name) == "" {
		return fmt.Errorf("name cannot be empty or only whitespace: %w", apperrors.ErrRequiredField)
	}
	if len(dto.Name) < 3 {
		return fmt.Errorf("name must be at least 3 characters: %w", apperrors.ErrValueTooShort)
	}
	if len(dto.Name) > 100 {
		return fmt.Errorf("name must be at most 100 characters: %w", apperrors.ErrValueTooLong)
	}
	return nil
}

func (s *subjectService) validateUpdateSubjectDTO(dto *UpdateSubjectDTO) error {
	// TODO: check if teacher with this ID exists in DB
	if dto.Name != nil {
		if strings.TrimSpace(*dto.Name) == "" {
			return fmt.Errorf("name cannot be only whitespace: %w", apperrors.ErrRequiredField)
		}
		if len(*dto.Name) < 3 {
			return fmt.Errorf("name must be at least 3 characters: %w", apperrors.ErrValueTooShort)
		}
		if len(*dto.Name) > 100 {
			return fmt.Errorf("name must be at most 100 characters: %w", apperrors.ErrValueTooLong)
		}
	}
	return nil
}

func SubjectToDTO(subject *model.Subject) *SubjectResponseDTO {
	return &SubjectResponseDTO{
		ID:        subject.ID,
		CreatedAt: subject.CreatedAt.Format(time.RFC3339),
		UpdatedAt: subject.UpdatedAt.Format(time.RFC3339),
		Name:      subject.Name,
	}
}
