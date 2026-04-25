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

type PostHandler struct {
	postService         service.PostService
	teacherService      service.TeacherService
	parentService       service.ParentService
	studentGroupService service.StudentGroupService
	studentService      service.StudentService
	log                 logger.Logger
}

func NewPostHandler(postService service.PostService, teacherService service.TeacherService, parentService service.ParentService, studentGroupService service.StudentGroupService, studentService service.StudentService, log logger.Logger) *PostHandler {
	return &PostHandler{
		postService:         postService,
		teacherService:      teacherService,
		parentService:       parentService,
		studentGroupService: studentGroupService,
		studentService:      studentService,
		log:                 log,
	}
}

func (h *PostHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helpers.MethodNotAllowedError(h.log, w)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 15<<20) // 15 MB
	// TODO: add cleaning of temporary data (ParseMultipartForm, r.MultipartForm, etc)
	if err := r.ParseMultipartForm(15 << 20); err != nil {
		h.log.Error("failed to parse multipart/formdata form")
		helpers.BadRequestError(h.log, w)
		return
	}
	// Get name field
	nameFields := r.PostForm["name"]
	if len(nameFields) > 1 {
		h.log.Error(fmt.Sprintf("failed to parse form: too many name fields (%d)", len(nameFields)))
		helpers.TooManyFieldsError(h.log, w, "name")
		return
	} else if len(nameFields) == 0 {
		h.log.Error("failed to parse form: name field cannot be empty")
		helpers.FieldRequiredError(h.log, w, "name")
		return
	}
	name := nameFields[0]
	// Get and convert user ID
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		h.log.Error("failed to get userID from context and convert it to UUID")
		helpers.InternalError(h.log, w)
		return
	}
	// Pre-assemble DTO
	dto := service.CreatePostDTO{
		Name:     name,
		AuthorID: userID,
	}
	// Get description (optional field)
	if descriptionFields := r.PostForm["description"]; len(descriptionFields) == 1 {
		dto.Description = descriptionFields[0]
	} else if len(descriptionFields) != 0 {
		h.log.Error("failed to parse form: too many description values")
		helpers.TooManyFieldsError(h.log, w, "description")
		return
	}
	// Get post photo (optional field)
	formFiles := r.MultipartForm.File["photo"]
	if len(formFiles) > 1 {
		h.log.Error("failed to parse form: too many post photo files")
		helpers.TooManyFieldsError(h.log, w, "photo")
		return
	} else if len(formFiles) == 1 {
		dto.Photo = formFiles[0]
	}
	// Check if user can verify posts
	userPermissions, ok := r.Context().Value(middleware.UserPermissionsKey).([]string)
	if !ok {
		h.log.Error("failed to get user permissions from the context")
		helpers.InternalError(h.log, w)
		return
	}
	canVerifyPost := slices.Contains(userPermissions, permissions.PostVerify)
	// Create post
	postResponse, err := h.postService.CreatePost(r.Context(), dto, canVerifyPost)
	if err != nil {
		helpers.HandleServiceError(h.log, w, fmt.Errorf("failed to create the post: %w", err))
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
	// Get and convert post ID
	postID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		h.log.Error("cannot convert post id to uuid")
		helpers.BadRequestFieldError(h.log, w, "id")
		return
	}
	// Get user permissions
	userPermissions, ok := r.Context().Value(middleware.UserPermissionsKey).([]string)
	if !ok {
		h.log.Error("failed to get user permissions from the context")
		helpers.InternalError(h.log, w)
		return
	}
	// Check if user updating his own post
	if slices.Contains(userPermissions, permissions.PostUpdateOwn) && !slices.Contains(userPermissions, permissions.PostUpdateAny) {
		// Get and convert user ID
		userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
		if !ok {
			h.log.Error("failed to get userID from context and convert it to UUID")
			helpers.InternalError(h.log, w)
			return
		}
		// Get post
		post, err := h.postService.GetPostByID(r.Context(), postID)
		if err != nil || post == nil {
			helpers.HandleServiceError(h.log, w, fmt.Errorf("failed to find the post by ID: %w", err))
			return
		}
		// Check if the post belongs to the user
		if userID != post.Author.ID {
			h.log.Error("forbidden: you do not have permission to update this post")
			helpers.ForbiddenError(h.log, w)
			return
		}
	}
	// DTO (all fields are optional)
	dto := service.UpdatePostDTO{}
	if nameFields := r.PostForm["name"]; len(nameFields) == 1 {
		dto.Name = &nameFields[0]
	} else if len(nameFields) != 0 {
		h.log.Error("failed to parse form: too many name values")
		helpers.TooManyFieldsError(h.log, w, "name")
		return
	}
	if descriptionFields := r.PostForm["description"]; len(descriptionFields) == 1 {
		dto.Description = &descriptionFields[0]
	} else if len(descriptionFields) != 0 {
		h.log.Error("failed to parse form: too many description values")
		helpers.TooManyFieldsError(h.log, w, "description")
		return
	}
	// Update post
	postResponse, err := h.postService.UpdatePost(r.Context(), postID, dto)
	if err != nil {
		helpers.HandleServiceError(h.log, w, fmt.Errorf("failed to update the post: %w", err))
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
		helpers.MethodNotAllowedError(h.log, w)
		return
	}
	// Get and convert post ID
	postID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		h.log.Error("cannot convert post id to uuid")
		helpers.BadRequestFieldError(h.log, w, "id")
		return
	}
	// Get user permissions
	userPermissions, ok := r.Context().Value(middleware.UserPermissionsKey).([]string)
	if !ok {
		h.log.Error("failed to get user permissions from the context")
		helpers.InternalError(h.log, w)
		return
	}
	// Check if user deleting his own post
	if slices.Contains(userPermissions, permissions.PostDeleteOwn) && !slices.Contains(userPermissions, permissions.PostDeleteAny) {
		// Get and convert user ID
		userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
		if !ok {
			h.log.Error("failed to get userID from context and convert it to UUID")
			helpers.InternalError(h.log, w)
			return
		}
		// Get post
		post, err := h.postService.GetPostByID(r.Context(), postID)
		if err != nil || post == nil {
			helpers.HandleServiceError(h.log, w, fmt.Errorf("failed to find the post by ID: %w", err))
			return
		}
		// Check if the post belongs to the user
		if userID != post.Author.ID {
			h.log.Error("forbidden: you do not have permission to delete this post")
			helpers.ForbiddenError(h.log, w)
			return
		}
	}
	// Delete post
	if err := h.postService.DeletePost(r.Context(), postID); err != nil {
		helpers.HandleServiceError(h.log, w, fmt.Errorf("failed to delete the post: %w", err))
		return
	}
	// Return response
	helpers.JsonResponse(w, map[string]interface{}{}, http.StatusNoContent)
}

