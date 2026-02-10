package models

import (
	"time"

	"github.com/google/uuid"
)

// UserRole defines the role of a user in the organization
type UserRole string

const (
	RoleAdmin      UserRole = "admin"
	RoleSupervisor UserRole = "supervisor"
	RoleAgent      UserRole = "agent"
)

// UserStatus defines the online status of a user
type UserStatus string

const (
	StatusOnline  UserStatus = "online"
	StatusAway    UserStatus = "away"
	StatusOffline UserStatus = "offline"
)

// User represents a team member (agent, supervisor, or admin)
type User struct {
	ID             uuid.UUID  `json:"id"`
	OrganizationID uuid.UUID  `json:"organization_id"`
	Email          string     `json:"email"`
	PasswordHash   string     `json:"-"` // Never expose password hash
	Name           string     `json:"name"`
	Role           UserRole   `json:"role"`
	AvatarURL      string     `json:"avatar_url,omitempty"`
	Status         UserStatus `json:"status"`
	LastSeenAt     *time.Time `json:"last_seen_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// UserPublic is the public representation of a user (for API responses)
type UserPublic struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Role      UserRole   `json:"role"`
	AvatarURL string     `json:"avatar_url,omitempty"`
	Status    UserStatus `json:"status"`
}

// ToPublic converts User to UserPublic
func (u *User) ToPublic() *UserPublic {
	return &UserPublic{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Role:      u.Role,
		AvatarURL: u.AvatarURL,
		Status:    u.Status,
	}
}

// RegisterInput for user registration
type RegisterInput struct {
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=8"`
	Name            string `json:"name" binding:"required,min=2,max=100"`
	OrganizationName string `json:"organization_name" binding:"required,min=2,max=100"`
}

// LoginInput for user login
type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// InviteUserInput for inviting a new team member
type InviteUserInput struct {
	Email string   `json:"email" binding:"required,email"`
	Name  string   `json:"name" binding:"required,min=2,max=100"`
	Role  UserRole `json:"role" binding:"required,oneof=admin supervisor agent"`
}

// UpdateUserInput for updating user profile
type UpdateUserInput struct {
	Name      string `json:"name,omitempty"`
	AvatarURL string `json:"avatar_url,omitempty"`
}

// AuthResponse returned after successful authentication
type AuthResponse struct {
	User         *UserPublic   `json:"user"`
	Organization *Organization `json:"organization"`
	AccessToken  string        `json:"access_token"`
	ExpiresIn    int64         `json:"expires_in"`
}
