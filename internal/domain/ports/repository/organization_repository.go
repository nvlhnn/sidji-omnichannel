package repository

import (
	"github.com/google/uuid"
)

// OrganizationRepository defines the outbound port for organization data
type OrganizationRepository interface {
	GetPlan(orgID uuid.UUID) (string, error)
	GetChannelCount(orgID uuid.UUID) (int, error)
}
