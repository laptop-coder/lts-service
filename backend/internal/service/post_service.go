package service

import (
	"backend/internal/model"
	"backend/internal/permissions"
	"backend/internal/repository"
	"backend/pkg/appcontext"
	"backend/pkg/apperrors"
	"backend/pkg/imghash"
	"backend/pkg/logger"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/valkey-io/valkey-go"
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
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"
)

type PostService interface {
	CreatePost(ctx context.Context, dto CreatePostDTO, canVerifyPost bool) (*PostResponseDTO, error)
	UpdatePost(ctx context.Context, id uuid.UUID, dto UpdatePostDTO) (*PostResponseDTO, error)
	DeletePost(ctx context.Context, id uuid.UUID) error
	RemovePhoto(ctx context.Context, postID uuid.UUID) error
	UpdatePhoto(ctx context.Context, postID uuid.UUID, dto *multipart.FileHeader) error
	GetPostByID(ctx context.Context, id uuid.UUID) (*PostResponseDTO, error)
	GetPosts(ctx context.Context, filter repository.PostFilter) ([]PostResponseDTO, error)
	ChangePostModerationStatus(ctx context.Context, id uuid.UUID, newStatus model.ModerationStatus, moderatorID *uuid.UUID, rejectReason *string) (*PostModerationResponseDTO, error)
	ReturnToOwner(ctx context.Context, id uuid.UUID) (*PostResponseDTO, error)
	GetSimilar(ctx context.Context, dto *GetSimilarDTO) ([]PostResponseDTO, error)
	CalcAllPhotosHashes(ctx context.Context) error
	ModeratePost(ctx context.Context, postID uuid.UUID) error
	ModerateAllPosts(ctx context.Context) error
	GetOldestPendingPost(ctx context.Context) (*PostResponseDTO, error)
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
	ID                   uuid.UUID                 `json:"id"`
	CreatedAt            string                    `json:"createdAt"`
	UpdatedAt            string                    `json:"updatedAt"`
	Name                 string                    `json:"name"`
	Description          string                    `json:"description,omitempty"`
	ThingReturnedToOwner bool                      `json:"thingReturnedToOwner"`
	HasPhoto             bool                      `json:"hasPhoto"`
	Author               UserResponseDTO           `json:"author"`
	Moderation           PostModerationResponseDTO `json:"moderation"`
}

type PostModerationResponseDTO struct {
	PostID    uuid.UUID `json:"postId"`
	CreatedAt string    `json:"createdAt"`
	UpdatedAt string    `json:"updatedAt"`

	Status        model.ModerationStatus `json:"status"`
	ModeratorID   *uuid.UUID             `json:"moderatorId,omitempty"`
	ModeratorUser *UserResponseDTO       `json:"moderatorUser,omitempty"`
	RejectReason  *string                `json:"rejectReason,omitempty"`
}

type GetSimilarDTO struct {
	ID          *uuid.UUID            `json:"id,omitempty"`
	Name        *string               `json:"name,omitempty"`
	Description *string               `json:"description,omitempty"`
	Photo       *multipart.FileHeader `form:"photo,omitempty"` // post photo file
	HasPhoto    bool                  `form:"hasPhoto"`
}

type postService struct {
	postRepo           repository.PostRepository
	postModerationRepo repository.PostModerationRepository
	userService        UserService
	hashCalc           imghash.HashCalculator
	db                 *gorm.DB
	client             valkey.Client
	config             PostServiceConfig
	log                logger.Logger
}

func NewPostService(
	postRepo repository.PostRepository,
	postModerationRepo repository.PostModerationRepository,
	userService UserService,
	hashCalc imghash.HashCalculator,
	db *gorm.DB,
	client valkey.Client,
	config PostServiceConfig,
	log logger.Logger,
) PostService {
	return &postService{
		postRepo:           postRepo,
		postModerationRepo: postModerationRepo,
		userService:        userService,
		hashCalc:           hashCalc,
		db:                 db,
		client:             client,
		config:             config,
		log:                log,
	}
}

