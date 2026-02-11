package services

import (
	"github.com/google/uuid"
	"github.com/sidji-omnichannel/internal/domain/ports/repository"
	"github.com/sidji-omnichannel/internal/models"
)

// MessageService handles message operations
type MessageService struct {
	messageRepo      repository.MessageRepository
	conversationRepo repository.ConversationRepository
	userRepo         repository.UserRepository
}

// NewMessageService creates a new message service
func NewMessageService(
	mr repository.MessageRepository,
	cr repository.ConversationRepository,
	ur repository.UserRepository,
) *MessageService {
	return &MessageService{
		messageRepo:      mr,
		conversationRepo: cr,
		userRepo:         ur,
	}
}

// List returns messages for a conversation
func (s *MessageService) List(convID uuid.UUID, filter *models.MessageFilter) (*models.MessageList, error) {
	return s.messageRepo.List(convID, filter)
}

// Create creates a new message
func (s *MessageService) Create(msg *models.Message) error {
	err := s.messageRepo.Create(msg)
	if err != nil {
		return err
	}

	// Update conversation's last_message_at
	_ = s.conversationRepo.UpdateLastMessage(msg.ConversationID, msg.CreatedAt)

	// If sender is AGENT, update last_human_reply_at
	if msg.SenderType == models.SenderAgent {
		_ = s.conversationRepo.UpdateLastHumanReply(msg.ConversationID, msg.CreatedAt)
	}

	return nil
}

// UpdateStatus updates message status (sent, delivered, read) by ExternalID
func (s *MessageService) UpdateStatus(externalID string, status models.MessageStatus) error {
	return s.messageRepo.UpdateStatus(externalID, status)
}

// UpdateStatusByID updates message status by internal UUID
func (s *MessageService) UpdateStatusByID(id uuid.UUID, status models.MessageStatus) error {
	return s.messageRepo.UpdateStatusByID(id, status)
}

// MarkAsRead marks all contact messages in a conversation as read
func (s *MessageService) MarkAsRead(convID uuid.UUID) error {
	return s.messageRepo.MarkAsRead(convID)
}

// CountUnread counts unread messages from contact in a conversation
func (s *MessageService) CountUnread(convID uuid.UUID) (int, error) {
	return s.messageRepo.CountUnread(convID)
}

// GetByExternalID finds a message by its external ID (WhatsApp/Instagram message ID)
func (s *MessageService) GetByExternalID(externalID string) (*models.Message, error) {
	return s.messageRepo.GetByExternalID(externalID)
}

// GetUserPublic fetches public info for a user
func (s *MessageService) GetUserPublic(userID uuid.UUID) (*models.UserPublic, error) {
	return s.userRepo.GetPublic(userID)
}
