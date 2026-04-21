package handler

import (
	"backend/internal/service"
	"backend/pkg/helpers"
	"backend/pkg/logger"
	"backend/pkg/middleware"
	"github.com/google/uuid"
	"net/http"
	"strings"
)

type ConversationHandler struct {
	conversationService service.ConversationService
	log                 logger.Logger
}

func NewConversationHandler(conversationService service.ConversationService, log logger.Logger) *ConversationHandler {
	return &ConversationHandler{
		conversationService: conversationService,
		log:                 log,
	}
}

// Create conversation and send first message
func (h *ConversationHandler) CreateConversation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helpers.MethodNotAllowedError(h.log, w)
		return
	}

	// Get post ID from URL
	postIDStr := r.PathValue("postId")
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		h.log.Error("invalid post id", "postId", postIDStr)
		helpers.BadRequestFieldError(h.log, w, "postId")
		return
	}

	// Get requesterID from the context
	requesterID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		h.log.Error("failed to get userID from context and convert it to UUID")
		helpers.InternalError(h.log, w)
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

	// Get message
	message := ""
	if messageFields := r.PostForm["message"]; len(messageFields) != 1 {
		h.log.Error("failed to parse form: message field must be specified exactly once")
		helpers.FieldExactlyOneError(h.log, w, "message")
		return
	} else {
		message = messageFields[0]
		if strings.TrimSpace(message) == "" {
			h.log.Error("failed to parse form: message cannot be empty or only whitespace")
			helpers.FieldRequiredError(h.log, w, "message")
			return
		}
	}

	// Create or get conversation
	conv, err := h.conversationService.CreateOrGet(r.Context(), postID, requesterID)
	if err != nil || conv == nil {
		h.log.Error("Failed to create or get conversation",
			"postId", postID,
			"requesterId", requesterID,
			"error", err.Error())

		helpers.HandleServiceError(h.log, w, err)
		return
	}

	// Send message
	msg, err := h.conversationService.SendMessage(r.Context(), conv.ID, requesterID, &message)
	if err != nil || msg == nil {
		h.log.Error("Failed to send message",
			"conversationId", conv.ID,
			"senderId", requesterID,
			"error", err.Error())

		helpers.HandleServiceError(h.log, w, err)
		return
	}

	helpers.SuccessResponse(w, map[string]interface{}{
		"conversationId": conv.ID,
		"messageId":      msg.ID,
	})
}

// Send message to the existing conversation
func (h *ConversationHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helpers.MethodNotAllowedError(h.log, w)
		return
	}

	// Get conversationId from URL
	convIDStr := r.PathValue("conversationId")
	convID, err := uuid.Parse(convIDStr)
	if err != nil {
		h.log.Error("invalid conversation id", "conversationId", convIDStr)
		helpers.BadRequestFieldError(h.log, w, "conversationId")
		return
	}

	// Get senderID from the context
	senderID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		h.log.Error("failed to get userID from context and convert it to UUID")
		helpers.InternalError(h.log, w)
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

	// Get message
	message := ""
	if messageFields := r.PostForm["message"]; len(messageFields) != 1 {
		h.log.Error("failed to parse form: message field must be specified exactly once")
		helpers.FieldExactlyOneError(h.log, w, "message")
		return
	} else {
		message = messageFields[0]
		if strings.TrimSpace(message) == "" {
			h.log.Error("failed to parse form: message cannot be empty or only whitespace")
			helpers.FieldRequiredError(h.log, w, "message")
			return
		}
	}

	// Send message
	msg, err := h.conversationService.SendMessage(r.Context(), convID, senderID, &message)
	if err != nil || msg == nil {
		h.log.Error("Failed to send message",
			"conversationId", convID,
			"senderId", senderID,
			"error", err.Error())
		helpers.HandleServiceError(h.log, w, err)
		return
	}

	helpers.SuccessResponse(w, map[string]interface{}{
		"message": *msg,
	})
}

// Get conversation with messages
func (h *ConversationHandler) GetConversation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		helpers.MethodNotAllowedError(h.log, w)
		return
	}

	// Get conversationId from URL
	convIDStr := r.PathValue("conversationId")
	convID, err := uuid.Parse(convIDStr)
	if err != nil {
		h.log.Error("invalid conversation id", "conversationId", convIDStr)
		helpers.BadRequestFieldError(h.log, w, "conversationId")
		return
	}

	// Get userID from the context
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		h.log.Error("failed to get userID from context and convert it to UUID")
		helpers.InternalError(h.log, w)
		return
	}

	// Get conversation
	conv, err := h.conversationService.GetConversation(r.Context(), convID, userID)
	if err != nil || conv == nil {
		h.log.Error("Failed to get conversation",
			"conversationId", convID,
			"userId", userID,
			"error", err.Error())

		helpers.HandleServiceError(h.log, w, err)
		return
	}

	// Mark messages as read
	if err := h.conversationService.MarkAsRead(r.Context(), convID, userID); err != nil {
		helpers.HandleServiceError(h.log, w, err)
		return
	}

	helpers.SuccessResponse(w, map[string]interface{}{
		"conversation": *conv,
	})
}

// Get the list of conversations of the current user
func (h *ConversationHandler) GetMyConversations(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		helpers.MethodNotAllowedError(h.log, w)
		return
	}

	// Get userID from the context
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		h.log.Error("failed to get userID from context and convert it to UUID")
		helpers.InternalError(h.log, w)
		return
	}

	// Get the list of conversations
	conversations, err := h.conversationService.GetUserConversations(r.Context(), userID, 0, 0)
	if err != nil {
		helpers.HandleServiceError(h.log, w, err)
		return
	}

	helpers.SuccessResponse(w, map[string]interface{}{
		"conversations": conversations,
	})
}

// Mark all messages in conversation as read
func (h *ConversationHandler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		helpers.MethodNotAllowedError(h.log, w)
		return
	}

	// Get conversationId from URL
	convIDStr := r.PathValue("conversationId")
	convID, err := uuid.Parse(convIDStr)
	if err != nil {
		h.log.Error("invalid conversation id", "conversationId", convIDStr)
		helpers.BadRequestFieldError(h.log, w, "conversationId")
		return
	}

	// Get userID from the context
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		h.log.Error("failed to get userID from context and convert it to UUID")
		helpers.InternalError(h.log, w)
		return
	}

	// Mark as read
	if err := h.conversationService.MarkAsRead(r.Context(), convID, userID); err != nil {
		helpers.HandleServiceError(h.log, w, err)
		return
	}

	helpers.JsonResponse(w, map[string]interface{}{}, http.StatusNoContent)
}

// Get total count of unread messages
func (h *ConversationHandler) GetTotalUnreadCount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		helpers.MethodNotAllowedError(h.log, w)
		return
	}

	// Get userID from the context
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		h.log.Error("failed to get userID from context and convert it to UUID")
		helpers.InternalError(h.log, w)
		return
	}

	// Get count
	count, err := h.conversationService.GetTotalUnreadCount(r.Context(), userID)
	if err != nil {
		helpers.HandleServiceError(h.log, w, err)
		return
	}

	helpers.SuccessResponse(w, map[string]interface{}{
		"unreadCount": count,
	})
}
