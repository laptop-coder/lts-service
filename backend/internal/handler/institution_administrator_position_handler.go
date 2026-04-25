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

type InstitutionAdministratorPositionHandler struct {
	institutionAdministratorPositionService service.InstitutionAdministratorPositionService
	log                                     logger.Logger
}

func NewInstitutionAdministratorPositionHandler(institutionAdministratorPositionService service.InstitutionAdministratorPositionService, log logger.Logger) *InstitutionAdministratorPositionHandler {
	return &InstitutionAdministratorPositionHandler{
		institutionAdministratorPositionService: institutionAdministratorPositionService,
		log:                                     log,
	}
}

func (h *InstitutionAdministratorPositionHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helpers.MethodNotAllowedError(h.log, w)
		return
	}
	// Restrictions
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MB
	// Parse form
	if err := r.ParseForm(); err != nil {
		h.log.Error("failed to parse x-www-form-urlencoded form")
		helpers.BadRequestError(h.log, w)
		return
	}
	// Get name
	nameFields := r.PostForm["name"]
	if len(nameFields) == 0 {
		helpers.FieldRequiredError(h.log, w, "name")
		return
	} else if len(nameFields) > 1 {
		h.log.Error(fmt.Sprintf("failed to parse form: too many name values (%d)", len(nameFields)))
		helpers.TooManyFieldsError(h.log, w, "name")
		return
	}
	name := nameFields[0]
	// Assemble DTO
	dto := service.CreateInstitutionAdministratorPositionDTO{
		Name: name,
	}
	// Create institutionAdministratorPosition
	institutionAdministratorPositionResponse, err := h.institutionAdministratorPositionService.CreatePosition(r.Context(), dto)
	if err != nil {
		helpers.HandleServiceError(h.log, w, fmt.Errorf("failed to create the institutionAdministratorPosition: %w", err))
		return
	}
	helpers.JsonResponse(w, map[string]interface{}{
		"institutionAdministratorPosition": institutionAdministratorPositionResponse,
	},
		http.StatusCreated,
	)
}

func (h *InstitutionAdministratorPositionHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodGet {
		helpers.MethodNotAllowedError(h.log, w)
		return
	}
	// Parse query parameters (for filter)
	limitString := r.URL.Query().Get("limit")
	offsetString := r.URL.Query().Get("offset")
	// Pre-assemble filter (fill with default values)
	filter := repository.InstitutionAdministratorPositionFilter{
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
			helpers.BadRequestFieldError(h.log, w, "limit")
			return
		}
	}
	// Parse offset if passed
	if offsetString != "" {
		if offset, err := strconv.Atoi(offsetString); err == nil && offset >= 0 {
			filter.Offset = offset
		} else {
			h.log.Error("invalid offset")
			helpers.BadRequestFieldError(h.log, w, "offset")
			return
		}
	}
	// Get institutionAdministratorPositions
	institutionAdministratorPositions, err := h.institutionAdministratorPositionService.GetPositions(r.Context(), filter)
	if err != nil {
		helpers.HandleServiceError(h.log, w, fmt.Errorf("failed to get institutionAdministratorPositions: %w", err))
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"institutionAdministratorPositions": institutionAdministratorPositions,
	})
}

func (h *InstitutionAdministratorPositionHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodPatch {
		helpers.MethodNotAllowedError(h.log, w)
		return
	}
	// Restrictions
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MB
	// Parse form
	if err := r.ParseForm(); err != nil {
		h.log.Error("failed to parse x-www-form-urlencoded form")
		helpers.BadRequestError(h.log, w)
		return
	}
	// Get and convert institutionAdministratorPosition ID
	institutionAdministratorPositionID64, err := strconv.ParseUint(r.PathValue("id"), 10, 16)
	if err != nil {
		h.log.Error("cannot convert institutionAdministratorPosition ID from string to uint64")
		helpers.BadRequestFieldError(h.log, w, "id")
		return
	}
	institutionAdministratorPositionID := uint16(institutionAdministratorPositionID64)
	// Assemble DTO (all fields are optional)
	dto := service.UpdateInstitutionAdministratorPositionDTO{}
	if nameFields := r.PostForm["name"]; len(nameFields) == 1 {
		dto.Name = &nameFields[0]
	} else if len(nameFields) != 0 {
		h.log.Error(fmt.Sprintf("failed to parse form: too many name values (%d)", len(nameFields)))
		helpers.TooManyFieldsError(h.log, w, "name")
		return
	}
	// Update institutionAdministratorPosition
	institutionAdministratorPositionResponse, err := h.institutionAdministratorPositionService.UpdatePosition(r.Context(), institutionAdministratorPositionID, dto)
	if err != nil {
		helpers.HandleServiceError(h.log, w, fmt.Errorf("failed to update the institutionAdministratorPosition: %w", err))
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"institutionAdministratorPosition": institutionAdministratorPositionResponse,
	})
}

func (h *InstitutionAdministratorPositionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodDelete {
		helpers.MethodNotAllowedError(h.log, w)
		return
	}
	// Get and convert institutionAdministratorPosition ID
	institutionAdministratorPositionID64, err := strconv.ParseUint(r.PathValue("id"), 10, 16)
	if err != nil {
		h.log.Error("cannot convert institutionAdministratorPosition ID from string to uint64")
		helpers.BadRequestFieldError(h.log, w, "id")
		return
	}
	institutionAdministratorPositionID := uint16(institutionAdministratorPositionID64)
	// Delete institutionAdministratorPosition
	if err := h.institutionAdministratorPositionService.DeletePosition(r.Context(), institutionAdministratorPositionID); err != nil {
		helpers.HandleServiceError(h.log, w, fmt.Errorf("failed to delete the institutionAdministratorPosition: %w", err))
		return
	}
	// Return response
	helpers.JsonResponse(w, map[string]interface{}{}, http.StatusNoContent)
}
