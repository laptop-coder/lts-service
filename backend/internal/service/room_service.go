package service

import (
	// "backend/pkg/apperrors" // TODO
	"backend/internal/model"
	"backend/internal/repository"
	"backend/pkg/logger"
	"context"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type RoomService interface {
	CreateRoom(ctx context.Context, dto CreateRoomDTO) (*RoomResponseDTO, error)
	GetRooms(ctx context.Context, filter repository.RoomFilter) ([]RoomResponseDTO, error)
	UpdateRoom(ctx context.Context, id uint8, dto UpdateRoomDTO) (*RoomResponseDTO, error)
	DeleteRoom(ctx context.Context, id uint8) error
}

type CreateRoomDTO struct {
	Name      string     `form:"name" validate:"required,min=1,max=20"`
	TeacherID *uuid.UUID `form:"teacherID,omitempty"`
}

type UpdateRoomDTO struct {
	Name      *string    `form:"name,omitempty" validate:"max=20"`
	TeacherID *uuid.UUID `form:"teacherID,omitempty"`
}

type RoomResponseDTO struct {
	ID        uint8      `json:"id"`
	CreatedAt string     `json:"createdAt"`
	UpdatedAt string     `json:"updatedAt"`
	Name      string     `json:"name"`
	TeacherID *uuid.UUID `json:"teacherID"`
}

type roomService struct {
	roomRepo repository.RoomRepository
	db       *gorm.DB
	log      logger.Logger
}

func NewRoomService(
	roomRepo repository.RoomRepository,
	db *gorm.DB,
	log logger.Logger,
) RoomService {
	return &roomService{
		roomRepo: roomRepo,
		db:       db,
		log:      log,
	}
}

func (s *roomService) CreateRoom(ctx context.Context, dto CreateRoomDTO) (*RoomResponseDTO, error) {
	// Input data validation
	if err := s.validateCreateRoomDTO(&dto); err != nil {
		return nil, fmt.Errorf("validation error during room creation: %w", err)
	}
	// TODO: check name uniqueness like in subject service
	// Creating model object
	room := &model.Room{
		Name:      dto.Name,
		TeacherID: dto.TeacherID,
	}
	// Create room
	if err := s.roomRepo.Create(ctx, room); err != nil {
		return nil, fmt.Errorf("failed to create room: %w", err)
	}
	// Get created room for response
	createdRoom, err := s.roomRepo.FindByID(ctx, &room.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch created room: %w", err)
	}
	return RoomToDTO(createdRoom), nil
}

func (s *roomService) GetRooms(ctx context.Context, filter repository.RoomFilter) ([]RoomResponseDTO, error) {
	rooms, err := s.roomRepo.FindAll(ctx, &filter)
	if err != nil {
		s.log.Error(
			"failed to get rooms from repository",
			"limit",
			filter.Limit,
			"offset",
			filter.Offset,
			"error",
			err,
		)
		return nil, fmt.Errorf(
			"failed to get rooms from repository (limit: %d, offset: %d): %w",
			filter.Limit,
			filter.Offset,
			err,
		)
	}
	roomDTOs := make([]RoomResponseDTO, len(rooms))
	for i, room := range rooms {
		roomDTOs[i] = *RoomToDTO(&room)
	}
	s.log.Info("successfully received the list of rooms")
	return roomDTOs, nil
}

func (s *roomService) UpdateRoom(ctx context.Context, id uint8, dto UpdateRoomDTO) (*RoomResponseDTO, error) {
	// Input data validation
	if err := s.validateUpdateRoomDTO(&dto); err != nil {
		return nil, fmt.Errorf("validation error during room updating: %w", err)
	}
	// Getting existing room
	room, err := s.roomRepo.FindByID(ctx, &id)
	if err != nil {
		return nil, fmt.Errorf("failed to get room for update: %w", err)
	}
	// Updating fields
	// TODO: add checks like in subject service
	updatedFieldsCount := 0 // TODO: perhaps this is unnecessary
	if dto.Name != nil && *dto.Name != room.Name {
		room.Name = *dto.Name
		updatedFieldsCount++
	}
	if dto.TeacherID != nil && *dto.TeacherID != *room.TeacherID {
		room.TeacherID = dto.TeacherID
		updatedFieldsCount++
	}
	// Updating room in DB
	if err := s.roomRepo.Update(ctx, room); err != nil {
		s.log.Error("failed to update the room")
		return nil, fmt.Errorf("failed to update the room: %w", err)
	}
	// Get updated room for response
	updatedRoom, err := s.roomRepo.FindByID(ctx, &room.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated room: %w", err)
	}
	return RoomToDTO(updatedRoom), nil
}

func (s *roomService) DeleteRoom(ctx context.Context, id uint8) error {
	s.log.Info("Starting room deletion...")
	if err := s.roomRepo.Delete(ctx, &id); err != nil {
		return fmt.Errorf("failed to delete the room: %w", err)
	}
	s.log.Info("room deleted successfully")
	return nil
}

func (s *roomService) validateCreateRoomDTO(dto *CreateRoomDTO) error {
	// TODO: check if teacher with this ID exists in DB
	return nil
}

func (s *roomService) validateUpdateRoomDTO(dto *UpdateRoomDTO) error {
	// TODO: check if teacher with this ID exists in DB
	return nil
}

func RoomToDTO(room *model.Room) *RoomResponseDTO {
	return &RoomResponseDTO{
		ID:        room.ID,
		CreatedAt: room.CreatedAt.Format(time.RFC3339),
		UpdatedAt: room.UpdatedAt.Format(time.RFC3339),
		Name:      room.Name,
		TeacherID: room.TeacherID,
	}
}
