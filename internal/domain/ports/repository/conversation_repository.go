package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/sidji-omnichannel/internal/models"
)

var (
	ErrNotFound = errors.New("resource not found")
)

// ConversationRepository defines the outbound port for conversation data
type ConversationRepository interface {
	List(orgID uuid.UUID, filter *models.ConversationFilter) ([]*models.ConversationListItem, int, error)
	GetByID(orgID, convID uuid.UUID) (*models.Conversation, error)
	Assign(orgID, convID, userID uuid.UUID) error
	UpdateStatus(orgID, convID uuid.UUID, status models.ConversationStatus) error
	FindOpen(orgID, channelID, contactID uuid.UUID) (uuid.UUID, error)
	Create(conv *models.Conversation) error
	Delete(orgID, convID uuid.UUID) error
	UpdateLastMessage(convID uuid.UUID, lastMessageAt interface{}) error
	UpdateLastHumanReply(convID uuid.UUID, lastReplyAt interface{}) error
}
