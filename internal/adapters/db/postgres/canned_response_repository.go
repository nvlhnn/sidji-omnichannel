package postgres

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/sidji-omnichannel/internal/domain/ports/repository"
	"github.com/sidji-omnichannel/internal/models"
)

type cannedResponseRepository struct {
	db *sql.DB
}

func NewCannedResponseRepository(db *sql.DB) repository.CannedResponseRepository {
	return &cannedResponseRepository{db: db}
}

func (r *cannedResponseRepository) List(orgID uuid.UUID) ([]*models.CannedResponse, error) {
	rows, err := r.db.Query(`SELECT id, organization_id, shortcut, title, content, created_by, created_at FROM canned_responses WHERE organization_id = $1 ORDER BY shortcut ASC`, orgID)
	if err != nil { return nil, err }
	defer rows.Close()
	var responses []*models.CannedResponse
	for rows.Next() {
		resp := &models.CannedResponse{}
		if err := rows.Scan(&resp.ID, &resp.OrganizationID, &resp.Shortcut, &resp.Title, &resp.Content, &resp.CreatedBy, &resp.CreatedAt); err != nil { return nil, err }
		responses = append(responses, resp)
	}
	return responses, nil
}

func (r *cannedResponseRepository) GetByShortcut(orgID uuid.UUID, shortcut string) (*models.CannedResponse, error) {
	resp := &models.CannedResponse{}
	err := r.db.QueryRow(`SELECT id, organization_id, shortcut, title, content, created_by, created_at FROM canned_responses WHERE organization_id = $1 AND shortcut = $2`, orgID, shortcut).Scan(&resp.ID, &resp.OrganizationID, &resp.Shortcut, &resp.Title, &resp.Content, &resp.CreatedBy, &resp.CreatedAt)
	if err == sql.ErrNoRows { return nil, repository.ErrNotFound }
	return resp, err
}

func (r *cannedResponseRepository) GetByID(orgID, responseID uuid.UUID) (*models.CannedResponse, error) {
	resp := &models.CannedResponse{}
	err := r.db.QueryRow(`SELECT id, organization_id, shortcut, title, content, created_by, created_at FROM canned_responses WHERE id = $1 AND organization_id = $2`, responseID, orgID).Scan(&resp.ID, &resp.OrganizationID, &resp.Shortcut, &resp.Title, &resp.Content, &resp.CreatedBy, &resp.CreatedAt)
	if err == sql.ErrNoRows { return nil, repository.ErrNotFound }
	return resp, err
}

func (r *cannedResponseRepository) Create(resp *models.CannedResponse) error {
	_, err := r.db.Exec(`INSERT INTO canned_responses (id, organization_id, shortcut, title, content, created_by) VALUES ($1, $2, $3, $4, $5, $6)`, resp.ID, resp.OrganizationID, resp.Shortcut, resp.Title, resp.Content, resp.CreatedBy)
	return err
}

func (r *cannedResponseRepository) Update(orgID, responseID uuid.UUID, input *models.CreateCannedResponseInput) error {
	res, err := r.db.Exec(`UPDATE canned_responses SET shortcut = $1, title = $2, content = $3, updated_at = NOW() WHERE id = $4 AND organization_id = $5`, input.Shortcut, input.Title, input.Content, responseID, orgID)
	if err != nil { return err }
	rows, _ := res.RowsAffected()
	if rows == 0 { return repository.ErrNotFound }
	return nil
}

func (r *cannedResponseRepository) Delete(orgID, responseID uuid.UUID) error {
	res, err := r.db.Exec(`DELETE FROM canned_responses WHERE id = $1 AND organization_id = $2`, responseID, orgID)
	if err != nil { return err }
	rows, _ := res.RowsAffected()
	if rows == 0 { return repository.ErrNotFound }
	return nil
}

func (r *cannedResponseRepository) Search(orgID uuid.UUID, query string) ([]*models.CannedResponse, error) {
	rows, err := r.db.Query(`SELECT id, organization_id, shortcut, title, content, created_by, created_at FROM canned_responses WHERE organization_id = $1 AND (shortcut ILIKE $2 OR title ILIKE $2 OR content ILIKE $2) ORDER BY shortcut ASC LIMIT 10`, orgID, "%"+query+"%")
	if err != nil { return nil, err }
	defer rows.Close()
	var responses []*models.CannedResponse
	for rows.Next() {
		resp := &models.CannedResponse{}
		if err := rows.Scan(&resp.ID, &resp.OrganizationID, &resp.Shortcut, &resp.Title, &resp.Content, &resp.CreatedBy, &resp.CreatedAt); err != nil { return nil, err }
		responses = append(responses, resp)
	}
	return responses, nil
}
