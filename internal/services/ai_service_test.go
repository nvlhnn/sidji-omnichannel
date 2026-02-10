package services

import (
	"context"
	"testing"
	"time"

	"github.com/sidji-omnichannel/internal/models"
	"github.com/sidji-omnichannel/internal/testutil"
)

type MockAIProvider struct {
	NameFunc          func() string
	EmbedTextFunc     func(ctx context.Context, text string) ([]float32, error)
	GenerateReplyFunc func(ctx context.Context, config *models.AIConfig, userQuery string, contextDocs []string, history []models.Message) (string, error)
}

func (m *MockAIProvider) Name() string {
	if m.NameFunc != nil {
		return m.NameFunc()
	}
	return "mock-provider"
}

func (m *MockAIProvider) EmbedText(ctx context.Context, text string) ([]float32, error) {
	if m.EmbedTextFunc != nil {
		return m.EmbedTextFunc(ctx, text)
	}
	return make([]float32, 3072), nil
}

func (m *MockAIProvider) GenerateReply(ctx context.Context, config *models.AIConfig, userQuery string, contextDocs []string, history []models.Message) (string, error) {
	if m.GenerateReplyFunc != nil {
		return m.GenerateReplyFunc(ctx, config, userQuery, contextDocs, history)
	}
	return "Mock reply", nil
}

func TestAIService_Config(t *testing.T) {
	db := testutil.TestDB(t)
	defer db.Close()
	defer testutil.CleanupTestData(t, db)

	org := testutil.CreateTestOrganization(t, db)
	channel := testutil.CreateTestChannel(t, db, org.ID)

	service := &AIService{db: db}

	// 1. Get default config
	config, err := service.GetConfig(channel.ID)
	if err != nil {
		t.Fatalf("Failed to get config: %v", err)
	}
	if config.IsEnabled {
		t.Error("Default config should be disabled")
	}
	if config.Mode != models.AIModeManual {
		t.Errorf("Expected mode %s, got %s", models.AIModeManual, config.Mode)
	}

	// 2. Update config
	config.IsEnabled = true
	config.Mode = models.AIModeHybrid
	config.Persona = "Test persona"
	config.HandoverTimeoutMinutes = 10

	err = service.UpdateConfig(config)
	if err != nil {
		t.Fatalf("Failed to update config: %v", err)
	}

	// 3. Verify update
	updated, err := service.GetConfig(channel.ID)
	if err != nil {
		t.Fatalf("Failed to get updated config: %v", err)
	}
	if !updated.IsEnabled {
		t.Error("Expected config to be enabled")
	}
	if updated.Mode != models.AIModeHybrid {
		t.Errorf("Expected mode %s, got %s", models.AIModeHybrid, updated.Mode)
	}
	if updated.HandoverTimeoutMinutes != 10 {
		t.Errorf("Expected timeout 10, got %d", updated.HandoverTimeoutMinutes)
	}
}

func TestAIService_DetermineAction(t *testing.T) {
	db := testutil.TestDB(t)
	defer db.Close()
	defer testutil.CleanupTestData(t, db)

	org := testutil.CreateTestOrganization(t, db)
	channel := testutil.CreateTestChannel(t, db, org.ID)
	service := &AIService{db: db}

	now := time.Now()
	recent := now.Add(-5 * time.Minute)
	old := now.Add(-30 * time.Minute)

	tests := []struct {
		name              string
		isEnabled         bool
		mode              models.AIMode
		timeout           int
		lastHumanReply    *time.Time
		wantAutoReply     bool
		wantDraft         bool
	}{
		{
			name:          "disabled",
			isEnabled:     false,
			wantAutoReply: false,
			wantDraft:     false,
		},
		{
			name:          "manual mode",
			isEnabled:     true,
			mode:          models.AIModeManual,
			wantAutoReply: false,
			wantDraft:     true,
		},
		{
			name:          "auto mode",
			isEnabled:     true,
			mode:          models.AIModeAuto,
			wantAutoReply: true,
			wantDraft:     false,
		},
		{
			name:           "hybrid mode - under timeout",
			isEnabled:      true,
			mode:           models.AIModeHybrid,
			timeout:        15,
			lastHumanReply: &recent,
			wantAutoReply:  false,
			wantDraft:      false,
		},
		{
			name:           "hybrid mode - over timeout",
			isEnabled:      true,
			mode:           models.AIModeHybrid,
			timeout:        15,
			lastHumanReply: &old,
			wantAutoReply:  true,
			wantDraft:      false,
		},
		{
			name:           "hybrid mode - no human reply",
			isEnabled:      true,
			mode:           models.AIModeHybrid,
			lastHumanReply: nil,
			wantAutoReply:  true,
			wantDraft:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Update settings for this test
			_ = service.UpdateConfig(&models.AIConfig{
				ChannelID:              channel.ID,
				IsEnabled:              tt.isEnabled,
				Mode:                   tt.mode,
				HandoverTimeoutMinutes: tt.timeout,
			})

			auto, draft, _, err := service.DetermineAction(channel.ID, tt.lastHumanReply)
			if err != nil {
				t.Fatalf("DetermineAction failed: %v", err)
			}
			if auto != tt.wantAutoReply {
				t.Errorf("Expected autoReply=%v, got %v", tt.wantAutoReply, auto)
			}
			if draft != tt.wantDraft {
				t.Errorf("Expected draft=%v, got %v", tt.wantDraft, draft)
			}
		})
	}
}

