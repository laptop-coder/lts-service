package handler

import (
	"backend/internal/service"
	"backend/pkg/helpers"
	"backend/pkg/logger"
	"fmt"
	"github.com/google/uuid"
	"net/http"
)

type PostHandler struct {
	postService service.PostService
	log         logger.Logger
}

func NewPostHandler(postService service.PostService, log logger.Logger) *PostHandler {
	return &PostHandler{
		postService: postService,
		log:         log,
	}
}

func (h *PostHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helpers.ErrorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20) // 10 MB
	// TODO: add cleaning of temporary data (ParseMultipartForm, r.MultipartForm, etc)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		helpers.ErrorResponse(w, "failed to parse multipart/formdata form", http.StatusBadRequest)
		return
	}
	// Get name field
	nameFields := r.PostForm["name"]
	if len(nameFields) > 1 {
		helpers.ErrorResponse(w, fmt.Sprintf("failed to parse form: too much name fields (%d)", len(nameFields)), http.StatusBadRequest)
		return
	} else if len(nameFields) == 0 {
		helpers.ErrorResponse(w, "failed to parse form: name field cannot be empty", http.StatusBadRequest)
		return
	}
	name := nameFields[0]
	// Get authorID field and convert the value to UUID
	authorIDFields := r.PostForm["authorID"]
	if len(authorIDFields) > 1 {
		helpers.ErrorResponse(w, fmt.Sprintf("failed to parse form: too much author ID fields (%d)", len(authorIDFields)), http.StatusBadRequest)
		return
	} else if len(nameFields) == 0 {
		helpers.ErrorResponse(w, "failed to parse form: author ID field cannot be empty", http.StatusBadRequest)
		return
	}
	authorID, err := uuid.Parse(authorIDFields[0])
	if err != nil {
		helpers.ErrorResponse(w, "cannot convert author id to uuid", http.StatusBadRequest)
	}
	// Pre-assemble DTO
	dto := service.CreatePostDTO{
		Name:     name,
		AuthorID: authorID,
	}
	// Get description (optional field)
	if descriptionFields := r.PostForm["description"]; len(descriptionFields) == 1 {
		dto.Description = descriptionFields[0]
	} else if len(descriptionFields) != 0 {
		helpers.ErrorResponse(w, "failed to parse form: to much description values", http.StatusBadRequest)
	}
	// Get post photo (optional field)
	formFiles := r.MultipartForm.File["photo"]
	if len(formFiles) > 1 {
		helpers.ErrorResponse(w, "failed to parse form: to much post photo files", http.StatusBadRequest)
		return
	} else if len(formFiles) == 1 {
		dto.Photo = formFiles[0]
	}
	postResponse, err := h.postService.CreatePost(r.Context(), dto)
	if err != nil {
		helpers.HandleServiceError(w, fmt.Errorf("failed to create the post: %w", err))
		return
	}
	helpers.JsonResponse(w, map[string]interface{}{
		"post": postResponse,
	},
		http.StatusCreated,
	)
}

func (h *PostHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodPatch {
		helpers.ErrorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Restrictions
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MB
	// Parse form
	if err := r.ParseForm(); err != nil {
		helpers.ErrorResponse(w, "failed to parse x-www-form-urlencoded form", http.StatusBadRequest)
		return
	}
	// Get and convert post ID
	postID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		helpers.ErrorResponse(w, "cannot convert post id to uuid", http.StatusBadRequest)
	}
	// DTO (all fields are optional)
	dto := service.UpdatePostDTO{}
	if nameFields := r.PostForm["name"]; len(nameFields) == 1 {
		dto.Name = &nameFields[0]
	} else if len(nameFields) != 0 {
		helpers.ErrorResponse(w, "failed to parse form: to much name values", http.StatusBadRequest)
	}
	if descriptionFields := r.PostForm["description"]; len(descriptionFields) == 1 {
		dto.Description = &descriptionFields[0]
	} else if len(descriptionFields) != 0 {
		helpers.ErrorResponse(w, "failed to parse form: to much description values", http.StatusBadRequest)
	}
	// Update post
	postResponse, err := h.postService.UpdatePost(r.Context(), postID, dto)
	if err != nil {
		helpers.HandleServiceError(w, fmt.Errorf("failed to update the post: %w", err))
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"post": postResponse,
	})
}

func (h *PostHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodDelete {
		helpers.ErrorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Get and convert post ID
	postID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		helpers.ErrorResponse(w, "cannot convert post id to uuid", http.StatusBadRequest)
	}
	// Delete post
	if err := h.postService.DeletePost(r.Context(), postID); err != nil {
		helpers.HandleServiceError(w, fmt.Errorf("failed to delete the post: %w", err))
		return
	}
	// Return response
	helpers.JsonResponse(w, map[string]interface{}{}, http.StatusNoContent)
}

func (h *PostHandler) RemovePhoto(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodDelete {
		helpers.ErrorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Get and convert post ID
	postID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		helpers.ErrorResponse(w, "cannot convert post id to uuid", http.StatusBadRequest)
	}
	// Remove post photo file
	if err := h.postService.RemovePhoto(r.Context(), postID); err != nil {
		helpers.HandleServiceError(w, fmt.Errorf("failed to remove post photo file: %w", err))
		return
	}
	// Return response
	helpers.JsonResponse(w, map[string]interface{}{}, http.StatusNoContent)
}
