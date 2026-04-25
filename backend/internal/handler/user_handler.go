package handler

import (
	"backend/internal/permissions"
	"backend/internal/repository"
	"backend/internal/service"
	"backend/pkg/helpers"
	"backend/pkg/logger"
	"backend/pkg/middleware"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"slices"
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
	// Get and convert user ID
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		h.log.Error("failed to get userID from context and convert it to UUID")
		helpers.InternalError(h.log, w)
		return
	}
	// DTO (all fields are optional)
	dto := service.UpdateUserDTO{}
	if firstNameFields := r.PostForm["firstName"]; len(firstNameFields) == 1 {
		dto.FirstName = &firstNameFields[0]
	} else if len(firstNameFields) != 0 {
		h.log.Error("failed to parse form: too many firstName values")
		helpers.TooManyFieldsError(h.log, w, "firstName")
		return
	}
	if middleNameFields := r.PostForm["middleName"]; len(middleNameFields) == 1 {
		dto.MiddleName = &middleNameFields[0]
	} else if len(middleNameFields) != 0 {
		h.log.Error("failed to parse form: too many middleName values")
		helpers.TooManyFieldsError(h.log, w, "middleName")
		return
	}
	if lastNameFields := r.PostForm["lastName"]; len(lastNameFields) == 1 {
		dto.LastName = &lastNameFields[0]
	} else if len(lastNameFields) != 0 {
		h.log.Error("failed to parse form: too many lastName values")
		helpers.TooManyFieldsError(h.log, w, "lastName")
		return
	}
	// Update user
	userResponse, err := h.userService.UpdateUser(r.Context(), userID, dto)
	if err != nil {
		helpers.HandleServiceError(h.log, w, fmt.Errorf("failed to update the user profile: %w", err))
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
		helpers.MethodNotAllowedError(h.log, w)
		return
	}
	// Get and convert user ID
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		h.log.Error("failed to get userID from context and convert it to UUID")
		helpers.InternalError(h.log, w)
		return
	}
	// Remove user avatar file
	if err := h.userService.RemoveAvatar(r.Context(), userID); err != nil {
		helpers.HandleServiceError(h.log, w, fmt.Errorf("failed to remove user avatar file: %w", err))
		return
	}
	// Return response
	helpers.JsonResponse(w, map[string]interface{}{}, http.StatusNoContent)
}

