package postgres

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/sidji-omnichannel/internal/domain/ports/repository"
	"github.com/sidji-omnichannel/internal/models"
)

type authRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) repository.AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) ExistsByEmail(email string) (bool, error) {
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", email).Scan(&exists)
	return exists, err
}

func (r *authRepository) CreateRegisterTransaction(org *models.Organization, user *models.User) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(
		"INSERT INTO organizations (id, name, slug, plan, subscription_status, ai_credits_limit, ai_credits_used, message_usage_limit, message_usage_used, billing_cycle_start) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)",
		org.ID, org.Name, org.Slug, org.Plan, org.SubscriptionStatus, org.AICreditsLimit, org.AICreditsUsed, org.MessageUsageLimit, org.MessageUsageUsed, org.BillingCycleStart,
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		`INSERT INTO users (id, organization_id, email, password_hash, name, role, status) 
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		user.ID, user.OrganizationID, user.Email, user.PasswordHash, user.Name, user.Role, user.Status,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *authRepository) GetAuthDataByEmail(email string) (*models.User, *models.Organization, error) {
	user := &models.User{}
	org := &models.Organization{}
	var avatarURL sql.NullString

	err := r.db.QueryRow(`
		SELECT u.id, u.organization_id, u.email, u.password_hash, u.name, u.role, u.status, u.avatar_url,
		       o.id, o.name, o.slug, o.plan, o.subscription_status, o.ai_credits_limit, o.ai_credits_used, o.message_usage_limit, o.message_usage_used, o.billing_cycle_start
		FROM users u
		JOIN organizations o ON o.id = u.organization_id
		WHERE u.email = $1
	`, email).Scan(
		&user.ID, &user.OrganizationID, &user.Email, &user.PasswordHash, &user.Name, &user.Role, &user.Status, &avatarURL,
		&org.ID, &org.Name, &org.Slug, &org.Plan, &org.SubscriptionStatus, &org.AICreditsLimit, &org.AICreditsUsed, &org.MessageUsageLimit, &org.MessageUsageUsed, &org.BillingCycleStart,
	)

	if err == sql.ErrNoRows {
		return nil, nil, repository.ErrNotFound
	}
	if err != nil {
		return nil, nil, err
	}

	if avatarURL.Valid {
		user.AvatarURL = avatarURL.String
	}

	return user, org, nil
}

func (r *authRepository) GetAuthDataByID(userID uuid.UUID) (*models.User, *models.Organization, error) {
	user := &models.User{}
	org := &models.Organization{}
	var avatarURL sql.NullString

	err := r.db.QueryRow(`
		SELECT u.id, u.organization_id, u.email, u.name, u.role, u.status, u.avatar_url, u.created_at,
		       o.id, o.name, o.slug, o.plan, o.subscription_status, o.ai_credits_limit, o.ai_credits_used, o.message_usage_limit, o.message_usage_used, o.billing_cycle_start, o.created_at, o.updated_at
		FROM users u
		JOIN organizations o ON o.id = u.organization_id
		WHERE u.id = $1
	`, userID).Scan(
		&user.ID, &user.OrganizationID, &user.Email, &user.Name, &user.Role, &user.Status, &avatarURL, &user.CreatedAt,
		&org.ID, &org.Name, &org.Slug, &org.Plan, &org.SubscriptionStatus, &org.AICreditsLimit, &org.AICreditsUsed, &org.MessageUsageLimit, &org.MessageUsageUsed, &org.BillingCycleStart, &org.CreatedAt, &org.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil, repository.ErrNotFound
	}
	if err != nil {
		return nil, nil, err
	}

	if avatarURL.Valid {
		user.AvatarURL = avatarURL.String
	}

	return user, org, nil
}

func (r *authRepository) UpdateLastSeen(userID uuid.UUID) error {
	_, err := r.db.Exec("UPDATE users SET last_seen_at = NOW(), status = 'online' WHERE id = $1", userID)
	return err
}

func (r *authRepository) GetUserByID(userID uuid.UUID) (*models.User, error) {
	user := &models.User{}
	var avatarURL sql.NullString

	err := r.db.QueryRow(`
		SELECT id, organization_id, email, name, role, status, avatar_url, created_at
		FROM users WHERE id = $1
	`, userID).Scan(
		&user.ID, &user.OrganizationID, &user.Email, &user.Name, &user.Role, &user.Status, &avatarURL, &user.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, repository.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	if avatarURL.Valid {
		user.AvatarURL = avatarURL.String
	}
	return user, nil
}

func (r *authRepository) GetOrganizationCounts(orgID uuid.UUID) (userCount, channelCount int, err error) {
	err = r.db.QueryRow("SELECT COUNT(*) FROM users WHERE organization_id = $1", orgID).Scan(&userCount)
	if err != nil {
		return 0, 0, err
	}
	err = r.db.QueryRow("SELECT COUNT(*) FROM channels WHERE organization_id = $1", orgID).Scan(&channelCount)
	if err != nil {
		return 0, 0, err
	}
	return userCount, channelCount, nil
}
