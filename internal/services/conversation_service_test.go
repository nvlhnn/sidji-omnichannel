package services

import (
	"testing"

	"github.com/google/uuid"
	"github.com/sidji-omnichannel/internal/models"
	"github.com/sidji-omnichannel/internal/testutil"
)

func TestConversationService_List(t *testing.T) {
	db := testutil.TestDB(t)
	defer db.Close()
	defer testutil.CleanupTestData(t, db)

	// Setup test data
	org := testutil.CreateTestOrganization(t, db)
	channel := testutil.CreateTestChannel(t, db, org.ID)
	contact := testutil.CreateTestContact(t, db, org.ID)
	_ = testutil.CreateTestConversation(t, db, org.ID, channel.ID, contact.ID)

	service := NewConversationService(db)

	tests := []struct {
		name           string
		orgID          uuid.UUID
		filter         *models.ConversationFilter
		wantCount      int
		wantErr        bool
	}{
		{
			name:  "list all conversations",
			orgID: org.ID,
			filter: &models.ConversationFilter{
				Page:  1,
				Limit: 20,
			},
			wantCount: 1,
			wantErr:   false,
		},
		{
			name:  "filter by open status",
			orgID: org.ID,
			filter: &models.ConversationFilter{
				Status: "open",
				Page:   1,
				Limit:  20,
			},
			wantCount: 1,
			wantErr:   false,
		},
		{
			name:  "filter by closed status (no results)",
			orgID: org.ID,
			filter: &models.ConversationFilter{
				Status: "closed",
				Page:   1,
				Limit:  20,
			},
			wantCount: 0,
			wantErr:   false,
		},
		{
			name:  "different organization (no results)",
			orgID: uuid.New(),
			filter: &models.ConversationFilter{
				Page:  1,
				Limit: 20,
			},
			wantCount: 0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conversations, total, err := service.List(tt.orgID, tt.filter)

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if len(conversations) != tt.wantCount {
				t.Errorf("Expected %d conversations, got %d", tt.wantCount, len(conversations))
			}

			if total < tt.wantCount {
				t.Errorf("Expected total >= %d, got %d", tt.wantCount, total)
			}
		})
	}
}

func TestConversationService_GetByID(t *testing.T) {
	db := testutil.TestDB(t)
	defer db.Close()
	defer testutil.CleanupTestData(t, db)

	// Setup test data
	org := testutil.CreateTestOrganization(t, db)
	channel := testutil.CreateTestChannel(t, db, org.ID)
	contact := testutil.CreateTestContact(t, db, org.ID)
	conv := testutil.CreateTestConversation(t, db, org.ID, channel.ID, contact.ID)

	service := NewConversationService(db)

	tests := []struct {
		name    string
		orgID   uuid.UUID
		convID  uuid.UUID
		wantErr error
	}{
		{
			name:    "existing conversation",
			orgID:   org.ID,
			convID:  conv.ID,
			wantErr: nil,
		},
		{
			name:    "non-existent conversation",
			orgID:   org.ID,
			convID:  uuid.New(),
			wantErr: ErrConversationNotFound,
		},
		{
			name:    "wrong organization",
			orgID:   uuid.New(),
			convID:  conv.ID,
			wantErr: ErrConversationNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.GetByID(tt.orgID, tt.convID)

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

			if result.ID != tt.convID {
				t.Errorf("Expected conversation ID %s, got %s", tt.convID, result.ID)
			}
		})
	}
}

