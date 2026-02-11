package services

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sidji-omnichannel/internal/adapters/db/postgres"
	"github.com/sidji-omnichannel/internal/models"
	"github.com/sidji-omnichannel/internal/testutil"
)

func TestMessageService_Create(t *testing.T) {
	db := testutil.TestDB(t)
	defer db.Close()
	defer testutil.CleanupTestData(t, db)

	// Setup test data
	org := testutil.CreateTestOrganization(t, db)
	channel := testutil.CreateTestChannel(t, db, org.ID)
	contact := testutil.CreateTestContact(t, db, org.ID)
	conv := testutil.CreateTestConversation(t, db, org.ID, channel.ID, contact.ID)
	user := testutil.CreateTestUser(t, db, org.ID, models.RoleAgent)

	mr := postgres.NewMessageRepository(db)
	cr := postgres.NewConversationRepository(db)
	ur := postgres.NewUserRepository(db)
	service := NewMessageService(mr, cr, ur)

	tests := []struct {
		name    string
		msg     *models.Message
		wantErr bool
	}{
		{
			name: "create text message from agent",
			msg: &models.Message{
				ConversationID: conv.ID,
				SenderType:     models.SenderAgent,
				SenderID:       user.ID,
				Content:        "Hello, how can I help you?",
				MessageType:    models.MessageTypeText,
				Status:         models.MessageStatusPending,
			},
			wantErr: false,
		},
		{
			name: "create text message from contact",
			msg: &models.Message{
				ConversationID: conv.ID,
				SenderType:     models.SenderContact,
				SenderID:       contact.ID,
				Content:        "I need help with my order",
				MessageType:    models.MessageTypeText,
				ExternalID:     "wamid.123456",
				Status:         models.MessageStatusDelivered,
			},
			wantErr: false,
		},
		{
			name: "create image message",
			msg: &models.Message{
				ConversationID: conv.ID,
				SenderType:     models.SenderAgent,
				SenderID:       user.ID,
				Content:        "Here is a screenshot",
				MessageType:    models.MessageTypeImage,
				MediaURL:       "https://example.com/image.jpg",
				MediaMimeType:  "image/jpeg",
				Status:         models.MessageStatusPending,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.Create(tt.msg)

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// Verify message was created
			if tt.msg.ID == uuid.Nil {
				t.Error("Expected message ID to be set")
			}
			if tt.msg.CreatedAt.IsZero() {
				t.Error("Expected CreatedAt to be set")
			}
		})
	}
}

func TestMessageService_List(t *testing.T) {
	db := testutil.TestDB(t)
	defer db.Close()
	defer testutil.CleanupTestData(t, db)

	// Setup test data
	org := testutil.CreateTestOrganization(t, db)
	channel := testutil.CreateTestChannel(t, db, org.ID)
	contact := testutil.CreateTestContact(t, db, org.ID)
	conv := testutil.CreateTestConversation(t, db, org.ID, channel.ID, contact.ID)
	user := testutil.CreateTestUser(t, db, org.ID, models.RoleAgent)

	mr := postgres.NewMessageRepository(db)
	cr := postgres.NewConversationRepository(db)
	ur := postgres.NewUserRepository(db)
	service := NewMessageService(mr, cr, ur)

	// Create some messages
	for i := 0; i < 5; i++ {
		msg := &models.Message{
			ConversationID: conv.ID,
			SenderType:     models.SenderAgent,
			SenderID:       user.ID,
			Content:        "Message " + string(rune('0'+i)),
			MessageType:    models.MessageTypeText,
			Status:         models.MessageStatusSent,
		}
		if err := service.Create(msg); err != nil {
			t.Fatalf("Failed to create message: %v", err)
		}
		// Small delay to ensure different timestamps
		time.Sleep(10 * time.Millisecond)
	}

	tests := []struct {
		name      string
		convID    uuid.UUID
		filter    *models.MessageFilter
		wantCount int
		wantErr   bool
	}{
		{
			name:   "list all messages",
			convID: conv.ID,
			filter: &models.MessageFilter{
				Limit: 50,
			},
			wantCount: 5,
			wantErr:   false,
		},
		{
			name:   "list with limit",
			convID: conv.ID,
			filter: &models.MessageFilter{
				Limit: 3,
			},
			wantCount: 3,
			wantErr:   false,
		},
		{
			name:   "different conversation (no results)",
			convID: uuid.New(),
			filter: &models.MessageFilter{
				Limit: 50,
			},
			wantCount: 0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.List(tt.convID, tt.filter)

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if len(result.Messages) != tt.wantCount {
				t.Errorf("Expected %d messages, got %d", tt.wantCount, len(result.Messages))
			}
		})
	}
}

