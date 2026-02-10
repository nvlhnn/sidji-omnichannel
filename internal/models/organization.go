package models

import (
	"time"

	"github.com/google/uuid"
)

// Organization represents a team/company using the platform
type Organization struct {
	ID                 uuid.UUID `json:"id"`
	Name               string    `json:"name"`
	Slug               string    `json:"slug"`
	Plan               string    `json:"plan"`
	SubscriptionStatus string    `json:"subscription_status"`
	AICreditsLimit      int       `json:"ai_credits_limit"`
	AICreditsUsed       int       `json:"ai_credits_used"`
	MessageUsageLimit   int       `json:"message_usage_limit"` // -1 for unlimited
	MessageUsageUsed    int       `json:"message_usage_used"`
	BillingCycleStart   time.Time `json:"billing_cycle_start"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`

	// UI helper fields (computed)
	IsOverLimit  bool `json:"is_over_limit"`
	UserCount    int  `json:"user_count"`
	ChannelCount int  `json:"channel_count"`
}

// CreateOrganizationInput for creating a new organization
type CreateOrganizationInput struct {
	Name string `json:"name" binding:"required,min=2,max=100"`
	Slug string `json:"slug" binding:"required,min=2,max=50,alphanum"`
}
