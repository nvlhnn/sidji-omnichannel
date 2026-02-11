package testutil

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/sidji-omnichannel/internal/config"
	"github.com/sidji-omnichannel/internal/models"
)

// TestDB returns a test database connection
func TestDB(t *testing.T) *sql.DB {
	t.Helper()

	// Use test database from environment or default
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://sidji:sidji123@localhost:5433/sidji_test?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to ping test database: %v", err)
	}

	return db
}

// TestConfig returns a test configuration
func TestConfig() *config.Config {
	return &config.Config{
		App: config.AppConfig{
			Env:    "test",
			Port:   8080,
			Secret: "test-secret-key-for-jwt-signing",
		},
		Database: config.DatabaseConfig{
			Host:     "localhost",
			Port:     5433,
			User:     "sidji",
			Password: "sidji123",
			Name:     "sidji_test",
		},
	}
}

// SetupTestRouter creates a test Gin router
func SetupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

// CreateTestOrganization creates a test organization in the database
func CreateTestOrganization(t *testing.T, db *sql.DB) *models.Organization {
	t.Helper()

	org := &models.Organization{
		ID:   uuid.New(),
		Name: fmt.Sprintf("Test Org %s", uuid.New().String()[:8]),
		Slug: fmt.Sprintf("test-org-%s", uuid.New().String()[:8]),
	}

	_, err := db.Exec(
		"INSERT INTO organizations (id, name, slug, plan, subscription_status) VALUES ($1, $2, $3, 'enterprise', 'active')",
		org.ID, org.Name, org.Slug,
	)
	if err != nil {
		t.Fatalf("Failed to create test organization: %v", err)
	}

	return org
}

// CreateTestUser creates a test user in the database
func CreateTestUser(t *testing.T, db *sql.DB, orgID uuid.UUID, role models.UserRole) *models.User {
	t.Helper()

	user := &models.User{
		ID:             uuid.New(),
		OrganizationID: orgID,
		Email:          fmt.Sprintf("test-%s@example.com", uuid.New().String()[:8]),
		PasswordHash:   "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // "password"
		Name:           "Test User",
		Role:           role,
		Status:         models.StatusOnline,
	}

	_, err := db.Exec(
		`INSERT INTO users (id, organization_id, email, password_hash, name, role, status)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		user.ID, user.OrganizationID, user.Email, user.PasswordHash, user.Name, user.Role, user.Status,
	)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	return user
}

// CreateTestContact creates a test contact in the database
func CreateTestContact(t *testing.T, db *sql.DB, orgID uuid.UUID) *models.Contact {
	t.Helper()

	contact := &models.Contact{
		ID:             uuid.New(),
		OrganizationID: orgID,
		Name:           fmt.Sprintf("Test Contact %s", uuid.New().String()[:8]),
		Phone:          fmt.Sprintf("+1%s", uuid.New().String()[:10]),
		WhatsAppID:     fmt.Sprintf("wa_%s", uuid.New().String()[:8]),
	}

	_, err := db.Exec(
		`INSERT INTO contacts (id, organization_id, name, phone, whatsapp_id)
		 VALUES ($1, $2, $3, $4, $5)`,
		contact.ID, contact.OrganizationID, contact.Name, contact.Phone, contact.WhatsAppID,
	)
	if err != nil {
		t.Fatalf("Failed to create test contact: %v", err)
	}

	return contact
}

// CreateTestChannel creates a test channel in the database
func CreateTestChannel(t *testing.T, db *sql.DB, orgID uuid.UUID) *models.Channel {
	t.Helper()

	channel := &models.Channel{
		ID:             uuid.New(),
		OrganizationID: orgID,
		Type:           models.ChannelWhatsApp,
		Name:           fmt.Sprintf("Test Channel %s", uuid.New().String()[:8]),
		PhoneNumberID:  fmt.Sprintf("phone_%s", uuid.New().String()[:8]),
		Status:         models.ChannelStatusActive,
	}

	_, err := db.Exec(
		`INSERT INTO channels (id, organization_id, type, provider, name, phone_number_id, status)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		channel.ID, channel.OrganizationID, channel.Type, "meta", channel.Name, channel.PhoneNumberID, channel.Status,
	)
	if err != nil {
		t.Fatalf("Failed to create test channel: %v", err)
	}

	return channel
}

// CreateTestConversation creates a test conversation in the database
func CreateTestConversation(t *testing.T, db *sql.DB, orgID, channelID, contactID uuid.UUID) *models.Conversation {
	t.Helper()

	conv := &models.Conversation{
		ID:             uuid.New(),
		OrganizationID: orgID,
		ChannelID:      channelID,
		ContactID:      contactID,
		Status:         models.ConversationStatusOpen,
	}

	_, err := db.Exec(
		`INSERT INTO conversations (id, organization_id, channel_id, contact_id, status)
		 VALUES ($1, $2, $3, $4, $5)`,
		conv.ID, conv.OrganizationID, conv.ChannelID, conv.ContactID, conv.Status,
	)
	if err != nil {
		t.Fatalf("Failed to create test conversation: %v", err)
	}

	return conv
}

// CleanupTestData removes test data from the database
func CleanupTestData(t *testing.T, db *sql.DB) {
	t.Helper()

	tables := []string{
		"messages",
		"conversation_notes",
		"knowledge_base",
		"ai_configs",
		"conversations",
		"contacts",
		"channels",
		"canned_responses",
		"labels",
		"users",
		"organizations",
	}

	for _, table := range tables {
		_, err := db.Exec(fmt.Sprintf("DELETE FROM %s WHERE id IS NOT NULL", table))
		if err != nil {
			t.Logf("Warning: Failed to cleanup table %s: %v", table, err)
		}
	}
}

// GenerateTestToken generates a test JWT token
func GenerateTestToken(user *models.User, secret string) (string, error) {
	// This would generate a real token - simplified for testing
	return "test-token-" + user.ID.String(), nil
}
