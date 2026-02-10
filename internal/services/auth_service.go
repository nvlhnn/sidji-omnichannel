package services

import (
	"database/sql"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sidji-omnichannel/internal/config"
	"github.com/sidji-omnichannel/internal/models"
	"github.com/sidji-omnichannel/internal/subscription"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserExists         = errors.New("user with this email already exists")
	ErrUserNotFound       = errors.New("user not found")
)

// AuthService handles authentication logic
type AuthService struct {
	db  *sql.DB
	cfg *config.Config
}

// NewAuthService creates a new auth service
func NewAuthService(db *sql.DB, cfg *config.Config) *AuthService {
	return &AuthService{db: db, cfg: cfg}
}

// Register creates a new user and organization
func (s *AuthService) Register(input *models.RegisterInput) (*models.AuthResponse, error) {
	// Check if user exists
	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", input.Email).Scan(&exists)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrUserExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Start transaction
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Create organization
	org := &models.Organization{
		ID:                 uuid.New(),
		Name:               input.OrganizationName,
		Slug:               generateSlug(input.OrganizationName),
		Plan:               "starter",
		SubscriptionStatus: "active",
		AICreditsLimit:     10,
		AICreditsUsed:      0,
		BillingCycleStart:  time.Now(),
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	org.MessageUsageLimit = 1000
	org.MessageUsageUsed = 0

	_, err = tx.Exec(
		"INSERT INTO organizations (id, name, slug, plan, subscription_status, ai_credits_limit, ai_credits_used, message_usage_limit, message_usage_used, billing_cycle_start) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)",
		org.ID, org.Name, org.Slug, org.Plan, org.SubscriptionStatus, org.AICreditsLimit, org.AICreditsUsed, org.MessageUsageLimit, org.MessageUsageUsed, org.BillingCycleStart,
	)
	if err != nil {
		return nil, err
	}

	// Create user as admin
	user := &models.User{
		ID:             uuid.New(),
		OrganizationID: org.ID,
		Email:          input.Email,
		PasswordHash:   string(hashedPassword),
		Name:           input.Name,
		Role:           models.RoleAdmin,
		Status:         models.StatusOffline,
	}

	_, err = tx.Exec(
		`INSERT INTO users (id, organization_id, email, password_hash, name, role, status) 
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		user.ID, user.OrganizationID, user.Email, user.PasswordHash, user.Name, user.Role, user.Status,
	)
	if err != nil {
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// Generate token
	token, expiresIn, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	// Populate initial compliance (1st user)
	org.UserCount = 1
	org.ChannelCount = 0
	org.IsOverLimit = false

	return &models.AuthResponse{
		User:         user.ToPublic(),
		Organization: org,
		AccessToken:  token,
		ExpiresIn:    expiresIn,
	}, nil
}

// Login authenticates a user
func (s *AuthService) Login(input *models.LoginInput) (*models.AuthResponse, error) {
	user := &models.User{}
	org := &models.Organization{}

	// Get user and organization
	var avatarURL sql.NullString
	err := s.db.QueryRow(`
		SELECT u.id, u.organization_id, u.email, u.password_hash, u.name, u.role, u.status, u.avatar_url,
		       o.id, o.name, o.slug, o.plan, o.subscription_status, o.ai_credits_limit, o.ai_credits_used, o.message_usage_limit, o.message_usage_used, o.billing_cycle_start
		FROM users u
		JOIN organizations o ON o.id = u.organization_id
		WHERE u.email = $1
	`, input.Email).Scan(
		&user.ID, &user.OrganizationID, &user.Email, &user.PasswordHash, &user.Name, &user.Role, &user.Status, &avatarURL,
		&org.ID, &org.Name, &org.Slug, &org.Plan, &org.SubscriptionStatus, &org.AICreditsLimit, &org.AICreditsUsed, &org.MessageUsageLimit, &org.MessageUsageUsed, &org.BillingCycleStart,
	)

	if err == sql.ErrNoRows {
		return nil, ErrInvalidCredentials
	}
	if err != nil {
		return nil, err
	}

	if avatarURL.Valid {
		user.AvatarURL = avatarURL.String
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Update last seen
	_, _ = s.db.Exec("UPDATE users SET last_seen_at = NOW(), status = 'online' WHERE id = $1", user.ID)

	// Fetch counts for compliance
	err = s.db.QueryRow("SELECT COUNT(*) FROM users WHERE organization_id = $1", org.ID).Scan(&org.UserCount)
	if err != nil {
		return nil, err
	}
	err = s.db.QueryRow("SELECT COUNT(*) FROM channels WHERE organization_id = $1", org.ID).Scan(&org.ChannelCount)
	if err != nil {
		return nil, err
	}
	org.IsOverLimit = !subscription.IsCompliance(org.Plan, org.UserCount, org.ChannelCount)

	// Generate token
	token, expiresIn, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		User:         user.ToPublic(),
		Organization: org,
		AccessToken:  token,
		ExpiresIn:    expiresIn,
	}, nil
}

// GetUserByID retrieves a user by ID
func (s *AuthService) GetUserByID(userID uuid.UUID) (*models.User, error) {
	user := &models.User{}
	var avatarURL sql.NullString

	err := s.db.QueryRow(`
		SELECT id, organization_id, email, name, role, status, avatar_url, created_at
		FROM users WHERE id = $1
	`, userID).Scan(
		&user.ID, &user.OrganizationID, &user.Email, &user.Name, &user.Role, &user.Status, &avatarURL, &user.CreatedAt,
	)
	if err == nil && avatarURL.Valid {
		user.AvatarURL = avatarURL.String
	}
	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

// generateToken creates a JWT token for a user
func (s *AuthService) generateToken(user *models.User) (string, int64, error) {
	expiresIn := int64(24 * 60 * 60) // 24 hours in seconds

	claims := &models.Claims{
		UserID:         user.ID,
		OrganizationID: user.OrganizationID,
		Role:           string(user.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiresIn) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.cfg.App.Secret))
	if err != nil {
		return "", 0, err
	}

	return tokenString, expiresIn, nil
}

// generateSlug creates a URL-friendly slug from a name
func generateSlug(name string) string {
	slug := ""
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			slug += string(r)
		} else if r >= 'A' && r <= 'Z' {
			slug += string(r + 32) // lowercase
		} else if r == ' ' || r == '-' {
			slug += "-"
		}
	}
	return slug + "-" + uuid.New().String()[:8]
}

// GetMe retrieves the current user and their organization with compliance status
func (s *AuthService) GetMe(userID uuid.UUID) (*models.AuthResponse, error) {
	user := &models.User{}
	org := &models.Organization{}
	var avatarURL sql.NullString

	err := s.db.QueryRow(`
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
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	if avatarURL.Valid {
		user.AvatarURL = avatarURL.String
	}

	// Fetch counts for compliance
	err = s.db.QueryRow("SELECT COUNT(*) FROM users WHERE organization_id = $1", org.ID).Scan(&org.UserCount)
	if err != nil {
		return nil, err
	}
	err = s.db.QueryRow("SELECT COUNT(*) FROM channels WHERE organization_id = $1", org.ID).Scan(&org.ChannelCount)
	if err != nil {
		return nil, err
	}

	org.IsOverLimit = !subscription.IsCompliance(org.Plan, org.UserCount, org.ChannelCount)

	return &models.AuthResponse{
		User:         user.ToPublic(),
		Organization: org,
	}, nil
}
