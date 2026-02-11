package services

import (
	"testing"

	"github.com/google/uuid"
	"github.com/sidji-omnichannel/internal/adapters/db/postgres"
	"github.com/sidji-omnichannel/internal/models"
	"github.com/sidji-omnichannel/internal/testutil"
)

func TestCannedResponseService_List(t *testing.T) {
	db := testutil.TestDB(t)
	defer db.Close()
	defer testutil.CleanupTestData(t, db)

	// Setup test data
	org := testutil.CreateTestOrganization(t, db)
	user := testutil.CreateTestUser(t, db, org.ID, models.RoleAgent)

	service := NewCannedResponseService(postgres.NewCannedResponseRepository(db))

	// Create some canned responses
	for i := 0; i < 3; i++ {
		_, err := service.Create(org.ID, user.ID, &models.CreateCannedResponseInput{
			Shortcut: "/test" + string(rune('0'+i)),
			Content:  "Test response " + string(rune('0'+i)),
		})
		if err != nil {
			t.Fatalf("Failed to create canned response: %v", err)
		}
	}

	tests := []struct {
		name      string
		orgID     uuid.UUID
		wantCount int
		wantErr   bool
	}{
		{
			name:      "list all responses",
			orgID:     org.ID,
			wantCount: 3,
			wantErr:   false,
		},
		{
			name:      "different organization (no results)",
			orgID:     uuid.New(),
			wantCount: 0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			responses, err := service.List(tt.orgID)

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if len(responses) != tt.wantCount {
				t.Errorf("Expected %d responses, got %d", tt.wantCount, len(responses))
			}
		})
	}
}

func TestCannedResponseService_Search(t *testing.T) {
	db := testutil.TestDB(t)
	defer db.Close()
	defer testutil.CleanupTestData(t, db)

	// Setup test data
	org := testutil.CreateTestOrganization(t, db)
	user := testutil.CreateTestUser(t, db, org.ID, models.RoleAgent)

	service := NewCannedResponseService(postgres.NewCannedResponseRepository(db))

	// Create canned responses
	_, _ = service.Create(org.ID, user.ID, &models.CreateCannedResponseInput{
		Shortcut: "/hello",
		Content:  "Hello! How can I help you today?",
	})
	_, _ = service.Create(org.ID, user.ID, &models.CreateCannedResponseInput{
		Shortcut: "/goodbye",
		Content:  "Thank you for contacting us. Goodbye!",
	})
	_, _ = service.Create(org.ID, user.ID, &models.CreateCannedResponseInput{
		Shortcut: "/help",
		Content:  "Here are the available options...",
	})

	tests := []struct {
		name      string
		orgID     uuid.UUID
		query     string
		wantCount int
		wantErr   bool
	}{
		{
			name:      "search by shortcut",
			orgID:     org.ID,
			query:     "hello",
			wantCount: 1,
			wantErr:   false,
		},
		{
			name:      "search by content",
			orgID:     org.ID,
			query:     "thank",
			wantCount: 1,
			wantErr:   false,
		},
		{
			name:      "search with no results",
			orgID:     org.ID,
			query:     "nonexistent",
			wantCount: 0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			responses, err := service.Search(tt.orgID, tt.query)

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if len(responses) != tt.wantCount {
				t.Errorf("Expected %d responses, got %d", tt.wantCount, len(responses))
			}
		})
	}
}

func TestCannedResponseService_Create(t *testing.T) {
	db := testutil.TestDB(t)
	defer db.Close()
	defer testutil.CleanupTestData(t, db)

	// Setup test data
	org := testutil.CreateTestOrganization(t, db)
	user := testutil.CreateTestUser(t, db, org.ID, models.RoleAgent)

	service := NewCannedResponseService(postgres.NewCannedResponseRepository(db))

	tests := []struct {
		name    string
		input   *models.CreateCannedResponseInput
		wantErr bool
	}{
		{
			name: "create valid response",
			input: &models.CreateCannedResponseInput{
				Shortcut: "/greeting",
				Content:  "Hello! Welcome to our support.",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.Create(org.ID, user.ID, tt.input)

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if result.Shortcut != tt.input.Shortcut {
				t.Errorf("Expected shortcut %s, got %s", tt.input.Shortcut, result.Shortcut)
			}
			if result.Content != tt.input.Content {
				t.Errorf("Expected content %s, got %s", tt.input.Content, result.Content)
			}
		})
	}
}

func TestCannedResponseService_Delete(t *testing.T) {
	db := testutil.TestDB(t)
	defer db.Close()
	defer testutil.CleanupTestData(t, db)

	// Setup test data
	org := testutil.CreateTestOrganization(t, db)
	user := testutil.CreateTestUser(t, db, org.ID, models.RoleAgent)

	service := NewCannedResponseService(postgres.NewCannedResponseRepository(db))

	// Create a response
	response, err := service.Create(org.ID, user.ID, &models.CreateCannedResponseInput{
		Shortcut: "/todelete",
		Content:  "This will be deleted",
	})
	if err != nil {
		t.Fatalf("Failed to create response: %v", err)
	}

	tests := []struct {
		name       string
		orgID      uuid.UUID
		responseID uuid.UUID
		wantErr    error
	}{
		{
			name:       "delete existing response",
			orgID:      org.ID,
			responseID: response.ID,
			wantErr:    nil,
		},
		{
			name:       "delete non-existent response",
			orgID:      org.ID,
			responseID: uuid.New(),
			wantErr:    ErrCannedResponseNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.Delete(tt.orgID, tt.responseID)

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
		})
	}
}