func (s *postService) CreatePost(ctx context.Context, dto CreatePostDTO, canVerifyPost bool) (*PostResponseDTO, error) {
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
		if err := s.savePostPhoto(ctx, postID, dto.Photo); err != nil {
			return nil, fmt.Errorf("failed to save post photo to storage: %w", err)
		}
		hasPhoto = true
	}
	// Creating post model object
	post := &model.Post{
		ID:                   postID,
		Name:                 dto.Name,
		Description:          dto.Description,
		ThingReturnedToOwner: false,
		HasPhoto:             hasPhoto,
		AuthorID:             dto.AuthorID,
	}
	// Automatically verify post if user has permission to verify posts
	verified := canVerifyPost
	// Creating moderation model object
	moderation := &model.PostModeration{
		PostID: postID,
		Status: func() model.ModerationStatus {
			if verified {
				return model.ModerationStatusApproved
			}
			return model.ModerationStatusPending
		}(),
		ModeratorID: func() *uuid.UUID {
			if verified {
				return &(dto.AuthorID)
			}
			return nil
		}(),
	}
	// Transaction for creating post
	err := s.db.Transaction(func(tx *gorm.DB) error {
		txPostRepo := repository.NewPostRepository(tx, s.client, s.log)
		if err := txPostRepo.Create(ctx, post); err != nil {
			// Delete the saved post photo, if the transaction is rolled back
			if hasPhoto {
				if err := s.removePostPhoto(ctx, postID); err != nil {
					s.log.Error("failed to create post, the transaction is rolled back, and failed to delete saved post photo", "error", err.Error())
					return fmt.Errorf("failed to create post, the transaction is rolled back, and failed to delete saved post photo: %w", err)
				}
			}
			s.log.Error("failed to create post", "error", err.Error())
			return fmt.Errorf("failed to create post: %w", err)
		}
		txModerationRepo := repository.NewPostModerationRepository(tx, s.log)
		if err := txModerationRepo.Create(ctx, moderation); err != nil {
			s.log.Error("failed to create post moderation", "error", err.Error())
			return fmt.Errorf("failed to create post moderation: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("transaction failed: %w", err)
	}
	if !verified {
		// Add post to moderation queue
		err = s.client.Do(
			ctx,
			s.client.B().
				Rpush().
				Key("moderation:posts:queue").
				Element(postID.String()).
				Build(),
		).Error()
		if err != nil {
			return nil, fmt.Errorf("failed to push post id to moderation queue: %w", err)
		}
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
		s.log.Error("failed to update the post")
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
	s.log.Info("starting post deletion...")
	// Getting existing post
	post, err := s.postRepo.FindByID(ctx, &id)
	if err != nil {
		return fmt.Errorf("failed to get post for delete: %w", err)
	}
	// Transaction for post deletion
	err = s.db.Transaction(func(tx *gorm.DB) error {
		txRepo := repository.NewPostRepository(tx, s.client, s.log)
		if post.HasPhoto {
			s.log.Info("removing post photo file...")
			if err := s.removePostPhoto(ctx, id); err != nil {
				s.log.Error("failed to delete post photo file", "error", err.Error())
				return fmt.Errorf("failed to delete post photo file: %w", err)
			}
		}
		if err := txRepo.Delete(ctx, &id); err != nil {
			s.log.Error("failed to delete the post")
			return fmt.Errorf("failed to delete the post: %w", err)
		}
		s.log.Info("post deleted successfully")
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
		return fmt.Errorf("post not found: %s: %w", err.Error(), apperrors.ErrPostNotFound)
	}
	// Transaction
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		if post.HasPhoto {
			// Change photo existence status in the database
			post.HasPhoto = false
			if err := s.postRepo.Update(ctx, post); err != nil {
				return fmt.Errorf("failed to delete post photo: %w", err)
			}
			s.log.Info("removing post photo file...")
			if err := s.removePostPhoto(ctx, postID); err != nil {
				s.log.Error("failed to delete post photo file", "error", err.Error())
				return fmt.Errorf("failed to delete post photo file: %w", err)
			}
		}
		return nil
	}); err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}
	s.log.Info("post photo file was successfully removed")
	return nil
}

func (s *postService) UpdatePhoto(ctx context.Context, postID uuid.UUID, photo *multipart.FileHeader) error {
	post, err := s.postRepo.FindByID(ctx, &postID)
	if err != nil {
		s.log.Error("post not found", "error", err.Error())
		return fmt.Errorf("post not found: %s: %w", err.Error(), apperrors.ErrPostNotFound)
	}
	// Validating the file
	if err := s.validatePostPhoto(photo); err != nil {
		s.log.Error("failed to validate the file", "error", err.Error())
		return err
	}
	// Saving new photo
	if err := s.savePostPhoto(ctx, postID, photo); err != nil {
		s.log.Error("failed to save new photo", "error", err.Error())
		return err
	}
	// Mark existence of the photo in the database
	post.HasPhoto = true
	if err := s.postRepo.Update(ctx, post); err != nil {
		// Rollback file saving in the case of error
		if err := s.removePostPhoto(ctx, postID); err != nil {
			s.log.Error("failed to update post photo and failed to rollback file saving", "error", err.Error())
			return fmt.Errorf("failed to update post photo and failed to rollback file saving: %w", err)
		}
		s.log.Error("failed to update post photo", "error", err.Error())
		return fmt.Errorf("failed to update post photo: %w", err)
	}
	return nil
}

func (s *postService) validatePostPhoto(fileHeader *multipart.FileHeader) error {
	// Check file size
	if fileHeader.Size > s.config.PhotoMaxSize {
		return fmt.Errorf("file size exceeds limit of %d bytes: %w", s.config.PhotoMaxSize, apperrors.ErrFileTooLarge)
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
		return fmt.Errorf("unsupported file type: %s (allowed: %v): %w", mimeType, s.config.PhotoAllowedMIMETypes, apperrors.ErrInvalidFileType)
	}
	return nil
}

func (s *postService) savePostPhoto(ctx context.Context, postID uuid.UUID, fileHeader *multipart.FileHeader) error {
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
		s.log.Error("failed to decode image", "format", format, "error", err.Error())
		return fmt.Errorf("failed to decode image (format: %s): %w", format, err)
	}
	s.log.Info("decoded image (post photo)", "format", format)
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
	// Calculate image hash
	hash, err := s.hashCalc.PerceptualHash(dst)
	if err != nil {
		s.log.Error("failed to calculate image hash")
		return fmt.Errorf("failed to calculate image hash: %w", err)
	}
	// Save hash
	if err := s.postRepo.UpdatePhotoHash(ctx, postID, hash); err != nil {
		s.log.Error("failed to save hash of the post photo", "error", err.Error())
		return fmt.Errorf("failed to save hash of the post photo: %w", err)
	}
	return nil
}

