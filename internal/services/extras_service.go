package services

import (
	"errors"

	"github.com/google/uuid"
	"github.com/sidji-omnichannel/internal/domain/ports/repository"
	"github.com/sidji-omnichannel/internal/models"
)

var (
	ErrCannedResponseNotFound = errors.New("canned response not found")
	ErrLabelNotFound          = errors.New("label not found")
)

// CannedResponseService handles canned response operations
type CannedResponseService struct {
	repo repository.CannedResponseRepository
}

// NewCannedResponseService creates a new canned response service
func NewCannedResponseService(repo repository.CannedResponseRepository) *CannedResponseService {
	return &CannedResponseService{repo: repo}
}

// List returns all canned responses for an organization
func (s *CannedResponseService) List(orgID uuid.UUID) ([]*models.CannedResponse, error) {
	return s.repo.List(orgID)
}

// GetByShortcut retrieves a canned response by shortcut
func (s *CannedResponseService) GetByShortcut(orgID uuid.UUID, shortcut string) (*models.CannedResponse, error) {
	resp, err := s.repo.GetByShortcut(orgID, shortcut)
	if err == repository.ErrNotFound {
		return nil, ErrCannedResponseNotFound
	}
	return resp, err
}

// Create creates a new canned response
func (s *CannedResponseService) Create(orgID, userID uuid.UUID, input *models.CreateCannedResponseInput) (*models.CannedResponse, error) {
	resp := &models.CannedResponse{
		ID:             uuid.New(),
		OrganizationID: orgID,
		Shortcut:       input.Shortcut,
		Title:          input.Title,
		Content:        input.Content,
		CreatedBy:      userID,
	}

	if err := s.repo.Create(resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// Update updates a canned response
func (s *CannedResponseService) Update(orgID, responseID uuid.UUID, input *models.CreateCannedResponseInput) (*models.CannedResponse, error) {
	err := s.repo.Update(orgID, responseID, input)
	if err == repository.ErrNotFound {
		return nil, ErrCannedResponseNotFound
	}
	if err != nil {
		return nil, err
	}

	return s.repo.GetByID(orgID, responseID)
}

// Delete deletes a canned response
func (s *CannedResponseService) Delete(orgID, responseID uuid.UUID) error {
	err := s.repo.Delete(orgID, responseID)
	if err == repository.ErrNotFound {
		return ErrCannedResponseNotFound
	}
	return err
}

// Search searches canned responses by shortcut or content
func (s *CannedResponseService) Search(orgID uuid.UUID, query string) ([]*models.CannedResponse, error) {
	return s.repo.Search(orgID, query)
}

// ============================================
// Label Service
// ============================================

// LabelService handles label operations
type LabelService struct {
	repo repository.LabelRepository
}

// NewLabelService creates a new label service
func NewLabelService(repo repository.LabelRepository) *LabelService {
	return &LabelService{repo: repo}
}

// List returns all labels for an organization
func (s *LabelService) List(orgID uuid.UUID) ([]*models.Label, error) {
	return s.repo.List(orgID)
}

// Create creates a new label
func (s *LabelService) Create(orgID uuid.UUID, input *models.CreateLabelInput) (*models.Label, error) {
	label := &models.Label{
		ID:             uuid.New(),
		OrganizationID: orgID,
		Name:           input.Name,
		Color:          input.Color,
	}

	if err := s.repo.Create(label); err != nil {
		return nil, err
	}

	return label, nil
}

// Update updates a label
func (s *LabelService) Update(orgID, labelID uuid.UUID, input *models.CreateLabelInput) (*models.Label, error) {
	err := s.repo.Update(orgID, labelID, input)
	if err == repository.ErrNotFound {
		return nil, ErrLabelNotFound
	}
	if err != nil {
		return nil, err
	}

	return s.repo.GetByID(orgID, labelID)
}

// Delete deletes a label
func (s *LabelService) Delete(orgID, labelID uuid.UUID) error {
	err := s.repo.Delete(orgID, labelID)
	if err == repository.ErrNotFound {
		return ErrLabelNotFound
	}
	return err
}

// AddToConversation adds a label to a conversation
func (s *LabelService) AddToConversation(convID, labelID uuid.UUID) error {
	return s.repo.AddToConversation(convID, labelID)
}

// RemoveFromConversation removes a label from a conversation
func (s *LabelService) RemoveFromConversation(convID, labelID uuid.UUID) error {
	return s.repo.RemoveFromConversation(convID, labelID)
}

// GetConversationLabels returns all labels for a conversation
func (s *LabelService) GetConversationLabels(convID uuid.UUID) ([]*models.Label, error) {
	return s.repo.GetConversationLabels(convID)
}
