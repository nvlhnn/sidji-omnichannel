package postgres

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/sidji-omnichannel/internal/domain/ports/repository"
	"github.com/sidji-omnichannel/internal/models"
)

type teamRepository struct {
	db *sql.DB
}

func NewTeamRepository(db *sql.DB) repository.TeamRepository {
	return &teamRepository{db: db}
}

func (r *teamRepository) ListMembers(orgID uuid.UUID) ([]*models.User, error) {
	rows, err := r.db.Query(`
		SELECT id, organization_id, email, name, role, avatar_url, status, last_seen_at, created_at
		FROM users WHERE organization_id = $1 ORDER BY created_at ASC
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
		if err := rows.Scan(&user.ID, &user.OrganizationID, &user.Email, &user.Name, &user.Role, &avatarURL, &user.Status, &lastSeenAt, &user.CreatedAt); err != nil {
			return nil, err
		}
		if avatarURL.Valid { user.AvatarURL = avatarURL.String }
		if lastSeenAt.Valid { user.LastSeenAt = &lastSeenAt.Time }
		members = append(members, user)
	}
	return members, nil
}

func (r *teamRepository) GetMember(orgID, userID uuid.UUID) (*models.User, error) {
	user := &models.User{}
	var avatarURL sql.NullString
	var lastSeenAt sql.NullTime
	err := r.db.QueryRow(`
		SELECT id, organization_id, email, name, role, avatar_url, status, last_seen_at, created_at
		FROM users WHERE id = $1 AND organization_id = $2
	`, userID, orgID).Scan(&user.ID, &user.OrganizationID, &user.Email, &user.Name, &user.Role, &avatarURL, &user.Status, &lastSeenAt, &user.CreatedAt)
	if err == sql.ErrNoRows { return nil, repository.ErrNotFound }
	if err != nil { return nil, err }
	if avatarURL.Valid { user.AvatarURL = avatarURL.String }
	if lastSeenAt.Valid { user.LastSeenAt = &lastSeenAt.Time }
	return user, nil
}

func (r *teamRepository) GetByEmail(email string) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRow(`SELECT id, organization_id, email, name, role FROM users WHERE email = $1`, email).Scan(&user.ID, &user.OrganizationID, &user.Email, &user.Name, &user.Role)
	if err == sql.ErrNoRows { return nil, repository.ErrNotFound }
	return user, err
}

func (r *teamRepository) ExistsByEmail(email string) (bool, error) {
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", email).Scan(&exists)
	return exists, err
}

func (r *teamRepository) CreateMember(user *models.User, passwordHash string) error {
	_, err := r.db.Exec(`
		INSERT INTO users (id, organization_id, email, password_hash, name, role, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, user.ID, user.OrganizationID, user.Email, passwordHash, user.Name, user.Role, user.Status)
	return err
}

func (r *teamRepository) UpdateMember(orgID, userID uuid.UUID, input *models.UpdateUserInput) error {
	res, err := r.db.Exec(`
		UPDATE users SET updated_at = NOW(), name = COALESCE(NULLIF($1, ''), name), avatar_url = COALESCE(NULLIF($2, ''), avatar_url)
		WHERE id = $3 AND organization_id = $4
	`, input.Name, input.AvatarURL, userID, orgID)
	if err != nil { return err }
	rows, _ := res.RowsAffected()
	if rows == 0 { return repository.ErrNotFound }
	return nil
}

func (r *teamRepository) UpdateMemberRole(orgID, userID uuid.UUID, role models.UserRole) error {
	res, err := r.db.Exec(`UPDATE users SET role = $1, updated_at = NOW() WHERE id = $2 AND organization_id = $3`, role, userID, orgID)
	if err != nil { return err }
	rows, _ := res.RowsAffected()
	if rows == 0 { return repository.ErrNotFound }
	return nil
}

func (r *teamRepository) DeleteMember(orgID, userID uuid.UUID) error {
	res, err := r.db.Exec(`DELETE FROM users WHERE id = $1 AND organization_id = $2`, userID, orgID)
	if err != nil { return err }
	rows, _ := res.RowsAffected()
	if rows == 0 { return repository.ErrNotFound }
	return nil
}

func (r *teamRepository) UpdateStatus(userID uuid.UUID, status models.UserStatus) error {
	_, err := r.db.Exec(`UPDATE users SET status = $1, last_seen_at = NOW(), updated_at = NOW() WHERE id = $2`, status, userID)
	return err
}

func (r *teamRepository) GetAdminCount(orgID uuid.UUID) (int, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM users WHERE organization_id = $1 AND role = 'admin'`, orgID).Scan(&count)
	return count, err
}

func (r *teamRepository) GetOrganization(orgID uuid.UUID) (*models.Organization, error) {
	org := &models.Organization{}
	err := r.db.QueryRow(`
		SELECT id, name, slug, plan, subscription_status, ai_credits_limit, ai_credits_used, message_usage_limit, message_usage_used, billing_cycle_start, created_at, updated_at 
		FROM organizations WHERE id = $1
	`, orgID).Scan(&org.ID, &org.Name, &org.Slug, &org.Plan, &org.SubscriptionStatus, &org.AICreditsLimit, &org.AICreditsUsed, &org.MessageUsageLimit, &org.MessageUsageUsed, &org.BillingCycleStart, &org.CreatedAt, &org.UpdatedAt)
	if err != nil { return nil, err }
	r.db.QueryRow("SELECT COUNT(*) FROM users WHERE organization_id = $1", orgID).Scan(&org.UserCount)
	r.db.QueryRow("SELECT COUNT(*) FROM channels WHERE organization_id = $1", orgID).Scan(&org.ChannelCount)
	return org, nil
}

func (r *teamRepository) UpdateOrganization(orgID uuid.UUID, name, plan string, aiLimit, msgLimit int) error {
	query := "UPDATE organizations SET updated_at = NOW()"
	args := []interface{}{}
	if name != "" {
		query += ", name = $1"
		args = append(args, name)
	}
	if plan != "" {
		argIdx := len(args) + 1
		query += ", plan = $" + string(rune('0'+argIdx)) + ", ai_credits_limit = $" + string(rune('0'+argIdx+1)) + ", message_usage_limit = $" + string(rune('0'+argIdx+2)) + ", billing_cycle_start = NOW(), ai_credits_used = 0, message_usage_used = 0"
		args = append(args, plan, aiLimit, msgLimit)
	}
	query += " WHERE id = $" + string(rune('0'+len(args)+1))
	args = append(args, orgID)
	_, err := r.db.Exec(query, args...)
	return err
}
