package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sidji-omnichannel/internal/config"
	"github.com/sidji-omnichannel/internal/domain/ports/repository"
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
	repo repository.AuthRepository
	cfg  *config.Config
}

// NewAuthService creates a new auth service
func NewAuthService(repo repository.AuthRepository, cfg *config.Config) *AuthService {
	return &AuthService{repo: repo, cfg: cfg}
}

// Register creates a new user and organization
func (s *AuthService) Register(input *models.RegisterInput) (*models.AuthResponse, error) {
	// Check if user exists
	exists, err := s.repo.ExistsByEmail(input.Email)
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

	// Create objects
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
		MessageUsageLimit:  1000,
		MessageUsageUsed:   0,
	}

	user := &models.User{
		ID:             uuid.New(),
		OrganizationID: org.ID,
		Email:          input.Email,
		PasswordHash:   string(hashedPassword),
		Name:           input.Name,
		Role:           models.RoleAdmin,
		Status:         models.StatusOffline,
	}

	if err := s.repo.CreateRegisterTransaction(org, user); err != nil {
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
	user, org, err := s.repo.GetAuthDataByEmail(input.Email)
	if err == repository.ErrNotFound {
		return nil, ErrInvalidCredentials
	}
	if err != nil {
		return nil, err
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Update last seen
	_ = s.repo.UpdateLastSeen(user.ID)

	// Fetch counts for compliance
	userCount, channelCount, err := s.repo.GetOrganizationCounts(org.ID)
	if err != nil {
		return nil, err
	}
	org.UserCount = userCount
	org.ChannelCount = channelCount
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
	user, err := s.repo.GetUserByID(userID)
	if err == repository.ErrNotFound {
		return nil, ErrUserNotFound
	}
	return user, err
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
	user, org, err := s.repo.GetAuthDataByID(userID)
	if err == repository.ErrNotFound {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	// Fetch counts for compliance
	userCount, channelCount, err := s.repo.GetOrganizationCounts(org.ID)
	if err != nil {
		return nil, err
	}
	org.UserCount = userCount
	org.ChannelCount = channelCount
	org.IsOverLimit = !subscription.IsCompliance(org.Plan, org.UserCount, org.ChannelCount)

	return &models.AuthResponse{
		User:         user.ToPublic(),
		Organization: org,
	}, nil
}

// GoogleLogin handles user login or registration via Google OAuth
func (s *AuthService) GoogleLogin(info *models.GoogleUserInfo) (*models.AuthResponse, error) {
	// 1. Try to fetch user by email
	user, org, err := s.repo.GetAuthDataByEmail(info.Email)
	
	if err == repository.ErrNotFound {
		// User doesn't exist yet, we must auto-register them
		// Generate random secure password, as they login via Google
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(uuid.New().String()), bcrypt.DefaultCost)
		
		orgName := info.GivenName + "'s Workspace"
		if info.GivenName == "" {
			orgName = info.Name + "'s Workspace"
		}

		org = &models.Organization{
			ID:                 uuid.New(),
			Name:               orgName,
			Slug:               generateSlug(orgName),
			Plan:               "starter",
			SubscriptionStatus: "active",
			AICreditsLimit:     10,
			AICreditsUsed:      0,
			BillingCycleStart:  time.Now(),
			CreatedAt:          time.Now(),
			UpdatedAt:          time.Now(),
			MessageUsageLimit:  1000,
			MessageUsageUsed:   0,
		}

		user = &models.User{
			ID:             uuid.New(),
			OrganizationID: org.ID,
			Email:          info.Email,
			PasswordHash:   string(hashedPassword),
			Name:           info.Name,
			Role:           models.RoleAdmin,
			Status:         models.StatusOffline,
			AvatarURL:      info.Picture,
		}

		if err := s.repo.CreateRegisterTransaction(org, user); err != nil {
			return nil, err
		}

		// Initial compliance
		org.UserCount = 1
		org.ChannelCount = 0
		org.IsOverLimit = false
	} else if err != nil {
		return nil, err
	} else {
		// User already exists, update their last seen and avatar
		_ = s.repo.UpdateLastSeen(user.ID)
		
		// If they don't have an avatar but Google provided one, we could update it
		// For simplicity, we just fetch counts for compliance
		userCount, channelCount, err := s.repo.GetOrganizationCounts(org.ID)
		if err == nil {
			org.UserCount = userCount
			org.ChannelCount = channelCount
			org.IsOverLimit = !subscription.IsCompliance(org.Plan, org.UserCount, org.ChannelCount)
		}
	}

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