func (s *postService) removePostPhoto(ctx context.Context, postID uuid.UUID) error {
	// Remove photo (soft-delete)
	filename := fmt.Sprintf("%s.jpeg", postID.String())
	filePath := filepath.Join(
		s.config.PhotoUploadPath,
		filename,
	)
	fileDeletePath := filepath.Join(
		s.config.PhotoDeletePath,
		filename,
	)
	if err := os.Rename(filePath, fileDeletePath); err != nil {
		s.log.Error("failed to move post photo to trash", "error", err.Error())
		return fmt.Errorf("failed to move post photo to trash: %w", err)
	}
	// Remove photo hash
	if err := s.postRepo.DeletePhotoHash(ctx, postID); err != nil {
		s.log.Error("failed to delete post photo hash", "error", err.Error())
		return fmt.Errorf("failed to delete post photo hash: %w", err)
	}
	return nil
}

// TODO: what is it? Does it duplicate UpdatePost?
func (s *postService) UpdatePostPhoto(ctx context.Context, postID uuid.UUID, postPhoto *multipart.FileHeader) error {
	post, err := s.postRepo.FindByID(ctx, &postID)
	if err != nil {
		return fmt.Errorf("post not found: %s: %w", err.Error(), apperrors.ErrPostNotFound)
	}
	// Validating the file
	if err := s.validatePostPhoto(postPhoto); err != nil {
		return err
	}
	// Saving the new photo
	if err := s.savePostPhoto(ctx, postID, postPhoto); err != nil {
		return err
	}
	// Mark existence of the photo in the database
	post.HasPhoto = true
	if err := s.postRepo.Update(ctx, post); err != nil {
		// Rollback file saving in the case of error
		if err := s.removePostPhoto(ctx, postID); err != nil {
			s.log.Error("failed to update post photo and failed to rollback file saving", "error", err.Error())
			return fmt.Errorf("failed to update post photo and failed to rollback file saving: %w", err)
		}
		return fmt.Errorf("failed to update post photo: %w", err)
	}
	return nil
}