func TestMessageService_UpdateStatus(t *testing.T) {
	db := testutil.TestDB(t)
	defer db.Close()
	defer testutil.CleanupTestData(t, db)

	// Setup test data
	org := testutil.CreateTestOrganization(t, db)
	channel := testutil.CreateTestChannel(t, db, org.ID)
	contact := testutil.CreateTestContact(t, db, org.ID)
	conv := testutil.CreateTestConversation(t, db, org.ID, channel.ID, contact.ID)
	user := testutil.CreateTestUser(t, db, org.ID, models.RoleAgent)

	mr := postgres.NewMessageRepository(db)
	cr := postgres.NewConversationRepository(db)
	ur := postgres.NewUserRepository(db)
	service := NewMessageService(mr, cr, ur)

	// Create a message with external ID
	msg := &models.Message{
		ConversationID: conv.ID,
		SenderType:     models.SenderAgent,
		SenderID:       user.ID,
		Content:        "Test message",
		MessageType:    models.MessageTypeText,
		ExternalID:     "wamid.test123",
		Status:         models.MessageStatusSent,
	}
	if err := service.Create(msg); err != nil {
		t.Fatalf("Failed to create message: %v", err)
	}

	// Test status updates
	tests := []struct {
		name       string
		externalID string
		newStatus  models.MessageStatus
		wantErr    bool
	}{
		{
			name:       "update to delivered",
			externalID: "wamid.test123",
			newStatus:  models.MessageStatusDelivered,
			wantErr:    false,
		},
		{
			name:       "update to read",
			externalID: "wamid.test123",
			newStatus:  models.MessageStatusRead,
			wantErr:    false,
		},
		{
			name:       "non-existent message (no error, just no rows affected)",
			externalID: "non-existent",
			newStatus:  models.MessageStatusRead,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.UpdateStatus(tt.externalID, tt.newStatus)

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
		})
	}

	// Verify final status
	fetchedMsg, err := service.GetByExternalID("wamid.test123")
	if err != nil {
		t.Fatalf("Failed to get message: %v", err)
	}
	if fetchedMsg == nil {
		t.Fatalf("Message wamid.test123 not found after update")
	}
	if fetchedMsg.Status != models.MessageStatusRead {
		t.Errorf("Expected status %s, got %s", models.MessageStatusRead, fetchedMsg.Status)
	}
}

func TestMessageService_MarkAsRead(t *testing.T) {
	db := testutil.TestDB(t)
	defer db.Close()
	defer testutil.CleanupTestData(t, db)

	// Setup test data
	org := testutil.CreateTestOrganization(t, db)
	channel := testutil.CreateTestChannel(t, db, org.ID)
	contact := testutil.CreateTestContact(t, db, org.ID)
	conv := testutil.CreateTestConversation(t, db, org.ID, channel.ID, contact.ID)

	mr := postgres.NewMessageRepository(db)
	cr := postgres.NewConversationRepository(db)
	ur := postgres.NewUserRepository(db)
	service := NewMessageService(mr, cr, ur)

	// Create messages from contact
	for i := 0; i < 3; i++ {
		msg := &models.Message{
			ConversationID: conv.ID,
			SenderType:     models.SenderContact,
			SenderID:       contact.ID,
			Content:        "Contact message " + string(rune('0'+i)),
			MessageType:    models.MessageTypeText,
			ExternalID:     "wamid.contact" + string(rune('0'+i)),
			Status:         models.MessageStatusDelivered,
		}
		if err := service.Create(msg); err != nil {
			t.Fatalf("Failed to create message: %v", err)
		}
	}

	// Mark as read
	err := service.MarkAsRead(conv.ID)
	if err != nil {
		t.Fatalf("Failed to mark as read: %v", err)
	}

	// Verify all messages are now read
	result, err := service.List(conv.ID, &models.MessageFilter{Limit: 50})
	if err != nil {
		t.Fatalf("Failed to list messages: %v", err)
	}

	for _, msg := range result.Messages {
		if msg.SenderType == models.SenderContact && msg.Status != models.MessageStatusRead {
			t.Errorf("Expected message %s to be read, got %s", msg.ID, msg.Status)
		}
	}
}

func TestMessageService_GetByExternalID(t *testing.T) {
	db := testutil.TestDB(t)
	defer db.Close()
	defer testutil.CleanupTestData(t, db)

	// Setup test data
	org := testutil.CreateTestOrganization(t, db)
	channel := testutil.CreateTestChannel(t, db, org.ID)
	contact := testutil.CreateTestContact(t, db, org.ID)
	conv := testutil.CreateTestConversation(t, db, org.ID, channel.ID, contact.ID)

	mr := postgres.NewMessageRepository(db)
	cr := postgres.NewConversationRepository(db)
	ur := postgres.NewUserRepository(db)
	service := NewMessageService(mr, cr, ur)

	// Create a message
	msg := &models.Message{
		ConversationID: conv.ID,
		SenderType:     models.SenderContact,
		SenderID:       contact.ID,
		Content:        "Test message",
		MessageType:    models.MessageTypeText,
		ExternalID:     "wamid.unique123",
		Status:         models.MessageStatusDelivered,
	}
	if err := service.Create(msg); err != nil {
		t.Fatalf("Failed to create message: %v", err)
	}

	tests := []struct {
		name       string
		externalID string
		wantNil    bool
		wantErr    bool
	}{
		{
			name:       "existing message",
			externalID: "wamid.unique123",
			wantNil:    false,
			wantErr:    false,
		},
		{
			name:       "non-existent message",
			externalID: "non-existent",
			wantNil:    true,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.GetByExternalID(tt.externalID)

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if tt.wantNil && result != nil {
				t.Error("Expected nil result")
			}
			if !tt.wantNil && result == nil {
				t.Error("Expected non-nil result")
			}
		})
	}
}
