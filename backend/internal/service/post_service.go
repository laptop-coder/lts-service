package service

import (
	"backend/internal/model"
	"backend/internal/repository"
	"backend/pkg/logger"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/image/draw"
	_ "golang.org/x/image/webp"
	"gorm.io/gorm"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"time"
)

var (
	ErrPostNotFound = errors.New("post not found")
)

type PostService interface {
	CreatePost(ctx context.Context, dto CreatePostDTO) (*PostResponseDTO, error)
	UpdatePost(ctx context.Context, id uuid.UUID, dto UpdatePostDTO) (*PostResponseDTO, error)
	DeletePost(ctx context.Context, id uuid.UUID) error
	RemovePhoto(ctx context.Context, postID uuid.UUID) error
	UpdatePhoto(ctx context.Context, postID uuid.UUID, dto *multipart.FileHeader) error
	GetPostByID(ctx context.Context, id uuid.UUID) (*PostResponseDTO, error)
	GetPosts(ctx context.Context, filter repository.PostFilter) ([]PostResponseDTO, error)
	VerifyPost(ctx context.Context, id uuid.UUID) (*PostResponseDTO, error)
	ReturnToOwner(ctx context.Context, id uuid.UUID) (*PostResponseDTO, error)
}

type CreatePostDTO struct {
	Name        string                `form:"name" validate:"required,min=2,max=50"`
	Description string                `form:"description,omitempty" validate:"max=1000"`
	Photo       *multipart.FileHeader `form:"photo,omitempty"` // post photo file
	AuthorID    uuid.UUID             `form:"authorID" validate:"required"`
}

type UpdatePostDTO struct {
	Name        *string `form:"name,omitempty" validate:"max=50"`
	Description *string `form:"description,omitempty" validate:"max=1000"`
}

type PostResponseDTO struct {
	ID                   uuid.UUID       `json:"id"`
	CreatedAt            string          `json:"createdAt"`
	UpdatedAt            string          `json:"updatedAt"`
	Name                 string          `json:"name"`
	Description          string          `json:"description,omitempty"`
	Verified             bool            `json:"verified"`
	ThingReturnedToOwner bool            `json:"thingReturnedToOwner"`
	HasPhoto             bool            `json:"hasPhoto"`
	Author               UserResponseDTO `json:"author"`
}

type postService struct {
	postRepo repository.PostRepository
	db       *gorm.DB
	config   PostServiceConfig
	log      logger.Logger
}

func NewPostService(
	postRepo repository.PostRepository,
	db *gorm.DB,
	config PostServiceConfig,
	log logger.Logger,
) PostService {
	return &postService{
		postRepo: postRepo,
		db:       db,
		config:   config,
		log:      log,
	}
}

