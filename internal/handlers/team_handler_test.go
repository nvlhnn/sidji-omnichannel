package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sidji-omnichannel/internal/adapters/db/postgres"
	"github.com/sidji-omnichannel/internal/models"
	"github.com/sidji-omnichannel/internal/services"
	"github.com/sidji-omnichannel/internal/testutil"
)

func setupTeamTestRouter(t *testing.T, db *sql.DB) (*gin.Engine, *models.Organization, *models.User) {
	// Setup test data
	org := testutil.CreateTestOrganization(t, db)
	// Create admin user
	admin := testutil.CreateTestUser(t, db, org.ID, models.RoleAdmin)

	// Create services
	teamService := services.NewTeamService(postgres.NewTeamRepository(db))

	handler := NewTeamHandler(teamService)

	router := gin.New()
	api := router.Group("/api")
	api.Use(func(c *gin.Context) {
		c.Set("user_id", admin.ID)
		c.Set("organization_id", org.ID)
		c.Set("role", admin.Role)
		c.Next()
	})

	api.GET("/team/members", handler.ListMembers)
	api.GET("/team/members/:id", handler.GetMember)
	api.POST("/team/members", handler.InviteMember)
	api.PATCH("/team/members/:id", handler.UpdateMember)
	api.PATCH("/team/members/:id/role", handler.UpdateMemberRole)
	api.DELETE("/team/members/:id", handler.RemoveMember)
	api.GET("/team/organization", handler.GetOrganization)
	api.PATCH("/team/organization", handler.UpdateOrganization)

	return router, org, admin
}

func TestTeamHandler_ListMembers(t *testing.T) {
	db := testutil.TestDB(t)
	defer db.Close()
	defer testutil.CleanupTestData(t, db)

	router, _, _ := setupTeamTestRouter(t, db)

	tests := []struct {
		name         string
		expectedCode int
	}{
		{
			name:         "list all members",
			expectedCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/api/team/members", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.expectedCode, w.Code, w.Body.String())
			}
		})
	}
}

func TestTeamHandler_GetMember(t *testing.T) {
	db := testutil.TestDB(t)
	defer db.Close()
	defer testutil.CleanupTestData(t, db)

	router, _, admin := setupTeamTestRouter(t, db)

	tests := []struct {
		name         string
		memberID     string
		expectedCode int
	}{
		{
			name:         "get existing member",
			memberID:     admin.ID.String(),
			expectedCode: http.StatusOK,
		},
		{
			name:         "non-existent member",
			memberID:     uuid.New().String(),
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "invalid UUID",
			memberID:     "invalid-uuid",
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/api/team/members/"+tt.memberID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.expectedCode, w.Code, w.Body.String())
			}
		})
	}
}

func TestTeamHandler_InviteMember(t *testing.T) {
	db := testutil.TestDB(t)
	defer db.Close()
	defer testutil.CleanupTestData(t, db)

	router, _, _ := setupTeamTestRouter(t, db)

	tests := []struct {
		name         string
		body         interface{}
		expectedCode int
	}{
		{
			name: "invite valid member",
			body: models.InviteUserInput{
				Email: "invited@example.com",
				Name:  "Invited User",
				Role:  models.RoleAgent,
			},
			expectedCode: http.StatusCreated,
		},
		{
			name: "missing email",
			body: models.InviteUserInput{
				Name: "No Email User",
				Role: models.RoleAgent,
			},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			req, _ := http.NewRequest("POST", "/api/team/members", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.expectedCode, w.Code, w.Body.String())
			}
		})
	}
}

func TestTeamHandler_UpdateMember(t *testing.T) {
	db := testutil.TestDB(t)
	defer db.Close()
	defer testutil.CleanupTestData(t, db)

	router, _, admin := setupTeamTestRouter(t, db)

	tests := []struct {
		name         string
		memberID     string
		body         interface{}
		expectedCode int
	}{
		{
			name:     "update existing member",
			memberID: admin.ID.String(),
			body: models.UpdateUserInput{
				Name: "Updated Admin Name",
			},
			expectedCode: http.StatusOK,
		},
		{
			name:     "update non-existent member",
			memberID: uuid.New().String(),
			body: models.UpdateUserInput{
				Name: "Ghost User",
			},
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			req, _ := http.NewRequest("PATCH", "/api/team/members/"+tt.memberID, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.expectedCode, w.Code, w.Body.String())
			}
		})
	}
}

func TestTeamHandler_UpdateMemberRole(t *testing.T) {
	db := testutil.TestDB(t)
	defer db.Close()
	defer testutil.CleanupTestData(t, db)

	router, _, admin := setupTeamTestRouter(t, db)

	tests := []struct {
		name         string
		memberID     string
		role         models.UserRole
		expectedCode int
	}{
		{
			name:         "update role valid",
			memberID:     admin.ID.String(),
			role:         models.RoleSupervisor,
			expectedCode: http.StatusOK,
		},
		{
			name:         "non-existent member",
			memberID:     uuid.New().String(),
			role:         models.RoleAgent,
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(map[string]interface{}{"role": tt.role})
			req, _ := http.NewRequest("PATCH", "/api/team/members/"+tt.memberID+"/role", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.expectedCode, w.Code, w.Body.String())
			}
		})
	}
}

func TestTeamHandler_RemoveMember(t *testing.T) {
	db := testutil.TestDB(t)
	defer db.Close()
	defer testutil.CleanupTestData(t, db)

	router, org, _ := setupTeamTestRouter(t, db)

	// Create a secondary user to delete
	user := testutil.CreateTestUser(t, db, org.ID, models.RoleAgent)

	tests := []struct {
		name         string
		memberID     string
		expectedCode int
	}{
		{
			name:         "remove existing member",
			memberID:     user.ID.String(),
			expectedCode: http.StatusOK,
		},
		{
			name:         "remove non-existent member",
			memberID:     uuid.New().String(),
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("DELETE", "/api/team/members/"+tt.memberID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.expectedCode, w.Code, w.Body.String())
			}
		})
	}
}

func TestTeamHandler_GetOrganization(t *testing.T) {
	db := testutil.TestDB(t)
	defer db.Close()
	defer testutil.CleanupTestData(t, db)

	router, _, _ := setupTeamTestRouter(t, db)

	req, _ := http.NewRequest("GET", "/api/team/organization", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusOK, w.Code, w.Body.String())
	}
}

func TestTeamHandler_UpdateOrganization(t *testing.T) {
	db := testutil.TestDB(t)
	defer db.Close()
	defer testutil.CleanupTestData(t, db)

	router, _, _ := setupTeamTestRouter(t, db)

	tests := []struct {
		name         string
		body         interface{}
		expectedCode int
	}{
		{
			name:         "update name",
			body:         map[string]string{"name": "Updated Org Name"},
			expectedCode: http.StatusOK,
		},
		{
			name:         "invalid name",
			body:         map[string]string{"name": "A"}, // too short
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			req, _ := http.NewRequest("PATCH", "/api/team/organization", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.expectedCode, w.Code, w.Body.String())
			}
		})
	}
}
