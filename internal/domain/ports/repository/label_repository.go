package repository

import (
	"github.com/google/uuid"
	"github.com/sidji-omnichannel/internal/models"
)

// LabelRepository defines the outbound port for label data
type LabelRepository interface {
	List(orgID uuid.UUID) ([]*models.Label, error)
	GetByID(orgID, labelID uuid.UUID) (*models.Label, error)
	Create(label *models.Label) error
	Update(orgID, labelID uuid.UUID, input *models.CreateLabelInput) error
	Delete(orgID, labelID uuid.UUID) error
	AddToConversation(convID, labelID uuid.UUID) error
	RemoveFromConversation(convID, labelID uuid.UUID) error
	GetConversationLabels(convID uuid.UUID) ([]*models.Label, error)
}
