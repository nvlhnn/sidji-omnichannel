package repository

import (
	"github.com/google/uuid"
	"github.com/sidji-omnichannel/internal/models"
)

// TeamRepository defines the outbound port for team/user data
type TeamRepository interface {
	ListMembers(orgID uuid.UUID) ([]*models.User, error)
	GetMember(orgID, userID uuid.UUID) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	ExistsByEmail(email string) (bool, error)
	CreateMember(user *models.User, passwordHash string) error
	UpdateMember(orgID, userID uuid.UUID, input *models.UpdateUserInput) error
	UpdateMemberRole(orgID, userID uuid.UUID, role models.UserRole) error
	DeleteMember(orgID, userID uuid.UUID) error
	UpdateStatus(userID uuid.UUID, status models.UserStatus) error
	GetAdminCount(orgID uuid.UUID) (int, error)
	GetOrganization(orgID uuid.UUID) (*models.Organization, error)
	UpdateOrganization(orgID uuid.UUID, name, plan string, aiLimit, msgLimit int) error
}
