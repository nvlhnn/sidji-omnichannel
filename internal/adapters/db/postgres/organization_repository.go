package postgres

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/sidji-omnichannel/internal/domain/ports/repository"
)

type organizationRepository struct {
	db *sql.DB
}

// NewOrganizationRepository creates a new postgres organization repository
func NewOrganizationRepository(db *sql.DB) repository.OrganizationRepository {
	return &organizationRepository{db: db}
}

func (r *organizationRepository) GetPlan(orgID uuid.UUID) (string, error) {
	var plan string
	err := r.db.QueryRow("SELECT plan FROM organizations WHERE id = $1", orgID).Scan(&plan)
	return plan, err
}

func (r *organizationRepository) GetChannelCount(orgID uuid.UUID) (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM channels WHERE organization_id = $1", orgID).Scan(&count)
	return count, err
}