func TestAIService_Knowledge(t *testing.T) {
	db := testutil.TestDB(t)
	defer db.Close()
	defer testutil.CleanupTestData(t, db)

	org := testutil.CreateTestOrganization(t, db)
	channel := testutil.CreateTestChannel(t, db, org.ID)
	
	mockProvider := &MockAIProvider{}
	service := &AIService{
		db:       db,
		provider: mockProvider,
	}

	ctx := context.Background()

	// 1. Add knowledge
	content := "Test knowledge content"
	err := service.AddKnowledge(channel.ID, content)
	if err != nil {
		t.Fatalf("Failed to add knowledge: %v", err)
	}

	// 2. List knowledge
	items, err := service.ListKnowledge(channel.ID)
	if err != nil {
		t.Fatalf("Failed to list knowledge: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("Expected 1 item, got %d", len(items))
	}
	if items[0].Content != content {
		t.Errorf("Expected content %s, got %s", content, items[0].Content)
	}

	knowledgeID := items[0].ID

	// 3. Update knowledge
	newContent := "Updated content"
	err = service.UpdateKnowledge(knowledgeID, newContent)
	if err != nil {
		t.Fatalf("Failed to update knowledge: %v", err)
	}

	// 4. Verify update
	updatedItems, _ := service.ListKnowledge(channel.ID)
	if updatedItems[0].Content != newContent {
		t.Errorf("Expected updated content %s, got %s", newContent, updatedItems[0].Content)
	}

	// 5. Search
	results, err := service.SearchKnowledge(channel.ID, make([]float32, 3072), 1)
	if err != nil {
		t.Fatalf("Failed to search: %v", err)
	}
	if len(results) != 1 || results[0] != newContent {
		t.Error("Search failed to return updated content")
	}

	// 6. Delete
	err = service.DeleteKnowledge(knowledgeID)
	if err != nil {
		t.Fatalf("Failed to delete knowledge: %v", err)
	}

	// 7. Verify deletion
	finalItems, _ := service.ListKnowledge(channel.ID)
	if len(finalItems) != 0 {
		t.Error("Knowledge base should be empty")
	}

	// 8. Generate Reply (test orchestration)
	reply, err := service.GenerateReply(ctx, &models.AIConfig{}, "hi", []string{content}, nil)
	if err != nil {
		t.Fatalf("GenerateReply failed: %v", err)
	}
	if reply != "Mock reply" {
		t.Errorf("Expected Mock reply, got %s", reply)
	}
}

func TestAIService_Credits(t *testing.T) {
	db := testutil.TestDB(t)
	defer db.Close()
	defer testutil.CleanupTestData(t, db)

	org := testutil.CreateTestOrganization(t, db)
	// Ensure default credits are set (since testutil might be minimal)
	// Check if credits are 10
	var limit int
	err := db.QueryRow("SELECT ai_credits_limit FROM organizations WHERE id = $1", org.ID).Scan(&limit)
	if err != nil {
		t.Fatalf("Failed to check limit: %v", err)
	}
	// If 0, it means testutil inserted without defaults or default didn't trigger?
	// Postgres defaults trigger on INSERT if column is omitted.
	
	channel := testutil.CreateTestChannel(t, db, org.ID)
	service := &AIService{db: db}

	// 1. Initial State: 10 credits (default)
	err = service.CheckCredits(channel.ID)
	if err != nil {
		t.Fatalf("Expected sufficient credits, got error: %v", err)
	}

	// 2. Consume all 10 credits
	for i := 0; i < 10; i++ {
		err := service.DeductCredit(channel.ID)
		if err != nil {
			t.Fatalf("Failed to deduct credit at %d: %v", i, err)
		}
	}

	// 3. Verify used count
	var used int
	db.QueryRow("SELECT ai_credits_used FROM organizations WHERE id = $1", org.ID).Scan(&used)
	if used != 10 {
		t.Errorf("Expected 10 credits used, got %d", used)
	}

	// 4. Check should fail now
	err = service.CheckCredits(channel.ID)
	if err == nil {
		t.Error("Expected error for insufficient credits, got nil")
	}

	// 5. Increase limit manually
	_, err = db.Exec("UPDATE organizations SET ai_credits_limit = 20 WHERE id = $1", org.ID)
	if err != nil {
		t.Fatalf("Failed to update limit: %v", err)
	}

	// 6. Check should pass
	err = service.CheckCredits(channel.ID)
	if err != nil {
		t.Errorf("Expected sufficient credits after limit increase, got error: %v", err)
	}
}
