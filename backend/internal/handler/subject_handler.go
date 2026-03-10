package handler

import (
	"backend/internal/service"
	"backend/internal/repository"
	"backend/pkg/helpers"
	"backend/pkg/logger"
	"fmt"
	"net/http"
	"strconv"
)

type SubjectHandler struct {
	subjectService service.SubjectService
	log            logger.Logger
}

func NewSubjectHandler(subjectService service.SubjectService, log logger.Logger) *SubjectHandler {
	return &SubjectHandler{
		subjectService: subjectService,
		log:            log,
	}
}

func (h *SubjectHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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
	// Get name
	nameFields := r.PostForm["name"]
	if len(nameFields) == 0 {
		helpers.ErrorResponse(w, "failed to parse form: name field is required", http.StatusBadRequest)
	} else if len(nameFields) > 1 {
		helpers.ErrorResponse(w, fmt.Sprintf("failed to parse form: to much name values (%d)", len(nameFields)), http.StatusBadRequest)
	}
	name := nameFields[0]
	// Assemble DTO
	dto := service.CreateSubjectDTO{
		Name: name,
	}
	// Create subject
	subjectResponse, err := h.subjectService.CreateSubject(r.Context(), dto)
	if err != nil {
		helpers.HandleServiceError(w, fmt.Errorf("failed to create the subject: %w", err))
		return
	}
	helpers.JsonResponse(w, map[string]interface{}{
		"subject": subjectResponse,
	},
		http.StatusCreated,
	)
}


func (h *SubjectHandler) GetSubjects(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodGet {
		helpers.ErrorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Parse query parameters (for filter)
	limitString := r.URL.Query().Get("limit")
	offsetString := r.URL.Query().Get("offset")
    // Pre-assemble filter (fill with default values)
	filter := repository.SubjectFilter {
		Limit: 20,
		Offset: 0,
	}
	// Parse limit if passed
	if limitString != "" {
		if limit, err := strconv.Atoi(limitString); err == nil && limit > 0 {
			if limit > 100 {
				limit = 100 // max value
			}
			filter.Limit = limit
		} else {
			h.log.Error("invalid limit")
			helpers.ErrorResponse(w, "invalid limit", http.StatusBadRequest)
			return
		}
	}
	// Parse offset if passed
	if offsetString != "" {
		if offset, err := strconv.Atoi(offsetString); err == nil && offset >= 0 {
			filter.Offset = offset
		} else {
			h.log.Error("invalid offset")
			helpers.ErrorResponse(w, "invalid offset", http.StatusBadRequest)
			return
		}
	}
	// Get subjects
	subjects, err := h.subjectService.GetSubjects(r.Context(), filter)
	if err != nil {
		helpers.HandleServiceError(w, fmt.Errorf("failed to get subjects: %w", err))
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"subjects": subjects,
	})
}


func (h *SubjectHandler) Update(w http.ResponseWriter, r *http.Request) {
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
	// Get and convert subject ID
	subjectID64, err := strconv.ParseUint(r.PathValue("id"), 10, 8)
	if err != nil {
		helpers.ErrorResponse(w, "cannot convert subject ID from string to uint64", http.StatusInternalServerError)
		return
	}
	subjectID := uint8(subjectID64)
	// Assemble DTO (all fields are optional)
	dto := service.UpdateSubjectDTO{}
	if nameFields := r.PostForm["name"]; len(nameFields) == 1 {
		dto.Name = &nameFields[0]
	} else if len(nameFields) != 0 {
		helpers.ErrorResponse(w, fmt.Sprintf("failed to parse form: to much name values (%d)", len(nameFields)), http.StatusBadRequest)
	}
	// Update subject
	subjectResponse, err := h.subjectService.UpdateSubject(r.Context(), subjectID, dto)
	if err != nil {
		helpers.HandleServiceError(w, fmt.Errorf("failed to update the subject: %w", err))
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"subject": subjectResponse,
	})
}

func (h *SubjectHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodDelete {
		helpers.ErrorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Get and convert subject ID
	subjectID64, err := strconv.ParseUint(r.PathValue("id"), 10, 8)
	if err != nil {
		helpers.ErrorResponse(w, "cannot convert subject ID from string to uint64", http.StatusInternalServerError)
		return
	}
	subjectID := uint8(subjectID64)
	// Delete subject
	if err := h.subjectService.DeleteSubject(r.Context(), subjectID); err != nil {
		helpers.HandleServiceError(w, fmt.Errorf("failed to delete the subject: %w", err))
		return
	}
	// Return response
	helpers.JsonResponse(w, map[string]interface{}{}, http.StatusNoContent)
}
