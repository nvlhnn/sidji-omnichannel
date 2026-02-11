package services

import (
	"errors"

	"github.com/google/uuid"
	"github.com/sidji-omnichannel/internal/domain/ports/repository"
	"github.com/sidji-omnichannel/internal/models"
	"github.com/sidji-omnichannel/internal/subscription"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrTeamMemberNotFound = errors.New("team member not found")
)

// TeamService handles team/user management operations
type TeamService struct {
	repo repository.TeamRepository
}

// NewTeamService creates a new team service
func NewTeamService(repo repository.TeamRepository) *TeamService {
	return &TeamService{repo: repo}
}

// ListMembers returns all team members for an organization
func (s *TeamService) ListMembers(orgID uuid.UUID) ([]*models.User, error) {
	return s.repo.ListMembers(orgID)
}

// GetMember retrieves a team member by ID
func (s *TeamService) GetMember(orgID, userID uuid.UUID) (*models.User, error) {
	user, err := s.repo.GetMember(orgID, userID)
	if err == repository.ErrNotFound {
		return nil, ErrTeamMemberNotFound
	}
	return user, err
}

// InviteMember creates a new team member
func (s *TeamService) InviteMember(orgID uuid.UUID, input *models.InviteUserInput) (*models.User, error) {
	// Check subscription limits
	org, err := s.GetOrganization(orgID)
	if err != nil {
		return nil, err
	}

	if err := subscription.CheckLimit(org.Plan, org.UserCount, "user"); err != nil {
		return nil, err
	}

	// Check if email already exists
	exists, err := s.repo.ExistsByEmail(input.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrUserExists
	}

	// Generate a temporary password
	tempPassword := uuid.New().String()[:8]
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(tempPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:             uuid.New(),
		OrganizationID: orgID,
		Email:          input.Email,
		Name:           input.Name,
		Role:           input.Role,
		Status:         models.StatusOffline,
	}

	if err := s.repo.CreateMember(user, string(hashedPassword)); err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateMember updates a team member's information
func (s *TeamService) UpdateMember(orgID, userID uuid.UUID, input *models.UpdateUserInput) (*models.User, error) {
	err := s.repo.UpdateMember(orgID, userID, input)
	if err == repository.ErrNotFound {
		return nil, ErrTeamMemberNotFound
	}
	if err != nil {
		return nil, err
	}
	return s.GetMember(orgID, userID)
}

// UpdateMemberRole updates a team member's role
func (s *TeamService) UpdateMemberRole(orgID, userID uuid.UUID, role models.UserRole) error {
	err := s.repo.UpdateMemberRole(orgID, userID, role)
	if err == repository.ErrNotFound {
		return ErrTeamMemberNotFound
	}
	return err
}

// RemoveMember removes a team member from the organization
func (s *TeamService) RemoveMember(orgID, userID uuid.UUID) error {
	// Don't allow removing the last admin
	adminCount, err := s.repo.GetAdminCount(orgID)
	if err != nil {
		return err
	}

	targetUser, err := s.repo.GetMember(orgID, userID)
	if err == repository.ErrNotFound {
		return ErrTeamMemberNotFound
	}
	if err != nil {
		return err
	}

	if targetUser.Role == models.RoleAdmin && adminCount <= 1 {
		return errors.New("cannot remove the last admin")
	}

	err = s.repo.DeleteMember(orgID, userID)
	if err == repository.ErrNotFound {
		return ErrTeamMemberNotFound
	}
	return err
}

// UpdateStatus updates a user's online status
func (s *TeamService) UpdateStatus(userID uuid.UUID, status models.UserStatus) error {
	return s.repo.UpdateStatus(userID, status)
}

// GetOrganization retrieves organization details with counts and limit status
func (s *TeamService) GetOrganization(orgID uuid.UUID) (*models.Organization, error) {
	org, err := s.repo.GetOrganization(orgID)
	if err != nil {
		return nil, err
	}

	// Compute compliance
	org.IsOverLimit = !subscription.IsCompliance(org.Plan, org.UserCount, org.ChannelCount)

	return org, nil
}

// UpdateOrganization updates organization details
func (s *TeamService) UpdateOrganization(orgID uuid.UUID, name string, plan string) (*models.Organization, error) {
	var aiLimit, msgLimit int
	if plan != "" {
		limits := subscription.GetSubscriptionLimits(plan)
		aiLimit = limits.MaxAIReply
		msgLimit = limits.MaxMessages
	}

	err := s.repo.UpdateOrganization(orgID, name, plan, aiLimit, msgLimit)
	if err != nil {
		return nil, err
	}
	return s.GetOrganization(orgID)
}