func (h *PostHandler) RemovePhoto(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodDelete {
		helpers.MethodNotAllowedError(h.log, w)
		return
	}
	// Get and convert post ID
	postID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		h.log.Error("cannot convert post id to uuid")
		helpers.BadRequestFieldError(h.log, w, "id")
		return
	}
	// Get user permissions
	userPermissions, ok := r.Context().Value(middleware.UserPermissionsKey).([]string)
	if !ok {
		h.log.Error("failed to get user permissions from the context")
		helpers.InternalError(h.log, w)
		return
	}
	// Check if user deleting photo of his own post
	if slices.Contains(userPermissions, permissions.PostPhotoDeleteOwn) && !slices.Contains(userPermissions, permissions.PostPhotoDeleteAny) {
		// Get and convert user ID
		userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
		if !ok {
			h.log.Error("failed to get userID from context and convert it to UUID")
			helpers.InternalError(h.log, w)
			return
		}
		// Get post
		post, err := h.postService.GetPostByID(r.Context(), postID)
		if err != nil || post == nil {
			helpers.HandleServiceError(h.log, w, fmt.Errorf("failed to find the post by ID: %w", err))
			return
		}
		// Check if the post belongs to the user
		if userID != post.Author.ID {
			h.log.Error("forbidden: you do not have permission to delete photo of this post")
			helpers.ForbiddenError(h.log, w)
			return
		}
	}
	// Remove post photo file
	if err := h.postService.RemovePhoto(r.Context(), postID); err != nil {
		helpers.HandleServiceError(h.log, w, fmt.Errorf("failed to remove post photo file: %w", err))
		return
	}
	// Return response
	helpers.JsonResponse(w, map[string]interface{}{}, http.StatusNoContent)
}

