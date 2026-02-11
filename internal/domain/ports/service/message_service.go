package service

import (
	"github.com/google/uuid"
	"github.com/sidji-omnichannel/internal/models"
)

// MessageService defines the inbound port for message operations
type MessageService interface {
	List(convID uuid.UUID, filter *models.MessageFilter) (*models.MessageList, error)
	Create(msg *models.Message) error
	UpdateStatus(externalID string, status models.MessageStatus) error
	UpdateStatusByID(id uuid.UUID, status models.MessageStatus) error
	MarkAsRead(convID uuid.UUID) error
	CountUnread(convID uuid.UUID) (int, error)
	GetByExternalID(externalID string) (*models.Message, error)
	GetUserPublic(userID uuid.UUID) (*models.UserPublic, error)
}
