package models

import (
	"time"

	"github.com/google/uuid"
)

// SenderType defines who sent the message
// @Description who sent the message
// @Enum contact,agent,system,ai,note
type SenderType string

const (
	SenderContact SenderType = "contact"
	SenderAgent   SenderType = "agent"
	SenderSystem  SenderType = "system"
	SenderAI      SenderType = "ai"
	SenderNote    SenderType = "note"
)

// MessageType defines the type of message content
type MessageType string

const (
	MessageTypeText     MessageType = "text"
	MessageTypeImage    MessageType = "image"
	MessageTypeVideo    MessageType = "video"
	MessageTypeAudio    MessageType = "audio"
	MessageTypeDocument MessageType = "document"
	MessageTypeSticker  MessageType = "sticker"
	MessageTypeTemplate MessageType = "template"
	MessageTypeReaction MessageType = "reaction"
)

// MessageStatus defines the delivery status of a message
type MessageStatus string

const (
	MessageStatusPending   MessageStatus = "pending"
	MessageStatusSent      MessageStatus = "sent"
	MessageStatusDelivered MessageStatus = "delivered"
	MessageStatusRead      MessageStatus = "read"
	MessageStatusFailed    MessageStatus = "failed"
)

// Message represents a single message in a conversation
type Message struct {
	ID             uuid.UUID     `json:"id"`
	ConversationID uuid.UUID     `json:"conversation_id"`
	SenderType     SenderType    `json:"sender_type"`
	SenderID       uuid.UUID     `json:"sender_id"` // Contact ID or User ID
	Content        string        `json:"content"`
	MessageType    MessageType   `json:"message_type"`
	MediaURL       string        `json:"media_url,omitempty"`
	MediaMimeType  string        `json:"media_mime_type,omitempty"`
	MediaFileName  string        `json:"media_file_name,omitempty"`
	ExternalID     string        `json:"external_id,omitempty"` // WhatsApp/IG message ID
	ReplyToID      *uuid.UUID    `json:"reply_to_id,omitempty"`
	Status         MessageStatus `json:"status"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`

	// Relations
	Sender  *UserPublic `json:"sender,omitempty"`
	ReplyTo *Message    `json:"reply_to,omitempty"`
}

// SendMessageInput for sending a new message
type SendMessageInput struct {
	Content     string      `json:"content" binding:"required_without=MediaURL,max=500"`
	MessageType MessageType `json:"message_type" binding:"required,oneof=text image video audio document"`
	MediaURL    string      `json:"media_url,omitempty"`
	ReplyToID   *uuid.UUID  `json:"reply_to_id,omitempty"`
	IsNote      bool        `json:"is_note,omitempty"`
}

// MessageList for paginated message list
type MessageList struct {
	Messages   []*Message `json:"messages"`
	TotalCount int        `json:"total_count"`
	HasMore    bool       `json:"has_more"`
}

// MessageFilter for filtering messages
type MessageFilter struct {
	Before *time.Time `form:"before,omitempty"`
	After  *time.Time `form:"after,omitempty"`
	Limit  int        `form:"limit,default=50"`
}
