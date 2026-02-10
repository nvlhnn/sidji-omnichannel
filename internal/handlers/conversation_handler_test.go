package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sidji-omnichannel/internal/models"
	"github.com/sidji-omnichannel/internal/services"
	"github.com/sidji-omnichannel/internal/testutil"
	"github.com/sidji-omnichannel/internal/websocket"
)

func setupConversationTestRouter(t *testing.T) (*gin.Engine, *sql.DB, *models.Organization, *models.User, *models.Conversation) {
	db := testutil.TestDB(t)

	// Setup test data
	org := testutil.CreateTestOrganization(t, db)
	user := testutil.CreateTestUser(t, db, org.ID, models.RoleAgent)
	channel := testutil.CreateTestChannel(t, db, org.ID)
	contact := testutil.CreateTestContact(t, db, org.ID)
	conv := testutil.CreateTestConversation(t, db, org.ID, channel.ID, contact.ID)

	// Create services
	convService := services.NewConversationService(db)
	msgService := services.NewMessageService(db)
	channelService := services.NewChannelService(db, testutil.TestConfig()) // nil meta client for tests
	contactService := services.NewContactService(db)
	hub := websocket.NewHub()

	handler := NewConversationHandler(convService, msgService, channelService, contactService, nil, hub)

	router := gin.New()
	api := router.Group("/api")
	api.Use(func(c *gin.Context) {
		c.Set("user_id", user.ID)
		c.Set("organization_id", org.ID)
		c.Next()
	})

	api.GET("/conversations", handler.List)
	api.GET("/conversations/:id", handler.Get)
	api.POST("/conversations/:id/assign", handler.Assign)
	api.PATCH("/conversations/:id/status", handler.UpdateStatus)
	api.GET("/conversations/:id/messages", handler.GetMessages)
	api.POST("/conversations/:id/messages", handler.SendMessage)

	return router, db, org, user, conv
}

func TestConversationHandler_List(t *testing.T) {
	router, db, _, _, _ := setupConversationTestRouter(t)
	defer db.Close()
	defer testutil.CleanupTestData(t, db)

	tests := []struct {
		name         string
		queryParams  string
		expectedCode int
	}{
		{
			name:         "list all conversations",
			queryParams:  "",
			expectedCode: http.StatusOK,
		},
		{
			name:         "list with status filter",
			queryParams:  "?status=open",
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
			req, _ := http.NewRequest("GET", "/api/conversations"+tt.queryParams, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.expectedCode, w.Code, w.Body.String())
			}
		})
	}
}

func TestConversationHandler_Get(t *testing.T) {
	router, db, _, _, conv := setupConversationTestRouter(t)
	defer db.Close()
	defer testutil.CleanupTestData(t, db)

	tests := []struct {
		name         string
		convID       string
		expectedCode int
	}{
		{
			name:         "existing conversation",
			convID:       conv.ID.String(),
			expectedCode: http.StatusOK,
		},
		{
			name:         "non-existent conversation",
			convID:       uuid.New().String(),
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "invalid UUID",
			convID:       "invalid-uuid",
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/api/conversations/"+tt.convID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.expectedCode, w.Code, w.Body.String())
			}
		})
	}
}

func TestConversationHandler_Assign(t *testing.T) {
	router, db, _, user, conv := setupConversationTestRouter(t)
	defer db.Close()
	defer testutil.CleanupTestData(t, db)

	tests := []struct {
		name         string
		convID       string
		body         interface{}
		expectedCode int
	}{
		{
			name:   "assign conversation",
			convID: conv.ID.String(),
			body: models.AssignConversationInput{
				UserID: user.ID,
			},
			expectedCode: http.StatusOK,
		},
		{
			name:   "non-existent conversation",
			convID: uuid.New().String(),
			body: models.AssignConversationInput{
				UserID: user.ID,
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "invalid request body",
			convID:       conv.ID.String(),
			body:         map[string]string{"invalid": "body"},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			req, _ := http.NewRequest("POST", "/api/conversations/"+tt.convID+"/assign", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.expectedCode, w.Code, w.Body.String())
			}
		})
	}
}

func TestConversationHandler_UpdateStatus(t *testing.T) {
	router, db, _, _, conv := setupConversationTestRouter(t)
	defer db.Close()
	defer testutil.CleanupTestData(t, db)

	tests := []struct {
		name         string
		convID       string
		body         interface{}
		expectedCode int
	}{
		{
			name:   "update to resolved",
			convID: conv.ID.String(),
			body: models.UpdateConversationStatusInput{
				Status: models.ConversationStatusResolved,
			},
			expectedCode: http.StatusOK,
		},
		{
			name:   "update to closed",
			convID: conv.ID.String(),
			body: models.UpdateConversationStatusInput{
				Status: models.ConversationStatusClosed,
			},
			expectedCode: http.StatusOK,
		},
		{
			name:   "non-existent conversation",
			convID: uuid.New().String(),
			body: models.UpdateConversationStatusInput{
				Status: models.ConversationStatusClosed,
			},
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			req, _ := http.NewRequest("PATCH", "/api/conversations/"+tt.convID+"/status", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.expectedCode, w.Code, w.Body.String())
			}
		})
	}
}

func TestConversationHandler_GetMessages(t *testing.T) {
	router, db, _, _, conv := setupConversationTestRouter(t)
	defer db.Close()
	defer testutil.CleanupTestData(t, db)

	tests := []struct {
		name         string
		convID       string
		queryParams  string
		expectedCode int
	}{
		{
			name:         "get messages",
			convID:       conv.ID.String(),
			queryParams:  "",
			expectedCode: http.StatusOK,
		},
		{
			name:         "get messages with pagination",
			convID:       conv.ID.String(),
			queryParams:  "?limit=10",
			expectedCode: http.StatusOK,
		},
		{
			name:         "non-existent conversation",
			convID:       uuid.New().String(),
			queryParams:  "",
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", fmt.Sprintf("/api/conversations/%s/messages%s", tt.convID, tt.queryParams), nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.expectedCode, w.Code, w.Body.String())
			}
		})
	}
}

func TestConversationHandler_SendMessage(t *testing.T) {
	router, db, _, _, conv := setupConversationTestRouter(t)
	defer db.Close()
	defer testutil.CleanupTestData(t, db)

	tests := []struct {
		name         string
		convID       string
		body         interface{}
		expectedCode int
	}{
		{
			name:   "send text message",
			convID: conv.ID.String(),
			body: models.SendMessageInput{
				Content:     "Hello, how can I help you?",
				MessageType: models.MessageTypeText,
			},
			expectedCode: http.StatusCreated,
		},
		{
			name:   "send image message",
			convID: conv.ID.String(),
			body: models.SendMessageInput{
				Content:     "Here is an image",
				MessageType: models.MessageTypeImage,
				MediaURL:    "https://example.com/image.jpg",
			},
			expectedCode: http.StatusCreated,
		},
		{
			name:   "non-existent conversation",
			convID: uuid.New().String(),
			body: models.SendMessageInput{
				Content:     "Test message",
				MessageType: models.MessageTypeText,
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "empty message",
			convID:       conv.ID.String(),
			body:         map[string]string{},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			req, _ := http.NewRequest("POST", "/api/conversations/"+tt.convID+"/messages", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.expectedCode, w.Code, w.Body.String())
			}
		})
	}
}
