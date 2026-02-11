package repository

import (
	"github.com/google/uuid"
	"github.com/sidji-omnichannel/internal/models"
)

// CannedResponseRepository defines the outbound port for canned response data
type CannedResponseRepository interface {
	List(orgID uuid.UUID) ([]*models.CannedResponse, error)
	GetByShortcut(orgID uuid.UUID, shortcut string) (*models.CannedResponse, error)
	GetByID(orgID, responseID uuid.UUID) (*models.CannedResponse, error)
	Create(resp *models.CannedResponse) error
	Update(orgID, responseID uuid.UUID, input *models.CreateCannedResponseInput) error
	Delete(orgID, responseID uuid.UUID) error
	Search(orgID uuid.UUID, query string) ([]*models.CannedResponse, error)
}
