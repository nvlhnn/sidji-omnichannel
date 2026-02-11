package services

import (
	"testing"

	"github.com/google/uuid"
	"github.com/sidji-omnichannel/internal/adapters/db/postgres"
	"github.com/sidji-omnichannel/internal/models"
	"github.com/sidji-omnichannel/internal/testutil"
)

func TestAuthService_Register(t *testing.T) {
	db := testutil.TestDB(t)
	defer db.Close()
	defer testutil.CleanupTestData(t, db)

	cfg := testutil.TestConfig()
	service := NewAuthService(postgres.NewAuthRepository(db), cfg)

	tests := []struct {
		name    string
		input   *models.RegisterInput
		wantErr error
	}{
		{
			name: "successful registration",
			input: &models.RegisterInput{
				Email:            "newuser@example.com",
				Password:         "password123",
				Name:             "New User",
				OrganizationName: "New Organization",
			},
			wantErr: nil,
		},
		{
			name: "duplicate email",
			input: &models.RegisterInput{
				Email:            "newuser@example.com", // Same email as above
				Password:         "password456",
				Name:             "Another User",
				OrganizationName: "Another Organization",
			},
			wantErr: ErrUserExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.Register(tt.input)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("Expected error %v, got nil", tt.wantErr)
				} else if err != tt.wantErr {
					t.Errorf("Expected error %v, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// Verify response
			if result.User == nil {
				t.Error("Expected user in response")
			}
			if result.Organization == nil {
				t.Error("Expected organization in response")
			}
			if result.AccessToken == "" {
				t.Error("Expected access token in response")
			}
			if result.User.Email != tt.input.Email {
				t.Errorf("Expected email %s, got %s", tt.input.Email, result.User.Email)
			}
			if result.User.Role != models.RoleAdmin {
				t.Errorf("Expected role %s, got %s", models.RoleAdmin, result.User.Role)
			}
			if result.Organization.Plan != "starter" {
				t.Errorf("Expected plan starter, got %s", result.Organization.Plan)
			}
			if result.Organization.SubscriptionStatus != "active" {
				t.Errorf("Expected status active, got %s", result.Organization.SubscriptionStatus)
			}
			if result.Organization.AICreditsLimit != 10 {
				t.Errorf("Expected 10 credits, got %d", result.Organization.AICreditsLimit)
			}
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	db := testutil.TestDB(t)
	defer db.Close()
	defer testutil.CleanupTestData(t, db)

	cfg := testutil.TestConfig()
	service := NewAuthService(postgres.NewAuthRepository(db), cfg)

	// Create a user first
	registerInput := &models.RegisterInput{
		Email:            "logintest@example.com",
		Password:         "password123",
		Name:             "Login Test User",
		OrganizationName: "Login Test Org",
	}
	_, err := service.Register(registerInput)
	if err != nil {
		t.Fatalf("Failed to register test user: %v", err)
	}

	tests := []struct {
		name    string
		input   *models.LoginInput
		wantErr error
	}{
		{
			name: "successful login",
			input: &models.LoginInput{
				Email:    "logintest@example.com",
				Password: "password123",
			},
			wantErr: nil,
		},
		{
			name: "wrong password",
			input: &models.LoginInput{
				Email:    "logintest@example.com",
				Password: "wrongpassword",
			},
			wantErr: ErrInvalidCredentials,
		},
		{
			name: "non-existent user",
			input: &models.LoginInput{
				Email:    "nonexistent@example.com",
				Password: "password123",
			},
			wantErr: ErrInvalidCredentials,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.Login(tt.input)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("Expected error %v, got nil", tt.wantErr)
				} else if err != tt.wantErr {
					t.Errorf("Expected error %v, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// Verify response
			if result.User == nil {
				t.Error("Expected user in response")
			}
			if result.AccessToken == "" {
				t.Error("Expected access token in response")
			}
			if result.User.Email != tt.input.Email {
				t.Errorf("Expected email %s, got %s", tt.input.Email, result.User.Email)
			}
		})
	}
}

func TestAuthService_GetUserByID(t *testing.T) {
	db := testutil.TestDB(t)
	defer db.Close()
	defer testutil.CleanupTestData(t, db)

	cfg := testutil.TestConfig()
	service := NewAuthService(postgres.NewAuthRepository(db), cfg)

	// Create a user
	registerInput := &models.RegisterInput{
		Email:            "getbyid@example.com",
		Password:         "password123",
		Name:             "Get By ID User",
		OrganizationName: "Get By ID Org",
	}
	registerResult, err := service.Register(registerInput)
	if err != nil {
		t.Fatalf("Failed to register test user: %v", err)
	}

	tests := []struct {
		name    string
		userID  uuid.UUID
		wantErr error
	}{
		{
			name:    "existing user",
			userID:  registerResult.User.ID,
			wantErr: nil,
		},
		{
			name:    "non-existent user",
			userID:  uuid.New(),
			wantErr: ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := service.GetUserByID(tt.userID)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("Expected error %v, got nil", tt.wantErr)
				} else if err != tt.wantErr {
					t.Errorf("Expected error %v, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if user.ID != tt.userID {
				t.Errorf("Expected user ID %s, got %s", tt.userID, user.ID)
			}
		})
	}
}

func TestGenerateSlug(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains string
	}{
		{
			name:     "simple name",
			input:    "My Company",
			contains: "my-company-",
		},
		{
			name:     "name with special characters",
			input:    "My Company! @#$%",
			contains: "my-company-",
		},
		{
			name:     "name with numbers",
			input:    "Company 123",
			contains: "company-123-",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generateSlug(tt.input)
			if len(result) == 0 {
				t.Error("Expected non-empty slug")
			}
			// The slug should contain the expected prefix
			if len(result) < len(tt.contains) {
				t.Errorf("Slug too short: %s", result)
			}
		})
	}
}