func TestConversationService_Assign(t *testing.T) {
	db := testutil.TestDB(t)
	defer db.Close()
	defer testutil.CleanupTestData(t, db)

	// Setup test data
	org := testutil.CreateTestOrganization(t, db)
	channel := testutil.CreateTestChannel(t, db, org.ID)
	contact := testutil.CreateTestContact(t, db, org.ID)
	conv := testutil.CreateTestConversation(t, db, org.ID, channel.ID, contact.ID)
	user := testutil.CreateTestUser(t, db, org.ID, models.RoleAgent)

	service := NewConversationService(db)

	tests := []struct {
		name    string
		orgID   uuid.UUID
		convID  uuid.UUID
		userID  uuid.UUID
		wantErr error
	}{
		{
			name:    "assign to agent",
			orgID:   org.ID,
			convID:  conv.ID,
			userID:  user.ID,
			wantErr: nil,
		},
		{
			name:    "non-existent conversation",
			orgID:   org.ID,
			convID:  uuid.New(),
			userID:  user.ID,
			wantErr: ErrConversationNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.Assign(tt.orgID, tt.convID, tt.userID)

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

			// Verify assignment
			result, err := service.GetByID(tt.orgID, tt.convID)
			if err != nil {
				t.Fatalf("Failed to get conversation: %v", err)
			}
			if result.AssignedTo == nil || *result.AssignedTo != tt.userID {
				t.Errorf("Expected assigned_to %s, got %v", tt.userID, result.AssignedTo)
			}
		})
	}
}

func TestConversationService_UpdateStatus(t *testing.T) {
	db := testutil.TestDB(t)
	defer db.Close()
	defer testutil.CleanupTestData(t, db)

	// Setup test data
	org := testutil.CreateTestOrganization(t, db)
	channel := testutil.CreateTestChannel(t, db, org.ID)
	contact := testutil.CreateTestContact(t, db, org.ID)
	
	service := NewConversationService(db)

	tests := []struct {
		name      string
		setup     func() uuid.UUID
		newStatus models.ConversationStatus
		wantErr   error
	}{
		{
			name: "update to resolved",
			setup: func() uuid.UUID {
				c := testutil.CreateTestConversation(t, db, org.ID, channel.ID, contact.ID)
				return c.ID
			},
			newStatus: models.ConversationStatusResolved,
			wantErr:   nil,
		},
		{
			name: "update to closed",
			setup: func() uuid.UUID {
				c := testutil.CreateTestConversation(t, db, org.ID, channel.ID, contact.ID)
				return c.ID
			},
			newStatus: models.ConversationStatusClosed,
			wantErr:   nil,
		},
		{
			name: "non-existent conversation",
			setup: func() uuid.UUID {
				return uuid.New()
			},
			newStatus: models.ConversationStatusClosed,
			wantErr:   ErrConversationNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			convID := tt.setup()
			err := service.UpdateStatus(org.ID, convID, tt.newStatus)

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

			// Verify status update
			result, err := service.GetByID(org.ID, convID)
			if err != nil {
				t.Fatalf("Failed to get conversation: %v", err)
			}
			if result.Status != tt.newStatus {
				t.Errorf("Expected status %s, got %s", tt.newStatus, result.Status)
			}
		})
	}
}

func TestConversationService_FindOrCreate(t *testing.T) {
	db := testutil.TestDB(t)
	defer db.Close()
	defer testutil.CleanupTestData(t, db)

	// Setup test data
	org := testutil.CreateTestOrganization(t, db)
	channel := testutil.CreateTestChannel(t, db, org.ID)
	contact := testutil.CreateTestContact(t, db, org.ID)

	service := NewConversationService(db)

	// First call should create
	conv1, err := service.FindOrCreate(org.ID, channel.ID, contact.ID)
	if err != nil {
		t.Fatalf("Failed to create conversation: %v", err)
	}
	if conv1 == nil {
		t.Fatal("Expected conversation, got nil")
	}

	// Second call should find existing
	conv2, err := service.FindOrCreate(org.ID, channel.ID, contact.ID)
	if err != nil {
		t.Fatalf("Failed to find conversation: %v", err)
	}
	if conv2 == nil {
		t.Fatal("Expected conversation, got nil")
	}

	// Should be the same conversation
	if conv1.ID != conv2.ID {
		t.Errorf("Expected same conversation ID, got %s and %s", conv1.ID, conv2.ID)
	}
}
