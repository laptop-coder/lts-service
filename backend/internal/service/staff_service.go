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
	// Position
	AssignPosition(ctx context.Context, userID uuid.UUID, positionID uint8) error
	GetPosition(ctx context.Context, userID uuid.UUID) (*StaffPositionResponseDTO, error)
}

type StaffResponseDTO struct {
	User     UserResponseDTO
	Position StaffPositionResponseDTO
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

func (s *staffService) AssignPosition(ctx context.Context, userID uuid.UUID, positionID uint8) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Get staff
		var staff model.Staff
		if err := tx.WithContext(ctx).
			First(&staff, "user_id = ?", userID).Error; err != nil {
			return fmt.Errorf("staff with user ID %s was not found: %w", userID, err)
		}
		// Check position existence
		var count int64
		if err := tx.WithContext(ctx).
			Model(&model.StaffPosition{}).
			Where("id = ?", positionID).
			Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			return fmt.Errorf("position with ID %d was not found", positionID)
		}
		staff.PositionID = positionID
		if err := tx.WithContext(ctx).Save(&staff).Error; err != nil {
			return fmt.Errorf("failed to assign position to staff: %w", err)
		}
		s.log.Info("Position was successfully assigned to the staff")
		return nil
	})
}

func (s *staffService) GetPosition(ctx context.Context, userID uuid.UUID) (*StaffPositionResponseDTO, error) {
	staff, err := s.staffRepo.FindByID(ctx, &userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("staff with user id %s was not found: %w", userID, err)
		}
		return nil, fmt.Errorf("failed to get staff: %w", err)
	}
	if staff.Position.ID == 0 {
		return nil, fmt.Errorf("staff has no position assigned")
	}
	return StaffPositionToDTO(&staff.Position), nil
}

func StaffToDTO(staff *model.Staff) *StaffResponseDTO {
	return &StaffResponseDTO{
		User:     *UserToDTO(&staff.User),
		Position: *StaffPositionToDTO(&staff.Position),
	}
}

