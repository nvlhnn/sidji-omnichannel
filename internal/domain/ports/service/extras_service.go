package service

import (
	"github.com/google/uuid"
	"github.com/sidji-omnichannel/internal/models"
)

// CannedResponseService defines the inbound port for canned response operations
type CannedResponseService interface {
	List(orgID uuid.UUID) ([]*models.CannedResponse, error)
	GetByShortcut(orgID uuid.UUID, shortcut string) (*models.CannedResponse, error)
	Create(orgID, userID uuid.UUID, input *models.CreateCannedResponseInput) (*models.CannedResponse, error)
	Update(orgID, responseID uuid.UUID, input *models.CreateCannedResponseInput) (*models.CannedResponse, error)
	Delete(orgID, responseID uuid.UUID) error
	Search(orgID uuid.UUID, query string) ([]*models.CannedResponse, error)
}

// LabelService defines the inbound port for label operations
type LabelService interface {
	List(orgID uuid.UUID) ([]*models.Label, error)
	Create(orgID uuid.UUID, input *models.CreateLabelInput) (*models.Label, error)
	Update(orgID, labelID uuid.UUID, input *models.CreateLabelInput) (*models.Label, error)
	Delete(orgID, labelID uuid.UUID) error
	AddToConversation(convID, labelID uuid.UUID) error
	RemoveFromConversation(convID, labelID uuid.UUID) error
	GetConversationLabels(convID uuid.UUID) ([]*models.Label, error)
}
