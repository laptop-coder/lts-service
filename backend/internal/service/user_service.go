// Package service provides business logic and use cases.
package service

import (
	"backend/internal/model"
	"backend/internal/repository"
	"backend/pkg/apperrors"
	"backend/pkg/logger"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/image/draw"
	_ "golang.org/x/image/webp"
	"gorm.io/gorm"
	"image"
	"image/color"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"time"
)

type UserService interface {
	CreateUser(ctx context.Context, createUserDTO CreateUserDTO, userExtensionsDTO UserExtensionsDTO) (*UserResponseDTO, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*UserResponseDTO, error)
	GetUserByEmail(ctx context.Context, email string) (*UserResponseDTO, error)
	GetUsers(ctx context.Context, filter repository.UserFilter) ([]UserResponseDTO, error)
	UpdateUser(ctx context.Context, id uuid.UUID, dto UpdateUserDTO) (*UserResponseDTO, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
	UpdateAvatar(ctx context.Context, userID uuid.UUID, dto *multipart.FileHeader) error
	RemoveAvatar(ctx context.Context, userID uuid.UUID) error
	// Roles
	GetUserRoles(ctx context.Context, userID uuid.UUID) ([]RoleResponseDTO, error)
	AssignRolesToUser(ctx context.Context, userID uuid.UUID, dto UserExtensionsDTO, roleIDs []uint16) error         // replace old roles with new ones
	AssignNonAdminRolesToUser(ctx context.Context, userID uuid.UUID, dto UserExtensionsDTO, roleIDs []uint16) error // replace old roles with new ones
	AddRolesToUser(ctx context.Context, userID uuid.UUID, dto UserExtensionsDTO, roleIDs []uint16) error
	RemoveRolesFromUser(ctx context.Context, userID uuid.UUID, roleIDs []uint16) error
	RemoveRoleFromUser(ctx context.Context, userID uuid.UUID, roleID uint16) error
	AssignExtensionsToUser(ctx context.Context, userID uuid.UUID, dto UserExtensionsDTO) error
}

type PermissionResponseDTO struct {
	ID        uint16 `json:"id"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	Name      string `json:"name"`
}

func PermissionToDTO(permission *model.Permission) *PermissionResponseDTO {
	return &PermissionResponseDTO{
		ID:        permission.ID,
		CreatedAt: permission.CreatedAt.Format(time.RFC3339),
		UpdatedAt: permission.UpdatedAt.Format(time.RFC3339),
		Name:      permission.Name,
	}
}

type CreateUserDTO struct {
	Email string `form:"email" validate:"required,email,min=5"`
	// TODO: add upper bound for password. Take into account:
	// 1. Cyrillic characters or emojis take up more bytes
	// 2. Bcrypt restrictions
	Password   string                `form:"password" validate:"required,min=8"`
	FirstName  string                `form:"firstName" validate:"required,min=2"`
	MiddleName *string               `form:"middleName,omitempty"`
	LastName   string                `form:"lastName" validate:"required,min=2"`
	RoleIDs    []uint16              `form:"roleIds,omitempty"`
	Avatar     *multipart.FileHeader `form:"avatar,omitempty"` // avatar file
}

type UserExtensionsDTO struct {
	TeacherClassroomID                 *uint16     `form:"teacherClassroomId,omitempty"`
	TeacherSubjectIDs                  []uint16    `form:"teacherSubjectIds,omitempty"`
	TeacherStudentGroupIDs             []uint16    `form:"teacherStudentGroupIds,omitempty"`
	StudentGroupID                     *uint16     `form:"studentGroupId,omitempty"`
	StaffPositionID                    *uint16     `form:"staffPositionId,omitempty"`
	InstitutionAdministratorPositionID *uint16     `form:"instituionAdministratorPositionId,omitempty"`
	ParentStudentIDs                   []uuid.UUID `form:"parentStudentIds,omitempty"`
}

type UserResponseDTO struct {
	ID         uuid.UUID         `json:"id"`
	CreatedAt  string            `json:"createdAt"`
	UpdatedAt  string            `json:"updatedAt"`
	Email      string            `json:"email"`
	FirstName  string            `json:"firstName"`
	MiddleName *string           `json:"middleName,omitempty"`
	LastName   string            `json:"lastName"`
	HasAvatar  bool              `json:"hasAvatar"`
	Roles      []RoleResponseDTO `json:"roles"`
}

type ChangePasswordDTO struct {
	OldPassword string `form:"oldPassword" validate:"required"`
	// TODO: move password min length to config
	NewPassword string `form:"newPassword" validate:"required,min=8"`
}

type UpdateUserDTO struct {
	FirstName  *string `form:"firstName,omitempty" validate:"min=2"`
	MiddleName *string `form:"middleName,omitempty"`
	LastName   *string `form:"lastName,omitempty" validate:"min=2"`
}

type userService struct {
	userRepo    repository.UserRepository
	studentRepo repository.StudentRepository
	roomRepo    repository.RoomRepository
	db          *gorm.DB
	config      UserServiceConfig
	log         logger.Logger
}

func NewUserService(
	userRepo repository.UserRepository,
	studentRepo repository.StudentRepository,
	roomRepo repository.RoomRepository,
	db *gorm.DB,
	config UserServiceConfig,
	log logger.Logger,
) UserService {
	return &userService{
		userRepo:    userRepo,
		studentRepo: studentRepo,
		roomRepo:    roomRepo,
		db:          db,
		config:      config,
		log:         log,
	}
}

func (s *userService) CreateUser(ctx context.Context, createUserDTO CreateUserDTO, userExtensionsDTO UserExtensionsDTO) (*UserResponseDTO, error) {
	// Input data validation
	if err := s.validateCreateUserDTO(&createUserDTO); err != nil {
		return nil, fmt.Errorf("validation error during user creation: %w", err)
	}
	// Check email uniqueness
	existingUser, err := s.userRepo.FindByEmail(ctx, &createUserDTO.Email)
	if err == nil && existingUser != nil {
		s.log.Error("user with this email already exists")
		return nil, fmt.Errorf("user with this email already exists: %w", apperrors.ErrUserWithThisEmailAlreadyExists)
	}
	if err != nil && !errors.Is(err, apperrors.ErrUserNotFound) {
		return nil, fmt.Errorf("failed to check email uniqueness: %w", err)
	}
	// Password hashing
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(createUserDTO.Password), s.config.BcryptCost)
	if err != nil {
		s.log.Error("failed to hash password", "error", err.Error())
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	// Generating ID for user
	userID := uuid.New()
	// Avatar processing (if passed)
	hasAvatar := false
	if createUserDTO.Avatar != nil {
		// Validating
		if err := s.validateAvatarFile(createUserDTO.Avatar); err != nil {
			s.log.Error("avatar validation failed", "error", err.Error())
			return nil, fmt.Errorf("avatar validation failed: %w", err)
		}
		// Saving to storage
		if err := s.saveAvatarFile(userID, createUserDTO.Avatar); err != nil {
			s.log.Error("failed to save avatar to storage", "error", err.Error())
			return nil, fmt.Errorf("failed to save avatar to storage: %w", err)
		}
		hasAvatar = true
	}
	// Creating model object
	user := &model.User{
		ID:         userID,
		Email:      createUserDTO.Email,
		Password:   string(passwordHash),
		FirstName:  createUserDTO.FirstName,
		MiddleName: createUserDTO.MiddleName,
		LastName:   createUserDTO.LastName,
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
			s.log.Error("failed to create user", "error", err.Error())
			return fmt.Errorf("failed to create user: %w", err)
		}
		return nil
	})
	if err != nil {
		s.log.Error("transaction failed", "error", err.Error())
		return nil, fmt.Errorf("transaction failed: %w", err)
	}
	// TODO: add transaction for roles assigning. If user was created, but roles
	// was not assigned (e.g. not all of them exists, so it causes error), user
	// must be deleted
	// Assign roles to user
	if err := s.assignRolesToUser(ctx, userID, userExtensionsDTO, createUserDTO.RoleIDs); err != nil {
		s.log.Error("failed to assign roles to user", "error", err.Error())
		return nil, fmt.Errorf("failed to assign roles to user: %w", err)
	}
	// Get created user for response
	createdUser, err := s.userRepo.FindByID(ctx, &user.ID)
	if err != nil {
		s.log.Error("failed to fetch created user", "error", err.Error())
		return nil, fmt.Errorf("failed to fetch created user: %w", err)
	}
	return UserToDTO(createdUser), nil
}

func (s *userService) UpdateUser(ctx context.Context, id uuid.UUID, dto UpdateUserDTO) (*UserResponseDTO, error) {
	// Input data validation
	if err := s.validateUpdateUserDTO(&dto); err != nil {
		return nil, fmt.Errorf("validation error during user updating: %w", err)
	}
	// Getting existing user
	user, err := s.userRepo.FindByID(ctx, &id)
	if err != nil {
		s.log.Error("failed to get user for update", "user_id", id, "error", err)
		return nil, fmt.Errorf("failed to get user for update: %w", err)
	}
	// Updating fields
	updatedFieldsCount := 0
	if dto.FirstName != nil && *dto.FirstName != user.FirstName {
		user.FirstName = *dto.FirstName
		updatedFieldsCount++
	}
	if dto.MiddleName != nil && (user.MiddleName == nil || *dto.MiddleName != *user.MiddleName) {
		user.MiddleName = dto.MiddleName
		updatedFieldsCount++
	}
	if dto.LastName != nil && *dto.LastName != user.LastName {
		user.LastName = *dto.LastName
		updatedFieldsCount++
	}
	// Updating user in DB
	if err := s.userRepo.Update(ctx, user); err != nil {
		s.log.Error("failed to update the user")
		return nil, fmt.Errorf("failed to update the user: %w", err)
	}
	// Get updated user for response
	updatedUser, err := s.userRepo.FindByID(ctx, &user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated user: %w", err)
	}
	return UserToDTO(updatedUser), nil
}

func (s *userService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	s.log.Info("Starting user deletion...")
	// Getting existing user
	user, err := s.userRepo.FindByID(ctx, &id)
	if err != nil {
		return fmt.Errorf("failed to get user for delete: %w", err)
	}
	// Transaction for user deletion
	err = s.db.Transaction(func(tx *gorm.DB) error {
		txRepo := repository.NewUserRepository(tx, s.log)
		if user.HasAvatar {
			s.log.Info("Removing user avatar file...")
			s.removeAvatarFile(id)
		}
		if err := txRepo.Delete(ctx, &id); err != nil {
			s.log.Error("failed to delete the user")
			return fmt.Errorf("failed to delete the user: %w", err)
		}
		s.log.Info("user deleted successfully")
		return nil
	})
	if err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}
	return nil
}

func (s *userService) GetUserByID(ctx context.Context, id uuid.UUID) (*UserResponseDTO, error) {
	user, err := s.userRepo.FindByID(ctx, &id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return UserToDTO(user), nil
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*UserResponseDTO, error) {
	user, err := s.userRepo.FindByEmail(ctx, &email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return UserToDTO(user), nil
}

func (s *userService) GetUsers(ctx context.Context, filter repository.UserFilter) ([]UserResponseDTO, error) {
	users, err := s.userRepo.FindAll(ctx, &filter)
	if err != nil {
		roleID := ""
		if filter.RoleID != nil {
			roleID = fmt.Sprintf("%d", *filter.RoleID)
		}
		s.log.Error(
			"failed to get users from repository",
			"role id",
			roleID,
			"limit",
			filter.Limit,
			"offset",
			filter.Offset,
			"error",
			err,
		)
		return nil, fmt.Errorf(
			"failed to get users from repository (role id: %s, limit: %d, offset: %d): %w",
			roleID,
			filter.Limit,
			filter.Offset,
			err,
		)
	}
	userDTOs := make([]UserResponseDTO, len(users))
	for i, user := range users {
		userDTOs[i] = *UserToDTO(&user)
	}
	s.log.Info("successfully received the list of users")
	return userDTOs, nil
}

func (s *userService) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]RoleResponseDTO, error) {
	// Get user by ID
	var user model.User
	if err := s.db.WithContext(ctx).
		Preload("Roles").
		Preload("Roles.Permissions").
		First(&user, userID).Error; err != nil {
		return nil, fmt.Errorf("user not found: %s: %w", err.Error(), apperrors.ErrUserNotFound)
	}
	// Get user's roles, convert them to DTO
	dtos := make([]RoleResponseDTO, len(user.Roles))
	for i, role := range user.Roles {
		dtos[i] = *RoleToDTO(&role)
	}
	// Return response
	return dtos, nil
}

func (s *userService) AssignRolesToUser(ctx context.Context, userID uuid.UUID, dto UserExtensionsDTO, roleIDs []uint16) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Check if required fields in DTO are filled in
		// TODO: move to validation function
		if slices.Contains(roleIDs, 3) && (dto.InstitutionAdministratorPositionID == nil) {
			return fmt.Errorf("bad request: required special fields for the institution administrator role cannot be empty: %w", apperrors.ErrRequiredField)
		}
		if slices.Contains(roleIDs, 4) && (dto.StaffPositionID == nil) {
			return fmt.Errorf("bad request: required special fields for the staff role cannot be empty: %w", apperrors.ErrRequiredField)
		}
		if slices.Contains(roleIDs, 5) && (len(dto.TeacherSubjectIDs) == 0) {
			return fmt.Errorf("bad request: required special fields for the teacher role cannot be empty: %w", apperrors.ErrRequiredField)
		}
		if slices.Contains(roleIDs, 7) && (dto.StudentGroupID == nil) {
			return fmt.Errorf("bad request: required special fields for the student role cannot be empty: %w", apperrors.ErrRequiredField)
		}
		// Get user by ID
		var user model.User
		if err := tx.First(&user, userID).Error; err != nil {
			return fmt.Errorf("user not found: %s: %w", err.Error(), apperrors.ErrUserNotFound)
		}
		// Get roles by IDs
		var roles []model.Role
		if len(roleIDs) > 0 {
			if err := tx.Where("id IN (?)", roleIDs).Find(&roles).Error; err != nil {
				return fmt.Errorf("failed to fetch roles: %w", err)
			}
			if len(roles) != len(roleIDs) {
				return fmt.Errorf("some roles not found: %w", apperrors.ErrNotFound)
			}
		}
		// Return response
		return s.assignRolesToUser(ctx, userID, dto, roleIDs)
	})
}

func (s *userService) AssignNonAdminRolesToUser(ctx context.Context, userID uuid.UUID, dto UserExtensionsDTO, roleIDs []uint16) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Filter IDs: skip superadmin and admin roles
		var filteredIDs []uint16
		for _, id := range roleIDs {
			if id != 1 && id != 2 {
				filteredIDs = append(filteredIDs, id)
			}
		}
		roleIDs = filteredIDs
		// TODO: move to the validation function
		// Check if required fields in DTO are filled in
		if slices.Contains(roleIDs, 3) && (dto.InstitutionAdministratorPositionID == nil) {
			return fmt.Errorf("bad request: required special fields for the institution administrator role cannot be empty: %w", apperrors.ErrRequiredField)
		}
		if slices.Contains(roleIDs, 4) && (dto.StaffPositionID == nil) {
			return fmt.Errorf("bad request: required special fields for the staff role cannot be empty: %w", apperrors.ErrRequiredField)
		}
		if slices.Contains(roleIDs, 5) && (len(dto.TeacherSubjectIDs) == 0) {
			return fmt.Errorf("bad request: required special fields for the teacher role cannot be empty: %w", apperrors.ErrRequiredField)
		}
		if slices.Contains(roleIDs, 7) && (dto.StudentGroupID == nil) {
			return fmt.Errorf("bad request: required special fields for the student role cannot be empty: %w", apperrors.ErrRequiredField)
		}
		// Get user by ID
		var user model.User
		if err := tx.Preload("Roles").First(&user, userID).Error; err != nil {
			return fmt.Errorf("user not found: %s: %w", err.Error(), apperrors.ErrUserNotFound)
		}
		// Get roles by IDs
		var roles []model.Role
		if len(roleIDs) > 0 {
			if err := tx.Where("id IN (?)", roleIDs).Find(&roles).Error; err != nil {
				return fmt.Errorf("failed to fetch roles: %w", err)
			}
			if len(roles) != len(roleIDs) {
				return fmt.Errorf("some roles not found: %w", apperrors.ErrNotFound)
			}
		}
		// Save superadmin/admin roles if user has them
		// TODO: maybe it is better to return error if user has superadmin role? Superadmin cannot have any other role.
		var adminRoleIDs []uint16
		for _, role := range user.Roles {
			if role.ID == 1 || role.ID == 2 {
				adminRoleIDs = append(adminRoleIDs, role.ID)
			}
		}
		// Append these roles to the new roles
		roleIDs = append(adminRoleIDs, roleIDs...)
		// Return response
		return s.assignRolesToUser(ctx, userID, dto, roleIDs)
	})
}

func (s *userService) AddRolesToUser(ctx context.Context, userID uuid.UUID, dto UserExtensionsDTO, roleIDs []uint16) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Check if required fields in DTO are filled in
		if slices.Contains(roleIDs, 3) && (dto.InstitutionAdministratorPositionID == nil) {
			return fmt.Errorf("bad request: required special fields for the institution administrator role cannot be empty: %w", apperrors.ErrRequiredField)
		}
		if slices.Contains(roleIDs, 4) && (dto.StaffPositionID == nil) {
			return fmt.Errorf("bad request: required special fields for the staff role cannot be empty: %w", apperrors.ErrRequiredField)
		}
		if slices.Contains(roleIDs, 5) && (len(dto.TeacherSubjectIDs) == 0) {
			return fmt.Errorf("bad request: required special fields for the teacher role cannot be empty: %w", apperrors.ErrRequiredField)
		}
		if slices.Contains(roleIDs, 7) && (dto.StudentGroupID == nil) {
			return fmt.Errorf("bad request: required special fields for the student role cannot be empty: %w", apperrors.ErrRequiredField)
		}
		// Get roles by IDs
		var roles []model.Role
		if err := tx.Where("id IN (?)", roleIDs).Find(&roles).Error; err != nil {
			return fmt.Errorf("failed to fetch roles: %w", err)
		}
		// Add roles to user, return response
		return s.addRolesToUser(ctx, userID, dto, roleIDs)
	})
}

func (s *userService) RemoveRolesFromUser(ctx context.Context, userID uuid.UUID, roleIDs []uint16) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Get roles by ID
		var roles []model.Role
		if err := tx.Where("id IN (?)", roleIDs).Find(&roles).Error; err != nil {
			return fmt.Errorf("failed to fetch roles: %w", err)
		}
		// Delete user info from extension tables
		for _, roleID := range roleIDs {
			if err := s.removeUserFromExtensionTable(tx, userID, roleID); err != nil {
				return fmt.Errorf("cannot remove user from extension tables: %w", err)
			}
		}
		// Delete roles from user, return response
		var user model.User
		user.ID = userID
		return tx.Model(&user).Association("Roles").Delete(&roles)
	})
}

func (s *userService) RemoveRoleFromUser(ctx context.Context, userID uuid.UUID, roleID uint16) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Get role by ID
		var role model.Role
		if err := tx.Where("id = ?", roleID).First(&role).Error; err != nil {
			return fmt.Errorf("failed to fetch role: %w", err)
		}
		// Delete user info from extension tables
		if err := s.removeUserFromExtensionTable(tx, userID, roleID); err != nil {
			return fmt.Errorf("cannot remove user from extension tables: %w", err)
		}
		// Delete role from user, return response
		var user model.User
		user.ID = userID
		return tx.Model(&user).Association("Roles").Delete(&role)
	})
}

type RoleResponseDTO struct {
	ID          uint16                  `json:"id"`
	CreatedAt   string                  `json:"createdAt"`
	UpdatedAt   string                  `json:"updatedAt"`
	Name        string                  `json:"name"`
	Permissions []PermissionResponseDTO `json:"permissions"`
}

func RoleToDTO(role *model.Role) *RoleResponseDTO {
	var permissions []PermissionResponseDTO
	for _, permission := range role.Permissions {
		permissions = append(permissions, *PermissionToDTO(&permission))
	}
	return &RoleResponseDTO{
		ID:          role.ID,
		CreatedAt:   role.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   role.UpdatedAt.Format(time.RFC3339),
		Name:        role.Name,
		Permissions: permissions,
	}
}

func (s *userService) validateAvatarFile(fileHeader *multipart.FileHeader) error {
	// Check file size
	if fileHeader.Size > s.config.AvatarMaxSize {
		return fmt.Errorf("file size exceeds limit of %d bytes: %w", s.config.AvatarMaxSize, apperrors.ErrFileTooLarge)
	}
	// read info
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
		return fmt.Errorf("unsupported file type: %s (allowed: %v): %w", mimeType, s.config.AvatarAllowedMIMETypes, apperrors.ErrInvalidFileType)
	}
	return nil
}

func (s *userService) saveAvatarFile(userID uuid.UUID, fileHeader *multipart.FileHeader) error {
	// Creating directory (if not exists)
	if err := os.MkdirAll(s.config.AvatarUploadPath, 0755); err != nil {
		s.log.Error("failed to create upload directory for avatars", "error", err.Error())
		return fmt.Errorf("failed to create upload directory for avatars: %w", err)
	}
	// Opening source file
	srcFile, err := fileHeader.Open()
	if err != nil {
		s.log.Error("failed to open uploaded file", "error", err.Error())
		return fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer srcFile.Close()
	// Decode image
	img, format, err := image.Decode(srcFile)
	if err != nil {
		s.log.Error("failed to decode image", "format", format, "error", err.Error())
		return fmt.Errorf("failed to decode image (format: %s): %w", format, err)
	}
	s.log.Info("Decoded image (avatar)", "format", format)
	// Convert to RGBA
	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, bounds, img, bounds.Min, draw.Src)
	// Resize to 200x200
	const avatarSize = 200
	dst := image.NewRGBA(image.Rect(0, 0, avatarSize, avatarSize))
	draw.Draw(dst, dst.Bounds(), image.NewUniform(color.White), image.Point{}, draw.Src) // white background (instead of transparent background in PNG)
	draw.ApproxBiLinear.Scale(dst, dst.Bounds(), rgba, rgba.Bounds(), draw.Over, nil)
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
	// Encode as JPEG with 80% quality
	opts := jpeg.Options{Quality: 80}
	if err := jpeg.Encode(dstFile, dst, &opts); err != nil {
		os.Remove(filePath)
		return fmt.Errorf("failed to encode image: %w", err)
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
		return fmt.Errorf("user not found: %s: %w", err.Error(), apperrors.ErrUserNotFound)
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

func (s *userService) RemoveAvatar(ctx context.Context, userID uuid.UUID) error {
	// Getting user
	user, err := s.userRepo.FindByID(ctx, &userID)
	if err != nil {
		return fmt.Errorf("user not found: %s: %w", err.Error(), apperrors.ErrUserNotFound)
	}
	// Transaction
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		if user.HasAvatar {
			// Change avatar existence status in the database
			user.HasAvatar = false
			if err := s.userRepo.Update(ctx, user); err != nil {
				return fmt.Errorf("failed to delete user avatar: %w", err)
			}
			s.log.Info("removing user avatar file...")
			s.removeAvatarFile(userID)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}
	s.log.Info("user avatar file was successfully removed")
	return nil
}

func (s *userService) validateCreateUserDTO(dto *CreateUserDTO) error {
	if dto.Email == "" {
		return fmt.Errorf("email is required: %w", apperrors.ErrEmptyEmail)
	}
	if len(dto.Password) < 8 {
		return fmt.Errorf("password must be at least 8 characters: %w", apperrors.ErrPasswordTooShort)
	}
	if len(dto.Password) > 72 {
		// TODO: add this restriction on 72 characters in other places
		return fmt.Errorf("password must be not more than 72 characters: %w", apperrors.ErrPasswordTooLong)
	}
	if len(dto.FirstName) < 2 {
		return fmt.Errorf("first name must be at least 2 characters: %w", apperrors.ErrValueTooShort)
	}
	if dto.MiddleName != nil && len(*dto.MiddleName) < 2 {
		return fmt.Errorf("middle name must be at least 2 characters or null: %w", apperrors.ErrValueTooShort)
	}
	if len(dto.LastName) < 2 {
		return fmt.Errorf("last name must be at least 2 characters: %w", apperrors.ErrValueTooShort)
	}
	if len(dto.RoleIDs) > 0 {
		// TODO: check if all reoles exist in DB
	}
	return nil
}

func (s *userService) validateUpdateUserDTO(dto *UpdateUserDTO) error {
	if dto.FirstName != nil && len(*dto.FirstName) < 2 {
		return fmt.Errorf("first name must be at least 2 characters or null: %w", apperrors.ErrValueTooShort)
	}
	if dto.MiddleName != nil && len(*dto.MiddleName) < 2 {
		return fmt.Errorf("middle name must be at least 2 characters or null: %w", apperrors.ErrValueTooShort)
	}
	if dto.LastName != nil && len(*dto.LastName) < 2 {
		return fmt.Errorf("last name must be at least 2 characters or null: %w", apperrors.ErrValueTooShort)
	}
	return nil
}

func UserToDTO(user *model.User) *UserResponseDTO {
	var roles []RoleResponseDTO
	for _, role := range user.Roles {
		roles = append(roles, *RoleToDTO(&role))
	}
	return &UserResponseDTO{
		ID:         user.ID,
		CreatedAt:  user.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  user.UpdatedAt.Format(time.RFC3339),
		Email:      user.Email,
		FirstName:  user.FirstName,
		MiddleName: user.MiddleName,
		LastName:   user.LastName,
		HasAvatar:  user.HasAvatar,
		Roles:      roles,
	}
}

func (s *userService) assignRolesToUser(ctx context.Context, userID uuid.UUID, dto UserExtensionsDTO, roleIDs []uint16) error {
	// TODO: now it is not supposed that length of roleIDs can be 0. So maybe in
	// this case all roles should be deleted: think about it.
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Block attempt to assign superadmin role
		if slices.Contains(roleIDs, 1) {
			return fmt.Errorf("you cannot assign superadmin role to user: %w", apperrors.ErrForbidden)
		}
		// Get user
		var user model.User
		if err := tx.WithContext(ctx).
			Preload("Roles").
			First(&user, "id = ?", userID).Error; err != nil {
			return fmt.Errorf("user with ID %s was not found: %s: %w", userID, err.Error(), apperrors.ErrUserNotFound)
		}
		// Get roles to assign
		var roles []model.Role
		if err := tx.WithContext(ctx).
			Where("id IN (?)", roleIDs).
			Find(&roles).Error; err != nil {
			return fmt.Errorf("failed to fetch roles for assigning: %w", err)
		}
		// Check if all roles were found
		if len(roles) != len(roleIDs) {
			return fmt.Errorf("%d role(-s) was(were) not found: %w", len(roleIDs)-len(roles), apperrors.ErrNotFound)
		}
		// Remove user from old tables-extensions ("teachers" for users with the
		// teacher role, e.g.; this tables contains specific information for
		// users)
		for _, role := range user.Roles {
			if err := s.removeUserFromExtensionTable(tx, userID, role.ID); err != nil {
				return fmt.Errorf("cannot remove user from extension table: %w", err)
			}
		}
		// Add user to new extention tables
		for _, role := range roles {
			if err := s.addUserToExtensionTable(ctx, tx, userID, dto, role.ID); err != nil {
				return fmt.Errorf("cannot add user to extension table: %w", err)
			}
		}
		// Replace old roles with new ones
		if err := tx.WithContext(ctx).
			Model(&user).
			Association("Roles").
			Replace(&roles); err != nil {
			return fmt.Errorf("failed to assign roles: %w", err)
		}
		return nil
	})
}

func (s *userService) addRolesToUser(ctx context.Context, userID uuid.UUID, dto UserExtensionsDTO, roleIDs []uint16) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Block attempt to add superadmin role
		if slices.Contains(roleIDs, 1) {
			return fmt.Errorf("you cannot add superadmin role to user: %w", apperrors.ErrForbidden)
		}
		// Get user
		var user model.User
		if err := tx.WithContext(ctx).
			Preload("Roles").
			First(&user, "id = ?", userID).Error; err != nil {
			return fmt.Errorf("user with ID %s was not found: %s: %w", userID, err.Error(), apperrors.ErrUserNotFound)
		}
		// Get roles to assign
		var roles []model.Role
		if err := tx.WithContext(ctx).
			Where("id IN (?)", roleIDs).
			Find(&roles).Error; err != nil {
			return fmt.Errorf("failed to fetch roles for assigning: %w", err)
		}
		// Check if all roles were found
		if len(roles) != len(roleIDs) {
			return fmt.Errorf("%d role(-s) was(were) not found: %w", len(roleIDs)-len(roles), apperrors.ErrNotFound)
		}
		// Add user to new extention tables
		for _, role := range roles {
			if err := s.addUserToExtensionTable(ctx, tx, userID, dto, role.ID); err != nil {
				return fmt.Errorf("cannot add user to extension table: %w", err)
			}
		}
		// Add new roles to old ones
		if err := tx.WithContext(ctx).
			Model(&user).
			Association("Roles").
			Append(&roles); err != nil {
			return fmt.Errorf("failed to append roles: %w", err)
		}
		return nil
	})
}

func (s *userService) addUserToExtensionTable(ctx context.Context, tx *gorm.DB, userID uuid.UUID, dto UserExtensionsDTO, roleID uint16) error {
	switch roleID {
	case 1: // superadmin
		return fmt.Errorf("you cannot add superadmin to extension table: %w", apperrors.ErrForbidden)
	case 2: // admin (there are no extension tables)
		return nil
	case 3: // institution_administrator
		return tx.Create(&model.InstitutionAdministrator{UserID: userID, PositionID: *dto.InstitutionAdministratorPositionID}).Error
	case 4: // staff
		return tx.Create(&model.Staff{UserID: userID, PositionID: *dto.StaffPositionID}).Error
	case 5: // teacher
		teacher := &model.Teacher{UserID: userID}
		if dto.TeacherClassroomID != nil {
			room, err := s.roomRepo.FindByID(ctx, dto.TeacherClassroomID)
			if err != nil {
				return err
			}
			teacher.Classroom = room
		}
		if err := tx.Create(teacher).Error; err != nil {
			return err
		}
		if len(dto.TeacherSubjectIDs) > 0 {
			var subjects []model.Subject
			if err := tx.Where("id IN (?)", dto.TeacherSubjectIDs).Find(&subjects).Error; err != nil {
				return err
			}
			if err := tx.Model(teacher).Association("Subjects").Append(&subjects); err != nil {
				return err
			}
		}
		if len(dto.TeacherStudentGroupIDs) > 0 {
			if err := tx.Model(&model.StudentGroup{}).
				Where("id IN (?)", dto.TeacherStudentGroupIDs).
				Update("group_advisor_id", userID).Error; err != nil {
				return err
			}
		}
		return nil
	case 6: // parent
		parent := &model.Parent{UserID: userID}
		if len(dto.ParentStudentIDs) > 0 {
			var students []model.Student
			if err := tx.Where("user_id IN (?)", dto.ParentStudentIDs).Find(&students).Error; err != nil {
				return fmt.Errorf("failed to find students: %w", err)
			}
			if len(students) != len(dto.ParentStudentIDs) {
				return fmt.Errorf("some students not found: %w", apperrors.ErrNotFound)
			}
			parent.Students = &students
		}
		return tx.Create(parent).Error
	case 7: // student
		return tx.Create(&model.Student{UserID: userID, StudentGroupID: *dto.StudentGroupID}).Error
	}
	return fmt.Errorf("role with id %d does not exist: %w", roleID, apperrors.ErrNotFound)
}

func (s *userService) removeUserFromExtensionTable(tx *gorm.DB, userID uuid.UUID, roleID uint16) error {
	switch roleID {
	case 1, 2: // superadmin, admin (there are no extension tables)
		return nil
	case 3: // institution_administrator
		return tx.Where("user_id = ?", userID).Delete(&model.InstitutionAdministrator{}).Error
	case 4: // staff
		return tx.Where("user_id = ?", userID).Delete(&model.Staff{}).Error
	case 5: // teacher
		return tx.Where("user_id = ?", userID).Delete(&model.Teacher{}).Error
	case 6: // parent
		return tx.Where("user_id = ?", userID).Delete(&model.Parent{}).Error
	case 7: // student
		return tx.Where("user_id = ?", userID).Delete(&model.Student{}).Error
	}
	return fmt.Errorf("role with id %d does not exist: %w", roleID, apperrors.ErrNotFound)
}

func (s *userService) AssignExtensionsToUser(ctx context.Context, userID uuid.UUID, dto UserExtensionsDTO) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Get user by ID
		var user model.User
		if err := tx.First(&user, userID).Error; err != nil {
			return fmt.Errorf("user not found: %s: %w", err.Error(), apperrors.ErrUserNotFound)
		}
		// Get role names
		roleNames := make([]string, len(user.Roles))
		for i, role := range user.Roles {
			roleNames[i] = role.Name
		}
		// Check if required fields in DTO are filled in
		if slices.Contains(roleNames, "institution_administrator") && (dto.InstitutionAdministratorPositionID == nil) {
			return fmt.Errorf("bad request: required extensions for the institution administrator role cannot be empty: %w", apperrors.ErrRequiredField)
		}
		if slices.Contains(roleNames, "staff") && (dto.StaffPositionID == nil) {
			return fmt.Errorf("bad request: required extensions for the staff role cannot be empty: %w", apperrors.ErrRequiredField)
		}
		if slices.Contains(roleNames, "teacher") && (len(dto.TeacherSubjectIDs) == 0) {
			return fmt.Errorf("bad request: required extensions for the teacher role cannot be empty: %w", apperrors.ErrRequiredField)
		}
		if slices.Contains(roleNames, "student") && (dto.StudentGroupID == nil) {
			return fmt.Errorf("bad request: required extensions for the student role cannot be empty: %w", apperrors.ErrRequiredField)
		}
		// Return response
		return s.assignExtensionsToUser(ctx, userID, dto)
	})
}

func (s *userService) assignExtensionsToUser(ctx context.Context, userID uuid.UUID, dto UserExtensionsDTO) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Get user
		var user model.User
		if err := tx.WithContext(ctx).
			Preload("Roles").
			First(&user, "id = ?", userID).Error; err != nil {
			return fmt.Errorf("user with ID %s was not found: %s: %w", userID, err.Error(), apperrors.ErrUserNotFound)
		}
		// Remove user from old tables-extensions
		for _, role := range user.Roles {
			if err := s.removeUserFromExtensionTable(tx, userID, role.ID); err != nil {
				return fmt.Errorf("cannot remove user from extension table: %w", err)
			}
		}
		// Add user to new extention tables
		for _, role := range user.Roles {
			if err := s.addUserToExtensionTable(ctx, tx, userID, dto, role.ID); err != nil {
				return fmt.Errorf("cannot add user to extension table: %w", err)
			}
		}
		return nil
	})
}
