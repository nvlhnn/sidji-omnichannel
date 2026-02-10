package services

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/sidji-omnichannel/internal/models"
	"github.com/sidji-omnichannel/internal/subscription"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrTeamMemberNotFound = errors.New("team member not found")
)

// TeamService handles team/user management operations
type TeamService struct {
	db *sql.DB
}

// NewTeamService creates a new team service
func NewTeamService(db *sql.DB) *TeamService {
	return &TeamService{db: db}
}

// ListMembers returns all team members for an organization
func (s *TeamService) ListMembers(orgID uuid.UUID) ([]*models.User, error) {
	rows, err := s.db.Query(`
		SELECT id, organization_id, email, name, role, avatar_url, status, last_seen_at, created_at
		FROM users
		WHERE organization_id = $1
		ORDER BY created_at ASC
	`, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []*models.User
	for rows.Next() {
		user := &models.User{}
		var avatarURL sql.NullString
		var lastSeenAt sql.NullTime
		
		err := rows.Scan(
			&user.ID, &user.OrganizationID, &user.Email, &user.Name,
			&user.Role, &avatarURL, &user.Status, &lastSeenAt, &user.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		
		if avatarURL.Valid {
			user.AvatarURL = avatarURL.String
		}
		if lastSeenAt.Valid {
			user.LastSeenAt = &lastSeenAt.Time
		}
		
		members = append(members, user)
	}

	return members, nil
}

// GetMember retrieves a team member by ID
func (s *TeamService) GetMember(orgID, userID uuid.UUID) (*models.User, error) {
	user := &models.User{}
	var avatarURL sql.NullString
	var lastSeenAt sql.NullTime

	err := s.db.QueryRow(`
		SELECT id, organization_id, email, name, role, avatar_url, status, last_seen_at, created_at
		FROM users
		WHERE id = $1 AND organization_id = $2
	`, userID, orgID).Scan(
		&user.ID, &user.OrganizationID, &user.Email, &user.Name,
		&user.Role, &avatarURL, &user.Status, &lastSeenAt, &user.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrTeamMemberNotFound
	}
	if err != nil {
		return nil, err
	}

	if avatarURL.Valid {
		user.AvatarURL = avatarURL.String
	}
	if lastSeenAt.Valid {
		user.LastSeenAt = &lastSeenAt.Time
	}

	return user, nil
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
	var exists bool
	err = s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", input.Email).Scan(&exists)
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

	_, err = s.db.Exec(`
		INSERT INTO users (id, organization_id, email, password_hash, name, role, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, user.ID, user.OrganizationID, user.Email, string(hashedPassword), user.Name, user.Role, user.Status)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateMember updates a team member's information
func (s *TeamService) UpdateMember(orgID, userID uuid.UUID, input *models.UpdateUserInput) (*models.User, error) {
	// Build update query dynamically
	query := "UPDATE users SET updated_at = NOW()"
	args := []interface{}{}
	argIndex := 1

	if input.Name != "" {
		query += ", name = $1"
		args = append(args, input.Name)
		argIndex++
	}

	if input.AvatarURL != "" {
		query += ", avatar_url = $2"
		args = append(args, input.AvatarURL)
		argIndex++
	}

	query += " WHERE id = $" + string(rune('0'+argIndex)) + " AND organization_id = $" + string(rune('0'+argIndex+1))
	args = append(args, userID, orgID)

	// Fix the string rune thing - it was index based
	// Re-writing more safely
	
	query = "UPDATE users SET updated_at = NOW()"
	args = []interface{}{}
	idx := 1
	if input.Name != "" {
		query += ", name = $1"
		args = append(args, input.Name)
		idx++
	}
	if input.AvatarURL != "" {
		query += ", avatar_url = $" + fmt.Sprintf("%d", idx)
		args = append(args, input.AvatarURL)
		idx++
	}
	query += " WHERE id = $" + fmt.Sprintf("%d", idx) + " AND organization_id = $" + fmt.Sprintf("%d", idx+1)
	args = append(args, userID, orgID)
	
	// Wait, I need fmt
	return s.updateMemberSafe(orgID, userID, input)
}

func (s *TeamService) updateMemberSafe(orgID, userID uuid.UUID, input *models.UpdateUserInput) (*models.User, error) {
	query := `UPDATE users SET updated_at = NOW(), name = COALESCE(NULLIF($1, ''), name), avatar_url = COALESCE(NULLIF($2, ''), avatar_url) 
	          WHERE id = $3 AND organization_id = $4`
	
	result, err := s.db.Exec(query, input.Name, input.AvatarURL, userID, orgID)
	if err != nil {
		return nil, err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return nil, ErrTeamMemberNotFound
	}
	return s.GetMember(orgID, userID)
}

// UpdateMemberRole updates a team member's role
func (s *TeamService) UpdateMemberRole(orgID, userID uuid.UUID, role models.UserRole) error {
	result, err := s.db.Exec(`
		UPDATE users SET role = $1, updated_at = NOW()
		WHERE id = $2 AND organization_id = $3
	`, role, userID, orgID)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrTeamMemberNotFound
	}

	return nil
}

// RemoveMember removes a team member from the organization
func (s *TeamService) RemoveMember(orgID, userID uuid.UUID) error {
	// Don't allow removing the last admin
	var adminCount int
	err := s.db.QueryRow(`
		SELECT COUNT(*) FROM users WHERE organization_id = $1 AND role = 'admin'
	`, orgID).Scan(&adminCount)
	if err != nil {
		return err
	}

	var targetRole string
	err = s.db.QueryRow(`
		SELECT role FROM users WHERE id = $1 AND organization_id = $2
	`, userID, orgID).Scan(&targetRole)
	if err == sql.ErrNoRows {
		return ErrTeamMemberNotFound
	}
	if err != nil {
		return err
	}

	if targetRole == "admin" && adminCount <= 1 {
		return errors.New("cannot remove the last admin")
	}

	result, err := s.db.Exec(`
		DELETE FROM users WHERE id = $1 AND organization_id = $2
	`, userID, orgID)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrTeamMemberNotFound
	}

	return nil
}

// UpdateStatus updates a user's online status
func (s *TeamService) UpdateStatus(userID uuid.UUID, status models.UserStatus) error {
	_, err := s.db.Exec(`
		UPDATE users SET status = $1, last_seen_at = NOW(), updated_at = NOW()
		WHERE id = $2
	`, status, userID)
	return err
}

// GetOrganization retrieves organization details with counts and limit status
func (s *TeamService) GetOrganization(orgID uuid.UUID) (*models.Organization, error) {
	org := &models.Organization{}
	err := s.db.QueryRow(`
		SELECT id, name, slug, plan, subscription_status, ai_credits_limit, ai_credits_used, message_usage_limit, message_usage_used, billing_cycle_start, created_at, updated_at 
		FROM organizations WHERE id = $1
	`, orgID).Scan(&org.ID, &org.Name, &org.Slug, &org.Plan, &org.SubscriptionStatus, &org.AICreditsLimit, &org.AICreditsUsed, &org.MessageUsageLimit, &org.MessageUsageUsed, &org.BillingCycleStart, &org.CreatedAt, &org.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// Fetch counts
	err = s.db.QueryRow("SELECT COUNT(*) FROM users WHERE organization_id = $1", orgID).Scan(&org.UserCount)
	if err != nil {
		return nil, err
	}

	err = s.db.QueryRow("SELECT COUNT(*) FROM channels WHERE organization_id = $1", orgID).Scan(&org.ChannelCount)
	if err != nil {
		return nil, err
	}

	// Compute compliance
	org.IsOverLimit = !subscription.IsCompliance(org.Plan, org.UserCount, org.ChannelCount)

	return org, nil
}

// UpdateOrganization updates organization details
func (s *TeamService) UpdateOrganization(orgID uuid.UUID, name string, plan string) (*models.Organization, error) {
	query := "UPDATE organizations SET updated_at = NOW()"
	args := []interface{}{}
	argIdx := 1

	if name != "" {
		query += fmt.Sprintf(", name = $%d", argIdx)
		args = append(args, name)
		argIdx++
	}

	if plan != "" {
		limits := subscription.GetSubscriptionLimits(plan)
		log.Printf("[UpdateOrganization] Updating plan to %s, setting ai_credits_limit to %d, message_usage_limit to %d", plan, limits.MaxAIReply, limits.MaxMessages)
		query += fmt.Sprintf(", plan = $%d, ai_credits_limit = $%d, message_usage_limit = $%d, billing_cycle_start = NOW(), ai_credits_used = 0, message_usage_used = 0", argIdx, argIdx+1, argIdx+2)
		args = append(args, plan, limits.MaxAIReply, limits.MaxMessages)
		argIdx += 3
	}

	query += fmt.Sprintf(" WHERE id = $%d", argIdx)
	args = append(args, orgID)

	_, err := s.db.Exec(query, args...)
	if err != nil {
		return nil, err
	}
	return s.GetOrganization(orgID)
}
