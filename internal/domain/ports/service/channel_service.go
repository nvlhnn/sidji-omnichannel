package service

import (
	"github.com/google/uuid"
	"github.com/sidji-omnichannel/internal/models"
)

// ChannelService defines the inbound port for channel operations
type ChannelService interface {
	List(orgID uuid.UUID) ([]*models.Channel, error)
	Get(channelID uuid.UUID) (*models.Channel, error)
	GetByID(orgID, channelID uuid.UUID) (*models.Channel, error)
	Create(orgID uuid.UUID, input *models.CreateChannelInput) (*models.Channel, error)
	Delete(orgID, channelID uuid.UUID) error
	UpdateStatus(channelID uuid.UUID, status models.ChannelStatus) error
	SendMessage(channel *models.Channel, contact *models.Contact, message *models.Message) (string, error)
	GetChannelByPhoneNumberID(phoneNumberID string) (*models.Channel, error)
	GetChannelByIGUserID(igUserID string) (*models.Channel, error)
	GetChannelByFacebookPageID(pageID string) (*models.Channel, error)
	GetChannelByTikTokOpenID(tiktokOpenID string) (*models.Channel, error)
	GetInstagramUserProfile(userID, accessToken string) (string, string, string, error)
	GetFacebookUserProfile(psid, accessToken string) (string, string, error)
	ConnectInstagram(orgID uuid.UUID, accessToken string, selectedID string) (*models.Channel, error)
	ConnectWhatsApp(orgID uuid.UUID, accessToken string, selectedID string) (*models.Channel, error)
	ConnectFacebook(orgID uuid.UUID, accessToken string, selectedID string) ([]*models.Channel, error)
	ConnectTikTok(orgID uuid.UUID, code string) (*models.Channel, error)
	DiscoverMetaAccounts(accessToken string) (map[string]interface{}, error)
	RefreshTokens() error
	StartTokenRefresher()
}
