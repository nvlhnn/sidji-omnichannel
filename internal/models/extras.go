package models

import (
	"time"

	"github.com/google/uuid"
)

// CannedResponse is a pre-saved reply template
type CannedResponse struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	Shortcut       string    `json:"shortcut"` // e.g., "/greeting"
	Title          string    `json:"title"`
	Content        string    `json:"content"`
	CreatedBy      uuid.UUID `json:"created_by"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// CreateCannedResponseInput for creating a canned response
type CreateCannedResponseInput struct {
	Shortcut string `json:"shortcut" binding:"required,min=2,max=50"`
	Title    string `json:"title" binding:"required,min=2,max=100"`
	Content  string `json:"content" binding:"required,min=1,max=2000"`
}

// Label for tagging conversations
type Label struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	Name           string    `json:"name"`
	Color          string    `json:"color"` // Hex color code
	CreatedAt      time.Time `json:"created_at"`
}

// CreateLabelInput for creating a label
type CreateLabelInput struct {
	Name  string `json:"name" binding:"required,min=1,max=50"`
	Color string `json:"color" binding:"required,hexcolor"`
}
