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

type StaffService interface {
	GetStaffByID(ctx context.Context, id uuid.UUID) (*StaffResponseDTO, error)
	GetStaff(ctx context.Context, filter repository.StaffFilter) ([]StaffResponseDTO, error)
}

type StaffResponseDTO struct {
	User     UserResponseDTO
	Position StaffPositionResponseDTO
}

type StaffPositionResponseDTO struct {
	ID        uint8  `json:"id"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	Name      string `json:"name"`
}

type staffService struct {
	staffRepo repository.StaffRepository
	userRepo  repository.UserRepository
	db        *gorm.DB
	log       logger.Logger
}

func NewStaffService(
	staffRepo repository.StaffRepository,
	userRepo repository.UserRepository,
	db *gorm.DB,
	log logger.Logger,
) StaffService {
	return &staffService{
		staffRepo: staffRepo,
		userRepo:  userRepo,
		db:        db,
		log:       log,
	}
}

func (s *staffService) GetStaffByID(ctx context.Context, id uuid.UUID) (*StaffResponseDTO, error) {
	staff, err := s.staffRepo.FindByID(ctx, &id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("staff with id %s was not found: %w", id, err)
		}
		return nil, fmt.Errorf("failed to get staff: %w", err)
	}
	return StaffToDTO(staff), nil
}

func (s *staffService) GetStaff(ctx context.Context, filter repository.StaffFilter) ([]StaffResponseDTO, error) {
	staffList, err := s.staffRepo.FindAll(ctx, &filter)
	if err != nil {
		s.log.Error(
			"failed to get staff from repository",
			"limit",
			filter.Limit,
			"offset",
			filter.Offset,
			"error",
			err,
		)
		return nil, fmt.Errorf(
			"failed to get staff from repository (limit: %d, offset: %d): %w",
			filter.Limit,
			filter.Offset,
			err,
		)
	}
	staffDTOs := make([]StaffResponseDTO, len(staffList))
	for i, staff := range staffList {
		staffDTOs[i] = *StaffToDTO(&staff)
	}
	s.log.Info("successfully received the list of staff")
	return staffDTOs, nil
}

func StaffToDTO(staff *model.Staff) *StaffResponseDTO {
	return &StaffResponseDTO{
		User:     *UserToDTO(&staff.User),
		Position: *StaffPositionToDTO(&staff.Position),
	}
}

func StaffPositionToDTO(position *model.StaffPosition) *StaffPositionResponseDTO {
	return &StaffPositionResponseDTO{
		ID:        position.ID,
		CreatedAt: position.CreatedAt.Format(time.RFC3339),
		UpdatedAt: position.UpdatedAt.Format(time.RFC3339),
		Name:      position.Name,
	}
}
