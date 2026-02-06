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
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"time"
)

type CreateUserDTO struct {
	Email string `form:"email" validate:"required,email,min=5"`
	// TODO: add upper bound for password. Take into account:
	// 1. Cyrillic characters or emojis take up more bytes
	// 2. Bcrypt restrictions
	Password   string                `form:"password" validate:"required,min=8"`
	FirstName  string                `form:"firstName" validate:"required,min=2"`
	MiddleName *string               `form:"middleName,omitempty"`
	LastName   string                `form:"lastName" validate:"required,min=2"`
	RoleID     *uint8                 `form:"roleId,omitempty"`
	Avatar     *multipart.FileHeader `form:"avatar,omitempty"` // avatar file
}

type UserResponseDTO struct {
	ID         uuid.UUID `json:"id"`
	Email      string    `json:"email"`
	FirstName  string    `json:"firstName"`
	MiddleName *string   `json:"middleName,omitempty"`
	LastName   string    `json:"lastName"`
	HasAvatar  bool      `json:"hasAvatar"`
	Roles      []RoleDTO `json:"roles"`
	CreatedAt  string    `json:"createdAt"`
}

type ChangePasswordDTO struct {
	OldPassword string `form:"oldPassword" validate:"required"`
	// TODO: move password min length to config
	NewPassword string `form:"newPassword" validate:"required,min=8"`
}

type UpdateUserDTO struct {
	Email      *string               `form:"email,omitempty" validate:"email,min=5"`
	FirstName  *string               `form:"firstName,omitempty" validate:"min=2"`
	MiddleName *string               `form:"middleName,omitempty"`
	LastName   *string               `form:"lastName,omitempty" validate:"min=2"`
	RoleID     *uint8                `form:"roleId,omitempty"`
	Avatar     *multipart.FileHeader `form:"avatar,omitempty"` // avatar file
}

type RoleDTO struct {
	ID   uint8  `json:"id"`
	Name string `json:"name"`
}

type userService struct {
	userRepo repository.UserRepository
	db       *gorm.DB
	config   UserServiceConfig
	log      logger.Logger
}