func (h *UserHandler) UpdateOwnAvatar(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodPut {
		helpers.MethodNotAllowedError(h.log, w)
		return
	}
	// Restrictions
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20) // 10 MB
	// Parse form
	// TODO: add cleaning of temporary data (ParseMultipartForm, r.MultipartForm, etc)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		h.log.Error("failed to parse multipart/formdata form")
		helpers.BadRequestError(h.log, w)
		return
	}
	// Get and convert user ID
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		h.log.Error("failed to get userID from context and convert it to UUID")
		helpers.InternalError(h.log, w)
		return
	}
	// Get avatar file from the request
	formFiles := r.MultipartForm.File["avatar"]
	if len(formFiles) > 1 {
		h.log.Error("failed to parse form: too many avatar files")
		helpers.TooManyFieldsError(h.log, w, "avatar")
		return
	} else if len(formFiles) == 0 {
		h.log.Error("failed to parse form: avatar cannot be empty")
		helpers.FieldRequiredError(h.log, w, "avatar")
		return
	}
	// Update avatar file
	if err := h.userService.UpdateAvatar(r.Context(), userID, formFiles[0]); err != nil {
		helpers.HandleServiceError(h.log, w, fmt.Errorf("failed to update the avatar: %w", err))
		return
	}
	// Return response
	helpers.JsonResponse(w, map[string]interface{}{}, http.StatusNoContent)
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodGet {
		helpers.MethodNotAllowedError(h.log, w)
		return
	}
	// Get and convert user ID
	userID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		h.log.Error("cannot convert user id to uuid")
		helpers.BadRequestFieldError(h.log, w, "id")
		return
	}
	// Get user
	response, err := h.userService.GetUserByID(r.Context(), userID)
	if err != nil {
		helpers.HandleServiceError(h.log, w, fmt.Errorf("failed to get user by id: %w", err))
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
		helpers.MethodNotAllowedError(h.log, w)
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
		roleID64, err := strconv.ParseUint(roleIDString, 10, 16)
		if err != nil {
			h.log.Error("cannot convert role ID from string to uint64")
			helpers.BadRequestFieldError(h.log, w, "roleId")
			return
		}
		// and to uint16
		roleID := uint16(roleID64)
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
	// Get users
	users, err := h.userService.GetUsers(r.Context(), filter)
	if err != nil {
		helpers.HandleServiceError(h.log, w, fmt.Errorf("failed to get users: %w", err))
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
		helpers.MethodNotAllowedError(h.log, w)
		return
	}
	// Get and convert user ID
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		h.log.Error("failed to get userID from context and convert it to UUID")
		helpers.InternalError(h.log, w)
		return
	}
	// Get user
	response, err := h.userService.GetUserByID(r.Context(), userID)
	if err != nil {
		helpers.HandleServiceError(h.log, w, fmt.Errorf("failed to get own user by id: %w", err))
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
		helpers.MethodNotAllowedError(h.log, w)
		return
	}
	// Get and convert user ID
	userID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		h.log.Error("cannot convert user id to uuid")
		helpers.BadRequestFieldError(h.log, w, "id")
		return
	}
	// Get roles
	roles, err := h.userService.GetUserRoles(r.Context(), userID)
	if err != nil {
		helpers.HandleServiceError(h.log, w, err)
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
		helpers.MethodNotAllowedError(h.log, w)
		return
	}
	// Get and convert user ID
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		h.log.Error("failed to get userID from context and convert it to UUID")
		helpers.InternalError(h.log, w)
		return
	}
	// Get roles
	roles, err := h.userService.GetUserRoles(r.Context(), userID)
	if err != nil {
		helpers.HandleServiceError(h.log, w, err)
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
	// Get and convert user ID
	userID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		h.log.Error("cannot convert user id to uuid")
		helpers.BadRequestFieldError(h.log, w, "id")
		return
	}
	// Get role IDs
	roleIDsFields := r.PostForm["roleId"]
	if len(roleIDsFields) == 0 {
		h.log.Error("the list of roles cannot be empty")
		helpers.AtLeastOneFieldError(h.log, w, "roleId")
		return
	}
	roleIDs := make([]uint16, len(roleIDsFields))
	for i, s := range roleIDsFields {
		val, err := strconv.ParseUint(s, 10, 16)
		if err != nil {
			h.log.Error("cannot convert IDs of roles from string to uint64")
			helpers.BadRequestFieldError(h.log, w, "roleId")
			return
		}
		roleIDs[i] = uint16(val)
	}
	// Get user permissions
	userPermissions, ok := r.Context().Value(middleware.UserPermissionsKey).([]string)
	if !ok {
		h.log.Error("failed to get user permissions from the context")
		helpers.InternalError(h.log, w)
		return
	}
	// Depending on whether assigning admin role (2) or user roles (3-7) require
	// different permissions
	// admin:
	if slices.Contains(roleIDs, 2) {
		if !slices.Contains(userPermissions, permissions.RoleAdminAssign) {
			h.log.Error("forbidden: you do not have permission to assign admin role")
			helpers.ForbiddenError(h.log, w)
			return
		}
	}
	// user:
	for _, roleID := range roleIDs {
		if slices.Contains([]uint16{3, 4, 5, 6, 7}, roleID) {
			if !slices.Contains(userPermissions, permissions.RoleUserAssign) {
				h.log.Error("forbidden: you do not have permission to assign user role")
				helpers.ForbiddenError(h.log, w)
				return
			}
			break
		}
	}
	// Get special fields (for user-extension tables)
	userExtensionsDTO := service.UserExtensionsDTO{}
	// TeacherClassroomID (special)
	if teacherClassroomIDFields := r.PostForm["teacherClassroomId"]; len(teacherClassroomIDFields) == 1 {
		// Convert to uint16
		teacherClassroomID64, err := strconv.ParseUint(teacherClassroomIDFields[0], 10, 16)
		if err != nil {
			h.log.Error("cannot convert teacher classroom ID from string to uint64")
			helpers.BadRequestFieldError(h.log, w, "teacherClassroomId")
			return
		}
		teacherClassroomID := uint16(teacherClassroomID64)
		userExtensionsDTO.TeacherClassroomID = &teacherClassroomID
	} else if len(teacherClassroomIDFields) != 0 {
		h.log.Error("failed to parse form: too many teacher classroom id values")
		helpers.TooManyFieldsError(h.log, w, "teacherClassroomId")
		return
	}
	// TeacherSubjectIDs (special)
	teacherSubjectIDsFields := r.PostForm["teacherSubjectId"]
	var teacherSubjectIDs = make([]uint16, len(teacherSubjectIDsFields))
	for i, subjectIDString := range teacherSubjectIDsFields {
		subjectID64, err := strconv.ParseUint(subjectIDString, 10, 16)
		if err != nil {
			h.log.Error("cannot convert teacher subject ID from string to uint64")
			helpers.BadRequestFieldError(h.log, w, "teacherSubjectId")
			return
		}
		subjectID8 := uint16(subjectID64)
		teacherSubjectIDs[i] = subjectID8
	}
	userExtensionsDTO.TeacherSubjectIDs = teacherSubjectIDs
	// TeacherStudentGroupIDs (special)
	if teacherStudentGroupIDsFields := r.PostForm["teacherStudentGroupId"]; len(teacherStudentGroupIDsFields) != 0 {
		var teacherStudentGroupIDs = make([]uint16, len(teacherStudentGroupIDsFields))
		for i, groupIDString := range teacherStudentGroupIDsFields {
			groupID64, err := strconv.ParseUint(groupIDString, 10, 16)
			if err != nil {
				h.log.Error("cannot convert teacher student group ID from string to uint64")
				helpers.BadRequestFieldError(h.log, w, "teacherStudentGroupId")
				return
			}
			groupID16 := uint16(groupID64)
			teacherStudentGroupIDs[i] = groupID16
		}
		userExtensionsDTO.TeacherStudentGroupIDs = teacherStudentGroupIDs
	}
	// StudentGroupID (special)
	if studentGroupIDFields := r.PostForm["studentGroupId"]; len(studentGroupIDFields) == 1 {
		// Convert to uint16
		studentGroupID64, err := strconv.ParseUint(studentGroupIDFields[0], 10, 16)
		if err != nil {
			h.log.Error("cannot convert student group ID from string to uint64")
			helpers.BadRequestFieldError(h.log, w, "studentGroupId")
			return
		}
		studentGroupID := uint16(studentGroupID64)
		userExtensionsDTO.StudentGroupID = &studentGroupID
	} else if len(studentGroupIDFields) != 0 {
		h.log.Error("failed to parse form: too many student group id values")
		helpers.TooManyFieldsError(h.log, w, "studentGroupId")
		return
	}
	// StaffPositionID (special)
	if staffPositionIDFields := r.PostForm["staffPositionId"]; len(staffPositionIDFields) == 1 {
		// Convert to uint16
		staffPositionID64, err := strconv.ParseUint(staffPositionIDFields[0], 10, 16)
		if err != nil {
			h.log.Error("cannot convert staff position ID from string to uint64")
			helpers.BadRequestFieldError(h.log, w, "staffPositionId")
			return
		}
		staffPositionID := uint16(staffPositionID64)
		userExtensionsDTO.StaffPositionID = &staffPositionID
	} else if len(staffPositionIDFields) != 0 {
		h.log.Error("failed to parse form: too many staff position id values")
		helpers.TooManyFieldsError(h.log, w, "staffPositionId")
		return
	}
	// InstitutionAdministratorPositionID (special)
	if institutionAdministratorPositionIDFields := r.PostForm["institutionAdministratorPositionId"]; len(institutionAdministratorPositionIDFields) == 1 {
		// Convert to uint16
		institutionAdministratorPositionID64, err := strconv.ParseUint(institutionAdministratorPositionIDFields[0], 10, 16)
		if err != nil {
			h.log.Error("cannot convert institution administrator position ID from string to uint64")
			helpers.BadRequestFieldError(h.log, w, "institutionAdministratorPositionId")
			return
		}
		institutionAdministratorPositionID := uint16(institutionAdministratorPositionID64)
		userExtensionsDTO.InstitutionAdministratorPositionID = &institutionAdministratorPositionID
	} else if len(institutionAdministratorPositionIDFields) != 0 {
		h.log.Error("failed to parse form: too many institution administrator position id values")
		helpers.TooManyFieldsError(h.log, w, "institutionAdministratorPositionId")
		return
	}
	// ParentStudentIDs (special)
	parentStudentIDsFields := r.PostForm["parentStudentId"]
	var parentStudentIDs = make([]uuid.UUID, len(parentStudentIDsFields))
	for i, parentStudentIDString := range parentStudentIDsFields {
		parentStudentID, err := uuid.Parse(parentStudentIDString)
		if err != nil {
			h.log.Error("cannot convert student id to uuid")
			helpers.BadRequestFieldError(h.log, w, "parentStudentId")
			return
		}
		parentStudentIDs[i] = parentStudentID
	}
	userExtensionsDTO.ParentStudentIDs = parentStudentIDs
	// Replace old roles with new ones
	if err := h.userService.AssignRolesToUser(ctx, userID, userExtensionsDTO, roleIDs); err != nil {
		helpers.HandleServiceError(h.log, w, err)
		return
	}
	// Get updated roles
	roles, err := h.userService.GetUserRoles(ctx, userID)
	if err != nil {
		helpers.HandleServiceError(h.log, w, err)
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"userID":  userID,
		"roles":   roles,
		"message": "roles updated successfully",
	})
}

