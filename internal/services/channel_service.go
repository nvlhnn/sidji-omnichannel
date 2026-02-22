package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sidji-omnichannel/internal/config"
	"github.com/sidji-omnichannel/internal/domain/ports/repository"
	"github.com/sidji-omnichannel/internal/integrations"
	"github.com/sidji-omnichannel/internal/integrations/meta"
	"github.com/sidji-omnichannel/internal/integrations/tiktok"
	"github.com/sidji-omnichannel/internal/models"
	"github.com/sidji-omnichannel/internal/subscription"
)

var (
	ErrChannelNotFound = errors.New("channel not found")
)

// ChannelService handles channel operations and message sending
type ChannelService struct {
	repo         repository.ChannelRepository
	orgRepo      repository.OrganizationRepository
	cfg          *config.Config
	metaClient   *meta.MetaClient
	tiktokClient *tiktok.TikTokClient
	providers    map[string]integrations.MessageProvider
}

// NewChannelService creates a new channel service
func NewChannelService(
	repo repository.ChannelRepository,
	orgRepo repository.OrganizationRepository,
	cfg *config.Config,
) *ChannelService {
	metaClient := meta.NewMetaClient(&cfg.Meta)
	tiktokClient := tiktok.NewTikTokClient(&cfg.TikTok)
	
	providers := make(map[string]integrations.MessageProvider)
	providers["meta"] = meta.NewMetaProvider(metaClient)
	providers["tiktok"] = tiktok.NewTikTokProvider(tiktokClient)

	return &ChannelService{
		repo:         repo,
		orgRepo:      orgRepo,
		cfg:          cfg,
		metaClient:   metaClient,
		tiktokClient: tiktokClient,
		providers:    providers,
	}
}

// List returns all channels for an organization
func (s *ChannelService) List(orgID uuid.UUID) ([]*models.Channel, error) {
	return s.repo.List(orgID)
}

// Get retrieves a channel by ID (without org check)
func (s *ChannelService) Get(channelID uuid.UUID) (*models.Channel, error) {
	ch, err := s.repo.Get(channelID)
	if err == repository.ErrNotFound {
		return nil, ErrChannelNotFound
	}
	return ch, err
}

// GetByID retrieves a channel by ID
func (s *ChannelService) GetByID(orgID, channelID uuid.UUID) (*models.Channel, error) {
	ch, err := s.repo.GetByID(orgID, channelID)
	if err == repository.ErrNotFound {
		return nil, ErrChannelNotFound
	}
	return ch, err
}

// Create creates a new channel
func (s *ChannelService) Create(orgID uuid.UUID, input *models.CreateChannelInput) (*models.Channel, error) {
	// Check subscription limits
	plan, err := s.orgRepo.GetPlan(orgID)
	if err != nil {
		return nil, err
	}

	channelCount, err := s.orgRepo.GetChannelCount(orgID)
	if err != nil {
		return nil, err
	}

	if err := subscription.CheckLimit(plan, channelCount, "channel"); err != nil {
		return nil, err
	}

	provider := input.Provider
	if provider == "" {
		provider = "meta"
	}

	channel := &models.Channel{
		ID:             uuid.New(),
		OrganizationID: orgID,
		Type:           input.Type,
		Provider:       provider,
		Name:           input.Name,
		AccessToken:    input.AccessToken,
		PhoneNumberID:  input.PhoneNumberID,
		IGUserID:       input.IGUserID,
		FacebookPageID: input.FacebookPageID,
		Status:         models.ChannelStatusPending,
	}

	// Build config based on channel type
	var configData interface{}
	if input.Type == models.ChannelWhatsApp {
		configData = models.WhatsAppConfig{
			BusinessAccountID: input.BusinessAccountID,
		}
	} else if input.Type == models.ChannelInstagram {
		configData = models.InstagramConfig{}
	} else if input.Type == models.ChannelFacebook {
		configData = models.FacebookConfig{
			PageID: input.FacebookPageID,
		}
	} else if input.Type == models.ChannelTikTok {
		configData = models.TikTokConfig{
			OpenID: input.TikTokOpenID,
		}
	}

	configJSON, err := json.Marshal(configData)
	if err != nil {
		return nil, err
	}
	channel.Config = configJSON

	if err := s.repo.Create(channel); err != nil {
		return nil, err
	}

	return channel, nil
}

