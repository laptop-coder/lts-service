package service

import (
	"backend/internal/model"
	"backend/internal/repository"
	"backend/pkg/logger"
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type SubjectService interface {
	CreateSubject(ctx context.Context, dto CreateSubjectDTO) (*SubjectResponseDTO, error)
	GetSubjects(ctx context.Context, filter repository.SubjectFilter) ([]SubjectResponseDTO, error)
	UpdateSubject(ctx context.Context, id uint8, dto UpdateSubjectDTO) (*SubjectResponseDTO, error)
	DeleteSubject(ctx context.Context, id uint8) error
}

type CreateSubjectDTO struct {
	Name string `form:"name" validate:"required,min=3,max=100"`
}

type UpdateSubjectDTO struct {
	Name *string `form:"name,omitempty" validate:"max=20"`
}

type SubjectResponseDTO struct {
	ID        uint8  `json:"id"`
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

func (s *subjectService) UpdateSubject(ctx context.Context, id uint8, dto UpdateSubjectDTO) (*SubjectResponseDTO, error) {
	// Input data validation
	if err := s.validateUpdateSubjectDTO(&dto); err != nil {
		return nil, fmt.Errorf("validation error during subject updating: %w", err)
	}
	// Getting existing subject
	subject, err := s.subjectRepo.FindByID(ctx, &id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.log.Error("Subject for update was not found by id", "subject id", id, "error", err)
			return nil, fmt.Errorf("subject with id %s was not found: %w", id, err)
		}
		s.log.Error("Failed to get subject for update", "subject id", id, "error", err)
		return nil, fmt.Errorf("failed to get subject for update: %w", err)
	}
	// Updating field
	if dto.Name != nil && *dto.Name != subject.Name {
		subject.Name = *dto.Name
	}
	// Updating subject in DB
	if err := s.subjectRepo.Update(ctx, subject); err != nil {
		s.log.Error("Failed to update the subject")
		return nil, fmt.Errorf("failed to update the subject: %w", err)
	}
	// Get updated subject for response
	updatedSubject, err := s.subjectRepo.FindByID(ctx, &subject.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated subject: %w", err)
	}
	return SubjectToDTO(updatedSubject), nil
}

func (s *subjectService) DeleteSubject(ctx context.Context, id uint8) error {
	s.log.Info("Starting subject deletion...")
	// Getting existing subject
	_, err := s.subjectRepo.FindByID(ctx, &id) // does it necessary?
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.log.Error("Subject for delete was not found by id", "subject id", id, "error", err)
			return fmt.Errorf("subject with id %s was not found: %w", id, err)
		}
		s.log.Error("Failed to get subject for delete", "subject id", id, "error", err)
		return fmt.Errorf("failed to get subject for delete: %w", err)
	}
	// Delete subject
	if err := s.subjectRepo.Delete(ctx, &id); err != nil {
		s.log.Error("Failed to delete the subject")
		return fmt.Errorf("failed to delete the subject: %w", err)
	}
	s.log.Info("Subject deleted successfully")
	return nil
}

func (s *subjectService) validateCreateSubjectDTO(dto *CreateSubjectDTO) error {
	// TODO: check if teacher with this ID exists in DB
	return nil
}

func (s *subjectService) validateUpdateSubjectDTO(dto *UpdateSubjectDTO) error {
	// TODO: check if teacher with this ID exists in DB
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
