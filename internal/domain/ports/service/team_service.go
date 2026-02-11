package service

import (
	"github.com/google/uuid"
	"github.com/sidji-omnichannel/internal/models"
)

// TeamService defines the inbound port for team/user operations
type TeamService interface {
	ListMembers(orgID uuid.UUID) ([]*models.User, error)
	GetMember(orgID, userID uuid.UUID) (*models.User, error)
	InviteMember(orgID uuid.UUID, input *models.InviteUserInput) (*models.User, error)
	UpdateMember(orgID, userID uuid.UUID, input *models.UpdateUserInput) (*models.User, error)
	UpdateMemberRole(orgID, userID uuid.UUID, role models.UserRole) error
	RemoveMember(orgID, userID uuid.UUID) error
	UpdateStatus(userID uuid.UUID, status models.UserStatus) error
	GetOrganization(orgID uuid.UUID) (*models.Organization, error)
	UpdateOrganization(orgID uuid.UUID, name string, plan string) (*models.Organization, error)
}
