package handler

import (
	"backend/internal/repository"
	"backend/internal/service"
	"backend/pkg/helpers"
	"backend/pkg/logger"
	"backend/pkg/middleware"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"strconv"
)

type UserHandler struct {
	userService service.UserService
	log         logger.Logger
}

func NewUserHandler(userService service.UserService, log logger.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		log:         log,
	}
}

func (h *UserHandler) UpdateOwnProfile(w http.ResponseWriter, r *http.Request) {
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
	// Get and convert user ID
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		helpers.ErrorResponse(w, "cannot convert user id to uuid", http.StatusUnauthorized)
		return
	}
	// DTO (all fields are optional)
	dto := service.UpdateUserDTO{}
	if firstNameFields := r.PostForm["firstName"]; len(firstNameFields) == 1 {
		dto.FirstName = &firstNameFields[0]
	} else if len(firstNameFields) != 0 {
		helpers.ErrorResponse(w, "failed to parse form: to much firstName values", http.StatusBadRequest)
		return
	}
	if middleNameFields := r.PostForm["middleName"]; len(middleNameFields) == 1 {
		dto.MiddleName = &middleNameFields[0]
	} else if len(middleNameFields) != 0 {
		helpers.ErrorResponse(w, "failed to parse form: to much middleName values", http.StatusBadRequest)
		return
	}
	if lastNameFields := r.PostForm["lastName"]; len(lastNameFields) == 1 {
		dto.LastName = &lastNameFields[0]
	} else if len(lastNameFields) != 0 {
		helpers.ErrorResponse(w, "failed to parse form: to much lastName values", http.StatusBadRequest)
		return
	}
	// Update user
	userResponse, err := h.userService.UpdateUser(r.Context(), userID, dto)
	if err != nil {
		helpers.HandleServiceError(w, fmt.Errorf("failed to update the user profile: %w", err))
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"user": userResponse,
	})
}

func (h *UserHandler) RemoveOwnAvatar(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodDelete {
		helpers.ErrorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Get and convert user ID
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		helpers.ErrorResponse(w, "cannot convert user id to uuid", http.StatusUnauthorized)
		return
	}
	// Remove user avatar file
	if err := h.userService.RemoveAvatar(r.Context(), userID); err != nil {
		helpers.HandleServiceError(w, fmt.Errorf("failed to remove user avatar file: %w", err))
		return
	}
	// Return response
	helpers.JsonResponse(w, map[string]interface{}{}, http.StatusNoContent)
}

func (h *UserHandler) UpdateOwnAvatar(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodPut {
		helpers.ErrorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Restrictions
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20) // 10 MB
	// Parse form
	// TODO: add cleaning of temporary data (ParseMultipartForm, r.MultipartForm, etc)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		helpers.ErrorResponse(w, "failed to parse multipart/formdata form", http.StatusBadRequest)
		return
	}
	// Get and convert user ID
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		helpers.ErrorResponse(w, "cannot convert user id to uuid", http.StatusUnauthorized)
		return
	}
	// Get avatar file from the request
	formFiles := r.MultipartForm.File["avatar"]
	if len(formFiles) > 1 {
		helpers.ErrorResponse(w, "failed to parse form: to much avatar files", http.StatusBadRequest)
		return
	} else if len(formFiles) == 0 {
		helpers.ErrorResponse(w, "failed to parse form: avatar cannot be empty", http.StatusBadRequest)
		return
	}
	// Update avatar file
	if err := h.userService.UpdateAvatar(r.Context(), userID, formFiles[0]); err != nil {
		helpers.HandleServiceError(w, fmt.Errorf("failed to update the avatar: %w", err))
		return
	}
	// Return response
	helpers.JsonResponse(w, map[string]interface{}{}, http.StatusNoContent)
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodGet {
		helpers.ErrorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Get and convert user ID
	userID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		helpers.ErrorResponse(w, "cannot convert user id to uuid", http.StatusBadRequest)
		return
	}
	// Get user
	response, err := h.userService.GetUserByID(r.Context(), userID)
	if err != nil {
		helpers.HandleServiceError(w, fmt.Errorf("failed to get user by id: %w", err))
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"user": response,
	})
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodGet {
		helpers.ErrorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Parse query parameters (for filter)
	roleIDString := r.URL.Query().Get("roleId")
	limitString := r.URL.Query().Get("limit")
	offsetString := r.URL.Query().Get("offset")
	// Pre-assemble filter (fill with default values)
	filter := repository.UserFilter{
		Limit:  20,
		Offset: 0,
	}
	// Parse role ID if passed
	if roleIDString != "" {
		// convert to uint64
		roleID64, err := strconv.ParseUint(roleIDString, 10, 8)
		if err != nil {
			h.log.Error("cannot convert role ID from string to uint64")
			helpers.ErrorResponse(w, "cannot convert role ID from string to uint64", http.StatusInternalServerError)
			return
		}
		// and to uint8
		roleID := uint8(roleID64)
		// Add to filter
		filter.RoleID = &roleID
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
	// Get users
	users, err := h.userService.GetUsers(r.Context(), filter)
	if err != nil {
		helpers.HandleServiceError(w, fmt.Errorf("failed to get users: %w", err))
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"users": users,
	})
}