// Handler to assign roles-depending fields that extends the "user" table.
func (h *UserHandler) AssignExtensionsOwn(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodPut {
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
	// Get and convert user ID
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		h.log.Error("failed to get userID from context and convert it to UUID")
		helpers.InternalError(h.log, w)
		return
	}
	// Get fields (for user-extension tables)
	userExtensionsDTO := service.UserExtensionsDTO{}
	// TeacherClassroomID
	if teacherClassroomIDFields := r.PostForm["teacherClassroomId"]; len(teacherClassroomIDFields) == 1 {
		// Convert to uint16
		teacherClassroomID64, err := strconv.ParseUint(teacherClassroomIDFields[0], 10, 16)
		if err != nil {
			h.log.Error("cannot convert teacher classroom ID from string to uint64")
			helpers.BadRequestFieldError(h.log, w, "teacherClassroomId")
			return
		}
		teacherClassroomID := uint16(teacherClassroomID64)
		userExtensionsDTO.TeacherClassroomID = &teacherClassroomID
	} else if len(teacherClassroomIDFields) != 0 {
		h.log.Error("failed to parse form: too many teacher classroom id values")
		helpers.TooManyFieldsError(h.log, w, "teacherClassroomId")
		return
	}
	// TODO: read only required fields (depending on roles)
	// TeacherSubjectIDs
	teacherSubjectIDsFields := r.PostForm["teacherSubjectId"]
	var teacherSubjectIDs = make([]uint16, len(teacherSubjectIDsFields))
	for i, subjectIDString := range teacherSubjectIDsFields {
		subjectID64, err := strconv.ParseUint(subjectIDString, 10, 16)
		if err != nil {
			h.log.Error("cannot convert teacher subject ID from string to uint64")
			helpers.BadRequestFieldError(h.log, w, "teacherSubjectId")
			return
		}
		subjectID8 := uint16(subjectID64)
		teacherSubjectIDs[i] = subjectID8
	}
	userExtensionsDTO.TeacherSubjectIDs = teacherSubjectIDs
	// TeacherStudentGroupIDs
	if teacherStudentGroupIDsFields := r.PostForm["teacherStudentGroupId"]; len(teacherStudentGroupIDsFields) != 0 {
		var teacherStudentGroupIDs = make([]uint16, len(teacherStudentGroupIDsFields))
		for i, groupIDString := range teacherStudentGroupIDsFields {
			groupID64, err := strconv.ParseUint(groupIDString, 10, 16)
			if err != nil {
				h.log.Error("cannot convert teacher student group ID from string to uint64")
				helpers.BadRequestFieldError(h.log, w, "teacherStudentGroupId")
				return
			}
			groupID16 := uint16(groupID64)
			teacherStudentGroupIDs[i] = groupID16
		}
		userExtensionsDTO.TeacherStudentGroupIDs = teacherStudentGroupIDs
	}
	// StudentGroupID
	if studentGroupIDFields := r.PostForm["studentGroupId"]; len(studentGroupIDFields) == 1 {
		// Convert to uint16
		studentGroupID64, err := strconv.ParseUint(studentGroupIDFields[0], 10, 16)
		if err != nil {
			h.log.Error("cannot convert student group ID from string to uint64")
			helpers.BadRequestFieldError(h.log, w, "studentGroupId")
			return
		}
		studentGroupID := uint16(studentGroupID64)
		userExtensionsDTO.StudentGroupID = &studentGroupID
	} else if len(studentGroupIDFields) != 0 {
		h.log.Error("failed to parse form: too many student group id values")
		helpers.TooManyFieldsError(h.log, w, "studentGroupId")
		return
	}
	// StaffPositionID
	if staffPositionIDFields := r.PostForm["staffPositionId"]; len(staffPositionIDFields) == 1 {
		// Convert to uint16
		staffPositionID64, err := strconv.ParseUint(staffPositionIDFields[0], 10, 16)
		if err != nil {
			h.log.Error("cannot convert staff position ID from string to uint64")
			helpers.BadRequestFieldError(h.log, w, "staffPositionId")
			return
		}
		staffPositionID := uint16(staffPositionID64)
		userExtensionsDTO.StaffPositionID = &staffPositionID
	} else if len(staffPositionIDFields) != 0 {
		h.log.Error("failed to parse form: too many staff position id values")
		helpers.TooManyFieldsError(h.log, w, "staffPositionId")
		return
	}
	// InstitutionAdministratorPositionID
	if institutionAdministratorPositionIDFields := r.PostForm["institutionAdministratorPositionId"]; len(institutionAdministratorPositionIDFields) == 1 {
		// Convert to uint16
		institutionAdministratorPositionID64, err := strconv.ParseUint(institutionAdministratorPositionIDFields[0], 10, 16)
		if err != nil {
			h.log.Error("cannot convert institution administrator position ID from string to uint64")
			helpers.BadRequestFieldError(h.log, w, "institutionAdministratorPositionId")
			return
		}
		institutionAdministratorPositionID := uint16(institutionAdministratorPositionID64)
		userExtensionsDTO.InstitutionAdministratorPositionID = &institutionAdministratorPositionID
	} else if len(institutionAdministratorPositionIDFields) != 0 {
		h.log.Error("failed to parse form: too many institution administrator position id values")
		helpers.TooManyFieldsError(h.log, w, "institutionAdministratorPositionId")
		return
	}
	// ParentStudentIDs
	parentStudentIDsFields := r.PostForm["parentStudentId"]
	var parentStudentIDs = make([]uuid.UUID, len(parentStudentIDsFields))
	for i, parentStudentIDString := range parentStudentIDsFields {
		parentStudentID, err := uuid.Parse(parentStudentIDString)
		if err != nil {
			h.log.Error("cannot convert student id to uuid")
			helpers.BadRequestFieldError(h.log, w, "parentStudentId")
			return
		}
		parentStudentIDs[i] = parentStudentID
	}
	userExtensionsDTO.ParentStudentIDs = parentStudentIDs
	// Replace old extensions with new ones
	if err := h.userService.AssignExtensionsToUser(r.Context(), userID, userExtensionsDTO); err != nil {
		helpers.HandleServiceError(h.log, w, err)
		return
	}
	helpers.JsonResponse(w, map[string]interface{}{}, http.StatusNoContent)
}

