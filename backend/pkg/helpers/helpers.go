package helpers

import (
	"backend/pkg/apperrors"
	"backend/pkg/logger"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func JsonResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	// Set JSON content type
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// Convert data to JSON format
	encodedData, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to encode response to JSON",
		}); err != nil {
			panic(err)
		}
		return
	}
	// Status code
	w.WriteHeader(statusCode)
	// Response
	w.Write(encodedData)
}

func ErrorResponse(log logger.Logger, w http.ResponseWriter, message string, statusCode int) {
	log.Error(message, "status_code", statusCode)
	JsonResponse(w, map[string]string{
		"error": message,
	}, statusCode)
}

func SuccessResponse(w http.ResponseWriter, data interface{}) {
	JsonResponse(w, data, http.StatusOK)
}

func HandleServiceError(log logger.Logger, w http.ResponseWriter, err error) {
	switch {

	case errors.Is(err, apperrors.ErrCannotContactOwnPost),
		errors.Is(err, apperrors.ErrDTOValidation),
		errors.Is(err, apperrors.ErrRequiredField),
		errors.Is(err, apperrors.ErrEmptyMessage),
		errors.Is(err, apperrors.ErrEmptyEmail),
		errors.Is(err, apperrors.ErrPasswordTooShort),
		errors.Is(err, apperrors.ErrPasswordTooLong),
		errors.Is(err, apperrors.ErrValueTooLong),
		errors.Is(err, apperrors.ErrValueTooShort),
		errors.Is(err, apperrors.ErrFileTooLarge),
		errors.Is(err, apperrors.ErrInvalidFileType):
		ErrorResponse(log, w, "Неверный запрос. Проверьте введённые данные", http.StatusBadRequest)

	case errors.Is(err, apperrors.ErrUnauthorized),
		errors.Is(err, apperrors.ErrInvalidCredentials),
		errors.Is(err, apperrors.ErrInvalidToken),
		errors.Is(err, apperrors.ErrTokenRevoked),
		errors.Is(err, apperrors.ErrTokenExpired):
		ErrorResponse(log, w, "Требуется вход в систему", http.StatusUnauthorized)

	case errors.Is(err, apperrors.ErrForbidden),
		errors.Is(err, apperrors.ErrNotConversationParticipant):
		ErrorResponse(log, w, "Доступ запрещён", http.StatusForbidden)

	case errors.Is(err, apperrors.ErrNotFound),
		errors.Is(err, apperrors.ErrPostNotFound),
		errors.Is(err, apperrors.ErrUserNotFound),
		errors.Is(err, apperrors.ErrConversationNotFound):
		ErrorResponse(log, w, "Запрошенный ресурс не найден", http.StatusNotFound)

	case errors.Is(err, apperrors.ErrSubjectAlreadyExists),
		errors.Is(err, apperrors.ErrRoomAlreadyExists),
		errors.Is(err, apperrors.ErrUserWithThisEmailAlreadyExists),
		errors.Is(err, apperrors.ErrInstitutionAdministratorPositionAlreadyExists),
		errors.Is(err, apperrors.ErrStaffPositionAlreadyExists),
		errors.Is(err, apperrors.ErrStudentGroupAlreadyExists),
		errors.Is(err, apperrors.ErrStudentGroupAlreadyHasAdvisor),
		errors.Is(err, apperrors.ErrRoomAlreadyHasTeacherAssignedToIt):
		ErrorResponse(log, w, "Ресурс уже существует. Используйте другие данные", http.StatusConflict)

	default:
		ErrorResponse(log, w, "Внутренняя ошибка сервера", http.StatusInternalServerError)

	}
}

func GetCookie(cookieKey string, r *http.Request) (string, error) {
	cookie, err := r.Cookie(cookieKey)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

// HTTP 400
func BadRequestError(log logger.Logger, w http.ResponseWriter) {
	ErrorResponse(log, w, "Неверный формат данных", http.StatusBadRequest)
}

func BadRequestFieldError(log logger.Logger, w http.ResponseWriter, field string) {
	ErrorResponse(log, w, fmt.Sprintf("Неверный формат поля %s", field), http.StatusBadRequest)
}

func FieldRequiredError(log logger.Logger, w http.ResponseWriter, field string) {
	ErrorResponse(log, w, fmt.Sprintf("Поле %s обязательно для заполнения", field), http.StatusBadRequest)
}

func TooManyFieldsError(log logger.Logger, w http.ResponseWriter, field string) {
	ErrorResponse(log, w, fmt.Sprintf("Слишком много значений %s", field), http.StatusBadRequest)
}

func FieldExactlyOneError(log logger.Logger, w http.ResponseWriter, field string) {
	ErrorResponse(log, w, fmt.Sprintf("Поле %s должно быть указано ровно один раз", field), http.StatusBadRequest)
}

func AtLeastOneFieldError(log logger.Logger, w http.ResponseWriter, field string) {
	ErrorResponse(log, w, fmt.Sprintf("Хотя бы одно поле %s должно быть указано", field), http.StatusBadRequest)
}

// HTTP 401
func UnauthorizedError(log logger.Logger, w http.ResponseWriter) {
	ErrorResponse(log, w, "Требуется вход в систему", http.StatusUnauthorized)
}

// HTTP 403
func ForbiddenError(log logger.Logger, w http.ResponseWriter) {
	ErrorResponse(log, w, "Доступ запрещён", http.StatusForbidden)
}

// HTTP 405
func MethodNotAllowedError(log logger.Logger, w http.ResponseWriter) {
	ErrorResponse(log, w, "Метод не поддерживается", http.StatusMethodNotAllowed)
}

// HTTP 500
func InternalError(log logger.Logger, w http.ResponseWriter) {
	ErrorResponse(log, w, "Внутренняя ошибка сервера. Попробуйте позже", http.StatusInternalServerError)
}