func (h *UserHandler) GetOwnUser(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodGet {
		helpers.ErrorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Get and convert user ID
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		helpers.ErrorResponse(w, "cannot convert own user id to uuid", http.StatusUnauthorized)
		return
	}
	// Get user
	response, err := h.userService.GetUserByID(r.Context(), userID)
	if err != nil {
		helpers.HandleServiceError(w, fmt.Errorf("failed to get own user by id: %w", err))
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"user": response,
	})
}

func (h *UserHandler) GetRoles(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodGet {
		helpers.ErrorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Get and convert user ID
	userID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		helpers.ErrorResponse(w, "cannot convert user id to uuid", http.StatusBadRequest)
		return
	}
	// Get roles
	roles, err := h.userService.GetUserRoles(r.Context(), userID)
	if err != nil {
		helpers.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"userID": userID,
		"roles":  roles,
	})
}

func (h *UserHandler) GetOwnRoles(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodGet {
		helpers.ErrorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Get and convert user ID
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		helpers.ErrorResponse(w, "cannot convert user id to uuid", http.StatusUnauthorized)
		return
	}
	// Get roles
	roles, err := h.userService.GetUserRoles(r.Context(), userID)
	if err != nil {
		helpers.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"userID": userID,
		"roles":  roles,
	})
}

func (h *UserHandler) AssignRoles(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// Check method
	if r.Method != http.MethodPut {
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
	// Get and convert user ID
	userID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		helpers.ErrorResponse(w, "cannot convert user id to uuid", http.StatusBadRequest)
		return
	}
	// Get role IDs
	roleIDsFields := r.PostForm["roleId"]
	if len(roleIDsFields) == 0 {
		h.log.Error("the list of roles cannot be empty")
		helpers.ErrorResponse(w, "the list of roles cannot be empty", http.StatusBadRequest)
		return
	}
	roleIDs := make([]uint8, len(roleIDsFields))
	for i, s := range roleIDsFields {
		val, err := strconv.ParseUint(s, 10, 8)
		if err != nil {
			h.log.Error("cannot convert IDs of roles from string to uint64")
			helpers.ErrorResponse(w, "cannot convert IDs of roles from string to uint64", http.StatusInternalServerError)
			return
		}
		roleIDs[i] = uint8(val)
	}
	// Get special fields (for user-extension tables)
	userRolesDTO := service.UserRolesDTO{}
	// TeacherClassroomID (special)
	if teacherClassroomIDFields := r.PostForm["teacherClassroomId"]; len(teacherClassroomIDFields) == 1 {
		// Convert to uint8
		teacherClassroomID64, err := strconv.ParseUint(teacherClassroomIDFields[0], 10, 8)
		if err != nil {
			helpers.ErrorResponse(w, "cannot convert teacher classroom ID from string to uint64", http.StatusInternalServerError)
			return
		}
		teacherClassroomID := uint8(teacherClassroomID64)
		userRolesDTO.TeacherClassroomID = &teacherClassroomID
	} else if len(teacherClassroomIDFields) != 0 {
		helpers.ErrorResponse(w, "failed to parse form: to much teacher classroom id values", http.StatusBadRequest)
		return
	}
	// TeacherSubjectIDs (special)
	teacherSubjectIDsFields := r.PostForm["teacherSubjectId"]
	var teacherSubjectIDs = make([]uint8, len(teacherSubjectIDsFields))
	for i, subjectIDString := range teacherSubjectIDsFields {
		subjectID64, err := strconv.ParseUint(subjectIDString, 10, 8)
		if err != nil {
			helpers.ErrorResponse(w, "cannot convert teacher subject ID from string to uint64", http.StatusInternalServerError)
			return
		}
		subjectID8 := uint8(subjectID64)
		teacherSubjectIDs[i] = subjectID8
	}
	userRolesDTO.TeacherSubjectIDs = teacherSubjectIDs
	// StudentGroupID (special)
	if studentGroupIDFields := r.PostForm["studentGroupId"]; len(studentGroupIDFields) == 1 {
		// Convert to uint16
		studentGroupID64, err := strconv.ParseUint(studentGroupIDFields[0], 10, 16)
		if err != nil {
			helpers.ErrorResponse(w, "cannot convert student group ID from string to uint64", http.StatusInternalServerError)
			return
		}
		studentGroupID := uint16(studentGroupID64)
		userRolesDTO.StudentGroupID = &studentGroupID
	} else if len(studentGroupIDFields) != 0 {
		helpers.ErrorResponse(w, "failed to parse form: to much student group id values", http.StatusBadRequest)
		return
	}
	// StaffPositionID (special)
	if staffPositionIDFields := r.PostForm["staffPositionId"]; len(staffPositionIDFields) == 1 {
		// Convert to uint8
		staffPositionID64, err := strconv.ParseUint(staffPositionIDFields[0], 10, 8)
		if err != nil {
			helpers.ErrorResponse(w, "cannot convert staff position ID from string to uint64", http.StatusInternalServerError)
			return
		}
		staffPositionID := uint8(staffPositionID64)
		userRolesDTO.StaffPositionID = &staffPositionID
	} else if len(staffPositionIDFields) != 0 {
		helpers.ErrorResponse(w, "failed to parse form: to much staff position id values", http.StatusBadRequest)
		return
	}
	// InstitutionAdministratorPositionID (special)
	if institutionAdministratorPositionIDFields := r.PostForm["institutionAdministratorPositionId"]; len(institutionAdministratorPositionIDFields) == 1 {
		// Convert to uint8
		institutionAdministratorPositionID64, err := strconv.ParseUint(institutionAdministratorPositionIDFields[0], 10, 8)
		if err != nil {
			helpers.ErrorResponse(w, "cannot convert institution administrator position ID from string to uint64", http.StatusInternalServerError)
			return
		}
		institutionAdministratorPositionID := uint8(institutionAdministratorPositionID64)
		userRolesDTO.InstitutionAdministratorPositionID = &institutionAdministratorPositionID
	} else if len(institutionAdministratorPositionIDFields) != 0 {
		helpers.ErrorResponse(w, "failed to parse form: to much institution administrator position id values", http.StatusBadRequest)
		return
	}
	// ParentStudentIDs (special)
	parentStudentIDsFields := r.PostForm["parentStudentId"]
	var parentStudentIDs = make([]uuid.UUID, len(parentStudentIDsFields))
	for i, parentStudentIDString := range parentStudentIDsFields {
		parentStudentID, err := uuid.Parse(parentStudentIDString)
		if err != nil {
			helpers.ErrorResponse(w, "cannot convert student id to uuid", http.StatusBadRequest)
			return
		}
		parentStudentIDs[i] = parentStudentID
	}
	userRolesDTO.ParentStudentIDs = parentStudentIDs
	// Replace old roles with new ones
	if err := h.userService.AssignRolesToUser(ctx, userID, userRolesDTO, roleIDs); err != nil {
		helpers.HandleServiceError(w, err)
		return
	}
	// Get updated roles
	roles, err := h.userService.GetUserRoles(ctx, userID)
	if err != nil {
		helpers.HandleServiceError(w, err)
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"userID":  userID,
		"roles":   roles,
		"message": "roles updated successfully",
	})
}