func (s *postService) CreatePost(ctx context.Context, dto CreatePostDTO) (*PostResponseDTO, error) {
	// Input data validation
	if err := s.validateCreatePostDTO(&dto); err != nil {
		return nil, fmt.Errorf("validation error during post creation: %w", err)
	}
	// Generating ID for post
	postID := uuid.New()
	// Photo processing (if passed)
	hasPhoto := false
	if dto.Photo != nil {
		// Validating
		if err := s.validatePostPhoto(dto.Photo); err != nil {
			return nil, fmt.Errorf("post photo validation failed: %w", err)
		}
		// Saving to storage
		if err := s.savePostPhoto(postID, dto.Photo); err != nil {
			return nil, fmt.Errorf("failed to save post photo to storage: %w", err)
		}
		hasPhoto = true
	}
	// Creating model object
	post := &model.Post{
		ID:                   postID,
		Name:                 dto.Name,
		Description:          dto.Description,
		Verified:             false,
		ThingReturnedToOwner: false,
		HasPhoto:             hasPhoto,
		AuthorID:             dto.AuthorID,
	}
	// Transaction for creating post
	err := s.db.Transaction(func(tx *gorm.DB) error {
		txRepo := repository.NewPostRepository(tx, s.log)
		if err := txRepo.Create(ctx, post); err != nil {
			// Delete the saved post photo, if the transaction is rolled back
			if hasPhoto {
				s.removePostPhoto(postID)
			}
			return fmt.Errorf("failed to create post: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("transaction failed: %w", err)
	}
	// Get created post for response
	createdPost, err := s.postRepo.FindByID(ctx, &post.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch created post: %w", err)
	}
	return PostToDTO(createdPost), nil
}

func (s *postService) UpdatePost(ctx context.Context, id uuid.UUID, dto UpdatePostDTO) (*PostResponseDTO, error) {
	// Input data validation
	if err := s.validateUpdatePostDTO(&dto); err != nil {
		return nil, fmt.Errorf("validation error during post updating: %w", err)
	}
	// Getting existing post
	post, err := s.postRepo.FindByID(ctx, &id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.log.Error("Post for update was not found by id", "post id", id, "error", err)
			return nil, fmt.Errorf("post with id %s was not found: %w", id, err)
		}
		s.log.Error("Failed to get post for update", "post id", id, "error", err)
		return nil, fmt.Errorf("failed to get post for update: %w", err)
	}
	// Updating fields
	updatedFieldsCount := 0
	if dto.Name != nil && *dto.Name != post.Name {
		post.Name = *dto.Name
		updatedFieldsCount++
	}
	if dto.Description != nil && *dto.Description != post.Description {
		post.Description = *dto.Description
		updatedFieldsCount++
	}
	// Updating post in DB
	if err := s.postRepo.Update(ctx, post); err != nil {
		s.log.Error("Failed to update the post")
		return nil, fmt.Errorf("failed to update the post: %w", err)
	}
	// Get updated post for response
	updatedPost, err := s.postRepo.FindByID(ctx, &post.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated post: %w", err)
	}
	return PostToDTO(updatedPost), nil
}

func (s *postService) DeletePost(ctx context.Context, id uuid.UUID) error {
	s.log.Info("Starting post deletion...")
	// Getting existing post
	post, err := s.postRepo.FindByID(ctx, &id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.log.Error("Post for delete was not found by id", "post id", id, "error", err)
			return fmt.Errorf("post with id %s was not found: %w", id, err)
		}
		s.log.Error("Failed to get post for delete", "post id", id, "error", err)
		return fmt.Errorf("failed to get post for delete: %w", err)
	}
	// Transaction for post deletion
	err = s.db.Transaction(func(tx *gorm.DB) error {
		txRepo := repository.NewPostRepository(tx, s.log)
		if post.HasPhoto {
			s.log.Info("Removing post photo file...")
			s.removePostPhoto(id)
		}
		if err := txRepo.Delete(ctx, &id); err != nil {
			s.log.Error("Failed to delete the post")
			return fmt.Errorf("failed to delete the post: %w", err)
		}
		s.log.Info("Post deleted successfully")
		return nil
	})
	if err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}
	return nil
}

func (s *postService) RemovePhoto(ctx context.Context, postID uuid.UUID) error {
	// Getting post
	post, err := s.postRepo.FindByID(ctx, &postID)
	if err != nil {
		return fmt.Errorf("post not found: %w", err)
	}
	// Transaction
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		if post.HasPhoto {
			// Change photo existence status in the database
			post.HasPhoto = false
			if err := s.postRepo.Update(ctx, post); err != nil {
				return fmt.Errorf("failed to delete post photo: %w", err)
			}
			s.log.Info("Removing post photo file...")
			s.removePostPhoto(postID)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}
	s.log.Info("Post photo file was successfully removed")
	return nil
}

func (s *postService) UpdatePhoto(ctx context.Context, postID uuid.UUID, photo *multipart.FileHeader) error {
	post, err := s.postRepo.FindByID(ctx, &postID)
	if err != nil {
		s.log.Error("Post not found", "error", err.Error())
		return fmt.Errorf("post not found: %w", err)
	}
	// Validating the file
	if err := s.validatePostPhoto(photo); err != nil {
		s.log.Error("Failed to validate the file", "error", err.Error())
		return err
	}
	// Saving new photo
	if err := s.savePostPhoto(postID, photo); err != nil {
		s.log.Error("Failed to save new photo", "error", err.Error())
		return err
	}
	// Mark existence of the photo in the database
	post.HasPhoto = true
	if err := s.postRepo.Update(ctx, post); err != nil {
		// Rollback file saving in the case of error
		s.removePostPhoto(postID)
		s.log.Error("Failed to update post photo", "error", err.Error())
		return fmt.Errorf("failed to update post photo: %w", err)
	}
	return nil
}

