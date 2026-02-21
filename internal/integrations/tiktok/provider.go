package tiktok

import (
	"context"
	"fmt"

	"github.com/sidji-omnichannel/internal/models"
)

// TikTokProvider implements the MessageProvider interface for TikTok
type TikTokProvider struct {
	client *TikTokClient
}

// NewTikTokProvider creates a new TikTok message provider
func NewTikTokProvider(client *TikTokClient) *TikTokProvider {
	return &TikTokProvider{client: client}
}

// GetProviderName returns the provider identifier
func (p *TikTokProvider) GetProviderName() string {
	return "tiktok"
}

// SendMessage sends a message through TikTok DM
func (p *TikTokProvider) SendMessage(ctx context.Context, channel *models.Channel, contact *models.Contact, message *models.Message) (string, error) {
	switch message.MessageType {
	case models.MessageTypeText:
		resp, err := p.client.SendText(
			channel.TikTokOpenID,
			contact.TikTokID,
			message.Content,
			channel.AccessToken,
		)
		if err != nil {
			return "", err
		}
		return resp.Data.MessageID, nil

	case models.MessageTypeImage:
		resp, err := p.client.SendImage(
			channel.TikTokOpenID,
			contact.TikTokID,
			message.MediaURL,
			channel.AccessToken,
		)
		if err != nil {
			return "", err
		}
		return resp.Data.MessageID, nil

	default:
		// Fallback: send as text
		resp, err := p.client.SendText(
			channel.TikTokOpenID,
			contact.TikTokID,
			message.Content,
			channel.AccessToken,
		)
		if err != nil {
			return "", err
		}
		return resp.Data.MessageID, nil
	}
}

// GetClient returns the underlying TikTok client for direct API calls
func (p *TikTokProvider) GetClient() *TikTokClient {
	return p.client
}

// GetUserProfile fetches a TikTok user's display name and avatar
func (p *TikTokProvider) GetUserProfile(accessToken string) (displayName, username, avatarURL string, err error) {
	userResp, err := p.client.GetUserInfo(accessToken)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to get tiktok user info: %w", err)
	}
	return userResp.Data.User.DisplayName, userResp.Data.User.Username, userResp.Data.User.AvatarURL, nil
}
