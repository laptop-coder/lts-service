package handler

import (
	"backend/internal/service"
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
		errorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20) // 10 MB
	// TODO: add cleaning of temporary data (ParseMultipartForm, r.MultipartForm, etc)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		errorResponse(w, "failed to parse multipart/formdata form", http.StatusBadRequest)
		return
	}
	// Get name field
	nameFields := r.PostForm["name"]
	if len(nameFields) > 1 {
		errorResponse(w, fmt.Sprintf("failed to parse form: too much name fields (%d)", len(nameFields)), http.StatusBadRequest)
		return
	} else if len(nameFields) == 0 {
		errorResponse(w, "failed to parse form: name field cannot be empty", http.StatusBadRequest)
		return
	}
	name := nameFields[0]
	// Get authorID field and convert the value to UUID
	authorIDFields := r.PostForm["authorID"]
	if len(authorIDFields) > 1 {
		errorResponse(w, fmt.Sprintf("failed to parse form: too much author ID fields (%d)", len(authorIDFields)), http.StatusBadRequest)
		return
	} else if len(nameFields) == 0 {
		errorResponse(w, "failed to parse form: author ID field cannot be empty", http.StatusBadRequest)
		return
	}
	authorID, err := uuid.Parse(authorIDFields[0])
	if err != nil {
		errorResponse(w, "cannot convert author id to uuid", http.StatusBadRequest)
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
		errorResponse(w, "failed to parse form: to much description values", http.StatusBadRequest)
	}
	// Get post photo (optional field)
	formFiles := r.MultipartForm.File["photo"]
	if len(formFiles) > 1 {
		errorResponse(w, "failed to parse form: to much post photo files", http.StatusBadRequest)
		return
	} else if len(formFiles) == 1 {
		dto.Photo = formFiles[0]
	}
	postResponse, err := h.postService.CreatePost(r.Context(), dto)
	if err != nil {
		handleServiceError(w, fmt.Errorf("failed to create the post: %w", err))
		return
	}
	jsonResponse(w, map[string]interface{}{
		"post": postResponse,
	},
		http.StatusCreated,
	)
}
