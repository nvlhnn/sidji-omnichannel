package integrations

import (
	"context"

	"github.com/sidji-omnichannel/internal/models"
)

// MessageProvider defines the interface for different messaging platforms (Twilio, Meta, etc.)
type MessageProvider interface {
	// SendMessage sends a message to a recipient
	SendMessage(ctx context.Context, channel *models.Channel, contact *models.Contact, message *models.Message) (string, error)
	
	// GetProviderName returns the unique identifier for this provider (e.g., "meta", "twilio")
	GetProviderName() string
}

// WhatsAppProvider is a specific interface for WhatsApp-related features
type WhatsAppProvider interface {
	MessageProvider
	// Add WhatsApp specific methods here if needed, like template management
}

// InstagramProvider is a specific interface for Instagram-related features
type InstagramProvider interface {
	MessageProvider
}
