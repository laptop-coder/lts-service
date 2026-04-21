package service

import (
	"fmt"
	"backend/internal/model"
	"backend/internal/repository"
	"backend/pkg/apperrors"
	"backend/pkg/logger"
	"context"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type ConversationService interface {
	CreateOrGet(ctx context.Context, postID uuid.UUID, requesterID uuid.UUID) (*ConversationResponseDTO, error)
	SendMessage(ctx context.Context, conversationID uuid.UUID, senderID uuid.UUID, content *string) (*MessageResponseDTO, error)
	GetConversation(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID) (*ConversationResponseDTO, error)
	GetUserConversations(ctx context.Context, userID uuid.UUID, limit, offset int) ([]ConversationListItemDTO, error)
	MarkAsRead(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID) error
	GetTotalUnreadCount(ctx context.Context, userID uuid.UUID) (int64, error)
}

type conversationService struct {
	convRepo     repository.ConversationRepository
	msgRepo      repository.MessageRepository
	postRepo     repository.PostRepository
	userRepo     repository.UserRepository
	emailService EmailService
	db           *gorm.DB
	log          logger.Logger
}

func NewConversationService(
	convRepo repository.ConversationRepository,
	msgRepo repository.MessageRepository,
	postRepo repository.PostRepository,
	userRepo repository.UserRepository,
	emailService EmailService,
	db *gorm.DB,
	log logger.Logger,
) ConversationService {
	return &conversationService{
		convRepo:     convRepo,
		msgRepo:      msgRepo,
		postRepo:     postRepo,
		userRepo:     userRepo,
		emailService: emailService,
		db:           db,
		log:          log,
	}
}

type MessageResponseDTO struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt string    `json:"createdAt"`
	UpdatedAt string    `json:"updatedAt"`
	SenderID  uuid.UUID `json:"senderId"`
	Content   string    `json:"content"`
	IsRead    bool      `json:"isRead"`
}

type ConversationResponseDTO struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt string    `json:"createdAt"`
	UpdatedAt string    `json:"updatedAt"`

	Post     PostResponseDTO      `json:"post"`
	Messages []MessageResponseDTO `json:"messages"`

	OtherUser UserResponseDTO `json:"otherUser"`
}

type ConversationListItemDTO struct {
	ID          uuid.UUID       `json:"id"`
	UpdatedAt   string          `json:"updatedAt"`
	PostID      uuid.UUID       `json:"postId"`
	PostName    string          `json:"postName"`
	UnreadCount int64           `json:"unreadCount"`
	LastMessage string          `json:"lastMessage"`
	OtherUser   UserResponseDTO `json:"otherUser"`
}

func ConversationToDTO(conversation *model.Conversation, otherUser *model.User) *ConversationResponseDTO {
	var messages []MessageResponseDTO
	for _, message := range conversation.Messages {
		messages = append(messages, *MessageToDTO(&message))
	}
	return &ConversationResponseDTO{
		ID:        conversation.ID,
		CreatedAt: conversation.CreatedAt.Format(time.RFC3339),
		UpdatedAt: conversation.UpdatedAt.Format(time.RFC3339),
		Post:      *PostToDTO(&conversation.Post),
		Messages:  messages,
		OtherUser: *UserToDTO(otherUser),
	}
}

func ConversationToListItemDTO(conversation *model.Conversation, otherUser *model.User, unreadCount int64, lastMessage string) *ConversationListItemDTO {
	return &ConversationListItemDTO{
		ID:          conversation.ID,
		UpdatedAt:   conversation.UpdatedAt.Format(time.RFC3339),
		PostID:      conversation.Post.ID,
		PostName:    conversation.Post.Name,
		UnreadCount: unreadCount,
		LastMessage: lastMessage,
		OtherUser:   *UserToDTO(otherUser),
	}
}

func MessageToDTO(message *model.Message) *MessageResponseDTO {
	return &MessageResponseDTO{
		ID:        message.ID,
		CreatedAt: message.CreatedAt.Format(time.RFC3339),
		UpdatedAt: message.UpdatedAt.Format(time.RFC3339),
		SenderID:  message.SenderID,
		Content:   message.Content,
		IsRead:    message.IsRead,
	}
}

func (s *conversationService) CreateOrGet(ctx context.Context, postID uuid.UUID, requesterID uuid.UUID) (*ConversationResponseDTO, error) {
	// Get post
	post, err := s.postRepo.FindByID(ctx, &postID)
	if err != nil || post == nil {
		return nil, fmt.Errorf("failed to get post by id: %w", err)
	}

	// Check if author is the requester
	if post.AuthorID == requesterID {
		return nil, fmt.Errorf("author of the post cannot be the requester: %w", apperrors.ErrCannotContactOwnPost)
	}

	// Get other user
	otherUser, err := s.userRepo.FindByID(ctx, &post.AuthorID)
	if err != nil || otherUser == nil {
		return nil, fmt.Errorf("failed to get other user: %w", err)
	}

	// Find existing conversation
	existingConv, err := s.convRepo.FindByPostAndUsers(ctx, postID, post.AuthorID, requesterID)
	if err != nil && !errors.Is(err, apperrors.ErrConversationNotFound) {
		return nil, err
	}

	if existingConv != nil {
		return ConversationToDTO(existingConv, otherUser), nil
	}

	// Get own user (i.e. requester)
	requester, err := s.userRepo.FindByID(ctx, &requesterID)
	if err != nil || requester == nil {
		return nil, fmt.Errorf("failed to get own user (requester): %w", err)
	}

	// Create new conversation
	conv := &model.Conversation{
		PostID:      postID,
		Post:        *post,
		AuthorID:    post.AuthorID,
		Author:      *otherUser,
		RequesterID: requesterID,
		Requester:   *requester,
		IsActive:    true,
	}

	if err := s.convRepo.Create(ctx, conv); err != nil {
		return nil, fmt.Errorf("failed to create conversation: %w", err)
	}

	return ConversationToDTO(conv, otherUser), nil
}

