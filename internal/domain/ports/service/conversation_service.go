package service

import (
	"github.com/google/uuid"
	"github.com/sidji-omnichannel/internal/models"
)

// ConversationService defines the inbound port for conversation operations
type ConversationService interface {
	List(orgID uuid.UUID, filter *models.ConversationFilter) ([]*models.ConversationListItem, int, error)
	GetByID(orgID, convID uuid.UUID) (*models.Conversation, error)
	Assign(orgID, convID, userID uuid.UUID) error
	UpdateStatus(orgID, convID uuid.UUID, status models.ConversationStatus) error
	FindOrCreate(orgID, channelID, contactID uuid.UUID) (*models.Conversation, error)
	Delete(orgID, convID uuid.UUID) error
}