// Delete deletes a channel
func (s *ChannelService) Delete(orgID, channelID uuid.UUID) error {
	err := s.repo.Delete(orgID, channelID)
	if err == repository.ErrNotFound {
		return ErrChannelNotFound
	}
	return err
}

// UpdateStatus updates channel status
func (s *ChannelService) UpdateStatus(channelID uuid.UUID, status models.ChannelStatus) error {
	return s.repo.UpdateStatus(channelID, status)
}

// SendMessage sends a message through a channel using the appropriate provider
func (s *ChannelService) SendMessage(channel *models.Channel, contact *models.Contact, message *models.Message) (string, error) {
	provider, ok := s.providers[channel.Provider]
	if !ok {
		// Fallback to meta if provider not set or unknown
		provider = s.providers["meta"]
	}

	return provider.SendMessage(context.Background(), channel, contact, message)
}

// GetChannelByPhoneNumberID finds a channel by WhatsApp phone number ID
func (s *ChannelService) GetChannelByPhoneNumberID(phoneNumberID string) (*models.Channel, error) {
	ch, err := s.repo.GetByPhoneNumberID(phoneNumberID)
	if err == repository.ErrNotFound {
		return nil, ErrChannelNotFound
	}
	return ch, err
}

// GetChannelByIGUserID finds a channel by Instagram user ID
func (s *ChannelService) GetChannelByIGUserID(igUserID string) (*models.Channel, error) {
	ch, err := s.repo.GetByIGUserID(igUserID)
	if err == repository.ErrNotFound {
		return nil, ErrChannelNotFound
	}
	return ch, err
}

// GetInstagramUserProfile fetches an Instagram user's profile
func (s *ChannelService) GetInstagramUserProfile(userID, accessToken string) (string, string, string, error) {
	user, err := s.metaClient.GetInstagramUserProfile(userID, accessToken)
	if err != nil {
		fmt.Printf("[ChannelService] Failed to fetch IG profile for %s: %v. Using fallback.\n", userID, err)
		return "", "", "", nil
	}
	return user.Name, user.Username, user.ProfilePicture, nil
}

// ConnectInstagram automatically connects an Instagram channel using a user access token
func (s *ChannelService) ConnectInstagram(orgID uuid.UUID, accessToken string, selectedID string) (*models.Channel, error) {
	if s.cfg.Meta.AppID != "" && s.cfg.Meta.AppSecret != "" {
		longLivedToken, expiresIn, err := s.metaClient.ExchangeForLongLivedToken(accessToken)
		if err == nil && longLivedToken != "" {
			fmt.Printf("Successfully exchanged for Long-Lived Token (expires in %d seconds)\n", expiresIn)
			accessToken = longLivedToken
		} else {
			fmt.Printf("Warning: Failed to exchange for Long-Lived Token: %v. Using original token.\n", err)
		}
	} else {
		fmt.Println("Warning: Meta App ID/Secret not configured. Skipping Long-Lived Token exchange.")
	}

	accounts, err := s.metaClient.GetMeAccounts(accessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch accounts: %w", err)
	}

	var connectedChannel *models.Channel

	for _, account := range accounts {
		igAccount, ok := account["instagram_business_account"].(map[string]interface{})
		if !ok || igAccount == nil {
			continue
		}

		igID, ok := igAccount["id"].(string)
		if !ok || igID == "" {
			continue
		}

		// Filter if selectedID is provided
		if selectedID != "" && igID != selectedID {
			continue
		}

		pageAccessToken, ok := account["access_token"].(string)
		if !ok || pageAccessToken == "" {
			pageAccessToken = accessToken
		}

		existing, err := s.GetChannelByIGUserID(igID)
		if err == nil {
			existing.AccessToken = pageAccessToken
			existing.Status = models.ChannelStatusActive
			if err := s.repo.Update(existing); err != nil {
				return nil, err
			}
			return existing, nil
		}

		channelName := fmt.Sprintf("%s (IG)", account["name"])
		if name, ok := account["name"].(string); ok {
			channelName = name + " Instagram"
		}

		input := &models.CreateChannelInput{
			Type:        models.ChannelInstagram,
			Provider:    "meta",
			Name:        channelName,
			AccessToken: pageAccessToken,
			IGUserID:    igID,
		}

		connectedChannel, err = s.Create(orgID, input)
		if err != nil {
			return nil, err
		}

		pageID, _ := account["id"].(string)
		if pageID != "" {
			if err := s.metaClient.SubscribePageToApp(pageID, pageAccessToken); err != nil {
				fmt.Printf("Failed to subscribe page %s to app: %v\n", pageID, err)
			} else {
				fmt.Printf("Successfully subscribed page %s to app webhooks\n", pageID)
			}
		}

		if err := s.UpdateStatus(connectedChannel.ID, models.ChannelStatusActive); err != nil {
			return nil, err
		}
		connectedChannel.Status = models.ChannelStatusActive

		return connectedChannel, nil
	}

	return nil, errors.New("no instagram professional account found linked to your facebook pages")
}