func (h *PostHandler) UpdatePhoto(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodPut {
		helpers.MethodNotAllowedError(h.log, w)
		return
	}
	// Restrictions
	r.Body = http.MaxBytesReader(w, r.Body, 15<<20) // 15 MB
	// Parse form
	// TODO: add cleaning of temporary data (ParseMultipartForm, r.MultipartForm, etc)
	if err := r.ParseMultipartForm(15 << 20); err != nil {
		h.log.Error("failed to parse multipart/formdata form")
		helpers.BadRequestError(h.log, w)
		return
	}
	// Get and convert post ID
	postID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		h.log.Error("cannot convert post id to uuid")
		helpers.BadRequestFieldError(h.log, w, "id")
		return
	}
	// Get user permissions
	userPermissions, ok := r.Context().Value(middleware.UserPermissionsKey).([]string)
	if !ok {
		h.log.Error("failed to get user permissions from the context")
		helpers.InternalError(h.log, w)
		return
	}
	// Check if user deleting photo of his own post
	if slices.Contains(userPermissions, permissions.PostPhotoUpdateOwn) && !slices.Contains(userPermissions, permissions.PostPhotoUpdateAny) {
		// Get and convert user ID
		userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
		if !ok {
			h.log.Error("failed to get userID from context and convert it to UUID")
			helpers.InternalError(h.log, w)
			return
		}
		// Get post
		post, err := h.postService.GetPostByID(r.Context(), postID)
		if err != nil || post == nil {
			helpers.HandleServiceError(h.log, w, fmt.Errorf("failed to find the post by ID: %w", err))
			return
		}
		// Check if the post belongs to the user
		if userID != post.Author.ID {
			h.log.Error("forbidden: you do not have permission to update photo of this post")
			helpers.ForbiddenError(h.log, w)
			return
		}
	}
	// Get photo file from the request
	formFiles := r.MultipartForm.File["photo"]
	if len(formFiles) > 1 {
		h.log.Error("failed to parse form: too many photo files")
		helpers.TooManyFieldsError(h.log, w, "photo")
		return
	} else if len(formFiles) == 0 {
		h.log.Error("failed to parse form: post photo cannot be empty")
		helpers.FieldRequiredError(h.log, w, "photo")
		return
	}
	// Update post photo
	if err := h.postService.UpdatePhoto(r.Context(), postID, formFiles[0]); err != nil {
		helpers.HandleServiceError(h.log, w, fmt.Errorf("failed to update post photo: %w", err))
		return
	}
	// Return response
	helpers.JsonResponse(w, map[string]interface{}{}, http.StatusNoContent)
}

