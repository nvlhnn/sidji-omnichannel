package repository

import (
	"github.com/google/uuid"
	"github.com/sidji-omnichannel/internal/models"
)

// ContactRepository defines the outbound port for contact data
type ContactRepository interface {
	List(orgID uuid.UUID, page, limit int, search string) ([]*models.Contact, int, error)
	GetByID(orgID, contactID uuid.UUID) (*models.Contact, error)
	GetByWhatsAppID(orgID uuid.UUID, whatsAppID string) (*models.Contact, error)
	GetByInstagramID(orgID uuid.UUID, instagramID string) (*models.Contact, error)
	GetByFacebookID(orgID uuid.UUID, facebookID string) (*models.Contact, error)
	Create(contact *models.Contact) error
	Update(orgID, contactID uuid.UUID, input *models.UpdateContactInput) (*models.Contact, error)
	UpdateNameAndAvatar(contactID uuid.UUID, name, avatarURL string) error
	Delete(orgID, contactID uuid.UUID) error
	GetConversations(orgID, contactID uuid.UUID) ([]*models.Conversation, error)
}