// StartTokenRefresher starts a background worker to refresh tokens daily
func (s *ChannelService) StartTokenRefresher() {
	go func() {
		fmt.Println("[TokenRefresher] Starting initial token refresh check...")
		if err := s.RefreshTokens(); err != nil {
			fmt.Printf("[TokenRefresher] Initial refresh failed: %v\n", err)
		}

		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			fmt.Println("[TokenRefresher] Starting daily token refresh check...")
			if err := s.RefreshTokens(); err != nil {
				fmt.Printf("[TokenRefresher] Daily refresh failed: %v\n", err)
			}
		}
	}()
}

// RefreshTokens refreshes access tokens for all active channels
func (s *ChannelService) RefreshTokens() error {
	if s.cfg.Meta.AppID == "" || s.cfg.Meta.AppSecret == "" {
		return nil 
	}

	channels, err := s.repo.ListActiveMeta()
	if err != nil {
		return err
	}

	for _, ch := range channels {
		newToken, expiresIn, err := s.metaClient.ExchangeForLongLivedToken(ch.AccessToken)
		if err != nil {
			fmt.Printf("[TokenRefresher] Failed to refresh token for channel %s (%s): %v\n", ch.Name, ch.ID, err)
			continue
		}

		ch.AccessToken = newToken
		if err := s.repo.Update(ch); err != nil {
			fmt.Printf("[TokenRefresher] Failed to update token in DB for likely channel %s (%s): %v\n", ch.Name, ch.ID, err)
			continue
		}

		fmt.Printf("[TokenRefresher] Successfully refreshed token for channel %s. Expires in ~%d days\n", ch.Name, expiresIn/86400)
	}

	return nil
}

// ConnectWhatsApp automatically connects a WhatsApp channel using a user access token from embedded signup
func (s *ChannelService) ConnectWhatsApp(orgID uuid.UUID, accessToken string, selectedID string) (*models.Channel, error) {
	if s.cfg.Meta.AppID != "" && s.cfg.Meta.AppSecret != "" {
		longLivedToken, _, err := s.metaClient.ExchangeForLongLivedToken(accessToken)
		if err == nil && longLivedToken != "" {
			accessToken = longLivedToken
		}
	}

	wabas, err := s.metaClient.GetMeWABAs(accessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch WABAs: %w", err)
	}

	if len(wabas) == 0 {
		return nil, errors.New("no whatsapp business accounts found")
	}

	for _, waba := range wabas {
		wabaID, _ := waba["id"].(string)
		wabaName, _ := waba["name"].(string)
		if wabaID == "" {
			continue
		}

		phoneNumbers, err := s.metaClient.GetWABAPhoneNumbers(wabaID, accessToken)
		if err != nil {
			fmt.Printf("Warning: Failed to fetch phone numbers for WABA %s: %v\n", wabaID, err)
			continue
		}

		for _, pn := range phoneNumbers {
			pnID, _ := pn["id"].(string)
			displayPN, _ := pn["display_phone_number"].(string)
			if pnID == "" {
				continue
			}

			// Filter if selectedID is provided
			if selectedID != "" && pnID != selectedID {
				continue
			}

			// Check if already exists
			existing, err := s.GetChannelByPhoneNumberID(pnID)
			if err == nil {
				// Update existing
				existing.AccessToken = accessToken
				existing.Status = models.ChannelStatusActive
				if err := s.repo.Update(existing); err != nil {
					return nil, err
				}
				return existing, nil
			}

			// Create new
			channelName := wabaName
			if displayPN != "" {
				channelName = fmt.Sprintf("WhatsApp (%s)", displayPN)
			}

			input := &models.CreateChannelInput{
				Type:              models.ChannelWhatsApp,
				Name:              channelName,
				Provider:          "meta",
				AccessToken:       accessToken,
				PhoneNumberID:     pnID,
				BusinessAccountID: wabaID,
			}

			channel, err := s.Create(orgID, input)
			if err != nil {
				return nil, err
			}

			// Activate immediately
			if err := s.UpdateStatus(channel.ID, models.ChannelStatusActive); err != nil {
				return nil, err
			}
			channel.Status = models.ChannelStatusActive

			return channel, nil
		}
	}

	return nil, errors.New("no phone numbers found in your whatsapp business accounts")
}