func (s *conversationService) SendMessage(ctx context.Context, conversationID uuid.UUID, senderID uuid.UUID, content *string) (*MessageResponseDTO, error) {
	if content == nil {
		return nil, fmt.Errorf("content cannot be empty: %w", apperrors.ErrEmptyMessage)
	}

	// Get conversation
	conv, err := s.convRepo.FindByID(ctx, conversationID)
	if err != nil {
		return nil, fmt.Errorf("failed to find conversation by id: %w", err)
	}

	// Check if sender is a conversation participant
	if conv.AuthorID != senderID && conv.RequesterID != senderID {
		return nil, fmt.Errorf("the user is not the conversation participant: %w", apperrors.ErrNotConversationParticipant)
	}

	// Create message
	msg := &model.Message{
		ConversationID: conversationID,
		SenderID:       senderID,
		Content:        *content,
		IsRead:         false,
	}

	if err := s.msgRepo.Create(ctx, msg); err != nil {
		return nil, fmt.Errorf("failed to create message: %w", err)
	}

	// Update UpdatedAt field
	s.convRepo.Update(ctx, conv)

	// Send email to the second participant
	errCh := make(chan error, 1)
	go func() {
		errCh <- s.notifyParticipant(ctx, conv, senderID, content)
	}()
	if err := <-errCh; err != nil {
		s.log.Error("Email (new message notification) was not sent", "error", err.Error())
	}

	return MessageToDTO(msg), nil
}

func (s *conversationService) GetConversation(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID) (*ConversationResponseDTO, error) {
	// Get conversation
	conv, err := s.convRepo.FindByID(ctx, conversationID)
	if err != nil {
		return nil, fmt.Errorf("failed to find conversation by id: %w", err)
	}

	// Check if the user is the conversation participant
	if conv.AuthorID != userID && conv.RequesterID != userID {
		return nil, fmt.Errorf("the user is not the conversation participant: %w", apperrors.ErrNotConversationParticipant)
	}

	// Identify second conversation participant
	var otherUser *model.User
	if conv.AuthorID == userID {
		otherUser = &conv.Requester
	} else {
		otherUser = &conv.Post.Author
	}

	return ConversationToDTO(conv, otherUser), nil
}

func (s *conversationService) GetUserConversations(ctx context.Context, userID uuid.UUID, limit, offset int) ([]ConversationListItemDTO, error) {
	conversations, err := s.convRepo.FindByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversations by user id: %w", err)
	}

	items := make([]ConversationListItemDTO, len(conversations))

	for i, conv := range conversations {
		// Identify second participant
		var otherUser model.User
		if conv.AuthorID == userID {
			otherUser = conv.Requester
		} else {
			otherUser = conv.Post.Author
		}

		// Last message
		lastMsg, err := s.msgRepo.FindLastMessage(ctx, conv.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to find last message: %w", err)
		}

		// Count of unread messages
		unreadCount, err := s.msgRepo.CountUnread(ctx, conv.ID, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to count unread messages: %w", err)
		}

		items[i] = *ConversationToListItemDTO(&conv, &otherUser, unreadCount, lastMsg.Content)
	}

	return items, nil
}

func (s *conversationService) MarkAsRead(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID) error {
	err := s.msgRepo.MarkAsRead(ctx, conversationID, userID)
	if err != nil {
		return fmt.Errorf("failed to mark all messages in conversation as read: %w", err)
	}
	return nil
}

func (s *conversationService) notifyParticipant(ctx context.Context, conv *model.Conversation, senderID uuid.UUID, content *string) error {
	if content == nil {
		return fmt.Errorf("content cannot be empty: %w", apperrors.ErrEmptyMessage)
	}

	// Identify recipient
	var recipientID uuid.UUID
	if conv.AuthorID == senderID {
		recipientID = conv.RequesterID
	} else {
		recipientID = conv.AuthorID
	}

	// Get recipient email
	recipient, err := s.userRepo.FindByID(ctx, &recipientID)
	if err != nil || recipient == nil {
		return fmt.Errorf("failed to get recipient: %w", err)
	}

	// Get sender
	sender, err := s.userRepo.FindByID(ctx, &senderID)
	if err != nil || sender == nil {
		return fmt.Errorf("failed to get sender: %w", err)
	}

	dto := NewMessageNotificationDTO{
		Post:           conv.Post,
		Recipient:      *recipient,
		Sender:         *sender,
		Message:        *content,
		ConversationID: conv.ID,
	}

	err = s.emailService.SendNewMessageNotification(ctx, &dto)
	if err != nil {
		return fmt.Errorf("failed to send notification about new message: %w", err)
	}
	return nil
}


func (s *conversationService) GetTotalUnreadCount(ctx context.Context, userID uuid.UUID) (int64, error) {
	return s.msgRepo.CountAllUnread(ctx, userID)
}
