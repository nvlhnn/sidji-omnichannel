package services

import (
	"errors"

	"github.com/google/uuid"
	"github.com/sidji-omnichannel/internal/domain/ports/repository"
	"github.com/sidji-omnichannel/internal/models"
)

var (
	ErrContactNotFound = errors.New("contact not found")
)

// ContactService handles contact operations
type ContactService struct {
	repo repository.ContactRepository
}

// NewContactService creates a new contact service
func NewContactService(repo repository.ContactRepository) *ContactService {
	return &ContactService{repo: repo}
}

// List returns all contacts for an organization with pagination
func (s *ContactService) List(orgID uuid.UUID, page, limit int, search string) ([]*models.Contact, int, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if page <= 0 {
		page = 1
	}
	return s.repo.List(orgID, page, limit, search)
}

// GetByID retrieves a contact by ID
func (s *ContactService) GetByID(orgID, contactID uuid.UUID) (*models.Contact, error) {
	contact, err := s.repo.GetByID(orgID, contactID)
	if err == repository.ErrNotFound {
		return nil, ErrContactNotFound
	}
	return contact, err
}

// Create creates a new contact
func (s *ContactService) Create(orgID uuid.UUID, input *models.CreateContactInput) (*models.Contact, error) {
	contact := &models.Contact{
		ID:             uuid.New(),
		OrganizationID: orgID,
		Name:           input.Name,
		Phone:          input.Phone,
		Email:          input.Email,
		WhatsAppID:     input.WhatsAppID,
		InstagramID:    input.InstagramID,
		FacebookID:     input.FacebookID,
	}

	if err := s.repo.Create(contact); err != nil {
		return nil, err
	}

	return contact, nil
}

// Update updates a contact
func (s *ContactService) Update(orgID, contactID uuid.UUID, input *models.UpdateContactInput) (*models.Contact, error) {
	contact, err := s.repo.Update(orgID, contactID, input)
	if err == repository.ErrNotFound {
		return nil, ErrContactNotFound
	}
	return contact, err
}

// Delete deletes a contact
func (s *ContactService) Delete(orgID, contactID uuid.UUID) error {
	err := s.repo.Delete(orgID, contactID)
	if err == repository.ErrNotFound {
		return ErrContactNotFound
	}
	return err
}

// FindOrCreateByWhatsAppID finds or creates a contact by WhatsApp ID
func (s *ContactService) FindOrCreateByWhatsAppID(orgID uuid.UUID, whatsAppID, name string) (*models.Contact, error) {
	contact, err := s.repo.GetByWhatsAppID(orgID, whatsAppID)
	if err == nil {
		if name != "" && contact.Name != name {
			s.repo.UpdateNameAndAvatar(contact.ID, name, "")
			contact.Name = name
		}
		return contact, nil
	}

	if err != repository.ErrNotFound {
		return nil, err
	}

	contact = &models.Contact{
		ID:             uuid.New(),
		OrganizationID: orgID,
		Name:           name,
		Phone:          whatsAppID,
		WhatsAppID:     whatsAppID,
	}

	if err := s.repo.Create(contact); err != nil {
		return nil, err
	}

	return contact, nil
}

// FindOrCreateByInstagramID finds or creates a contact by Instagram ID
func (s *ContactService) FindOrCreateByInstagramID(orgID uuid.UUID, instagramID, name, avatarURL string) (*models.Contact, error) {
	contact, err := s.repo.GetByInstagramID(orgID, instagramID)
	if err == nil {
		if (name != "" && contact.Name != name) || (avatarURL != "" && contact.AvatarURL != avatarURL) {
			s.repo.UpdateNameAndAvatar(contact.ID, name, avatarURL)
			if name != "" { contact.Name = name }
			if avatarURL != "" { contact.AvatarURL = avatarURL }
		}
		return contact, nil
	}

	if err != repository.ErrNotFound {
		return nil, err
	}

	contact = &models.Contact{
		ID:             uuid.New(),
		OrganizationID: orgID,
		Name:           name,
		InstagramID:    instagramID,
		AvatarURL:      avatarURL,
	}

	if err := s.repo.Create(contact); err != nil {
		return nil, err
	}

	return contact, nil
}

// FindOrCreateByFacebookID finds or creates a contact by Facebook ID
func (s *ContactService) FindOrCreateByFacebookID(orgID uuid.UUID, facebookID, name, avatarURL string) (*models.Contact, error) {
	contact, err := s.repo.GetByFacebookID(orgID, facebookID)
	if err == nil {
		if (name != "" && contact.Name != name) || (avatarURL != "" && contact.AvatarURL != avatarURL) {
			s.repo.UpdateNameAndAvatar(contact.ID, name, avatarURL)
			if name != "" { contact.Name = name }
			if avatarURL != "" { contact.AvatarURL = avatarURL }
		}
		return contact, nil
	}

	if err != repository.ErrNotFound {
		return nil, err
	}

	contact = &models.Contact{
		ID:             uuid.New(),
		OrganizationID: orgID,
		Name:           name,
		FacebookID:     facebookID,
		AvatarURL:      avatarURL,
	}

	if err := s.repo.Create(contact); err != nil {
		return nil, err
	}

	return contact, nil
}

// GetConversations returns all conversations for a contact
func (s *ContactService) GetConversations(orgID, contactID uuid.UUID) ([]*models.Conversation, error) {
	return s.repo.GetConversations(orgID, contactID)
}

// FindOrCreateByTikTokID finds or creates a contact by TikTok ID
func (s *ContactService) FindOrCreateByTikTokID(orgID uuid.UUID, tiktokID, name, avatarURL string) (*models.Contact, error) {
	contact, err := s.repo.GetByTikTokID(orgID, tiktokID)
	if err == nil {
		if (name != "" && contact.Name != name) || (avatarURL != "" && contact.AvatarURL != avatarURL) {
			s.repo.UpdateNameAndAvatar(contact.ID, name, avatarURL)
			if name != "" { contact.Name = name }
			if avatarURL != "" { contact.AvatarURL = avatarURL }
		}
		return contact, nil
	}

	if err != repository.ErrNotFound {
		return nil, err
	}

	contact = &models.Contact{
		ID:             uuid.New(),
		OrganizationID: orgID,
		Name:           name,
		TikTokID:       tiktokID,
		AvatarURL:      avatarURL,
	}

	if err := s.repo.Create(contact); err != nil {
		return nil, err
	}

	return contact, nil
}