func (h *UserHandler) AddRoles(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// Check method
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
	// Get and convert user ID
	userID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		helpers.ErrorResponse(w, "cannot convert user id to uuid", http.StatusBadRequest)
		return
	}
	// Get role IDs
	roleIDsFields := r.PostForm["roleId"]
	if len(roleIDsFields) == 0 {
		h.log.Error("the list of roles cannot be empty")
		helpers.ErrorResponse(w, "the list of roles cannot be empty", http.StatusBadRequest)
		return
	}
	roleIDs := make([]uint8, len(roleIDsFields))
	for i, s := range roleIDsFields {
		val, err := strconv.ParseUint(s, 10, 8)
		if err != nil {
			h.log.Error("cannot convert IDs of roles from string to uint64")
			helpers.ErrorResponse(w, "cannot convert IDs of roles from string to uint64", http.StatusInternalServerError)
			return
		}
		roleIDs[i] = uint8(val)
	}
	// Get special fields (for user-extension tables)
	userRolesDTO := service.UserRolesDTO{}
	// TeacherClassroomID (special)
	if teacherClassroomIDFields := r.PostForm["teacherClassroomId"]; len(teacherClassroomIDFields) == 1 {
		// Convert to uint8
		teacherClassroomID64, err := strconv.ParseUint(teacherClassroomIDFields[0], 10, 8)
		if err != nil {
			helpers.ErrorResponse(w, "cannot convert teacher classroom ID from string to uint64", http.StatusInternalServerError)
			return
		}
		teacherClassroomID := uint8(teacherClassroomID64)
		userRolesDTO.TeacherClassroomID = &teacherClassroomID
	} else if len(teacherClassroomIDFields) != 0 {
		helpers.ErrorResponse(w, "failed to parse form: to much teacher classroom id values", http.StatusBadRequest)
		return
	}
	// TeacherSubjectIDs (special)
	teacherSubjectIDsFields := r.PostForm["teacherSubjectId"]
	var teacherSubjectIDs = make([]uint8, len(teacherSubjectIDsFields))
	for i, subjectIDString := range teacherSubjectIDsFields {
		subjectID64, err := strconv.ParseUint(subjectIDString, 10, 8)
		if err != nil {
			helpers.ErrorResponse(w, "cannot convert teacher subject ID from string to uint64", http.StatusInternalServerError)
			return
		}
		subjectID8 := uint8(subjectID64)
		teacherSubjectIDs[i] = subjectID8
	}
	userRolesDTO.TeacherSubjectIDs = teacherSubjectIDs
	// StudentGroupID (special)
	if studentGroupIDFields := r.PostForm["studentGroupId"]; len(studentGroupIDFields) == 1 {
		// Convert to uint16
		studentGroupID64, err := strconv.ParseUint(studentGroupIDFields[0], 10, 16)
		if err != nil {
			helpers.ErrorResponse(w, "cannot convert student group ID from string to uint64", http.StatusInternalServerError)
			return
		}
		studentGroupID := uint16(studentGroupID64)
		userRolesDTO.StudentGroupID = &studentGroupID
	} else if len(studentGroupIDFields) != 0 {
		helpers.ErrorResponse(w, "failed to parse form: to much student group id values", http.StatusBadRequest)
		return
	}
	// StaffPositionID (special)
	if staffPositionIDFields := r.PostForm["staffPositionId"]; len(staffPositionIDFields) == 1 {
		// Convert to uint8
		staffPositionID64, err := strconv.ParseUint(staffPositionIDFields[0], 10, 8)
		if err != nil {
			helpers.ErrorResponse(w, "cannot convert staff position ID from string to uint64", http.StatusInternalServerError)
			return
		}
		staffPositionID := uint8(staffPositionID64)
		userRolesDTO.StaffPositionID = &staffPositionID
	} else if len(staffPositionIDFields) != 0 {
		helpers.ErrorResponse(w, "failed to parse form: to much staff position id values", http.StatusBadRequest)
		return
	}
	// InstitutionAdministratorPositionID (special)
	if institutionAdministratorPositionIDFields := r.PostForm["institutionAdministratorPositionId"]; len(institutionAdministratorPositionIDFields) == 1 {
		// Convert to uint8
		institutionAdministratorPositionID64, err := strconv.ParseUint(institutionAdministratorPositionIDFields[0], 10, 8)
		if err != nil {
			helpers.ErrorResponse(w, "cannot convert institution administrator position ID from string to uint64", http.StatusInternalServerError)
			return
		}
		institutionAdministratorPositionID := uint8(institutionAdministratorPositionID64)
		userRolesDTO.InstitutionAdministratorPositionID = &institutionAdministratorPositionID
	} else if len(institutionAdministratorPositionIDFields) != 0 {
		helpers.ErrorResponse(w, "failed to parse form: to much institution administrator position id values", http.StatusBadRequest)
		return
	}
	// ParentStudentIDs (special)
	parentStudentIDsFields := r.PostForm["parentStudentId"]
	var parentStudentIDs = make([]uuid.UUID, len(parentStudentIDsFields))
	for i, parentStudentIDString := range parentStudentIDsFields {
		parentStudentID, err := uuid.Parse(parentStudentIDString)
		if err != nil {
			helpers.ErrorResponse(w, "cannot convert student id to uuid", http.StatusBadRequest)
			return
		}
		parentStudentIDs[i] = parentStudentID
	}
	userRolesDTO.ParentStudentIDs = parentStudentIDs
	// Add new roles to the old ones
	if err := h.userService.AddRolesToUser(ctx, userID, userRolesDTO, roleIDs); err != nil {
		helpers.HandleServiceError(w, err)
		return
	}
	// Get updated roles
	roles, err := h.userService.GetUserRoles(ctx, userID)
	if err != nil {
		helpers.HandleServiceError(w, err)
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"userID":  userID,
		"roles":   roles,
		"message": "roles added successfully",
	})
}

func (h *UserHandler) RemoveRole(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// Check method
	if r.Method != http.MethodDelete {
		helpers.ErrorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Get and convert user ID
	userID, err := uuid.Parse(r.PathValue("userId"))
	if err != nil {
		helpers.ErrorResponse(w, "cannot convert user id to uuid", http.StatusBadRequest)
		return
	}
	// Get and convert role ID
	roleIDString := r.PathValue("roleId")
	roleID64, err := strconv.ParseUint(roleIDString, 10, 8)
	if err != nil {
		h.log.Error("cannot convert role ID from string to uint64")
		helpers.ErrorResponse(w, "cannot convert role ID from string to uint64", http.StatusInternalServerError)
		return
	}
	roleID := uint8(roleID64)
	// Remove user role
	if err := h.userService.RemoveRoleFromUser(ctx, userID, roleID); err != nil {
		helpers.HandleServiceError(w, err)
		return
	}
	helpers.JsonResponse(w, map[string]interface{}{}, http.StatusNoContent)
}
