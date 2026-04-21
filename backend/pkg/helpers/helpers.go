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

	// HTTP 400
	case errors.Is(err, apperrors.ErrCannotContactOwnPost):
		ErrorResponse(log, w, "Нельзя написать самому себе", http.StatusBadRequest)
	case errors.Is(err, apperrors.ErrRequiredField):
		ErrorResponse(log, w, "Поле обязательно для заполнения", http.StatusBadRequest)
	case errors.Is(err, apperrors.ErrEmptyMessage):
		ErrorResponse(log, w, "Сообщение не может быть пустым", http.StatusBadRequest)
	case errors.Is(err, apperrors.ErrEmptyEmail):
		ErrorResponse(log, w, "Email не может быть пустым", http.StatusBadRequest)
	case errors.Is(err, apperrors.ErrPasswordTooShort):
		ErrorResponse(log, w, "Пароль слишком короткий", http.StatusBadRequest)
	case errors.Is(err, apperrors.ErrPasswordTooLong):
		ErrorResponse(log, w, "Пароль слишком длинный", http.StatusBadRequest)
	case errors.Is(err, apperrors.ErrValueTooShort):
		ErrorResponse(log, w, "Введённая строка слишком короткая", http.StatusBadRequest)
	case errors.Is(err, apperrors.ErrValueTooLong):
		ErrorResponse(log, w, "Введённая строка слишком длинная", http.StatusBadRequest)
	case errors.Is(err, apperrors.ErrFileTooLarge):
		ErrorResponse(log, w, "Файл слишком большой", http.StatusBadRequest)
	case errors.Is(err, apperrors.ErrInvalidFileType):
		ErrorResponse(log, w, "Неверный формат файла", http.StatusBadRequest)
	case errors.Is(err, apperrors.ErrDTOValidation):
		ErrorResponse(log, w, "Неверный запрос. Проверьте введённые данные", http.StatusBadRequest)

	// HTTP 401
	case errors.Is(err, apperrors.ErrInvalidCredentials):
		ErrorResponse(log, w, "Неверный email или пароль", http.StatusUnauthorized)
	case errors.Is(err, apperrors.ErrUnauthorized),
		errors.Is(err, apperrors.ErrInvalidToken),
		errors.Is(err, apperrors.ErrTokenRevoked),
		errors.Is(err, apperrors.ErrTokenExpired):
		ErrorResponse(log, w, "Требуется вход в систему", http.StatusUnauthorized)

	// HTTP 403
	case errors.Is(err, apperrors.ErrForbidden),
		errors.Is(err, apperrors.ErrNotConversationParticipant):
		ErrorResponse(log, w, "Доступ запрещён", http.StatusForbidden)

	// HTTP 404
	case errors.Is(err, apperrors.ErrPostNotFound):
		ErrorResponse(log, w, "Объявление не найдено", http.StatusNotFound)
	case errors.Is(err, apperrors.ErrUserNotFound):
		ErrorResponse(log, w, "Пользователь не найден", http.StatusNotFound)
	case errors.Is(err, apperrors.ErrConversationNotFound):
		ErrorResponse(log, w, "Переписка не найдена", http.StatusNotFound)
	case errors.Is(err, apperrors.ErrNotFound):
		ErrorResponse(log, w, "Запрошенный ресурс не найден", http.StatusNotFound)

		// HTTP 409
	case errors.Is(err, apperrors.ErrSubjectAlreadyExists):
		ErrorResponse(log, w, "Предмет уже существует. Используйте другие данные", http.StatusConflict)
	case errors.Is(err, apperrors.ErrRoomAlreadyExists):
		ErrorResponse(log, w, "Комната уже существует. Используйте другие данные", http.StatusConflict)
	case errors.Is(err, apperrors.ErrUserWithThisEmailAlreadyExists):
		ErrorResponse(log, w, "Пользователь с таким email уже существует", http.StatusConflict)
	case errors.Is(err, apperrors.ErrInstitutionAdministratorPositionAlreadyExists):
		ErrorResponse(log, w, "Должность администрации ОУ уже существует. Используйте другие данные", http.StatusConflict)
	case errors.Is(err, apperrors.ErrStaffPositionAlreadyExists):
		ErrorResponse(log, w, "Должность сотрудника ОУ уже существует. Используйте другие данные", http.StatusConflict)
	case errors.Is(err, apperrors.ErrStudentGroupAlreadyExists):
		ErrorResponse(log, w, "Учебная группа/класс уже существует. Используйте другие данные", http.StatusConflict)
	case errors.Is(err, apperrors.ErrStudentGroupAlreadyHasAdvisor):
		ErrorResponse(log, w, "Учебная группа/класс уже имеет наставника/классного руководителя", http.StatusConflict)
	case errors.Is(err, apperrors.ErrRoomAlreadyHasTeacherAssignedToIt):
		ErrorResponse(log, w, "Комната уже назначена другому учителю", http.StatusConflict)

	// HTTP 500
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