func (s *postService) RemovePostPhoto(ctx context.Context, postID uuid.UUID) error {
	// Getting post
	post, err := s.postRepo.FindByID(ctx, &postID)
	if err != nil {
		return fmt.Errorf("post not found: %s: %w", err.Error(), apperrors.ErrPostNotFound)
	}
	// Transaction
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		if post.HasPhoto {
			// Change photo existence status in the database
			post.HasPhoto = false
			if err := s.postRepo.Update(ctx, post); err != nil {
				return fmt.Errorf("failed to delete post photo: %w", err)
			}
			s.log.Info("removing post photos...")
			if err := s.removePostPhoto(ctx, postID); err != nil {
				s.log.Error("failed to delete post photo", "error", err.Error())
				return fmt.Errorf("failed to delete post photo: %w", err)
			}
		}
		return nil
	}); err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}
	s.log.Info("post photo was successfully removed")
	return nil
}

func (s *postService) GetPostByID(ctx context.Context, id uuid.UUID) (*PostResponseDTO, error) {
	post, err := s.postRepo.FindByID(ctx, &id)
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}
	return PostToDTO(post), nil
}

func (s *postService) GetPosts(ctx context.Context, filter repository.PostFilter) ([]PostResponseDTO, error) {
	posts, err := s.postRepo.FindAll(ctx, &filter)
	if err != nil {
		authorIDs := []string{}
		if len(filter.AuthorIDs) > 0 {
			for _, id := range filter.AuthorIDs {
				authorIDs = append(authorIDs, id.String())
			}
		}
		thingReturnedToOwner := ""
		if filter.ThingReturnedToOwner != nil {
			thingReturnedToOwner = strconv.FormatBool(*filter.ThingReturnedToOwner)
		}
		s.log.Error(
			"failed to get posts from repository",
			"author ids",
			authorIDs,
			"moderation_statuses",
			fmt.Sprintf("%v", filter.ModerationStatuses),
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
			"failed to get posts from repository (author ids: %v, moderation_statuses: %v, thing returned to owner: %s, limit: %d, offset: %d): %w",
			authorIDs,
			filter.ModerationStatuses,
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

func (s *postService) ChangePostModerationStatus(ctx context.Context, postID uuid.UUID, newStatus model.ModerationStatus, moderatorID *uuid.UUID, rejectReason *string) (*PostModerationResponseDTO, error) {
	// If reject reason is specified, check that new post status is "rejected"
	if newStatus != model.ModerationStatusRejected &&
		newStatus != model.ModerationStatusAutoRejected &&
		rejectReason != nil {
		s.log.Error(
			"cannot specify reject reason for not rejected post",
			"current_new_status",
			string(newStatus),
			"required_new_status",
			fmt.Sprintf(
				"%s or %s",
				string(model.ModerationStatusRejected),
				string(model.ModerationStatusAutoRejected),
			),
		)
		return nil, fmt.Errorf("cannot specify reject reason for not rejected post (current new status is %s, but required %s or %s)", string(newStatus), string(model.ModerationStatusRejected), string(model.ModerationStatusAutoRejected)) // TODO: return Bad Request
	}
	// Getting existing post moderation
	moderation, err := s.postModerationRepo.FindByID(ctx, &postID)
	if err != nil {
		return nil, fmt.Errorf("failed to get post moderation for changing moderation status: %w", err)
	}
	// Updating fields
	moderation.Status = newStatus
	moderation.ModeratorID = moderatorID
	moderation.RejectReason = rejectReason
	// Updating post moderation in DB
	if err := s.postModerationRepo.Update(ctx, moderation); err != nil {
		s.log.Error("failed to change post moderation status")
		return nil, fmt.Errorf("failed to change post moderation status: %w", err)
	}
	// Get post moderation for response
	// TODO: refactor in the whole code, maybe re-use "moderation" variable instead
	//of using FindByID twice
	changedModeration, err := s.postModerationRepo.FindByID(ctx, &moderation.PostID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch changed post moderation: %w", err)
	}
	return ModerationToDTO(changedModeration), nil
}

func (s *postService) ReturnToOwner(ctx context.Context, id uuid.UUID) (*PostResponseDTO, error) {
	// Getting existing post
	post, err := s.postRepo.FindByID(ctx, &id)
	if err != nil || post == nil {
		return nil, fmt.Errorf("failed to get post for changing thing returning status: %w", err)
	}
	// Check if the post approved
	if post.Moderation.Status != model.ModerationStatusAutoApproved && post.Moderation.Status != model.ModerationStatusApproved {
		s.log.Error("failed to mark thing as returned to owner for not approved post", "post id", id)
		return nil, fmt.Errorf("failed to mark thing as returned to owner for not approved post: %w", apperrors.ErrForbidden)
	}
	// Updating field
	post.ThingReturnedToOwner = true
	// Updating post in DB
	if err := s.postRepo.Update(ctx, post); err != nil {
		s.log.Error("failed to change thing returning status")
		return nil, fmt.Errorf("failed to change thing returning status: %w", err)
	}
	// Get updated post for response
	updatedPost, err := s.postRepo.FindByID(ctx, &post.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch post with changed thing returning status: %w", err)
	}
	return PostToDTO(updatedPost), nil
}

func (s *postService) GetSimilar(ctx context.Context, dto *GetSimilarDTO) ([]PostResponseDTO, error) {
	if dto.ID == nil && dto.Name == nil && dto.Description == nil && dto.Photo == nil && !dto.HasPhoto {
		s.log.Error("failed to get similar posts: at least one search parameter must be specified")
		return nil, fmt.Errorf("at least one search parameter must be specified")
	}
	if dto.HasPhoto && dto.ID == nil {
		s.log.Error("failed to get similar posts: the ID must be specified if HasPhoto is true")
		return nil, fmt.Errorf("failed to get similar posts: the ID must be specified if HasPhoto is true")
	}
	var (
		imageMatches       []model.Post
		nameMatches        []model.Post
		descriptionMatches []model.Post
		err                error
	)
	// If post has photo, the photo by ID has the priority over the passed file
	if dto.HasPhoto && dto.ID != nil {
		// Read file
		file, err := os.Open(filepath.Join(s.config.PhotoUploadPath, fmt.Sprintf("%s.jpeg", (*dto.ID).String())))
		if err != nil || file == nil { // TODO: if file is nil, cannot use err.Error(), because err is nil (check the whole code and refactor)
			s.log.Error("failed to open file", "error", err.Error())
			return nil, fmt.Errorf("failed to open file: %w", err)
		}
		defer file.Close()
		// Decode as JPEG
		img, err := jpeg.Decode(file)
		if err != nil {
			s.log.Error("failed to decode JPEG image")
			return nil, fmt.Errorf("failed to decode JPEG image: %w", err)
		}
		// Calculate image hash
		hash, err := s.hashCalc.PerceptualHash(img)
		if err != nil {
			s.log.Error("failed to calculate image hash")
			return nil, fmt.Errorf("failed to calculate image hash: %w", err)
		}
		// Get matches
		imageMatches, err = s.postRepo.FindSimilarByImageHashDistance(ctx, hash, 25)
		if err != nil {
			s.log.Error("failed to find image matches")
			return nil, fmt.Errorf("failed to find image matches: %w", err)
		}
	} else if dto.Photo != nil {
		// Check file size
		if err := s.validatePostPhoto(dto.Photo); err != nil {
			return nil, fmt.Errorf("failed to validate post photo: %w", err)
		}
		// read info
		file, err := dto.Photo.Open()
		if err != nil {
			return nil, fmt.Errorf("failed to open file: %w", err)
		}
		defer file.Close()
		// Decode image
		img, format, err := image.Decode(file)
		if err != nil {
			s.log.Error("failed to decode image", "format", format, "error", err.Error())
			return nil, fmt.Errorf("failed to decode image (format: %s): %w", format, err)
		}
		s.log.Info("decoded image (post photo)", "format", format)
		// TODO: refactor (the code is duplicated)
		// Calculate image hash
		hash, err := s.hashCalc.PerceptualHash(img)
		if err != nil {
			s.log.Error("failed to calculate image hash")
			return nil, fmt.Errorf("failed to calculate image hash: %w", err)
		}
		// Get matches
		imageMatches, err = s.postRepo.FindSimilarByImageHashDistance(ctx, hash, 25)
		if err != nil {
			s.log.Error("failed to find image matches")
			return nil, fmt.Errorf("failed to find image matches: %w", err)
		}
	}
	if dto.Name != nil {
		// Get matches
		nameMatches, err = s.postRepo.FindSimilarByName(ctx, *dto.Name)
		if err != nil {
			s.log.Error("failed to find name matches")
			return nil, fmt.Errorf("failed to find name matches: %w", err)
		}
	}
	if dto.Description != nil {
		// Get matches
		descriptionMatches, err = s.postRepo.FindSimilarByDescription(ctx, *dto.Description)
		if err != nil {
			s.log.Error("failed to find description matches")
			return nil, fmt.Errorf("failed to find description matches: %w", err)
		}
	}
	// Collect all matches
	scores := make(map[uuid.UUID]float64)
	for _, p := range imageMatches {
		scores[p.ID] += 0.5
	}
	for _, p := range nameMatches {
		scores[p.ID] += 0.3
	}
	for _, p := range descriptionMatches {
		scores[p.ID] += 0.2
	}
	// Sort by scores
	type keyValue struct {
		Key   uuid.UUID
		Value float64
	}
	var sortedScores []keyValue
	for key, value := range scores {
		if dto.ID != nil && *dto.ID == key {
			continue
		}
		sortedScores = append(sortedScores, keyValue{key, value})
	}
	sort.Slice(sortedScores, func(i, j int) bool {
		return sortedScores[i].Value > sortedScores[j].Value
	})
	// Get posts' info
	// TODO: optimize, don't collect info about posts if threre are already 10 posts
	var postDTOs []PostResponseDTO
	for _, kv := range sortedScores {
		id := kv.Key
		// Get post by ID
		post, err := s.postRepo.FindByID(ctx, &id)
		if err != nil || post == nil {
			return nil, fmt.Errorf("failed to fetch matched post: %w", err)
		}
		// Convert to DTO
		postDTOs = append(postDTOs, *PostToDTO(post))
	}
	// Filter posts
	var filteredPosts []PostResponseDTO
	// Check if user is authorized
	userPermissions, ok := ctx.Value(appcontext.UserPermissionsKey).([]string)
	if ok {
		s.log.Debug("user is authorized")
		// Filter posts depends on user's permissions
		if slices.Contains(userPermissions, permissions.PostReadAny) {
			s.log.Debug("user can read any posts")
			filteredPosts = postDTOs
		} else if slices.Contains(userPermissions, permissions.PostReadOwn) {
			s.log.Debug("user can read own posts")
			userID, ok := ctx.Value(appcontext.UserIDKey).(uuid.UUID)
			if !ok {
				return nil, fmt.Errorf("failed to get user ID from the context and convert it to UUID")
			}
			for _, post := range postDTOs {
				if post.Author.ID == userID {
					filteredPosts = append(filteredPosts, post)
				}
			}
		} else {
			// Show only verified posts
			for _, post := range postDTOs {
				if post.Moderation.Status == model.ModerationStatusApproved ||
					post.Moderation.Status == model.ModerationStatusAutoApproved {
					filteredPosts = append(filteredPosts, post)
				}
			}
		}
	} else {
		s.log.Debug("user is not authorized")
		// Show only verified posts
		for _, post := range postDTOs {
			if post.Moderation.Status == model.ModerationStatusApproved ||
				post.Moderation.Status == model.ModerationStatusAutoApproved {
				filteredPosts = append(filteredPosts, post)
			}
		}
	}
	// Get first 10 posts
	if len(filteredPosts) > 10 {
		filteredPosts = filteredPosts[:10]
	}
	s.log.Info("successfully received the list of similar posts")
	return filteredPosts, nil
}

func (s *postService) CalcAllPhotosHashes(ctx context.Context) error {
	// Get photos for which to calculate hashes
	photoIDs, err := s.postRepo.FindPhotosWithoutHashes(ctx)
	if err != nil {
		return fmt.Errorf("failed to calculate hashes of photos for all posts: %w", err)
	}
	// Calculate hash for each photo and save it
	for _, id := range photoIDs {
		// Read file
		file, err := os.Open(filepath.Join(s.config.PhotoUploadPath, fmt.Sprintf("%s.jpeg", id.String())))
		if err != nil || file == nil {
			s.log.Error("failed to open file", "error", err.Error())
			return fmt.Errorf("failed to open file: %w", err)
		}
		defer file.Close()
		// Decode as JPEG
		img, err := jpeg.Decode(file)
		if err != nil {
			s.log.Error("failed to decode JPEG image")
			return fmt.Errorf("failed to decode JPEG image: %w", err)
		}
		// Calculate image hash
		hash, err := s.hashCalc.PerceptualHash(img)
		if err != nil {
			s.log.Error("failed to calculate image hash")
			return fmt.Errorf("failed to calculate image hash: %w", err)
		}
		// Save it
		if err := s.postRepo.UpdatePhotoHash(ctx, id, hash); err != nil {
			s.log.Error("failed to save hash of the post photo", "error", err.Error())
			return fmt.Errorf("failed to save hash of the post photo: %w", err)
		}
	}
	return nil
}

func (s *postService) ModeratePost(ctx context.Context, postID uuid.UUID) error {
	// Get post
	post, err := s.postRepo.FindByID(ctx, &postID)
	if err != nil {
		return fmt.Errorf("failed to get post by id (%s): %w", postID.String(), err)
	}
	if post == nil {
		return fmt.Errorf("post is nil")
	}
	// Return nil if moderation status of the post is not pending
	if post.Moderation.Status != model.ModerationStatusPending {
		return nil
	}
	// Get moderation bot user
	bot, err := s.userService.GetPostsModeratorBot(ctx)
	if err != nil {
		s.log.Error("failed to get moderation bot user", "error", err.Error())
		return fmt.Errorf("failed to get moderation bot user: %w", err)
	}
	if bot == nil {
		s.log.Error("bot is nil")
		return fmt.Errorf("bot is nil")
	}
	// Change moderation status to "in progress"
	s.ChangePostModerationStatus(
		ctx,
		postID,
		model.ModerationStatusInProgress,
		&bot.ID,
		nil,
	)
	// Send request and get response
	var description *string
	if strings.TrimSpace(post.Description) != "" {
		description = &post.Description
	}
	var id *uuid.UUID
	if post.HasPhoto {
		id = &post.ID
	}
	res, err := s.sendModerateRequest(post.Name, description, id)
	if err != nil {
		return err
	}
	if res == nil {
		return fmt.Errorf("response from moderation service is nil")
	}
	s.log.Info("successfully received new moderation status of the post", "post_id", postID.String(), "new_moderation_status", (*res).Status)
	// Parse moderation status from the response
	parsedStatus, err := model.ParseModerationStatus((*res).Status)
	if err != nil {
		s.log.Error("failed to parse moderation status", "error", err.Error())
		return fmt.Errorf("failed to parse moderation status: %w", err)
	}
	if parsedStatus == nil {
		s.log.Error("parsed status is nil")
		return fmt.Errorf("parsed status is nil")
	}
	// Change moderation status of the post in the DB
	s.ChangePostModerationStatus(
		ctx,
		postID,
		*parsedStatus,
		&bot.ID,
		func() *string {
			if *parsedStatus == model.ModerationStatusAutoRejected {
				reason := "Содержится неприемлемый контент"
				return &reason
			}
			return nil
		}(),
	)
	s.log.Info("the post was successfully moderated", "post_id", postID.String())
	return nil
}

func (s *postService) GetOldestPendingPost(ctx context.Context) (*PostResponseDTO, error) {
	var post model.Post
	result := s.db.WithContext(ctx).
		Model(&model.Post{}).
		Where("status = ?", string(model.ModerationStatusPending)).
		Order("created_at ASC").
		First(&post)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get the oldest post with pending moderation status")
	}
	return PostToDTO(&post), nil
}

func (s *postService) ModerateAllPosts(ctx context.Context) error {
	// Get IDs of the posts with pending status
	var postModerations []model.PostModeration
	result := s.db.WithContext(ctx).Model(&model.PostModeration{}).Where("status = ?", string(model.ModerationStatusPending)).Find(&postModerations)
	if result.Error != nil {
		return fmt.Errorf("failed to get IDs of the posts with pending moderation status")
	}
	// Moderate posts
	for _, moderation := range postModerations {
		s.log.Debug("trying to moderate post", "post_id", moderation.PostID.String())
		if err := s.ModeratePost(ctx, moderation.PostID); err != nil {
			s.log.Error("failed to moderate post", "post_id", moderation.PostID.String(), "error", err.Error())
			return fmt.Errorf("failed to moderate post with id %s: %w", moderation.PostID.String(), err)
		}
		s.log.Debug("the post was successfully moderated", "post_id", moderation.PostID.String())
	}
	return nil
}

func (s *postService) sendModerateRequest(postTitle string, postDescription *string, postID *uuid.UUID) (*ModerationResult, error) {
	data := url.Values{}

	data.Set("title", postTitle)
	if postDescription != nil {
		data.Set("description", *postDescription)
	}
	if postID != nil {
		data.Set("post_id", (*postID).String())
	}

	res, err := http.Post(
		"http://ml:4746/moderate",
		"application/x-www-form-urlencoded",
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		s.log.Error("failed to send POST request to moderation service", "error", err.Error())
		return nil, fmt.Errorf("failed to send POST request to moderation service: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		s.log.Error("failed to send POST request to moderation service: status code is not 200")
		return nil, fmt.Errorf("failed to send POST request to moderation service: status code is not 200")
	}
	// Parse the result
	var result ModerationResult
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		s.log.Error("failed to send POST request to moderation service: failed to parse result as JSON", "error", err.Error())
		return nil, fmt.Errorf("failed to send POST request to moderation service: failed to parse result as JSON: %w", err)
	}
	return &result, nil
}

type ModerationResult struct {
	Status string `json:"status"`
}

func (s *postService) validateCreatePostDTO(dto *CreatePostDTO) error {
	//TODO
	return nil
}

func (s *postService) validateUpdatePostDTO(dto *UpdatePostDTO) error {
	//TODO
	return nil
}

func ModerationToDTO(moderation *model.PostModeration) *PostModerationResponseDTO {
	return &PostModerationResponseDTO{
		PostID:        moderation.PostID,
		CreatedAt:     moderation.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     moderation.UpdatedAt.Format(time.RFC3339),
		Status:        moderation.Status,
		ModeratorID:   moderation.ModeratorID,
		ModeratorUser: UserToDTO(moderation.ModeratorUser),
		RejectReason:  moderation.RejectReason,
	}
}

func PostToDTO(post *model.Post) *PostResponseDTO {
	return &PostResponseDTO{
		ID:                   post.ID,
		CreatedAt:            post.CreatedAt.Format(time.RFC3339),
		UpdatedAt:            post.UpdatedAt.Format(time.RFC3339),
		Name:                 post.Name,
		Description:          post.Description,
		ThingReturnedToOwner: post.ThingReturnedToOwner,
		HasPhoto:             post.HasPhoto,
		Author:               *UserToDTO(&post.Author),
		Moderation:           *ModerationToDTO(&post.Moderation),
	}
}
