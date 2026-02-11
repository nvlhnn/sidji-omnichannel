package repository

import (
	"github.com/google/uuid"
	"github.com/sidji-omnichannel/internal/models"
)

// ChannelRepository defines the outbound port for channel data
type ChannelRepository interface {
	List(orgID uuid.UUID) ([]*models.Channel, error)
	Get(channelID uuid.UUID) (*models.Channel, error)
	GetByID(orgID, channelID uuid.UUID) (*models.Channel, error)
	GetByPhoneNumberID(phoneNumberID string) (*models.Channel, error)
	GetByIGUserID(igUserID string) (*models.Channel, error)
	GetByFacebookPageID(pageID string) (*models.Channel, error)
	Create(channel *models.Channel) error
	Update(channel *models.Channel) error
	UpdateStatus(channelID uuid.UUID, status models.ChannelStatus) error
	Delete(orgID, channelID uuid.UUID) error
	ListActiveMeta() ([]*models.Channel, error)
}
