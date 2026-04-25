package handler

import (
	"backend/internal/repository"
	"backend/internal/service"
	"backend/pkg/helpers"
	"backend/pkg/logger"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"strconv"
)

type RoomHandler struct {
	roomService service.RoomService
	log         logger.Logger
}

func NewRoomHandler(roomService service.RoomService, log logger.Logger) *RoomHandler {
	return &RoomHandler{
		roomService: roomService,
		log:         log,
	}
}

func (h *RoomHandler) Create(w http.ResponseWriter, r *http.Request) {
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
		h.log.Error("failed to parse form: name field is required")
		helpers.FieldRequiredError(h.log, w, "name")
		return
	} else if len(nameFields) > 1 {
		h.log.Error(fmt.Sprintf("failed to parse form: too many name values (%d)", len(nameFields)))
		helpers.TooManyFieldsError(h.log, w, "name")
		return
	}
	name := nameFields[0]
	// Get and convert teacher ID
	teacherIDFields := r.PostForm["teacherId"]
	var teacherID *uuid.UUID
	if len(teacherIDFields) > 1 {
		h.log.Error(fmt.Sprintf("failed to parse form: too many teacherID values (%d)", len(teacherIDFields)))
		helpers.TooManyFieldsError(h.log, w, "teacherId")
		return
	} else if len(teacherIDFields) == 1 {
		teacherIDUUID, err := uuid.Parse(teacherIDFields[0])
		if err != nil {
			h.log.Error("cannot convert teacher id to uuid")
			helpers.BadRequestFieldError(h.log, w, "teacherId")
			return
		}
		teacherID = &teacherIDUUID
	}
	// Assemble DTO
	dto := service.CreateRoomDTO{
		Name:      name,
		TeacherID: teacherID,
	}
	// Create room
	roomResponse, err := h.roomService.CreateRoom(r.Context(), dto)
	if err != nil {
		helpers.HandleServiceError(h.log, w, fmt.Errorf("failed to create the room: %w", err))
		return
	}
	helpers.JsonResponse(w, map[string]interface{}{
		"room": roomResponse,
	},
		http.StatusCreated,
	)
}

func (h *RoomHandler) GetRooms(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodGet {
		helpers.MethodNotAllowedError(h.log, w)
		return
	}
	// Parse query parameters (for filter)
	limitString := r.URL.Query().Get("limit")
	offsetString := r.URL.Query().Get("offset")
	// Pre-assemble filter (fill with default values)
	filter := repository.RoomFilter{
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
	// Get rooms
	rooms, err := h.roomService.GetRooms(r.Context(), filter)
	if err != nil {
		helpers.HandleServiceError(h.log, w, fmt.Errorf("failed to get rooms: %w", err))
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"rooms": rooms,
	})
}

func (h *RoomHandler) Update(w http.ResponseWriter, r *http.Request) {
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
	// Get and convert room ID
	roomID64, err := strconv.ParseUint(r.PathValue("id"), 10, 8)
	if err != nil {
		h.log.Error("cannot convert room ID from string to uint64")
		helpers.BadRequestFieldError(h.log, w, "id")
		return
	}
	roomID := uint16(roomID64)
	// Assemble DTO (all fields are optional)
	dto := service.UpdateRoomDTO{}
	if nameFields := r.PostForm["name"]; len(nameFields) == 1 {
		dto.Name = &nameFields[0]
	} else if len(nameFields) != 0 {
		h.log.Error(fmt.Sprintf("failed to parse form: too many name values (%d)", len(nameFields)))
		helpers.TooManyFieldsError(h.log, w, "name")
		return
	}
	if teacherIDFields := r.PostForm["teacherId"]; len(teacherIDFields) == 1 {
		// Convert teacher ID to UUID
		teacherID, err := uuid.Parse(teacherIDFields[0])
		if err != nil {
			h.log.Error("cannot convert teacher id to uuid")
			helpers.BadRequestFieldError(h.log, w, "teacherId")
			return
		}
		dto.TeacherID = &teacherID
	} else if len(teacherIDFields) != 0 {
		h.log.Error(fmt.Sprintf("failed to parse form: too many teacher ID values (%d)", len(teacherIDFields)))
		helpers.TooManyFieldsError(h.log, w, "teacherId")
		return
	}
	// Update room
	roomResponse, err := h.roomService.UpdateRoom(r.Context(), roomID, dto)
	if err != nil {
		helpers.HandleServiceError(h.log, w, fmt.Errorf("failed to update the room: %w", err))
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"room": roomResponse,
	})
}

func (h *RoomHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodDelete {
		helpers.MethodNotAllowedError(h.log, w)
		return
	}
	// Get and convert room ID
	roomID64, err := strconv.ParseUint(r.PathValue("id"), 10, 8)
	if err != nil {
		h.log.Error("cannot convert room ID from string to uint64")
		helpers.BadRequestFieldError(h.log, w, "id")
		return
	}
	roomID := uint16(roomID64)
	// Delete room
	if err := h.roomService.DeleteRoom(r.Context(), roomID); err != nil {
		helpers.HandleServiceError(h.log, w, fmt.Errorf("failed to delete the room: %w", err))
		return
	}
	// Return response
	helpers.JsonResponse(w, map[string]interface{}{}, http.StatusNoContent)
}
