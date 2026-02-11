package services

import (
	"errors"

	"github.com/google/uuid"
	"github.com/sidji-omnichannel/internal/domain/ports/repository"
	"github.com/sidji-omnichannel/internal/models"
)

var (
	ErrConversationNotFound = errors.New("conversation not found")
)

// ConversationService handles conversation operations
type ConversationService struct {
	repo repository.ConversationRepository
}

// NewConversationService creates a new conversation service
func NewConversationService(repo repository.ConversationRepository) *ConversationService {
	return &ConversationService{repo: repo}
}

// List returns conversations for an organization with filters
func (s *ConversationService) List(orgID uuid.UUID, filter *models.ConversationFilter) ([]*models.ConversationListItem, int, error) {
	return s.repo.List(orgID, filter)
}

// GetByID retrieves a single conversation
func (s *ConversationService) GetByID(orgID, convID uuid.UUID) (*models.Conversation, error) {
	conv, err := s.repo.GetByID(orgID, convID)
	if err == repository.ErrNotFound || conv == nil {
		return nil, ErrConversationNotFound
	}
	if err != nil {
		return nil, err
	}
	return conv, nil
}

// Assign assigns a conversation to an agent
func (s *ConversationService) Assign(orgID, convID, userID uuid.UUID) error {
	err := s.repo.Assign(orgID, convID, userID)
	if err == repository.ErrNotFound {
		return ErrConversationNotFound
	}
	return err
}

// UpdateStatus updates conversation status
func (s *ConversationService) UpdateStatus(orgID, convID uuid.UUID, status models.ConversationStatus) error {
	err := s.repo.UpdateStatus(orgID, convID, status)
	if err == repository.ErrNotFound {
		return ErrConversationNotFound
	}
	return err
}

// FindOrCreate finds or creates a conversation for a contact on a channel
func (s *ConversationService) FindOrCreate(orgID, channelID, contactID uuid.UUID) (*models.Conversation, error) {
	// Try to find existing open conversation
	convID, err := s.repo.FindOpen(orgID, channelID, contactID)

	if err == nil {
		return s.GetByID(orgID, convID)
	}

	// Create new conversation
	conv := &models.Conversation{
		ID:             uuid.New(),
		OrganizationID: orgID,
		ChannelID:      channelID,
		ContactID:      contactID,
		Status:         models.ConversationStatusOpen,
	}
	err = s.repo.Create(conv)
	if err != nil {
		return nil, err
	}

	return s.GetByID(orgID, conv.ID)
}

// Delete permanently removes a conversation
func (s *ConversationService) Delete(orgID, convID uuid.UUID) error {
	err := s.repo.Delete(orgID, convID)
	if err == repository.ErrNotFound {
		return ErrConversationNotFound
	}
	return err
}