func (h *UserHandler) AssignNonAdminRoles(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// Check method
	if r.Method != http.MethodPut {
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
	// Get and convert user ID
	userID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		h.log.Error("cannot convert user id to uuid")
		helpers.BadRequestFieldError(h.log, w, "id")
		return
	}
	// Get role IDs
	roleIDsFields := r.PostForm["roleId"]
	if len(roleIDsFields) == 0 {
		h.log.Error("the list of roles cannot be empty")
		helpers.AtLeastOneFieldError(h.log, w, "roleId")
		return
	}
	roleIDs := make([]uint16, len(roleIDsFields))
	for i, s := range roleIDsFields {
		val, err := strconv.ParseUint(s, 10, 16)
		if err != nil {
			h.log.Error("cannot convert IDs of roles from string to uint64")
			helpers.BadRequestFieldError(h.log, w, "roleId")
			return
		}
		roleIDs[i] = uint16(val)
	}
	// Get user permissions
	userPermissions, ok := r.Context().Value(middleware.UserPermissionsKey).([]string)
	if !ok {
		h.log.Error("failed to get user permissions from the context")
		helpers.InternalError(h.log, w)
		return
	}
	// Check permissions
	if !slices.Contains(userPermissions, permissions.RoleUserAssign) {
		h.log.Error("forbidden: you do not have permission to assign user role")
		helpers.ForbiddenError(h.log, w)
		return
	}
	// Get special fields (for user-extension tables)
	userExtensionsDTO := service.UserExtensionsDTO{}
	// TeacherClassroomID (special)
	if teacherClassroomIDFields := r.PostForm["teacherClassroomId"]; len(teacherClassroomIDFields) == 1 {
		// Convert to uint16
		teacherClassroomID64, err := strconv.ParseUint(teacherClassroomIDFields[0], 10, 16)
		if err != nil {
			h.log.Error("cannot convert teacher classroom ID from string to uint64")
			helpers.BadRequestFieldError(h.log, w, "teacherClassroomId")
			return
		}
		teacherClassroomID := uint16(teacherClassroomID64)
		userExtensionsDTO.TeacherClassroomID = &teacherClassroomID
	} else if len(teacherClassroomIDFields) != 0 {
		h.log.Error("failed to parse form: too many teacher classroom id values")
		helpers.TooManyFieldsError(h.log, w, "teacherClassroomId")
		return
	}
	// TeacherSubjectIDs (special)
	teacherSubjectIDsFields := r.PostForm["teacherSubjectId"]
	var teacherSubjectIDs = make([]uint16, len(teacherSubjectIDsFields))
	for i, subjectIDString := range teacherSubjectIDsFields {
		subjectID64, err := strconv.ParseUint(subjectIDString, 10, 16)
		if err != nil {
			h.log.Error("cannot convert teacher subject ID from string to uint64")
			helpers.BadRequestFieldError(h.log, w, "teacherSubjectId")
			return
		}
		subjectID8 := uint16(subjectID64)
		teacherSubjectIDs[i] = subjectID8
	}
	userExtensionsDTO.TeacherSubjectIDs = teacherSubjectIDs
	// TeacherStudentGroupIDs (special)
	if teacherStudentGroupIDsFields := r.PostForm["teacherStudentGroupId"]; len(teacherStudentGroupIDsFields) != 0 {
		var teacherStudentGroupIDs = make([]uint16, len(teacherStudentGroupIDsFields))
		for i, groupIDString := range teacherStudentGroupIDsFields {
			groupID64, err := strconv.ParseUint(groupIDString, 10, 16)
			if err != nil {
				h.log.Error("cannot convert teacher student group ID from string to uint64")
				helpers.BadRequestFieldError(h.log, w, "teacherStudentGroupId")
				return
			}
			groupID16 := uint16(groupID64)
			teacherStudentGroupIDs[i] = groupID16
		}
		userExtensionsDTO.TeacherStudentGroupIDs = teacherStudentGroupIDs
	}
	// StudentGroupID (special)
	if studentGroupIDFields := r.PostForm["studentGroupId"]; len(studentGroupIDFields) == 1 {
		// Convert to uint16
		studentGroupID64, err := strconv.ParseUint(studentGroupIDFields[0], 10, 16)
		if err != nil {
			h.log.Error("cannot convert student group ID from string to uint64")
			helpers.BadRequestFieldError(h.log, w, "studentGroupId")
			return
		}
		studentGroupID := uint16(studentGroupID64)
		userExtensionsDTO.StudentGroupID = &studentGroupID
	} else if len(studentGroupIDFields) != 0 {
		h.log.Error("failed to parse form: too many student group id values")
		helpers.TooManyFieldsError(h.log, w, "studentGroupId")
		return
	}
	// StaffPositionID (special)
	if staffPositionIDFields := r.PostForm["staffPositionId"]; len(staffPositionIDFields) == 1 {
		// Convert to uint16
		staffPositionID64, err := strconv.ParseUint(staffPositionIDFields[0], 10, 16)
		if err != nil {
			h.log.Error("cannot convert staff position ID from string to uint64")
			helpers.BadRequestFieldError(h.log, w, "staffPositionId")
			return
		}
		staffPositionID := uint16(staffPositionID64)
		userExtensionsDTO.StaffPositionID = &staffPositionID
	} else if len(staffPositionIDFields) != 0 {
		h.log.Error("failed to parse form: too many staff position id values")
		helpers.TooManyFieldsError(h.log, w, "staffPositionId")
		return
	}
	// InstitutionAdministratorPositionID (special)
	if institutionAdministratorPositionIDFields := r.PostForm["institutionAdministratorPositionId"]; len(institutionAdministratorPositionIDFields) == 1 {
		// Convert to uint16
		institutionAdministratorPositionID64, err := strconv.ParseUint(institutionAdministratorPositionIDFields[0], 10, 16)
		if err != nil {
			h.log.Error("cannot convert institution administrator position ID from string to uint64")
			helpers.BadRequestFieldError(h.log, w, "institutionAdministratorPositionId")
			return
		}
		institutionAdministratorPositionID := uint16(institutionAdministratorPositionID64)
		userExtensionsDTO.InstitutionAdministratorPositionID = &institutionAdministratorPositionID
	} else if len(institutionAdministratorPositionIDFields) != 0 {
		h.log.Error("failed to parse form: too many institution administrator position id values")
		helpers.TooManyFieldsError(h.log, w, "institutionAdministratorPositionId")
		return
	}
	// ParentStudentIDs (special)
	parentStudentIDsFields := r.PostForm["parentStudentId"]
	var parentStudentIDs = make([]uuid.UUID, len(parentStudentIDsFields))
	for i, parentStudentIDString := range parentStudentIDsFields {
		parentStudentID, err := uuid.Parse(parentStudentIDString)
		if err != nil {
			h.log.Error("cannot convert parent student id to uuid")
			helpers.BadRequestFieldError(h.log, w, "parentStudentId")
			return
		}
		parentStudentIDs[i] = parentStudentID
	}
	userExtensionsDTO.ParentStudentIDs = parentStudentIDs
	// Replace old roles with new ones
	if err := h.userService.AssignNonAdminRolesToUser(ctx, userID, userExtensionsDTO, roleIDs); err != nil {
		helpers.HandleServiceError(h.log, w, err)
		return
	}
	// Get updated roles
	roles, err := h.userService.GetUserRoles(ctx, userID)
	if err != nil {
		helpers.HandleServiceError(h.log, w, err)
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
	// Get and convert user ID
	userID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		h.log.Error("cannot convert user id to uuid")
		helpers.BadRequestFieldError(h.log, w, "id")
		return
	}
	// Get role IDs
	roleIDsFields := r.PostForm["roleId"]
	if len(roleIDsFields) == 0 {
		h.log.Error("the list of roles cannot be empty")
		helpers.AtLeastOneFieldError(h.log, w, "roleId")
		return
	}
	roleIDs := make([]uint16, len(roleIDsFields))
	for i, s := range roleIDsFields {
		val, err := strconv.ParseUint(s, 10, 16)
		if err != nil {
			h.log.Error("cannot convert IDs of roles from string to uint64")
			helpers.BadRequestFieldError(h.log, w, "roleId")
			return
		}
		roleIDs[i] = uint16(val)
	}
	// Get user permissions
	userPermissions, ok := r.Context().Value(middleware.UserPermissionsKey).([]string)
	if !ok {
		h.log.Error("failed to get user permissions from the context")
		helpers.InternalError(h.log, w)
		return
	}
	// Depending on whether adding admin role (2) or user roles (3-7) require
	// different permissions
	// admin:
	if slices.Contains(roleIDs, 2) {
		if !slices.Contains(userPermissions, permissions.RoleAdminAdd) {
			h.log.Error("forbidden: you do not have permission to add admin role")
			helpers.ForbiddenError(h.log, w)
			return
		}
	}
	// user:
	for _, roleID := range roleIDs {
		if slices.Contains([]uint16{3, 4, 5, 6, 7}, roleID) {
			if !slices.Contains(userPermissions, permissions.RoleUserAdd) {
				h.log.Error("forbidden: you do not have permission to add user role")
				helpers.ForbiddenError(h.log, w)
				return
			}
			break
		}
	}
	// Get special fields (for user-extension tables)
	userExtensionsDTO := service.UserExtensionsDTO{}
	// TeacherClassroomID (special)
	if teacherClassroomIDFields := r.PostForm["teacherClassroomId"]; len(teacherClassroomIDFields) == 1 {
		// Convert to uint16
		teacherClassroomID64, err := strconv.ParseUint(teacherClassroomIDFields[0], 10, 16)
		if err != nil {
			h.log.Error("cannot convert teacher classroom ID from string to uint64")
			helpers.BadRequestFieldError(h.log, w, "teacherClassroomId")
			return
		}
		teacherClassroomID := uint16(teacherClassroomID64)
		userExtensionsDTO.TeacherClassroomID = &teacherClassroomID
	} else if len(teacherClassroomIDFields) != 0 {
		h.log.Error("failed to parse form: too many teacher classroom id values")
		helpers.TooManyFieldsError(h.log, w, "teacherClassroomId")
		return
	}
	// TeacherSubjectIDs (special)
	teacherSubjectIDsFields := r.PostForm["teacherSubjectId"]
	var teacherSubjectIDs = make([]uint16, len(teacherSubjectIDsFields))
	for i, subjectIDString := range teacherSubjectIDsFields {
		subjectID64, err := strconv.ParseUint(subjectIDString, 10, 16)
		if err != nil {
			h.log.Error("cannot convert teacher subject ID from string to uint64")
			helpers.BadRequestFieldError(h.log, w, "teacherSubjectId")
			return
		}
		subjectID8 := uint16(subjectID64)
		teacherSubjectIDs[i] = subjectID8
	}
	userExtensionsDTO.TeacherSubjectIDs = teacherSubjectIDs
	// TeacherStudentGroupIDs (special)
	if teacherStudentGroupIDsFields := r.PostForm["teacherStudentGroupId"]; len(teacherStudentGroupIDsFields) != 0 {
		var teacherStudentGroupIDs = make([]uint16, len(teacherStudentGroupIDsFields))
		for i, groupIDString := range teacherStudentGroupIDsFields {
			groupID64, err := strconv.ParseUint(groupIDString, 10, 16)
			if err != nil {
				h.log.Error("cannot convert teacher student group ID from string to uint64")
				helpers.BadRequestFieldError(h.log, w, "teacherStudentGroupId")
				return
			}
			groupID16 := uint16(groupID64)
			teacherStudentGroupIDs[i] = groupID16
		}
		userExtensionsDTO.TeacherStudentGroupIDs = teacherStudentGroupIDs
	}
	// StudentGroupID (special)
	if studentGroupIDFields := r.PostForm["studentGroupId"]; len(studentGroupIDFields) == 1 {
		// Convert to uint16
		studentGroupID64, err := strconv.ParseUint(studentGroupIDFields[0], 10, 16)
		if err != nil {
			h.log.Error("cannot convert student group ID from string to uint64")
			helpers.BadRequestFieldError(h.log, w, "studentGroupId")
			return
		}
		studentGroupID := uint16(studentGroupID64)
		userExtensionsDTO.StudentGroupID = &studentGroupID
	} else if len(studentGroupIDFields) != 0 {
		h.log.Error("failed to parse form: too many student group id values")
		helpers.TooManyFieldsError(h.log, w, "studentGroupId")
		return
	}
	// StaffPositionID (special)
	if staffPositionIDFields := r.PostForm["staffPositionId"]; len(staffPositionIDFields) == 1 {
		// Convert to uint16
		staffPositionID64, err := strconv.ParseUint(staffPositionIDFields[0], 10, 16)
		if err != nil {
			h.log.Error("cannot convert staff position ID from string to uint64")
			helpers.BadRequestFieldError(h.log, w, "staffPositionId")
			return
		}
		staffPositionID := uint16(staffPositionID64)
		userExtensionsDTO.StaffPositionID = &staffPositionID
	} else if len(staffPositionIDFields) != 0 {
		h.log.Error("failed to parse form: too many staff position id values")
		helpers.TooManyFieldsError(h.log, w, "staffPositionId")
		return
	}
	// InstitutionAdministratorPositionID (special)
	if institutionAdministratorPositionIDFields := r.PostForm["institutionAdministratorPositionId"]; len(institutionAdministratorPositionIDFields) == 1 {
		// Convert to uint16
		institutionAdministratorPositionID64, err := strconv.ParseUint(institutionAdministratorPositionIDFields[0], 10, 16)
		if err != nil {
			h.log.Error("cannot convert institution administrator position ID from string to uint64")
			helpers.BadRequestFieldError(h.log, w, "institutionAdministratorPositionId")
			return
		}
		institutionAdministratorPositionID := uint16(institutionAdministratorPositionID64)
		userExtensionsDTO.InstitutionAdministratorPositionID = &institutionAdministratorPositionID
	} else if len(institutionAdministratorPositionIDFields) != 0 {
		h.log.Error("failed to parse form: too many institution administrator position id values")
		helpers.TooManyFieldsError(h.log, w, "institutionAdministratorPositionId")
		return
	}
	// ParentStudentIDs (special)
	parentStudentIDsFields := r.PostForm["parentStudentId"]
	var parentStudentIDs = make([]uuid.UUID, len(parentStudentIDsFields))
	for i, parentStudentIDString := range parentStudentIDsFields {
		parentStudentID, err := uuid.Parse(parentStudentIDString)
		if err != nil {
			h.log.Error("cannot convert parent student id to uuid")
			helpers.BadRequestFieldError(h.log, w, "parentStudentId")
			return
		}
		parentStudentIDs[i] = parentStudentID
	}
	userExtensionsDTO.ParentStudentIDs = parentStudentIDs
	// Add new roles to the old ones
	if err := h.userService.AddRolesToUser(ctx, userID, userExtensionsDTO, roleIDs); err != nil {
		helpers.HandleServiceError(h.log, w, err)
		return
	}
	// Get updated roles
	roles, err := h.userService.GetUserRoles(ctx, userID)
	if err != nil {
		helpers.HandleServiceError(h.log, w, err)
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
		helpers.MethodNotAllowedError(h.log, w)
		return
	}
	// Get and convert user ID
	userID, err := uuid.Parse(r.PathValue("userId"))
	if err != nil {
		h.log.Error("cannot convert user id to uuid")
		helpers.BadRequestFieldError(h.log, w, "userId")
		return
	}
	// Get and convert role ID
	roleIDString := r.PathValue("roleId")
	roleID64, err := strconv.ParseUint(roleIDString, 10, 16)
	if err != nil {
		h.log.Error("cannot convert role ID from string to uint64")
		helpers.BadRequestFieldError(h.log, w, "roleId")
		return
	}
	roleID := uint16(roleID64)
	// Get user permissions
	userPermissions, ok := r.Context().Value(middleware.UserPermissionsKey).([]string)
	if !ok {
		h.log.Error("failed to get user permissions from the context")
		helpers.InternalError(h.log, w)
		return
	}
	// Depending on whether adding admin role (2) or user roles (3-7) require
	// different permissions
	// admin:
	if roleID == 2 {
		if !slices.Contains(userPermissions, permissions.RoleAdminUnassign) {
			h.log.Error("forbidden: you do not have permission to unassign admin role")
			helpers.ForbiddenError(h.log, w)
			return
		}
	}
	// user:
	if slices.Contains([]uint16{3, 4, 5, 6, 7}, roleID) {
		if !slices.Contains(userPermissions, permissions.RoleUserUnassign) {
			h.log.Error("forbidden: you do not have permission to unassign user role")
			helpers.ForbiddenError(h.log, w)
			return
		}
	}
	// Remove user role
	if err := h.userService.RemoveRoleFromUser(ctx, userID, roleID); err != nil {
		helpers.HandleServiceError(h.log, w, err)
		return
	}
	helpers.JsonResponse(w, map[string]interface{}{}, http.StatusNoContent)
}
