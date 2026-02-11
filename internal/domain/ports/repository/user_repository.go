package repository

import (
	"github.com/google/uuid"
	"github.com/sidji-omnichannel/internal/models"
)

// UserRepository defines the outbound port for user data
type UserRepository interface {
	GetPublic(userID uuid.UUID) (*models.UserPublic, error)
}