func (h *PostHandler) GetPosts(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodGet {
		helpers.MethodNotAllowedError(h.log, w)
		return
	}
	// Parse query parameters (for filter)
	authorIDString := r.URL.Query().Get("authorId")
	verifiedString := r.URL.Query().Get("verified")
	thingReturnedToOwnerString := r.URL.Query().Get("thingReturnedToOwner")
	limitString := r.URL.Query().Get("limit")
	offsetString := r.URL.Query().Get("offset")
	// Pre-assemble filter (fill with default values)
	filter := repository.PostFilter{
		Limit:  20,
		Offset: 0,
	}
	// Parse author ID if passed
	if authorIDString != "" {
		// Convert to UUID
		authorID, err := uuid.Parse(authorIDString)
		if err != nil {
			h.log.Error("cannot convert author id (i.e. user id) to uuid")
			helpers.BadRequestFieldError(h.log, w, "authorId")
			return
		}
		// Add to filter
		filter.AuthorIDs = []uuid.UUID{authorID}
	}
	// Parse verification status if passed
	if verifiedString != "" {
		verified, err := strconv.ParseBool(verifiedString)
		if err != nil {
			h.log.Error("cannot convert verification status from string to boolean")
			helpers.BadRequestFieldError(h.log, w, "verified")
			return
		}
		filter.Verified = &verified
	}
	// Parse thing returning to owner status if passed
	if thingReturnedToOwnerString != "" {
		thingReturnedToOwner, err := strconv.ParseBool(thingReturnedToOwnerString)
		if err != nil {
			h.log.Error("cannot convert thing returning to owner status from string to boolean")
			helpers.BadRequestFieldError(h.log, w, "thingReturnedToOwner")
			return
		}
		filter.ThingReturnedToOwner = &thingReturnedToOwner
	}
	// Parse limit if passed
	// TODO: move limit and offset parsing to a separate helper
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
	// Get posts
	posts, err := h.postService.GetPosts(r.Context(), filter)
	if err != nil {
		helpers.HandleServiceError(h.log, w, fmt.Errorf("failed to get posts: %w", err))
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"posts": posts,
	})
}

func (h *PostHandler) GetPostByID(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodGet {
		helpers.MethodNotAllowedError(h.log, w)
		return
	}
	// Get and convert post ID
	postID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		h.log.Error("cannot convert post id to uuid")
		helpers.BadRequestFieldError(h.log, w, "id")
		return
	}
	// Get post
	post, err := h.postService.GetPostByID(r.Context(), postID)
	if err != nil {
		h.log.Error("Failed to get post by id", "error", err.Error())
		helpers.HandleServiceError(h.log, w, fmt.Errorf("failed to get post by id: %w", err))
		return
	}
	// Get user permissions
	userPermissions, ok := r.Context().Value(middleware.UserPermissionsKey).([]string)
	if !ok {
		h.log.Error("failed to get user permissions from the context")
		helpers.InternalError(h.log, w)
		return
	}
	// Get ID of the authorized user
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		h.log.Error("failed to get userID from context and convert it to UUID")
		helpers.InternalError(h.log, w)
		return
	}
	// Return post in three cases:
	// 1. if post verified (public access)
	// 2. if post was not verified, but the user is the author of this post
	// 3. if the user is not the author of the post, but he has permission to read any post
	if post.Verified || (slices.Contains(userPermissions, permissions.PostReadOwn) && (post.Author.ID == userID)) || slices.Contains(userPermissions, permissions.PostReadAny) {
		helpers.SuccessResponse(w, map[string]interface{}{
			"post": post,
		})
		return
	}
	h.log.Error("forbidden: you do not have permission to view this post")
	helpers.ForbiddenError(h.log, w)
}

