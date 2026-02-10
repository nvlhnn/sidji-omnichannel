package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sidji-omnichannel/internal/models"
	"github.com/sidji-omnichannel/internal/services"
	"github.com/sidji-omnichannel/internal/testutil"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestAuthHandler_Register(t *testing.T) {
	db := testutil.TestDB(t)
	defer db.Close()
	defer testutil.CleanupTestData(t, db)

	cfg := testutil.TestConfig()
	authService := services.NewAuthService(db, cfg)
	handler := NewAuthHandler(authService)

	tests := []struct {
		name         string
		body         interface{}
		expectedCode int
	}{
		{
			name: "successful registration",
			body: models.RegisterInput{
				Email:            "test@example.com",
				Password:         "password123",
				Name:             "Test User",
				OrganizationName: "Test Organization",
			},
			expectedCode: http.StatusCreated,
		},
		{
			name: "invalid request - missing fields",
			body: map[string]string{
				"email": "test2@example.com",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "duplicate email",
			body: models.RegisterInput{
				Email:            "test@example.com", // Same as first test
				Password:         "password456",
				Name:             "Another User",
				OrganizationName: "Another Org",
			},
			expectedCode: http.StatusConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.POST("/auth/register", handler.Register)

			body, _ := json.Marshal(tt.body)
			req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.expectedCode, w.Code, w.Body.String())
			}

			if tt.expectedCode == http.StatusCreated {
				var response models.AuthResponse
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}
				if response.AccessToken == "" {
					t.Error("Expected access token in response")
				}
			}
		})
	}
}

func TestAuthHandler_Login(t *testing.T) {
	db := testutil.TestDB(t)
	defer db.Close()
	defer testutil.CleanupTestData(t, db)

	cfg := testutil.TestConfig()
	authService := services.NewAuthService(db, cfg)
	handler := NewAuthHandler(authService)

	// Create a user first
	_, err := authService.Register(&models.RegisterInput{
		Email:            "logintest@example.com",
		Password:         "password123",
		Name:             "Login Test User",
		OrganizationName: "Login Test Org",
	})
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	tests := []struct {
		name         string
		body         interface{}
		expectedCode int
	}{
		{
			name: "successful login",
			body: models.LoginInput{
				Email:    "logintest@example.com",
				Password: "password123",
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "wrong password",
			body: models.LoginInput{
				Email:    "logintest@example.com",
				Password: "wrongpassword",
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "non-existent user",
			body: models.LoginInput{
				Email:    "nonexistent@example.com",
				Password: "password123",
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "invalid request",
			body: map[string]string{
				"email": "test@example.com",
			},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.POST("/auth/login", handler.Login)

			body, _ := json.Marshal(tt.body)
			req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.expectedCode, w.Code, w.Body.String())
			}

			if tt.expectedCode == http.StatusOK {
				var response models.AuthResponse
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}
				if response.AccessToken == "" {
					t.Error("Expected access token in response")
				}
			}
		})
	}
}

func TestAuthHandler_Me(t *testing.T) {
	db := testutil.TestDB(t)
	defer db.Close()
	defer testutil.CleanupTestData(t, db)

	cfg := testutil.TestConfig()
	authService := services.NewAuthService(db, cfg)
	handler := NewAuthHandler(authService)

	// Create a user
	result, err := authService.Register(&models.RegisterInput{
		Email:            "metest@example.com",
		Password:         "password123",
		Name:             "Me Test User",
		OrganizationName: "Me Test Org",
	})
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	tests := []struct {
		name         string
		setupContext func(*gin.Context)
		expectedCode int
	}{
		{
			name: "authenticated user",
			setupContext: func(c *gin.Context) {
				c.Set("user_id", result.User.ID)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "unauthenticated request",
			setupContext: func(c *gin.Context) {},
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.GET("/auth/me", func(c *gin.Context) {
				tt.setupContext(c)
				handler.Me(c)
			})

			req, _ := http.NewRequest("GET", "/auth/me", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.expectedCode, w.Code, w.Body.String())
			}

			if tt.expectedCode == http.StatusOK {
				var response models.AuthResponse
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}
				if response.User.Email != "metest@example.com" {
					t.Errorf("Expected email metest@example.com, got %s", response.User.Email)
				}
			}
		})
	}
}