func (s *postService) validatePostPhoto(fileHeader *multipart.FileHeader) error {
	// Check file size
	if fileHeader.Size > s.config.PhotoMaxSize {
		return fmt.Errorf("file size exceeds limit of %d bytes", s.config.PhotoMaxSize)
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
	if !slices.Contains(s.config.PhotoAllowedMIMETypes, mimeType) {
		return fmt.Errorf("unsupported file type: %s. Allowed: %v", mimeType, s.config.PhotoAllowedMIMETypes)
	}
	return nil
}

func (s *postService) savePostPhoto(postID uuid.UUID, fileHeader *multipart.FileHeader) error {
	// Creating directory (if not exists)
	if err := os.MkdirAll(s.config.PhotoUploadPath, 0755); err != nil {
		return fmt.Errorf("failed to create upload directory for post photos: %w", err)
	}
	// Opening source file
	srcFile, err := fileHeader.Open()
	if err != nil {
		return fmt.Errorf("failed to open uploaded file (post photo): %w", err)
	}
	defer srcFile.Close()
	// Decode image
	img, format, err := image.Decode(srcFile)
	if err != nil {
		s.log.Error("Failed to decode image", "format", format, "error", err.Error())
		return fmt.Errorf("failed to decode image (format: %s): %w", format, err)
	}
	s.log.Info("Decoded image (post photo)", "format", format)
	// Convert to RGBA
	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, bounds, img, bounds.Min, draw.Src)
	// Resize if too large (max width 1200px)
	bounds = rgba.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	maxWidth := 1200
	var dst image.Image = rgba
	if width > maxWidth {
		newHeight := height * maxWidth / width
		resized := image.NewRGBA(image.Rect(0, 0, maxWidth, newHeight))
		draw.ApproxBiLinear.Scale(resized, resized.Bounds(), rgba, bounds, draw.Over, nil)
		dst = resized
	}
	// Creating file path (where to save post photo)
	filePath := filepath.Join(
		s.config.PhotoUploadPath,
		fmt.Sprintf("%s.jpeg", postID.String()),
	)
	// Creating file in storage
	dstFile, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer dstFile.Close()
	// Encode as JPEG with 85% quality
	opts := jpeg.Options{Quality: 85}
	if err := jpeg.Encode(dstFile, dst, &opts); err != nil {
		os.Remove(filePath)
		return fmt.Errorf("failed to encode image: %w", err)
	}
	return nil
}

func (s *postService) removePostPhoto(postID uuid.UUID) {
	filePath := filepath.Join(
		s.config.PhotoUploadPath,
		fmt.Sprintf("%s.jpeg", postID.String()),
	)
	os.Remove(filePath)
}

func (s *postService) UpdatePostPhoto(ctx context.Context, postID uuid.UUID, postPhoto *multipart.FileHeader) error {
	post, err := s.postRepo.FindByID(ctx, &postID)
	if err != nil {
		return fmt.Errorf("post not found: %w", err)
	}
	// Validating the file
	if err := s.validatePostPhoto(postPhoto); err != nil {
		return err
	}
	// Saving the new photo
	if err := s.savePostPhoto(postID, postPhoto); err != nil {
		return err
	}
	// Mark existence of the photo in the database
	post.HasPhoto = true
	if err := s.postRepo.Update(ctx, post); err != nil {
		// Rollback file saving in the case of error
		s.removePostPhoto(postID)
		return fmt.Errorf("failed to update post photo: %w", err)
	}
	return nil
}

func (s *postService) RemovePostPhoto(ctx context.Context, postID uuid.UUID) error {
	// Getting post
	post, err := s.postRepo.FindByID(ctx, &postID)
	if err != nil {
		return fmt.Errorf("post not found: %w", err)
	}
	// Transaction
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		if post.HasPhoto {
			// Change photo existence status in the database
			post.HasPhoto = false
			if err := s.postRepo.Update(ctx, post); err != nil {
				return fmt.Errorf("failed to delete post photo: %w", err)
			}
			s.log.Info("Removing post photos...")
			s.removePostPhoto(postID)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}
	s.log.Info("Post photo was successfully removed")
	return nil
}

func (s *postService) GetPostByID(ctx context.Context, id uuid.UUID) (*PostResponseDTO, error) {
	post, err := s.postRepo.FindByID(ctx, &id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.log.Error("Post with id %s was not found", "error", err.Error())
			return nil, fmt.Errorf("post with id %s was not found: %w", id, err)
		}
		s.log.Error("Failed to get post", "error", err.Error())
		return nil, fmt.Errorf("failed to get post: %w", err)
	}
	return PostToDTO(post), nil
}

