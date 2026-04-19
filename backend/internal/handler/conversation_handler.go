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
		helpers.ErrorResponse(h.log, w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get post ID from URL
	postIDStr := r.PathValue("postId")
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		h.log.Error("Invalid post id", "postId", postIDStr)
		helpers.ErrorResponse(h.log, w, "invalid post id", http.StatusBadRequest)
		return
	}

	// Get requesterID from the context
	requesterID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		helpers.ErrorResponse(h.log, w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Restrictions
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MB
	// Parse form
	if err := r.ParseForm(); err != nil {
		helpers.ErrorResponse(h.log, w, "failed to parse x-www-form-urlencoded form", http.StatusBadRequest)
		return
	}

	// Get message
	message := ""
	if messageFields := r.PostForm["message"]; len(messageFields) != 1 {
		helpers.ErrorResponse(h.log, w, "failed to parse form: message field must be specified exactly once", http.StatusBadRequest)
		return
	} else {
		message = messageFields[0]
		if strings.TrimSpace(message) == "" {
			helpers.ErrorResponse(h.log, w, "failed to parse form: message cannot be empty or only whitespace", http.StatusBadRequest)
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
		helpers.ErrorResponse(h.log, w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get conversationId from URL
	convIDStr := r.PathValue("conversationId")
	convID, err := uuid.Parse(convIDStr)
	if err != nil {
		h.log.Error("Invalid conversation id", "conversationId", convIDStr)
		helpers.ErrorResponse(h.log, w, "invalid conversation id", http.StatusBadRequest)
		return
	}

	// Get senderID from the context
	senderID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		helpers.ErrorResponse(h.log, w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Restrictions
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MB
	// Parse form
	if err := r.ParseForm(); err != nil {
		helpers.ErrorResponse(h.log, w, "failed to parse x-www-form-urlencoded form", http.StatusBadRequest)
		return
	}

	// Get message
	message := ""
	if messageFields := r.PostForm["message"]; len(messageFields) != 1 {
		helpers.ErrorResponse(h.log, w, "failed to parse form: message field must be specified exactly once", http.StatusBadRequest)
		return
	} else {
		message = messageFields[0]
		if strings.TrimSpace(message) == "" {
			helpers.ErrorResponse(h.log, w, "failed to parse form: message cannot be empty or only whitespace", http.StatusBadRequest)
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
		helpers.ErrorResponse(h.log, w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get conversationId from URL
	convIDStr := r.PathValue("conversationId")
	convID, err := uuid.Parse(convIDStr)
	if err != nil {
		h.log.Error("Invalid conversation id", "conversationId", convIDStr)
		helpers.ErrorResponse(h.log, w, "invalid conversation id", http.StatusBadRequest)
		return
	}

	// Get userID from the context
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		helpers.ErrorResponse(h.log, w, "unauthorized", http.StatusUnauthorized)
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
	}

	helpers.SuccessResponse(w, map[string]interface{}{
		"conversation": *conv,
	})
}

// Get the list of conversations of the current user
func (h *ConversationHandler) GetMyConversations(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		helpers.ErrorResponse(h.log, w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get userID from the context
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		helpers.ErrorResponse(h.log, w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Get the list of conversations
	conversations, err := h.conversationService.GetUserConversations(r.Context(), userID, 0, 0)
	if err != nil {
		h.log.Error("Failed to get user conversations",
			"userId", userID,
			"error", err.Error())

		helpers.ErrorResponse(h.log, w, "failed to get user conversations", http.StatusInternalServerError)
		return
	}

	helpers.SuccessResponse(w, map[string]interface{}{
		"conversations": conversations,
	})
}

// Mark all messages in conversation as read
func (h *ConversationHandler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		helpers.ErrorResponse(h.log, w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get conversationId from URL
	convIDStr := r.PathValue("conversationId")
	convID, err := uuid.Parse(convIDStr)
	if err != nil {
		h.log.Error("Invalid conversation id", "conversationId", convIDStr)
		helpers.ErrorResponse(h.log, w, "invalid conversation id", http.StatusBadRequest)
		return
	}

	// Get userID from the context
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		helpers.ErrorResponse(h.log, w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Mark as read
	if err := h.conversationService.MarkAsRead(r.Context(), convID, userID); err != nil {
		h.log.Error("Failed to mark messages as read",
			"conversationId", convID,
			"userId", userID,
			"error", err.Error())

		helpers.ErrorResponse(h.log, w, "failed to mark messages as read", http.StatusInternalServerError)
		return
	}

	helpers.JsonResponse(w, map[string]interface{}{}, http.StatusNoContent)
}

// Get total count of unread messages
func (h *ConversationHandler) GetTotalUnreadCount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		helpers.ErrorResponse(h.log, w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get userID from the context
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		helpers.ErrorResponse(h.log, w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Get count
	count, err := h.conversationService.GetTotalUnreadCount(r.Context(), userID)
	if err != nil {
		h.log.Error("Failed to get count of all unread messages",
			"userId", userID,
			"error", err.Error())

		helpers.ErrorResponse(h.log, w, "failed to get count of all unread messages", http.StatusInternalServerError)
		return
	}

	helpers.SuccessResponse(w, map[string]interface{}{
		"unreadCount": count,
	})
}
