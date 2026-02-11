package services

import (
	"testing"

	"github.com/google/uuid"
	"github.com/sidji-omnichannel/internal/adapters/db/postgres"
	"github.com/sidji-omnichannel/internal/models"
	"github.com/sidji-omnichannel/internal/testutil"
)

func TestTeamService_ListMembers(t *testing.T) {
	db := testutil.TestDB(t)
	defer db.Close()
	defer testutil.CleanupTestData(t, db)

	// Setup test data
	org := testutil.CreateTestOrganization(t, db)
	testutil.CreateTestUser(t, db, org.ID, models.RoleAdmin)
	testutil.CreateTestUser(t, db, org.ID, models.RoleAgent)

	repo := postgres.NewTeamRepository(db)
	service := NewTeamService(repo)

	tests := []struct {
		name      string
		orgID     uuid.UUID
		wantCount int
		wantErr   bool
	}{
		{
			name:      "list all members",
			orgID:     org.ID,
			wantCount: 2,
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
			members, err := service.ListMembers(tt.orgID)

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if len(members) != tt.wantCount {
				t.Errorf("Expected %d members, got %d", tt.wantCount, len(members))
			}
		})
	}
}

func TestTeamService_GetMember(t *testing.T) {
	db := testutil.TestDB(t)
	defer db.Close()
	defer testutil.CleanupTestData(t, db)

	// Setup test data
	org := testutil.CreateTestOrganization(t, db)
	user := testutil.CreateTestUser(t, db, org.ID, models.RoleAgent)

	repo := postgres.NewTeamRepository(db)
	service := NewTeamService(repo)

	tests := []struct {
		name     string
		orgID    uuid.UUID
		memberID uuid.UUID
		wantErr  error
	}{
		{
			name:     "existing member",
			orgID:    org.ID,
			memberID: user.ID,
			wantErr:  nil,
		},
		{
			name:     "non-existent member",
			orgID:    org.ID,
			memberID: uuid.New(),
			wantErr:  ErrTeamMemberNotFound,
		},
		{
			name:     "wrong organization",
			orgID:    uuid.New(),
			memberID: user.ID,
			wantErr:  ErrTeamMemberNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.GetMember(tt.orgID, tt.memberID)

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

			if result.ID != tt.memberID {
				t.Errorf("Expected member ID %s, got %s", tt.memberID, result.ID)
			}
		})
	}
}

func TestTeamService_InviteMember(t *testing.T) {
	db := testutil.TestDB(t)
	defer db.Close()
	defer testutil.CleanupTestData(t, db)

	// Setup test data
	org := testutil.CreateTestOrganization(t, db)

	repo := postgres.NewTeamRepository(db)
	service := NewTeamService(repo)

	tests := []struct {
		name    string
		input   *models.InviteUserInput
		wantErr error
	}{
		{
			name: "invite new agent",
			input: &models.InviteUserInput{
				Email:    "newagent@example.com",
				Name:     "New Agent",
				Role:     models.RoleAgent,
			},
			wantErr: nil,
		},
		{
			name: "invite supervisor",
			input: &models.InviteUserInput{
				Email:    "supervisor@example.com",
				Name:     "New Supervisor",
				Role:     models.RoleSupervisor,
			},
			wantErr: nil,
		},
		{
			name: "duplicate email",
			input: &models.InviteUserInput{
				Email:    "newagent@example.com", // same as first
				Name:     "Another Agent",
				Role:     models.RoleAgent,
			},
			wantErr: ErrUserExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.InviteMember(org.ID, tt.input)

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

			if result.Email != tt.input.Email {
				t.Errorf("Expected email %s, got %s", tt.input.Email, result.Email)
			}
			if result.Role != tt.input.Role {
				t.Errorf("Expected role %s, got %s", tt.input.Role, result.Role)
			}
		})
	}
}