// Label Service Tests

func TestLabelService_List(t *testing.T) {
	db := testutil.TestDB(t)
	defer db.Close()
	defer testutil.CleanupTestData(t, db)

	// Setup test data
	org := testutil.CreateTestOrganization(t, db)

	service := NewLabelService(postgres.NewLabelRepository(db))

	// Create some labels
	for i := 0; i < 3; i++ {
		_, err := service.Create(org.ID, &models.CreateLabelInput{
			Name:  "Label " + string(rune('0'+i)),
			Color: "#ff000" + string(rune('0'+i)),
		})
		if err != nil {
			t.Fatalf("Failed to create label: %v", err)
		}
	}

	tests := []struct {
		name      string
		orgID     uuid.UUID
		wantCount int
		wantErr   bool
	}{
		{
			name:      "list all labels",
			orgID:     org.ID,
			wantCount: 3,
			wantErr:   false,
		},
		{
			name:      "different organization (no results)",
			orgID:     uuid.New(),
			wantCount: 0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			labels, err := service.List(tt.orgID)

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if len(labels) != tt.wantCount {
				t.Errorf("Expected %d labels, got %d", tt.wantCount, len(labels))
			}
		})
	}
}

func TestLabelService_Create(t *testing.T) {
	db := testutil.TestDB(t)
	defer db.Close()
	defer testutil.CleanupTestData(t, db)

	// Setup test data
	org := testutil.CreateTestOrganization(t, db)

	service := NewLabelService(postgres.NewLabelRepository(db))

	tests := []struct {
		name    string
		input   *models.CreateLabelInput
		wantErr bool
	}{
		{
			name: "create valid label",
			input: &models.CreateLabelInput{
				Name:  "Urgent",
				Color: "#ff0000",
			},
			wantErr: false,
		},
		{
			name: "create another label",
			input: &models.CreateLabelInput{
				Name:  "Follow-up",
				Color: "#00ff00",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.Create(org.ID, tt.input)

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if result.Name != tt.input.Name {
				t.Errorf("Expected name %s, got %s", tt.input.Name, result.Name)
			}
			if result.Color != tt.input.Color {
				t.Errorf("Expected color %s, got %s", tt.input.Color, result.Color)
			}
		})
	}
}

func TestLabelService_Update(t *testing.T) {
	db := testutil.TestDB(t)
	defer db.Close()
	defer testutil.CleanupTestData(t, db)

	// Setup test data
	org := testutil.CreateTestOrganization(t, db)

	service := NewLabelService(postgres.NewLabelRepository(db))

	// Create a label
	label, err := service.Create(org.ID, &models.CreateLabelInput{
		Name:  "Original",
		Color: "#000000",
	})
	if err != nil {
		t.Fatalf("Failed to create label: %v", err)
	}

	tests := []struct {
		name    string
		orgID   uuid.UUID
		labelID uuid.UUID
		input   *models.CreateLabelInput
		wantErr error
	}{
		{
			name:    "update label",
			orgID:   org.ID,
			labelID: label.ID,
			input: &models.CreateLabelInput{
				Name:  "Updated",
				Color: "#ffffff",
			},
			wantErr: nil,
		},
		{
			name:    "update non-existent label",
			orgID:   org.ID,
			labelID: uuid.New(),
			input: &models.CreateLabelInput{
				Name:  "Test",
				Color: "#000000",
			},
			wantErr: ErrLabelNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.Update(tt.orgID, tt.labelID, tt.input)

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

			if result.Name != tt.input.Name {
				t.Errorf("Expected name %s, got %s", tt.input.Name, result.Name)
			}
		})
	}
}

func TestLabelService_Delete(t *testing.T) {
	db := testutil.TestDB(t)
	defer db.Close()
	defer testutil.CleanupTestData(t, db)

	// Setup test data
	org := testutil.CreateTestOrganization(t, db)

	service := NewLabelService(postgres.NewLabelRepository(db))

	// Create a label
	label, err := service.Create(org.ID, &models.CreateLabelInput{
		Name:  "ToDelete",
		Color: "#123456",
	})
	if err != nil {
		t.Fatalf("Failed to create label: %v", err)
	}

	tests := []struct {
		name    string
		orgID   uuid.UUID
		labelID uuid.UUID
		wantErr error
	}{
		{
			name:    "delete existing label",
			orgID:   org.ID,
			labelID: label.ID,
			wantErr: nil,
		},
		{
			name:    "delete non-existent label",
			orgID:   org.ID,
			labelID: uuid.New(),
			wantErr: ErrLabelNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.Delete(tt.orgID, tt.labelID)

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
		})
	}
}
