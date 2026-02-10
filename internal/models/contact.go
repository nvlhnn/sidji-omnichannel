package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Contact represents a customer/contact from messaging channels
type Contact struct {
	ID             uuid.UUID       `json:"id"`
	OrganizationID uuid.UUID       `json:"organization_id"`
	Name           string          `json:"name"`
	Phone          string          `json:"phone,omitempty"`
	Email          string          `json:"email,omitempty"`
	AvatarURL      string          `json:"avatar_url,omitempty"`
	Metadata       json.RawMessage `json:"metadata,omitempty" swaggertype:"object"`
	// External IDs for different channels
	WhatsAppID  string `json:"whatsapp_id,omitempty"`
	InstagramID string `json:"instagram_id,omitempty"`
	FacebookID  string `json:"facebook_id,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateContactInput for creating a new contact
type CreateContactInput struct {
	Name        string `json:"name" binding:"required,min=1,max=200"`
	Phone       string `json:"phone,omitempty"`
	Email       string `json:"email,omitempty"`
	WhatsAppID  string `json:"whatsapp_id,omitempty"`
	InstagramID string `json:"instagram_id,omitempty"`
	FacebookID  string `json:"facebook_id,omitempty"`
}

// UpdateContactInput for updating contact info
type UpdateContactInput struct {
	Name  string `json:"name,omitempty"`
	Phone string `json:"phone,omitempty"`
	Email string `json:"email,omitempty"`
}