func (h *PostHandler) GetPostsPublic(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodGet {
		helpers.MethodNotAllowedError(h.log, w)
		return
	}
	// Parse query parameters (for filter)
	authorIDString := r.URL.Query().Get("authorId")
	thingReturnedToOwnerString := r.URL.Query().Get("thingReturnedToOwner")
	authorString := r.URL.Query().Get("author") // filter posts by owners
	limitString := r.URL.Query().Get("limit")
	offsetString := r.URL.Query().Get("offset")
	// Pre-assemble filter (fill with default values)
	filter := repository.PostFilter{
		Limit:  20,
		Offset: 0,
	}
	// Parse author ID if passed
	if authorIDString != "" {
		// Convert to UUID
		authorID, err := uuid.Parse(authorIDString)
		if err != nil {
			h.log.Error("cannot convert author id (i.e. user id) to uuid")
			helpers.BadRequestFieldError(h.log, w, "authorId")
			return
		}
		// Add to filter
		filter.AuthorIDs = []uuid.UUID{authorID}
	}
	// Show only verified posts
	verified := true
	filter.Verified = &verified
	// Parse thing returning to owner status if passed
	if thingReturnedToOwnerString != "" {
		thingReturnedToOwner, err := strconv.ParseBool(thingReturnedToOwnerString)
		if err != nil {
			h.log.Error("cannot convert thing returning to owner status from string to boolean")
			helpers.BadRequestFieldError(h.log, w, "thingReturnedToOwner")
			return
		}
		filter.ThingReturnedToOwner = &thingReturnedToOwner
	}
	// Get and convert user ID
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		h.log.Error("failed to get userID from context and convert it to UUID")
		helpers.InternalError(h.log, w)
		return
	}
	// Parse author if passed
	if authorString != "" {
		switch authorString {
		case "all":
		case "me":
			filter.AuthorIDs = []uuid.UUID{userID}
		case "students": // my students (for teacher role)
			teacher, err := h.teacherService.GetTeacherByID(r.Context(), userID)
			if err != nil {
				helpers.HandleServiceError(h.log, w, err)
				return
			}
			studentIDs := []uuid.UUID{}
			for _, group := range teacher.StudentGroups {
				for _, student := range group.Students {
					studentIDs = append(studentIDs, student.UserID)
				}
			}
			filter.AuthorIDs = studentIDs
		case "children": // my children (for parent role)
			parent, err := h.parentService.GetParentByID(r.Context(), userID)
			if err != nil {
				helpers.HandleServiceError(h.log, w, err)
				return
			}
			studentIDs := []uuid.UUID{}
			for _, student := range parent.Students {
				studentIDs = append(studentIDs, student.UserID)
			}
			filter.AuthorIDs = studentIDs
		case "children_groups": // my children student groups (for parent role)
			parent, err := h.parentService.GetParentByID(r.Context(), userID)
			if err != nil {
				helpers.HandleServiceError(h.log, w, err)
				return
			}
			studentIDs := []uuid.UUID{}
			usedGroupIDs := []uint16{}
			for _, student := range parent.Students {
				if !slices.Contains(usedGroupIDs, student.StudentGroup.ID) {
					group, err := h.studentGroupService.GetStudentGroupByID(r.Context(), student.StudentGroup.ID)
					if err != nil {
						helpers.HandleServiceError(h.log, w, err)
						return
					}
					for _, student := range group.Students {
						studentIDs = append(studentIDs, student.UserID)
					}
					usedGroupIDs = append(usedGroupIDs, group.ID)
				}
			}
			filter.AuthorIDs = studentIDs
		case "parents": // my parents (for student role)
			student, err := h.studentService.GetStudentByID(r.Context(), userID)
			if err != nil {
				helpers.HandleServiceError(h.log, w, err)
				return
			}
			parentIDs := []uuid.UUID{}
			for _, parent := range student.Parents {
				parentIDs = append(parentIDs, parent.UserID)
			}
			filter.AuthorIDs = parentIDs
		case "classmates": // my student group, i.e. my classmates (for student role)
			student, err := h.studentService.GetStudentByID(r.Context(), userID)
			if err != nil {
				helpers.HandleServiceError(h.log, w, err)
				return
			}
			classmateIDs := []uuid.UUID{}
			for _, classmate := range student.StudentGroup.Students {
				classmateIDs = append(classmateIDs, classmate.UserID)
			}
			filter.AuthorIDs = classmateIDs
		default:
			h.log.Error("failed to parse author query parameter")
			helpers.BadRequestFieldError(h.log, w, "author")
			return
		}
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
	// Get posts
	posts, err := h.postService.GetPosts(r.Context(), filter)
	if err != nil {
		helpers.HandleServiceError(h.log, w, fmt.Errorf("failed to get posts: %w", err))
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"posts": posts,
	})
}

