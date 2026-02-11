package repository

import (
	"github.com/google/uuid"
	"github.com/sidji-omnichannel/internal/models"
)

// MessageRepository defines the outbound port for message data
type MessageRepository interface {
	Create(msg *models.Message) error
	List(convID uuid.UUID, filter *models.MessageFilter) (*models.MessageList, error)
	GetByExternalID(externalID string) (*models.Message, error)
	UpdateStatus(externalID string, status models.MessageStatus) error
	UpdateStatusByID(id uuid.UUID, status models.MessageStatus) error
	MarkAsRead(convID uuid.UUID) error
	CountUnread(convID uuid.UUID) (int, error)
}
