package repository

import (
	"backend/internal/model"
	"backend/pkg/logger"
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type RoomRepository interface {
	Create(ctx context.Context, room *model.Room) error
	FindAll(ctx context.Context) ([]model.Room, error)
	FindByID(ctx context.Context, id *uint8) (*model.Room, error)
	Update(ctx context.Context, room *model.Room) error
	Delete(ctx context.Context, id *uint8) error
}

type roomRepository struct {
	db  *gorm.DB
	log logger.Logger
}

func NewRoomRepository(db *gorm.DB, log logger.Logger) RoomRepository {
	if db == nil {
		log.Error("DB is nil")
		panic("DB is nil")
	}
	return &roomRepository{db: db, log: log}
}

func (r *roomRepository) Create(ctx context.Context, room *model.Room) error {
	if room == nil {
		return fmt.Errorf("room cannot be nil")
	}

	result := r.db.WithContext(ctx).Create(room)
	if result.Error != nil {
		return fmt.Errorf("failed to create new room: %w", result.Error)
	}

	return nil
}

func (r *roomRepository) FindAll(ctx context.Context) ([]model.Room, error) {
	var rooms []model.Room

	err := r.db.WithContext(ctx).
		Model(&model.Room{}).
		Order("name").
		Find(&rooms).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch rooms list: %w", err)
	}

	return rooms, nil
}

func (r *roomRepository) FindByID(ctx context.Context, id *uint8) (*model.Room, error) {
	if id == nil {
		return nil, fmt.Errorf("room id cannot be nil")
	}

	var room model.Room
	result := r.db.WithContext(ctx).First(&room, *id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("room with id %d was not found: %w", *id, result.Error)
		}
		return nil, fmt.Errorf("failed to fetch room by id (%d): %w", *id, result.Error)
	}

	return &room, nil
}

func (r *roomRepository) Update(ctx context.Context, room *model.Room) error {
	if room == nil {
		return fmt.Errorf("room cannot be nil")
	}

	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.Room{}).
		Where("id = ?", room.ID).
		Count(&count).Error

	if err != nil {
		return fmt.Errorf("failed to check room existence: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("room with id %d was not found", room.ID)
	}

	result := r.db.WithContext(ctx).Save(room)
	if result.Error != nil {
		return fmt.Errorf("failed to update room: %w", result.Error)
	}

	return nil
}

func (r *roomRepository) Delete(ctx context.Context, id *uint8) error {
	result := r.db.WithContext(ctx).Unscoped().Delete(&model.Room{}, *id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete room with id %d: %w", *id, result.Error)
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