func (h *PostHandler) GetOwnPosts(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodGet {
		helpers.MethodNotAllowedError(h.log, w)
		return
	}
	// Parse query parameters (for filter)
	verifiedString := r.URL.Query().Get("verified")
	thingReturnedToOwnerString := r.URL.Query().Get("thingReturnedToOwner")
	limitString := r.URL.Query().Get("limit")
	offsetString := r.URL.Query().Get("offset")
	// Pre-assemble filter (fill with default values)
	filter := repository.PostFilter{
		Limit:  20,
		Offset: 0,
	}
	// Get and convert user ID
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		h.log.Error("failed to get userID from context and convert it to UUID")
		helpers.InternalError(h.log, w)
		return
	}
	// Set author ID to user ID
	filter.AuthorIDs = []uuid.UUID{userID}
	// Parse verification status if passed
	if verifiedString != "" {
		verified, err := strconv.ParseBool(verifiedString)
		if err != nil {
			h.log.Error("cannot convert verification status from string to boolean")
			helpers.BadRequestFieldError(h.log, w, "verified")
			return
		}
		filter.Verified = &verified
	}
	// Parse thing returning to owner status if passed
	if thingReturnedToOwnerString != "" {
		thingReturnedToOwner, err := strconv.ParseBool(thingReturnedToOwnerString)
		if err != nil {
			h.log.Error("cannot convert thing returning to owner status from string to boolean")
			helpers.BadRequestFieldError(h.log, w, "thingReturnedToOwner")
			return
		}
		filter.ThingReturnedToOwner = &thingReturnedToOwner
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
	// Get posts
	posts, err := h.postService.GetPosts(r.Context(), filter)
	if err != nil {
		helpers.HandleServiceError(h.log, w, fmt.Errorf("failed to get posts: %w", err))
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"posts": posts,
	})
}

func (h *PostHandler) Verify(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodPatch {
		helpers.MethodNotAllowedError(h.log, w)
		return
	}
	// Restrictions
	// TODO: think about the restrictions in the whole code
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MB
	// Parse form
	if err := r.ParseForm(); err != nil {
		h.log.Error("failed to parse x-www-form-urlencoded form")
		helpers.BadRequestError(h.log, w)
		return
	}
	// Get and convert post ID
	postID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		h.log.Error("cannot convert post id to uuid")
		helpers.BadRequestFieldError(h.log, w, "id")
		return
	}
	// Verify post
	postResponse, err := h.postService.VerifyPost(r.Context(), postID)
	if err != nil {
		helpers.HandleServiceError(h.log, w, fmt.Errorf("failed to change post verification status: %w", err))
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"post": postResponse,
	})
}

func (h *PostHandler) ReturnToOwner(w http.ResponseWriter, r *http.Request) {
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
	// Get and convert post ID
	postID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		h.log.Error("cannot convert post id to uuid")
		helpers.BadRequestFieldError(h.log, w, "id")
		return
	}
	// Get user permissions
	userPermissions, ok := r.Context().Value(middleware.UserPermissionsKey).([]string)
	if !ok {
		h.log.Error("failed to get user permissions from the context")
		helpers.InternalError(h.log, w)
		return
	}
	// Check if user updating his own post
	if slices.Contains(userPermissions, permissions.PostMarkReturnedOwn) && !slices.Contains(userPermissions, permissions.PostMarkReturnedAny) {
		// Get and convert user ID
		userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
		if !ok {
			h.log.Error("failed to get userID from context and convert it to UUID")
			helpers.InternalError(h.log, w)
			return
		}
		// Get post
		post, err := h.postService.GetPostByID(r.Context(), postID)
		if err != nil || post == nil {
			helpers.HandleServiceError(h.log, w, fmt.Errorf("failed to find the post by ID: %w", err))
			return
		}
		// Check if the post belongs to the user
		if userID != post.Author.ID {
			h.log.Error("forbidden: you do not have permission to change status of this post")
			helpers.ForbiddenError(h.log, w)
			return
		}
	}
	// Update post
	postResponse, err := h.postService.ReturnToOwner(r.Context(), postID)
	if err != nil {
		helpers.HandleServiceError(h.log, w, fmt.Errorf("failed to change post status: %w", err))
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"post": postResponse,
	})
}
