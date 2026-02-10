package services

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/sidji-omnichannel/internal/models"
)

var (
	ErrCannedResponseNotFound = errors.New("canned response not found")
	ErrLabelNotFound          = errors.New("label not found")
)

// CannedResponseService handles canned response operations
type CannedResponseService struct {
	db *sql.DB
}

// NewCannedResponseService creates a new canned response service
func NewCannedResponseService(db *sql.DB) *CannedResponseService {
	return &CannedResponseService{db: db}
}

// List returns all canned responses for an organization
func (s *CannedResponseService) List(orgID uuid.UUID) ([]*models.CannedResponse, error) {
	rows, err := s.db.Query(`
		SELECT id, organization_id, shortcut, title, content, created_by, created_at
		FROM canned_responses
		WHERE organization_id = $1
		ORDER BY shortcut ASC
	`, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var responses []*models.CannedResponse
	for rows.Next() {
		resp := &models.CannedResponse{}
		err := rows.Scan(
			&resp.ID, &resp.OrganizationID, &resp.Shortcut, &resp.Title,
			&resp.Content, &resp.CreatedBy, &resp.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		responses = append(responses, resp)
	}

	return responses, nil
}

// GetByShortcut retrieves a canned response by shortcut
func (s *CannedResponseService) GetByShortcut(orgID uuid.UUID, shortcut string) (*models.CannedResponse, error) {
	resp := &models.CannedResponse{}
	err := s.db.QueryRow(`
		SELECT id, organization_id, shortcut, title, content, created_by, created_at
		FROM canned_responses
		WHERE organization_id = $1 AND shortcut = $2
	`, orgID, shortcut).Scan(
		&resp.ID, &resp.OrganizationID, &resp.Shortcut, &resp.Title,
		&resp.Content, &resp.CreatedBy, &resp.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrCannedResponseNotFound
	}
	if err != nil {
		return nil, err
	}
	return resp, nil
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

	_, err := s.db.Exec(`
		INSERT INTO canned_responses (id, organization_id, shortcut, title, content, created_by)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, resp.ID, resp.OrganizationID, resp.Shortcut, resp.Title, resp.Content, resp.CreatedBy)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Update updates a canned response
func (s *CannedResponseService) Update(orgID, responseID uuid.UUID, input *models.CreateCannedResponseInput) (*models.CannedResponse, error) {
	result, err := s.db.Exec(`
		UPDATE canned_responses 
		SET shortcut = $1, title = $2, content = $3, updated_at = NOW()
		WHERE id = $4 AND organization_id = $5
	`, input.Shortcut, input.Title, input.Content, responseID, orgID)

	if err != nil {
		return nil, err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return nil, ErrCannedResponseNotFound
	}

	resp := &models.CannedResponse{}
	err = s.db.QueryRow(`
		SELECT id, organization_id, shortcut, title, content, created_by, created_at
		FROM canned_responses WHERE id = $1
	`, responseID).Scan(
		&resp.ID, &resp.OrganizationID, &resp.Shortcut, &resp.Title,
		&resp.Content, &resp.CreatedBy, &resp.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Delete deletes a canned response
func (s *CannedResponseService) Delete(orgID, responseID uuid.UUID) error {
	result, err := s.db.Exec(`
		DELETE FROM canned_responses WHERE id = $1 AND organization_id = $2
	`, responseID, orgID)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrCannedResponseNotFound
	}

	return nil
}

// Search searches canned responses by shortcut or content
func (s *CannedResponseService) Search(orgID uuid.UUID, query string) ([]*models.CannedResponse, error) {
	rows, err := s.db.Query(`
		SELECT id, organization_id, shortcut, title, content, created_by, created_at
		FROM canned_responses
		WHERE organization_id = $1 AND (shortcut ILIKE $2 OR title ILIKE $2 OR content ILIKE $2)
		ORDER BY shortcut ASC
		LIMIT 10
	`, orgID, "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var responses []*models.CannedResponse
	for rows.Next() {
		resp := &models.CannedResponse{}
		err := rows.Scan(
			&resp.ID, &resp.OrganizationID, &resp.Shortcut, &resp.Title,
			&resp.Content, &resp.CreatedBy, &resp.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		responses = append(responses, resp)
	}

	return responses, nil
}

// ============================================
// Label Service
// ============================================

// LabelService handles label operations
type LabelService struct {
	db *sql.DB
}

// NewLabelService creates a new label service
func NewLabelService(db *sql.DB) *LabelService {
	return &LabelService{db: db}
}

// List returns all labels for an organization
func (s *LabelService) List(orgID uuid.UUID) ([]*models.Label, error) {
	rows, err := s.db.Query(`
		SELECT id, organization_id, name, color, created_at
		FROM labels
		WHERE organization_id = $1
		ORDER BY name ASC
	`, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var labels []*models.Label
	for rows.Next() {
		label := &models.Label{}
		err := rows.Scan(
			&label.ID, &label.OrganizationID, &label.Name, &label.Color, &label.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		labels = append(labels, label)
	}

	return labels, nil
}

// Create creates a new label
func (s *LabelService) Create(orgID uuid.UUID, input *models.CreateLabelInput) (*models.Label, error) {
	label := &models.Label{
		ID:             uuid.New(),
		OrganizationID: orgID,
		Name:           input.Name,
		Color:          input.Color,
	}

	_, err := s.db.Exec(`
		INSERT INTO labels (id, organization_id, name, color)
		VALUES ($1, $2, $3, $4)
	`, label.ID, label.OrganizationID, label.Name, label.Color)

	if err != nil {
		return nil, err
	}

	return label, nil
}

// Update updates a label
func (s *LabelService) Update(orgID, labelID uuid.UUID, input *models.CreateLabelInput) (*models.Label, error) {
	result, err := s.db.Exec(`
		UPDATE labels SET name = $1, color = $2
		WHERE id = $3 AND organization_id = $4
	`, input.Name, input.Color, labelID, orgID)

	if err != nil {
		return nil, err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return nil, ErrLabelNotFound
	}

	label := &models.Label{}
	err = s.db.QueryRow(`
		SELECT id, organization_id, name, color, created_at FROM labels WHERE id = $1
	`, labelID).Scan(&label.ID, &label.OrganizationID, &label.Name, &label.Color, &label.CreatedAt)
	if err != nil {
		return nil, err
	}

	return label, nil
}

// Delete deletes a label
func (s *LabelService) Delete(orgID, labelID uuid.UUID) error {
	result, err := s.db.Exec(`
		DELETE FROM labels WHERE id = $1 AND organization_id = $2
	`, labelID, orgID)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrLabelNotFound
	}

	return nil
}

// AddToConversation adds a label to a conversation
func (s *LabelService) AddToConversation(convID, labelID uuid.UUID) error {
	_, err := s.db.Exec(`
		INSERT INTO conversation_labels (conversation_id, label_id)
		VALUES ($1, $2)
		ON CONFLICT (conversation_id, label_id) DO NOTHING
	`, convID, labelID)
	return err
}

// RemoveFromConversation removes a label from a conversation
func (s *LabelService) RemoveFromConversation(convID, labelID uuid.UUID) error {
	_, err := s.db.Exec(`
		DELETE FROM conversation_labels WHERE conversation_id = $1 AND label_id = $2
	`, convID, labelID)
	return err
}

// GetConversationLabels returns all labels for a conversation
func (s *LabelService) GetConversationLabels(convID uuid.UUID) ([]*models.Label, error) {
	rows, err := s.db.Query(`
		SELECT l.id, l.organization_id, l.name, l.color, l.created_at
		FROM labels l
		JOIN conversation_labels cl ON cl.label_id = l.id
		WHERE cl.conversation_id = $1
		ORDER BY l.name ASC
	`, convID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var labels []*models.Label
	for rows.Next() {
		label := &models.Label{}
		err := rows.Scan(
			&label.ID, &label.OrganizationID, &label.Name, &label.Color, &label.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		labels = append(labels, label)
	}

	return labels, nil
}