func TestTeamService_UpdateMemberRole(t *testing.T) {
	db := testutil.TestDB(t)
	defer db.Close()
	defer testutil.CleanupTestData(t, db)

	// Setup test data
	org := testutil.CreateTestOrganization(t, db)
	user := testutil.CreateTestUser(t, db, org.ID, models.RoleAgent)

	repo := postgres.NewTeamRepository(db)
	service := NewTeamService(repo)

	tests := []struct {
		name     string
		orgID    uuid.UUID
		memberID uuid.UUID
		newRole  models.UserRole
		wantErr  error
	}{
		{
			name:     "promote to supervisor",
			orgID:    org.ID,
			memberID: user.ID,
			newRole:  models.RoleSupervisor,
			wantErr:  nil,
		},
		{
			name:     "demote to agent",
			orgID:    org.ID,
			memberID: user.ID,
			newRole:  models.RoleAgent,
			wantErr:  nil,
		},
		{
			name:     "non-existent member",
			orgID:    org.ID,
			memberID: uuid.New(),
			newRole:  models.RoleAgent,
			wantErr:  ErrTeamMemberNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.UpdateMemberRole(tt.orgID, tt.memberID, tt.newRole)

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

			// Verify role update
			member, err := service.GetMember(tt.orgID, tt.memberID)
			if err != nil {
				t.Fatalf("Failed to get member: %v", err)
			}
			if member.Role != tt.newRole {
				t.Errorf("Expected role %s, got %s", tt.newRole, member.Role)
			}
		})
	}
}

func TestTeamService_RemoveMember(t *testing.T) {
	db := testutil.TestDB(t)
	defer db.Close()
	defer testutil.CleanupTestData(t, db)

	// Setup test data
	org := testutil.CreateTestOrganization(t, db)
	admin := testutil.CreateTestUser(t, db, org.ID, models.RoleAdmin)
	agent := testutil.CreateTestUser(t, db, org.ID, models.RoleAgent)

	repo := postgres.NewTeamRepository(db)
	service := NewTeamService(repo)

	tests := []struct {
		name     string
		orgID    uuid.UUID
		memberID uuid.UUID
		wantErr  bool
	}{
		{
			name:     "remove agent",
			orgID:    org.ID,
			memberID: agent.ID,
			wantErr:  false,
		},
		{
			name:     "remove last admin (should fail)",
			orgID:    org.ID,
			memberID: admin.ID,
			wantErr:  true, // Cannot remove last admin
		},
		{
			name:     "remove non-existent member",
			orgID:    org.ID,
			memberID: uuid.New(),
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.RemoveMember(tt.orgID, tt.memberID)

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// Verify deletion
			_, err = service.GetMember(tt.orgID, tt.memberID)
			if err != ErrTeamMemberNotFound {
				t.Errorf("Expected member to be deleted, got error: %v", err)
			}
		})
	}
}

func TestTeamService_GetOrganization(t *testing.T) {
	db := testutil.TestDB(t)
	defer db.Close()
	defer testutil.CleanupTestData(t, db)

	// Setup test data
	org := testutil.CreateTestOrganization(t, db)

	repo := postgres.NewTeamRepository(db)
	service := NewTeamService(repo)

	tests := []struct {
		name    string
		orgID   uuid.UUID
		wantErr bool
	}{
		{
			name:    "existing organization",
			orgID:   org.ID,
			wantErr: false,
		},
		{
			name:    "non-existent organization",
			orgID:   uuid.New(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.GetOrganization(tt.orgID)

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if result.ID != tt.orgID {
				t.Errorf("Expected org ID %s, got %s", tt.orgID, result.ID)
			}
		})
	}
}

func TestTeamService_UpdateOrganization(t *testing.T) {
	db := testutil.TestDB(t)
	defer db.Close()
	defer testutil.CleanupTestData(t, db)

	// Setup test data
	org := testutil.CreateTestOrganization(t, db)

	repo := postgres.NewTeamRepository(db)
	service := NewTeamService(repo)

	tests := []struct {
		name    string
		orgID   uuid.UUID
		newName string
		wantErr bool
	}{
		{
			name:    "update organization name",
			orgID:   org.ID,
			newName: "Updated Organization Name",
			wantErr: false,
		},
		{
			name:    "non-existent organization",
			orgID:   uuid.New(),
			newName: "New Name",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.UpdateOrganization(tt.orgID, tt.newName, "")

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if result.Name != tt.newName {
				t.Errorf("Expected name %s, got %s", tt.newName, result.Name)
			}
		})
	}
}
