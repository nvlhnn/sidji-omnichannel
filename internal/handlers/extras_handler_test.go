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

func setupExtrasTestRouter(t *testing.T) (*gin.Engine, *models.Organization, *models.User) {
	db := testutil.TestDB(t)

	org := testutil.CreateTestOrganization(t, db)
	user := testutil.CreateTestUser(t, db, org.ID, models.RoleAgent)

	cannedService := services.NewCannedResponseService(db)
	labelService := services.NewLabelService(db)

	cannedHandler := NewCannedResponseHandler(cannedService)
	labelHandler := NewLabelHandler(labelService)

	router := gin.New()
	api := router.Group("/api")
	api.Use(func(c *gin.Context) {
		c.Set("user_id", user.ID)
		c.Set("organization_id", org.ID)
		c.Next()
	})

	api.GET("/canned-responses", cannedHandler.List)
	api.GET("/canned-responses/search", cannedHandler.Search)
	api.POST("/canned-responses", cannedHandler.Create)
	api.PUT("/canned-responses/:id", cannedHandler.Update)
	api.DELETE("/canned-responses/:id", cannedHandler.Delete)

	api.GET("/labels", labelHandler.List)
	api.POST("/labels", labelHandler.Create)
	api.PUT("/labels/:id", labelHandler.Update)
	api.DELETE("/labels/:id", labelHandler.Delete)

	return router, org, user
}

func TestCannedResponseHandler_Create(t *testing.T) {
	router, _, _ := setupExtrasTestRouter(t)
	db := testutil.TestDB(t)
	defer db.Close()
	defer testutil.CleanupTestData(t, db)

	tests := []struct {
		name         string
		body         interface{}
		expectedCode int
	}{
		{
			name: "create valid response",
			body: models.CreateCannedResponseInput{
				Shortcut: "hello",
				Title:    "Hello Greeting",
				Content:  "Hello! How can I help you?",
			},
			expectedCode: http.StatusCreated,
		},
		{
			name: "missing shortcut",
			body: models.CreateCannedResponseInput{
				Content: "Content without shortcut",
			},
			expectedCode: http.StatusBadRequest, // assuming model validation
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			req, _ := http.NewRequest("POST", "/api/canned-responses", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.expectedCode, w.Code, w.Body.String())
			}
		})
	}
}

func TestCannedResponseHandler_List(t *testing.T) {
	router, _, _ := setupExtrasTestRouter(t)
	db := testutil.TestDB(t)
	defer db.Close()
	defer testutil.CleanupTestData(t, db)

	req, _ := http.NewRequest("GET", "/api/canned-responses", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusOK, w.Code, w.Body.String())
	}
}

func TestLabelHandler_Create(t *testing.T) {
	router, _, _ := setupExtrasTestRouter(t)
	db := testutil.TestDB(t)
	defer db.Close()
	defer testutil.CleanupTestData(t, db)

	tests := []struct {
		name         string
		body         interface{}
		expectedCode int
	}{
		{
			name: "create valid label",
			body: models.CreateLabelInput{
				Name:  "Bug",
				Color: "#FF0000",
			},
			expectedCode: http.StatusCreated,
		},
		{
			name: "missing name",
			body: models.CreateLabelInput{
				Color: "#00FF00",
			},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			req, _ := http.NewRequest("POST", "/api/labels", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.expectedCode, w.Code, w.Body.String())
			}
		})
	}
}

func TestLabelHandler_List(t *testing.T) {
	router, _, _ := setupExtrasTestRouter(t)
	db := testutil.TestDB(t)
	defer db.Close()
	defer testutil.CleanupTestData(t, db)

	req, _ := http.NewRequest("GET", "/api/labels", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusOK, w.Code, w.Body.String())
	}
}