// GetChannelByFacebookPageID finds a channel by Facebook page ID
func (s *ChannelService) GetChannelByFacebookPageID(pageID string) (*models.Channel, error) {
	ch, err := s.repo.GetByFacebookPageID(pageID)
	if err == repository.ErrNotFound {
		return nil, ErrChannelNotFound
	}
	return ch, err
}

// GetChannelByTikTokOpenID finds a channel by TikTok open ID
func (s *ChannelService) GetChannelByTikTokOpenID(tiktokOpenID string) (*models.Channel, error) {
	ch, err := s.repo.GetByTikTokOpenID(tiktokOpenID)
	if err == repository.ErrNotFound {
		return nil, ErrChannelNotFound
	}
	return ch, err
}

// ConnectFacebook automatically connects Facebook Pages
func (s *ChannelService) ConnectFacebook(orgID uuid.UUID, accessToken string, selectedID string) ([]*models.Channel, error) {
	if s.cfg.Meta.AppID != "" && s.cfg.Meta.AppSecret != "" {
		longLivedToken, _, err := s.metaClient.ExchangeForLongLivedToken(accessToken)
		if err == nil && longLivedToken != "" {
			accessToken = longLivedToken
		}
	}

	accounts, err := s.metaClient.GetMeAccounts(accessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch accounts: %w", err)
	}

	var connectedChannels []*models.Channel

	for _, account := range accounts {
		pageID, _ := account["id"].(string)
		pageName, _ := account["name"].(string)
		pageAccessToken, _ := account["access_token"].(string)

		if pageID == "" || pageAccessToken == "" {
			continue
		}

		// Filter if selectedID is provided
		if selectedID != "" && pageID != selectedID {
			continue
		}

		// Check if already exists
		existing, err := s.GetChannelByFacebookPageID(pageID)
		if err == nil {
			existing.AccessToken = pageAccessToken
			existing.Status = models.ChannelStatusActive
			if err := s.repo.Update(existing); err != nil {
				return nil, err
			}
			connectedChannels = append(connectedChannels, existing)
			continue
		}

		input := &models.CreateChannelInput{
			Type:           models.ChannelFacebook,
			Provider:       "meta",
			Name:           pageName + " (Facebook)",
			AccessToken:    pageAccessToken,
			FacebookPageID: pageID,
		}

		channel, err := s.Create(orgID, input)
		if err != nil {
			return nil, err
		}

		// Subscribe page to app webhooks
		if err := s.metaClient.SubscribePageToApp(pageID, pageAccessToken); err != nil {
			fmt.Printf("Failed to subscribe page %s to app: %v\n", pageID, err)
		}

		if err := s.UpdateStatus(channel.ID, models.ChannelStatusActive); err != nil {
			return nil, err
		}
		channel.Status = models.ChannelStatusActive

		connectedChannels = append(connectedChannels, channel)
	}

	return connectedChannels, nil
}

