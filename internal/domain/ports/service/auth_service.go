package service

import (
	"github.com/google/uuid"
	"github.com/sidji-omnichannel/internal/models"
)

// AuthService defines the inbound port for authentication operations
type AuthService interface {
	Register(input *models.RegisterInput) (*models.AuthResponse, error)
	Login(input *models.LoginInput) (*models.AuthResponse, error)
	GetUserByID(userID uuid.UUID) (*models.User, error)
	GetMe(userID uuid.UUID) (*models.AuthResponse, error)
	GoogleLogin(info *models.GoogleUserInfo) (*models.AuthResponse, error)
}
