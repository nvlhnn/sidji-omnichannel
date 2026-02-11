package service

import (
	"github.com/google/uuid"
	"github.com/sidji-omnichannel/internal/models"
)

// ContactService defines the inbound port for contact operations
type ContactService interface {
	List(orgID uuid.UUID, page, limit int, search string) ([]*models.Contact, int, error)
	GetByID(orgID, contactID uuid.UUID) (*models.Contact, error)
	Create(orgID uuid.UUID, input *models.CreateContactInput) (*models.Contact, error)
	Update(orgID, contactID uuid.UUID, input *models.UpdateContactInput) (*models.Contact, error)
	Delete(orgID, contactID uuid.UUID) error
	FindOrCreateByWhatsAppID(orgID uuid.UUID, whatsAppID, name string) (*models.Contact, error)
	FindOrCreateByInstagramID(orgID uuid.UUID, instagramID, name, avatarURL string) (*models.Contact, error)
	FindOrCreateByFacebookID(orgID uuid.UUID, facebookID, name, avatarURL string) (*models.Contact, error)
	GetConversations(orgID, contactID uuid.UUID) ([]*models.Conversation, error)
}