type UserService interface {
	// CRUD
	CreateUser(ctx context.Context, dto CreateUserDTO) (*UserResponseDTO, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*UserResponseDTO, error)
	GetUserByEmail(ctx context.Context, email string) (*UserResponseDTO, error)
	GetUsers(ctx context.Context, filter repository.UserFilter) ([]UserResponseDTO, error)
	UpdateUser(ctx context.Context, id uuid.UUID, dto UpdateUserDTO) (*UserResponseDTO, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error

	ChangePassword(ctx context.Context, id uuid.UUID, dto ChangePasswordDTO) error
	UpdateAvatar(ctx context.Context, userID uuid.UUID, dto *multipart.FileHeader) error
	RemoveAvatar(ctx context.Context, userID uuid.UUID) error

	GetStudentGroupAdvisorByGroupID(ctx context.Context, id uint16) (*UserResponseDTO, error)
}

func NewUserService(
	userRepo repository.UserRepository,
	db *gorm.DB,
	config UserServiceConfig,
	log logger.Logger,
) UserService {
	return &userService{
		userRepo: userRepo,
		db:       db,
		config:   config,
		log:      log,
	}
}

func (s *userService) CreateUser(ctx context.Context, dto CreateUserDTO) (*UserResponseDTO, error) {
	// Input data validation
	if err := validateCreateUserDTO(&dto); err != nil {
		return nil, fmt.Errorf("validation error during user creation: %w", err)
	}
	// Check email uniqueness
	existingUser, err := s.userRepo.FindByEmail(ctx, &dto.Email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("user with this email already exists")
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check email uniqueness: %w", err)
	}
	// Password hashing
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(dto.Password), s.config.BcryptCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	// Generating ID for user
	userID := uuid.New()
	// Avatar processing (if passed)
	hasAvatar := false
	if dto.Avatar != nil {
		// Validating
		if err := s.validateAvatarFile(dto.Avatar); err != nil {
			return nil, fmt.Errorf("avatar validation failed: %w", err)
		}
		// Saving to storage
		if err := s.saveAvatarFile(userID, dto.Avatar); err != nil {
			return nil, fmt.Errorf("failed to save avatar to storage: %w", err)
		}
		hasAvatar = true
	}
	// Creating model object
	user := &model.User{
		ID:         userID,
		Email:      dto.Email,
		Password:   string(passwordHash),
		FirstName:  dto.FirstName,
		MiddleName: dto.MiddleName,
		LastName:   dto.LastName,
		HasAvatar:  hasAvatar,
	}
	// Transaction for creating user
	err = s.db.Transaction(func(tx *gorm.DB) error {
		txRepo := repository.NewUserRepository(tx, s.log)
		if err := txRepo.Create(ctx, user); err != nil {
			// Delete the saved avatar, if the transaction is rolled back
			if hasAvatar {
				s.removeAvatarFile(userID)
			}
			return fmt.Errorf("failed to create user: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("transaction failed: %w", err)
	}
	// Get created user for response
	createdUser, err := s.userRepo.FindByID(ctx, &user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch created user: %w", err)
	}
	return s.userToDTO(createdUser), nil
}

func (s *userService) validateAvatarFile(fileHeader *multipart.FileHeader) error {
	// Check file size
	if fileHeader.Size > s.config.AvatarMaxSize {
		return fmt.Errorf("file size exceeds limit of %d bytes", s.config.AvatarMaxSize)
	}

	// read meta information
	file, err := fileHeader.Open()
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// to return to the start of the file after determiming the MIME type
	if seeker, ok := file.(io.Seeker); ok {
		defer seeker.Seek(0, io.SeekStart)
	}

	buffer := make([]byte, 512) // read first 512 bytes to determine MIME type
	_, err = file.Read(buffer)
	if err != nil && err != io.EOF {
		return fmt.Errorf("failed to read file: %w", err)
	}

	mimeType := http.DetectContentType(buffer)
	if !slices.Contains(s.config.AvatarAllowedMIMETypes, mimeType) {
		return fmt.Errorf("unsupported file type: %s. Allowed: %v", mimeType, s.config.AvatarAllowedMIMETypes)
	}

	return nil
}

func (s *userService) saveAvatarFile(userID uuid.UUID, fileHeader *multipart.FileHeader) error {
	// Creating directory (if not exists)
	if err := os.MkdirAll(s.config.AvatarUploadPath, 0755); err != nil {
		return fmt.Errorf("failed to create upload directory for avatars: %w", err)
	}
	// Opening source file
	srcFile, err := fileHeader.Open()
	if err != nil {
		return fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer srcFile.Close()
	// Creating file path (where to save avatar)
	filePath := filepath.Join(
		s.config.AvatarUploadPath,
		fmt.Sprintf("%s.jpeg", userID.String()),
	)
	// Creating file in storage
	dstFile, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer dstFile.Close()
	// Copying the content from the source file to the destination file
	if _, err = io.Copy(dstFile, srcFile); err != nil {
		// Deleting a partially filled file in the case of error
		os.Remove(filePath)
		return fmt.Errorf("failed to copy file: %w", err)
	}
	return nil
}

func (s *userService) removeAvatarFile(userID uuid.UUID) {
	filePath := filepath.Join(
		s.config.AvatarUploadPath,
		fmt.Sprintf("%s.jpeg", userID.String()),
	)
	os.Remove(filePath)
}

func (s *userService) UpdateAvatar(ctx context.Context, userID uuid.UUID, avatar *multipart.FileHeader) error {
	user, err := s.userRepo.FindByID(ctx, &userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}
	// Validating the file
	if err := s.validateAvatarFile(avatar); err != nil {
		return err
	}
	// Saving the new avatar
	if err := s.saveAvatarFile(userID, avatar); err != nil {
		return err
	}
	// Mark existence of the avatar in the database
	user.HasAvatar = true
	if err := s.userRepo.Update(ctx, user); err != nil {
		// Rollback file saving in the case of error
		s.removeAvatarFile(userID)
		return fmt.Errorf("failed to update user avatar: %w", err)
	}
	return nil
}

func validateCreateUserDTO(dto *CreateUserDTO) error {
	if dto.Email == "" {
		return fmt.Errorf("email is required")
	}
	if len(dto.Password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}
	if len(dto.FirstName) < 2 {
		return fmt.Errorf("first name must be at least 2 characters")
	}
	if len(dto.LastName) < 2 {
		return fmt.Errorf("last name must be at least 2 characters")
	}
	return nil
}

func (s *userService) userToDTO(user *model.User) *UserResponseDTO {
	var roles []RoleDTO
	for _, role := range user.Roles {
		roles = append(roles, RoleDTO{
			ID:   role.ID,
			Name: role.Name,
		})
	}

	return &UserResponseDTO{
		ID:         user.ID,
		Email:      user.Email,
		FirstName:  user.FirstName,
		MiddleName: user.MiddleName,
		LastName:   user.LastName,
		HasAvatar:  user.HasAvatar,
		Roles:      roles,
		CreatedAt:  user.CreatedAt.Format(time.RFC3339),
	}
}
