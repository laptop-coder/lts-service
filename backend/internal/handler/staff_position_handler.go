package handler

import (
	"backend/internal/repository"
	"backend/internal/service"
	"backend/pkg/helpers"
	"backend/pkg/logger"
	"fmt"
	"net/http"
	"strconv"
)

type StaffPositionHandler struct {
	staffPositionService service.StaffPositionService
	log                  logger.Logger
}

func NewStaffPositionHandler(staffPositionService service.StaffPositionService, log logger.Logger) *StaffPositionHandler {
	return &StaffPositionHandler{
		staffPositionService: staffPositionService,
		log:                  log,
	}
}

func (h *StaffPositionHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helpers.ErrorResponse(h.log, w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Restrictions
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MB
	// Parse form
	if err := r.ParseForm(); err != nil {
		helpers.ErrorResponse(h.log, w, "failed to parse x-www-form-urlencoded form", http.StatusBadRequest)
		return
	}
	// Get name
	nameFields := r.PostForm["name"]
	if len(nameFields) == 0 {
		helpers.ErrorResponse(h.log, w, "failed to parse form: name field is required", http.StatusBadRequest)
		return
	} else if len(nameFields) > 1 {
		helpers.ErrorResponse(h.log, w, fmt.Sprintf("failed to parse form: to much name values (%d)", len(nameFields)), http.StatusBadRequest)
		return
	}
	name := nameFields[0]
	// Assemble DTO
	dto := service.CreateStaffPositionDTO{
		Name: name,
	}
	// Create staffPosition
	staffPositionResponse, err := h.staffPositionService.CreatePosition(r.Context(), dto)
	if err != nil {
		helpers.HandleServiceError(h.log, w, fmt.Errorf("failed to create the staffPosition: %w", err))
		return
	}
	helpers.JsonResponse(w, map[string]interface{}{
		"staffPosition": staffPositionResponse,
	},
		http.StatusCreated,
	)
}

func (h *StaffPositionHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodGet {
		helpers.ErrorResponse(h.log, w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Parse query parameters (for filter)
	limitString := r.URL.Query().Get("limit")
	offsetString := r.URL.Query().Get("offset")
	// Pre-assemble filter (fill with default values)
	filter := repository.StaffPositionFilter{
		Limit:  20,
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
			helpers.ErrorResponse(h.log, w, "invalid limit", http.StatusBadRequest)
			return
		}
	}
	// Parse offset if passed
	if offsetString != "" {
		if offset, err := strconv.Atoi(offsetString); err == nil && offset >= 0 {
			filter.Offset = offset
		} else {
			h.log.Error("invalid offset")
			helpers.ErrorResponse(h.log, w, "invalid offset", http.StatusBadRequest)
			return
		}
	}
	// Get staffPositions
	staffPositions, err := h.staffPositionService.GetPositions(r.Context(), filter)
	if err != nil {
		helpers.HandleServiceError(h.log, w, fmt.Errorf("failed to get staffPositions: %w", err))
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"staffPositions": staffPositions,
	})
}

func (h *StaffPositionHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodPatch {
		helpers.ErrorResponse(h.log, w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Restrictions
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MB
	// Parse form
	if err := r.ParseForm(); err != nil {
		helpers.ErrorResponse(h.log, w, "failed to parse x-www-form-urlencoded form", http.StatusBadRequest)
		return
	}
	// Get and convert staffPosition ID
	staffPositionID64, err := strconv.ParseUint(r.PathValue("id"), 10, 8)
	if err != nil {
		helpers.ErrorResponse(h.log, w, "cannot convert staffPosition ID from string to uint64", http.StatusBadRequest)
		return
	}
	staffPositionID := uint8(staffPositionID64)
	// Assemble DTO (all fields are optional)
	dto := service.UpdateStaffPositionDTO{}
	if nameFields := r.PostForm["name"]; len(nameFields) == 1 {
		dto.Name = &nameFields[0]
	} else if len(nameFields) != 0 {
		helpers.ErrorResponse(h.log, w, fmt.Sprintf("failed to parse form: to much name values (%d)", len(nameFields)), http.StatusBadRequest)
		return
	}
	// Update staffPosition
	staffPositionResponse, err := h.staffPositionService.UpdatePosition(r.Context(), staffPositionID, dto)
	if err != nil {
		helpers.HandleServiceError(h.log, w, fmt.Errorf("failed to update the staffPosition: %w", err))
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"staffPosition": staffPositionResponse,
	})
}

func (h *StaffPositionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodDelete {
		helpers.ErrorResponse(h.log, w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Get and convert staffPosition ID
	staffPositionID64, err := strconv.ParseUint(r.PathValue("id"), 10, 8)
	if err != nil {
		helpers.ErrorResponse(h.log, w, "cannot convert staffPosition ID from string to uint64", http.StatusInternalServerError)
		return
	}
	staffPositionID := uint8(staffPositionID64)
	// Delete staffPosition
	if err := h.staffPositionService.DeletePosition(r.Context(), staffPositionID); err != nil {
		helpers.HandleServiceError(h.log, w, fmt.Errorf("failed to delete the staffPosition: %w", err))
		return
	}
	// Return response
	helpers.JsonResponse(w, map[string]interface{}{}, http.StatusNoContent)
}
