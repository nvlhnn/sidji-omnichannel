package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// ChannelType defines the type of messaging channel
type ChannelType string

const (
	ChannelWhatsApp  ChannelType = "whatsapp"
	ChannelInstagram ChannelType = "instagram"
	ChannelFacebook  ChannelType = "facebook"
)

// ChannelStatus defines the status of a channel connection
type ChannelStatus string

const (
	ChannelStatusActive       ChannelStatus = "active"
	ChannelStatusDisconnected ChannelStatus = "disconnected"
	ChannelStatusPending      ChannelStatus = "pending"
)

// Channel represents a connected messaging channel (WhatsApp/Instagram account)
type Channel struct {
	ID             uuid.UUID       `json:"id"`
	OrganizationID uuid.UUID       `json:"organization_id"`
	Type           ChannelType     `json:"type"`
	Provider       string          `json:"provider"` // e.g. "meta", "twilio"
	Name           string          `json:"name"`
	Config         json.RawMessage `json:"config,omitempty"` // Channel-specific config
	AccessToken    string          `json:"-"`                // Never expose access token
	PhoneNumberID  string          `json:"phone_number_id,omitempty"` // For WhatsApp
	IGUserID       string          `json:"ig_user_id,omitempty"`      // For Instagram
	FacebookPageID string          `json:"facebook_page_id,omitempty"` // For Facebook
	Status         ChannelStatus   `json:"status"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

// WhatsAppConfig holds WhatsApp-specific configuration
type WhatsAppConfig struct {
	BusinessAccountID string `json:"business_account_id"`
	TwilioAccountSID  string `json:"twilio_account_sid,omitempty"`
	DisplayPhone      string `json:"display_phone"`
	VerifiedName      string `json:"verified_name"`
}

// InstagramConfig holds Instagram-specific configuration
type InstagramConfig struct {
	Username    string `json:"username"`
	ProfilePic  string `json:"profile_pic"`
	AccountType string `json:"account_type"`
}

// FacebookConfig holds Facebook-specific configuration
type FacebookConfig struct {
	PageID      string `json:"page_id"`
	PageName    string `json:"page_name"`
	AccessToken string `json:"access_token"` // Page-scoped access token
}

// CreateChannelInput for connecting a new channel
type CreateChannelInput struct {
	Type        ChannelType `json:"type" binding:"required,oneof=whatsapp instagram facebook"`
	Provider    string      `json:"provider"` // defaults to meta
	Name        string      `json:"name" binding:"required,min=2,max=100"`
	AccessToken string      `json:"access_token"` // Required for Meta, SID/Token for Twilio
	// WhatsApp specific
	PhoneNumberID     string `json:"phone_number_id,omitempty"`
	BusinessAccountID string `json:"business_account_id,omitempty"`
	// Instagram specific
	IGUserID string `json:"ig_user_id,omitempty"`
	// Facebook specific
	FacebookPageID string `json:"facebook_page_id,omitempty"`
}

// ConnectInstagramInput for auto-connecting Instagram
type ConnectInstagramInput struct {
	AccessToken string `json:"access_token" binding:"required"`
	SelectedID  string `json:"selected_id,omitempty"`
}

// ConnectWhatsAppInput for auto-connecting WhatsApp
type ConnectWhatsAppInput struct {
	AccessToken string `json:"access_token" binding:"required"`
	SelectedID  string `json:"selected_id,omitempty"`
}

// ConnectFacebookInput for auto-connecting Facebook
type ConnectFacebookInput struct {
	AccessToken string `json:"access_token" binding:"required"`
	SelectedID  string `json:"selected_id,omitempty"`
}

// ChannelPublic is the public representation of a channel
type ChannelPublic struct {
	ID       uuid.UUID     `json:"id"`
	Type     ChannelType   `json:"type"`
	Provider string        `json:"provider"`
	Name     string        `json:"name"`
	Status   ChannelStatus `json:"status"`
}