// DiscoverMetaAccounts finds all available Meta accounts (FB Pages, IG, WhatsApp)
func (s *ChannelService) DiscoverMetaAccounts(accessToken string) (map[string]interface{}, error) {
	result := map[string]interface{}{
		"facebook_pages":   []interface{}{},
		"instagram_users":  []interface{}{},
		"whatsapp_numbers": []interface{}{},
	}

	// 1. Discovery Pages & Instagram
	accounts, err := s.metaClient.GetMeAccounts(accessToken)
	if err == nil {
		for _, acc := range accounts {
			// Add to Pages
			result["facebook_pages"] = append(result["facebook_pages"].([]interface{}), map[string]interface{}{
				"id":   acc["id"],
				"name": acc["name"],
			})

			// Add to Instagram if exists
			if ig, ok := acc["instagram_business_account"].(map[string]interface{}); ok && ig != nil {
				result["instagram_users"] = append(result["instagram_users"].([]interface{}), map[string]interface{}{
					"id":        ig["id"],
					"name":      fmt.Sprintf("%s (Instagram)", acc["name"]),
					"parent_id": acc["id"],
				})
			}
		}
	}

	// 2. Discovery WhatsApp
	wabas, err := s.metaClient.GetMeWABAs(accessToken)
	if err == nil {
		for _, waba := range wabas {
			wabaID, _ := waba["id"].(string)
			if wabaID == "" {
				continue
			}

			numbers, nerr := s.metaClient.GetWABAPhoneNumbers(wabaID, accessToken)
			if nerr == nil {
				for _, num := range numbers {
					result["whatsapp_numbers"] = append(result["whatsapp_numbers"].([]interface{}), map[string]interface{}{
						"id":           num["id"],
						"display_name": num["display_phone_number"],
						"waba_id":      wabaID,
						"waba_name":    waba["name"],
					})
				}
			}
		}
	}

	return result, nil
}

// GetFacebookUserProfile retrieves a Facebook user's profile
func (s *ChannelService) GetFacebookUserProfile(psid, accessToken string) (string, string, error) {
	profile, err := s.metaClient.GetFacebookUserProfile(psid, accessToken)
	if err != nil {
		return "", "", err
	}
	name := fmt.Sprintf("%s %s", profile.FirstName, profile.LastName)
	return name, profile.ProfilePicture, nil
}

// ConnectTikTok connects a TikTok channel via OAuth code exchange
func (s *ChannelService) ConnectTikTok(orgID uuid.UUID, code string) (*models.Channel, error) {
	// Exchange code for access token
	tokenResp, err := s.tiktokClient.ExchangeCodeForToken(code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange TikTok code: %w", err)
	}

	openID := tokenResp.GetOpenID()
	accessToken := tokenResp.GetAccessToken()
	refreshToken := tokenResp.GetRefreshToken()

	if openID == "" {
		return nil, errors.New("tiktok oauth returned empty open_id")
	}

	// Check if channel already exists
	existing, err := s.GetChannelByTikTokOpenID(openID)
	if err == nil {
		// Update existing channel
		existing.AccessToken = accessToken
		existing.Status = models.ChannelStatusActive
		// Update config with new refresh token
		tikTokProvider := s.providers["tiktok"].(*tiktok.TikTokProvider)
		displayName, username, avatarURL, _ := tikTokProvider.GetUserProfile(accessToken)
		configData := models.TikTokConfig{
			OpenID:       openID,
			Username:     username,
			DisplayName:  displayName,
			AvatarURL:    avatarURL,
			RefreshToken: refreshToken,
		}
		configJSON, _ := json.Marshal(configData)
		existing.Config = configJSON
		if err := s.repo.Update(existing); err != nil {
			return nil, err
		}
		return existing, nil
	}

	// Get user profile
	tikTokProvider := s.providers["tiktok"].(*tiktok.TikTokProvider)
	displayName, username, avatarURL, _ := tikTokProvider.GetUserProfile(accessToken)

	channelName := "TikTok"
	if displayName != "" {
		channelName = displayName + " (TikTok)"
	} else if username != "" {
		channelName = "@" + username + " (TikTok)"
	}

	// Build config
	configData := models.TikTokConfig{
		OpenID:       openID,
		Username:     username,
		DisplayName:  displayName,
		AvatarURL:    avatarURL,
		RefreshToken: refreshToken,
	}
	configJSON, err := json.Marshal(configData)
	if err != nil {
		return nil, err
	}

	channel := &models.Channel{
		ID:             uuid.New(),
		OrganizationID: orgID,
		Type:           models.ChannelTikTok,
		Provider:       "tiktok",
		Name:           channelName,
		AccessToken:    accessToken,
		TikTokOpenID:   openID,
		Config:         configJSON,
		Status:         models.ChannelStatusActive,
	}

	if err := s.repo.Create(channel); err != nil {
		return nil, err
	}

	return channel, nil
}
