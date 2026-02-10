package models

import (
	"time"

	"github.com/google/uuid"
)

// ConversationStatus defines the status of a conversation
type ConversationStatus string

const (
	ConversationStatusOpen     ConversationStatus = "open"
	ConversationStatusPending  ConversationStatus = "pending"
	ConversationStatusResolved ConversationStatus = "resolved"
	ConversationStatusClosed   ConversationStatus = "closed"
)

// Conversation represents a chat thread with a contact
type Conversation struct {
	ID             uuid.UUID          `json:"id"`
	OrganizationID uuid.UUID          `json:"organization_id"`
	ChannelID      uuid.UUID          `json:"channel_id"`
	ContactID      uuid.UUID          `json:"contact_id"`
	AssignedTo     *uuid.UUID         `json:"assigned_to,omitempty"`
	Status         ConversationStatus `json:"status"`
	Subject        string             `json:"subject,omitempty"`
	LastMessageAt  *time.Time         `json:"last_message_at,omitempty"`
	LastHumanReplyAt *time.Time       `json:"last_human_reply_at,omitempty"` // Added for AI Hybrid mode
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`

	// Relations (populated when needed)
	Channel      *ChannelPublic `json:"channel,omitempty"`
	Contact      *Contact       `json:"contact,omitempty"`
	AssignedUser *UserPublic    `json:"assigned_user,omitempty"`
	LastMessage  *Message       `json:"last_message,omitempty"`
	UnreadCount  int            `json:"unread_count,omitempty"`
}

// ConversationListItem for inbox list view
type ConversationListItem struct {
	ID            uuid.UUID          `json:"id"`
	Status        ConversationStatus `json:"status"`
	Channel       *ChannelPublic     `json:"channel"`
	Contact       *Contact           `json:"contact"`
	AssignedUser  *UserPublic        `json:"assigned_user,omitempty"`
	LastMessage   *Message           `json:"last_message,omitempty"`
	LastMessageAt *time.Time         `json:"last_message_at"`
	UnreadCount   int                `json:"unread_count"`
}

// AssignConversationInput for assigning a conversation to an agent
type AssignConversationInput struct {
	UserID uuid.UUID `json:"user_id" binding:"required"`
}

// UpdateConversationStatusInput for updating conversation status
type UpdateConversationStatusInput struct {
	Status ConversationStatus `json:"status" binding:"required,oneof=open pending resolved closed"`
}

// ConversationFilter for filtering conversations
type ConversationFilter struct {
	Status     ConversationStatus `form:"status,omitempty"`
	ChannelID  uuid.UUID          `form:"channel_id,omitempty"`
	AssignedTo uuid.UUID          `form:"assigned_to,omitempty"`
	Unassigned bool               `form:"unassigned,omitempty"`
	Search     string             `form:"search,omitempty"`
	Page       int                `form:"page,default=1"`
	Limit      int                `form:"limit,default=20"`
}
