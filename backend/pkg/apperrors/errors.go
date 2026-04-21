package apperrors

import (
	"errors"
)

// HTTP 400
var (
	ErrCannotContactOwnPost = errors.New("cannot contact your own post")

	ErrDTOValidation = errors.New("error during dto validation")

	ErrRequiredField = errors.New("required field")
	ErrEmptyMessage  = errors.New("message cannot be empty")
	ErrEmptyEmail    = errors.New("email cannot be empty")

	ErrPasswordTooShort = errors.New("password too short")
	ErrPasswordTooLong  = errors.New("password too long")
	ErrValueTooLong     = errors.New("value too long")
	ErrValueTooShort    = errors.New("value too short")
	ErrFileTooLarge     = errors.New("file too large")
	ErrInvalidFileType  = errors.New("invalid file type")
)

// HTTP 401
var (
	ErrUnauthorized       = errors.New("unauthorized")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
	ErrTokenRevoked       = errors.New("token revoked")
	ErrTokenExpired       = errors.New("token expired")
)

// HTTP 403
var (
	ErrForbidden                  = errors.New("forbidden")
	ErrNotConversationParticipant = errors.New("user is not a participant of this conversation")
)

// HTTP 404
var (
	ErrNotFound = errors.New("resource not found")

	ErrPostNotFound         = errors.New("post not found")
	ErrUserNotFound         = errors.New("user not found")
	ErrConversationNotFound = errors.New("conversation not found")
)

// HTTP 409
var (
	ErrSubjectAlreadyExists                          = errors.New("subject already exists")
	ErrRoomAlreadyExists                             = errors.New("room already exists")
	ErrUserWithThisEmailAlreadyExists                = errors.New("user with this email already exists")
	ErrInstitutionAdministratorPositionAlreadyExists = errors.New("institution administrator position already exists")
	ErrStaffPositionAlreadyExists                    = errors.New("staff position already exists")
	ErrStudentGroupAlreadyExists                     = errors.New("student group already exists")
	ErrStudentGroupAlreadyHasAdvisor                 = errors.New("student group already has an advisor")
	ErrRoomAlreadyHasTeacherAssignedToIt             = errors.New("room already has teacher assigned to it")
)
