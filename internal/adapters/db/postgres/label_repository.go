package postgres

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/sidji-omnichannel/internal/domain/ports/repository"
	"github.com/sidji-omnichannel/internal/models"
)

type labelRepository struct {
	db *sql.DB
}

func NewLabelRepository(db *sql.DB) repository.LabelRepository {
	return &labelRepository{db: db}
}

func (r *labelRepository) List(orgID uuid.UUID) ([]*models.Label, error) {
	rows, err := r.db.Query(`SELECT id, organization_id, name, color, created_at FROM labels WHERE organization_id = $1 ORDER BY name ASC`, orgID)
	if err != nil { return nil, err }
	defer rows.Close()
	var labels []*models.Label
	for rows.Next() {
		l := &models.Label{}
		if err := rows.Scan(&l.ID, &l.OrganizationID, &l.Name, &l.Color, &l.CreatedAt); err != nil { return nil, err }
		labels = append(labels, l)
	}
	return labels, nil
}

func (r *labelRepository) GetByID(orgID, labelID uuid.UUID) (*models.Label, error) {
	l := &models.Label{}
	err := r.db.QueryRow(`SELECT id, organization_id, name, color, created_at FROM labels WHERE id = $1 AND organization_id = $2`, labelID, orgID).Scan(&l.ID, &l.OrganizationID, &l.Name, &l.Color, &l.CreatedAt)
	if err == sql.ErrNoRows { return nil, repository.ErrNotFound }
	return l, err
}

func (r *labelRepository) Create(l *models.Label) error {
	_, err := r.db.Exec(`INSERT INTO labels (id, organization_id, name, color) VALUES ($1, $2, $3, $4)`, l.ID, l.OrganizationID, l.Name, l.Color)
	return err
}

func (r *labelRepository) Update(orgID, labelID uuid.UUID, input *models.CreateLabelInput) error {
	res, err := r.db.Exec(`UPDATE labels SET name = $1, color = $2 WHERE id = $3 AND organization_id = $4`, input.Name, input.Color, labelID, orgID)
	if err != nil { return err }
	rows, _ := res.RowsAffected()
	if rows == 0 { return repository.ErrNotFound }
	return nil
}

func (r *labelRepository) Delete(orgID, labelID uuid.UUID) error {
	res, err := r.db.Exec(`DELETE FROM labels WHERE id = $1 AND organization_id = $2`, labelID, orgID)
	if err != nil { return err }
	rows, _ := res.RowsAffected()
	if rows == 0 { return repository.ErrNotFound }
	return nil
}

func (r *labelRepository) AddToConversation(convID, labelID uuid.UUID) error {
	_, err := r.db.Exec(`INSERT INTO conversation_labels (conversation_id, label_id) VALUES ($1, $2) ON CONFLICT (conversation_id, label_id) DO NOTHING`, convID, labelID)
	return err
}

func (r *labelRepository) RemoveFromConversation(convID, labelID uuid.UUID) error {
	_, err := r.db.Exec(`DELETE FROM conversation_labels WHERE conversation_id = $1 AND label_id = $2`, convID, labelID)
	return err
}

func (r *labelRepository) GetConversationLabels(convID uuid.UUID) ([]*models.Label, error) {
	rows, err := r.db.Query(`SELECT l.id, l.organization_id, l.name, l.color, l.created_at FROM labels l JOIN conversation_labels cl ON cl.label_id = l.id WHERE cl.conversation_id = $1 ORDER BY l.name ASC`, convID)
	if err != nil { return nil, err }
	defer rows.Close()
	var labels []*models.Label
	for rows.Next() {
		l := &models.Label{}
		if err := rows.Scan(&l.ID, &l.OrganizationID, &l.Name, &l.Color, &l.CreatedAt); err != nil { return nil, err }
		labels = append(labels, l)
	}
	return labels, nil
}