func (s *postService) GetPosts(ctx context.Context, filter repository.PostFilter) ([]PostResponseDTO, error) {
	posts, err := s.postRepo.FindAll(ctx, &filter)
	if err != nil {
		authorID := ""
		if filter.AuthorID != nil {
			authorID = (*filter.AuthorID).String()
		}
		verified := ""
		if filter.Verified != nil {
			verified = strconv.FormatBool(*filter.Verified)
		}
		thingReturnedToOwner := ""
		if filter.ThingReturnedToOwner != nil {
			thingReturnedToOwner = strconv.FormatBool(*filter.ThingReturnedToOwner)
		}
		s.log.Error(
			"failed to get posts from repository",
			"author id",
			authorID,
			"verified",
			verified,
			"thing returned to owner",
			thingReturnedToOwner,
			"limit",
			filter.Limit,
			"offset",
			filter.Offset,
			"error",
			err,
		)
		return nil, fmt.Errorf(
			"failed to get posts from repository (author id: %s, verified: %s, thing returned to owner: %s, limit: %d, offset: %d): %w",
			authorID,
			verified,
			thingReturnedToOwner,
			filter.Limit,
			filter.Offset,
			err,
		)
	}
	postDTOs := make([]PostResponseDTO, len(posts))
	for i, post := range posts {
		postDTOs[i] = *PostToDTO(&post)
	}
	s.log.Info("successfully received the list of posts")
	return postDTOs, nil
}

func (s *postService) VerifyPost(ctx context.Context, id uuid.UUID) (*PostResponseDTO, error) {
	// Getting existing post
	post, err := s.postRepo.FindByID(ctx, &id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.log.Error("Post for verification was not found by id", "post id", id, "error", err)
			return nil, fmt.Errorf("post with id %s was not found: %w", id, err)
		}
		s.log.Error("Failed to get post for verification", "post id", id, "error", err)
		return nil, fmt.Errorf("failed to get post for verification: %w", err)
	}
	// Updating field
	post.Verified = true
	// Updating post in DB
	if err := s.postRepo.Update(ctx, post); err != nil {
		s.log.Error("Failed to change post verification status")
		return nil, fmt.Errorf("failed to change post verification status: %w", err)
	}
	// Get verified post for response
	// TODO: refactor in the whole code, maybe re-use "post" variable instead
	//of using FindByID twice
	verifiedPost, err := s.postRepo.FindByID(ctx, &post.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch verified post: %w", err)
	}
	return PostToDTO(verifiedPost), nil
}

func (s *postService) ReturnToOwner(ctx context.Context, id uuid.UUID) (*PostResponseDTO, error) {
	// Getting existing post
	post, err := s.postRepo.FindByID(ctx, &id)
	if err != nil || post == nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.log.Error("Post for changing thing returning status was not found by id", "post id", id, "error", err)
			return nil, fmt.Errorf("post with id %s was not found: %w", id, err)
		}
		s.log.Error("Failed to get post for changing thing returning status", "post id", id, "error", err)
		return nil, fmt.Errorf("failed to get post for changing thing returning status: %w", err)
	}
	// Check if the post verified
	if post.Verified != true {
		s.log.Error("Failed to mark thing as returned to owner for not verified post", "post id", id)
		return nil, fmt.Errorf("forbidden: failed to mark thing as returned to owner for not verified post")
	}
	// Updating field
	post.ThingReturnedToOwner = true
	// Updating post in DB
	if err := s.postRepo.Update(ctx, post); err != nil {
		s.log.Error("Failed to change thing returning status")
		return nil, fmt.Errorf("failed to change thing returning status: %w", err)
	}
	// Get updated post for response
	updatedPost, err := s.postRepo.FindByID(ctx, &post.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch post with changed thing returning status: %w", err)
	}
	return PostToDTO(updatedPost), nil
}

func (s *postService) validateCreatePostDTO(dto *CreatePostDTO) error {
	return nil
}

func (s *postService) validateUpdatePostDTO(dto *UpdatePostDTO) error {
	return nil
}

func PostToDTO(post *model.Post) *PostResponseDTO {
	return &PostResponseDTO{
		ID:                   post.ID,
		CreatedAt:            post.CreatedAt.Format(time.RFC3339),
		UpdatedAt:            post.UpdatedAt.Format(time.RFC3339),
		Name:                 post.Name,
		Description:          post.Description,
		Verified:             post.Verified,
		ThingReturnedToOwner: post.ThingReturnedToOwner,
		HasPhoto:             post.HasPhoto,
		Author:               *UserToDTO(&post.Author),
	}
}
