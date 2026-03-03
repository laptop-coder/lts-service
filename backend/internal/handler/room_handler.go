package handler

import (
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
	// Get and convert teacher ID
	teacherIDFields := r.PostForm["teacherID"]
	var teacherID *uuid.UUID
	if len(teacherIDFields) > 1 {
		helpers.ErrorResponse(w, fmt.Sprintf("failed to parse form: to much teacherID values (%d)", len(teacherIDFields)), http.StatusBadRequest)
	} else if len(teacherIDFields) == 1 {
		teacherIDUUID, err := uuid.Parse(teacherIDFields[0])
		if err != nil {
			helpers.ErrorResponse(w, "cannot convert teacher id to uuid", http.StatusBadRequest)
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
		helpers.HandleServiceError(w, fmt.Errorf("failed to create the room: %w", err))
		return
	}
	helpers.JsonResponse(w, map[string]interface{}{
		"room": roomResponse,
	},
		http.StatusCreated,
	)
}

func (h *RoomHandler) Update(w http.ResponseWriter, r *http.Request) {
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
	// Get and convert room ID
	roomID64, err := strconv.ParseUint(r.PathValue("id"), 10, 8)
	if err != nil {
		helpers.ErrorResponse(w, "cannot convert room ID from string to uint64", http.StatusInternalServerError)
		return
	}
	roomID := uint8(roomID64)
	// Assemble DTO (all fields are optional)
	dto := service.UpdateRoomDTO{}
	if nameFields := r.PostForm["name"]; len(nameFields) == 1 {
		dto.Name = &nameFields[0]
	} else if len(nameFields) != 0 {
		helpers.ErrorResponse(w, fmt.Sprintf("failed to parse form: to much name values (%d)", len(nameFields)), http.StatusBadRequest)
	}
	if teacherIDFields := r.PostForm["teacherID"]; len(teacherIDFields) == 1 {
		// Convert teacher ID to UUID
		teacherID, err := uuid.Parse(teacherIDFields[0])
		if err != nil {
			helpers.ErrorResponse(w, "cannot convert teacher id to uuid", http.StatusBadRequest)
		}
		dto.TeacherID = &teacherID
	} else if len(teacherIDFields) != 0 {
		helpers.ErrorResponse(w, fmt.Sprintf("failed to parse form: to much teacher ID values (%d)", len(teacherIDFields)), http.StatusBadRequest)
	}
	// Update room
	roomResponse, err := h.roomService.UpdateRoom(r.Context(), roomID, dto)
	if err != nil {
		helpers.HandleServiceError(w, fmt.Errorf("failed to update the room: %w", err))
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
		helpers.ErrorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Get and convert room ID
	roomID64, err := strconv.ParseUint(r.PathValue("id"), 10, 8)
	if err != nil {
		helpers.ErrorResponse(w, "cannot convert room ID from string to uint64", http.StatusInternalServerError)
		return
	}
	roomID := uint8(roomID64)
	// Delete room
	if err := h.roomService.DeleteRoom(r.Context(), roomID); err != nil {
		helpers.HandleServiceError(w, fmt.Errorf("failed to delete the room: %w", err))
		return
	}
	// Return response
	helpers.JsonResponse(w, map[string]interface{}{}, http.StatusNoContent)
}
