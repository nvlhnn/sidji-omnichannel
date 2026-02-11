package services

import (
	"testing"

	"github.com/google/uuid"
	"github.com/sidji-omnichannel/internal/adapters/db/postgres"
	"github.com/sidji-omnichannel/internal/models"
	"github.com/sidji-omnichannel/internal/testutil"
)

func TestContactService_List(t *testing.T) {
	db := testutil.TestDB(t)
	defer db.Close()
	defer testutil.CleanupTestData(t, db)

	// Setup test data
	org := testutil.CreateTestOrganization(t, db)
	testutil.CreateTestContact(t, db, org.ID)
	testutil.CreateTestContact(t, db, org.ID)

	repo := postgres.NewContactRepository(db)
	service := NewContactService(repo)

	tests := []struct {
		name      string
		orgID     uuid.UUID
		page      int
		limit     int
		search    string
		wantCount int
		wantErr   bool
	}{
		{
			name:      "list all contacts",
			orgID:     org.ID,
			page:      1,
			limit:     20,
			search:    "",
			wantCount: 2,
			wantErr:   false,
		},
		{
			name:      "different organization (no results)",
			orgID:     uuid.New(),
			page:      1,
			limit:     20,
			search:    "",
			wantCount: 0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			contacts, total, err := service.List(tt.orgID, tt.page, tt.limit, tt.search)

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if len(contacts) != tt.wantCount {
				t.Errorf("Expected %d contacts, got %d", tt.wantCount, len(contacts))
			}

			if total < tt.wantCount {
				t.Errorf("Expected total >= %d, got %d", tt.wantCount, total)
			}
		})
	}
}

func TestContactService_GetByID(t *testing.T) {
	db := testutil.TestDB(t)
	defer db.Close()
	defer testutil.CleanupTestData(t, db)

	// Setup test data
	org := testutil.CreateTestOrganization(t, db)
	contact := testutil.CreateTestContact(t, db, org.ID)

	repo := postgres.NewContactRepository(db)
	service := NewContactService(repo)

	tests := []struct {
		name      string
		orgID     uuid.UUID
		contactID uuid.UUID
		wantErr   error
	}{
		{
			name:      "existing contact",
			orgID:     org.ID,
			contactID: contact.ID,
			wantErr:   nil,
		},
		{
			name:      "non-existent contact",
			orgID:     org.ID,
			contactID: uuid.New(),
			wantErr:   ErrContactNotFound,
		},
		{
			name:      "wrong organization",
			orgID:     uuid.New(),
			contactID: contact.ID,
			wantErr:   ErrContactNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.GetByID(tt.orgID, tt.contactID)

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

			if result.ID != tt.contactID {
				t.Errorf("Expected contact ID %s, got %s", tt.contactID, result.ID)
			}
		})
	}
}

func TestContactService_Create(t *testing.T) {
	db := testutil.TestDB(t)
	defer db.Close()
	defer testutil.CleanupTestData(t, db)

	// Setup test data
	org := testutil.CreateTestOrganization(t, db)

	repo := postgres.NewContactRepository(db)
	service := NewContactService(repo)

	tests := []struct {
		name    string
		input   *models.CreateContactInput
		wantErr bool
	}{
		{
			name: "create contact with all fields",
			input: &models.CreateContactInput{
				Name:       "John Doe",
				Phone:      "+1234567890",
				Email:      "john@example.com",
				WhatsAppID: "wa_john123",
			},
			wantErr: false,
		},
		{
			name: "create contact with minimal fields",
			input: &models.CreateContactInput{
				Name: "Jane Doe",
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

			if result.ID == uuid.Nil {
				t.Error("Expected contact ID to be set")
			}
		})
	}
}

func TestContactService_Update(t *testing.T) {
	db := testutil.TestDB(t)
	defer db.Close()
	defer testutil.CleanupTestData(t, db)

	// Setup test data
	org := testutil.CreateTestOrganization(t, db)
	contact := testutil.CreateTestContact(t, db, org.ID)

	repo := postgres.NewContactRepository(db)
	service := NewContactService(repo)

	tests := []struct {
		name      string
		orgID     uuid.UUID
		contactID uuid.UUID
		input     *models.UpdateContactInput
		wantErr   error
	}{
		{
			name:      "update contact name",
			orgID:     org.ID,
			contactID: contact.ID,
			input: &models.UpdateContactInput{
				Name: "Updated Name",
			},
			wantErr: nil,
		},
		{
			name:      "non-existent contact",
			orgID:     org.ID,
			contactID: uuid.New(),
			input: &models.UpdateContactInput{
				Name: "Updated Name",
			},
			wantErr: ErrContactNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.Update(tt.orgID, tt.contactID, tt.input)

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

func TestContactService_Delete(t *testing.T) {
	db := testutil.TestDB(t)
	defer db.Close()
	defer testutil.CleanupTestData(t, db)

	// Setup test data
	org := testutil.CreateTestOrganization(t, db)
	contact := testutil.CreateTestContact(t, db, org.ID)

	repo := postgres.NewContactRepository(db)
	service := NewContactService(repo)

	tests := []struct {
		name      string
		orgID     uuid.UUID
		contactID uuid.UUID
		wantErr   error
	}{
		{
			name:      "delete existing contact",
			orgID:     org.ID,
			contactID: contact.ID,
			wantErr:   nil,
		},
		{
			name:      "delete non-existent contact",
			orgID:     org.ID,
			contactID: uuid.New(),
			wantErr:   ErrContactNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.Delete(tt.orgID, tt.contactID)

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

			// Verify deletion
			_, err = service.GetByID(tt.orgID, tt.contactID)
			if err != ErrContactNotFound {
				t.Errorf("Expected contact to be deleted, got error: %v", err)
			}
		})
	}
}
