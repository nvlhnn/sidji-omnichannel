package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func setupContactTestRouter(t *testing.T) (*gin.Engine, *models.Organization, *models.Contact) {
	db := testutil.TestDB(t)

	// Setup test data
	org := testutil.CreateTestOrganization(t, db)
	// Create user mostly for middleware context
	user := testutil.CreateTestUser(t, db, org.ID, models.RoleAgent)
	contact := testutil.CreateTestContact(t, db, org.ID)

	// Create services
	repo := postgres.NewContactRepository(db)
	contactService := services.NewContactService(repo)
	handler := NewContactHandler(contactService)

	router := gin.New()
	api := router.Group("/api")
	api.Use(func(c *gin.Context) {
		c.Set("user_id", user.ID)
		c.Set("organization_id", org.ID)
		c.Next()
	})

	api.GET("/contacts", handler.List)
	api.GET("/contacts/:id", handler.Get)
	api.POST("/contacts", handler.Create)
	api.PATCH("/contacts/:id", handler.Update)
	api.DELETE("/contacts/:id", handler.Delete)
	api.GET("/contacts/:id/conversations", handler.GetConversations)

	return router, org, contact
}

func TestContactHandler_List(t *testing.T) {
	router, _, _ := setupContactTestRouter(t)
	db := testutil.TestDB(t)
	defer db.Close()
	defer testutil.CleanupTestData(t, db)

	tests := []struct {
		name         string
		queryParams  string
		expectedCode int
	}{
		{
			name:         "list all contacts",
			queryParams:  "",
			expectedCode: http.StatusOK,
		},
		{
			name:         "list with pagination",
			queryParams:  "?page=1&limit=10",
			expectedCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/api/contacts"+tt.queryParams, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.expectedCode, w.Code, w.Body.String())
			}
		})
	}
}

func TestContactHandler_Get(t *testing.T) {
	router, _, contact := setupContactTestRouter(t)
	db := testutil.TestDB(t)
	defer db.Close()
	defer testutil.CleanupTestData(t, db)

	tests := []struct {
		name         string
		contactID    string
		expectedCode int
	}{
		{
			name:         "existing contact",
			contactID:    contact.ID.String(),
			expectedCode: http.StatusOK,
		},
		{
			name:         "non-existent contact",
			contactID:    uuid.New().String(),
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "invalid UUID",
			contactID:    "invalid-uuid",
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/api/contacts/"+tt.contactID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.expectedCode, w.Code, w.Body.String())
			}
		})
	}
}

func TestContactHandler_Create(t *testing.T) {
	router, _, _ := setupContactTestRouter(t)
	db := testutil.TestDB(t)
	defer db.Close()
	defer testutil.CleanupTestData(t, db)

	tests := []struct {
		name         string
		body         interface{}
		expectedCode int
	}{
		{
			name: "create valid contact",
			body: models.CreateContactInput{
				Name:  "New Contact",
				Phone: "+1234567890",
				Email: "new@example.com",
			},
			expectedCode: http.StatusCreated,
		},
		{
			name: "missing name",
			body: models.CreateContactInput{
				Phone: "+1234567890",
			},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			req, _ := http.NewRequest("POST", "/api/contacts", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.expectedCode, w.Code, w.Body.String())
			}
		})
	}
}

func TestContactHandler_Update(t *testing.T) {
	router, _, contact := setupContactTestRouter(t)
	db := testutil.TestDB(t)
	defer db.Close()
	defer testutil.CleanupTestData(t, db)

	tests := []struct {
		name         string
		contactID    string
		body         interface{}
		expectedCode int
	}{
		{
			name:      "update valid contact",
			contactID: contact.ID.String(),
			body: models.UpdateContactInput{
				Name: "Updated Name",
			},
			expectedCode: http.StatusOK,
		},
		{
			name:      "update non-existent contact",
			contactID: uuid.New().String(),
			body: models.UpdateContactInput{
				Name: "Updated Name",
			},
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			req, _ := http.NewRequest("PATCH", "/api/contacts/"+tt.contactID, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.expectedCode, w.Code, w.Body.String())
			}
		})
	}
}

func TestContactHandler_Delete(t *testing.T) {
	router, _, contact := setupContactTestRouter(t)
	db := testutil.TestDB(t)
	defer db.Close()
	defer testutil.CleanupTestData(t, db)

	tests := []struct {
		name         string
		contactID    string
		expectedCode int
	}{
		{
			name:         "delete existing contact",
			contactID:    contact.ID.String(),
			expectedCode: http.StatusOK,
		},
		{
			name:         "delete non-existent contact",
			contactID:    uuid.New().String(),
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("DELETE", "/api/contacts/"+tt.contactID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.expectedCode, w.Code, w.Body.String())
			}
		})
	}
}

func TestContactHandler_GetConversations(t *testing.T) {
	router, _, contact := setupContactTestRouter(t)
	db := testutil.TestDB(t)
	defer db.Close()
	defer testutil.CleanupTestData(t, db)

	tests := []struct {
		name         string
		contactID    string
		expectedCode int
	}{
		{
			name:         "get conversations for existing contact",
			contactID:    contact.ID.String(),
			expectedCode: http.StatusOK,
		},
		{
			name:         "non-existent contact",
			contactID:    uuid.New().String(), // Returns empty list with 200 usually, wait. Service returns empty list. Handler returns 200.
			expectedCode: http.StatusOK, // If contact not found service might return empty or error. Service `GetConversations` queries by `contactID` but doesn't check if contact exists first. So it returns empty list.
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", fmt.Sprintf("/api/contacts/%s/conversations", tt.contactID), nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.expectedCode, w.Code, w.Body.String())
			}
		})
	}
}
