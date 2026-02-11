package repository

import (
	"github.com/google/uuid"
	"github.com/sidji-omnichannel/internal/models"
)

// AuthRepository defines the outbound port for authentication-related data operations
type AuthRepository interface {
	ExistsByEmail(email string) (bool, error)
	CreateRegisterTransaction(org *models.Organization, user *models.User) error
	GetAuthDataByEmail(email string) (*models.User, *models.Organization, error)
	GetAuthDataByID(userID uuid.UUID) (*models.User, *models.Organization, error)
	UpdateLastSeen(userID uuid.UUID) error
	GetUserByID(userID uuid.UUID) (*models.User, error)
	GetOrganizationCounts(orgID uuid.UUID) (userCount, channelCount int, err error)
}
