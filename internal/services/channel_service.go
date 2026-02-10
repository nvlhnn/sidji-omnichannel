package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sidji-omnichannel/internal/config"
	"github.com/sidji-omnichannel/internal/integrations"
	"github.com/sidji-omnichannel/internal/integrations/meta"
	"github.com/sidji-omnichannel/internal/models"
	"github.com/sidji-omnichannel/internal/subscription"
)

var (
	ErrChannelNotFound = errors.New("channel not found")
)

// ChannelService handles channel operations and message sending
type ChannelService struct {
	db         *sql.DB
	cfg        *config.Config
	metaClient *meta.MetaClient
	providers  map[string]integrations.MessageProvider
}

// NewChannelService creates a new channel service
func NewChannelService(db *sql.DB, cfg *config.Config) *ChannelService {
	metaClient := meta.NewMetaClient(&cfg.Meta)
	
	providers := make(map[string]integrations.MessageProvider)
	providers["meta"] = meta.NewMetaProvider(metaClient)

	return &ChannelService{
		db:         db,
		cfg:        cfg,
		metaClient: metaClient,
		providers:  providers,
	}
}

// List returns all channels for an organization
func (s *ChannelService) List(orgID uuid.UUID) ([]*models.Channel, error) {
	rows, err := s.db.Query(`
		SELECT id, organization_id, type, provider, name, config, phone_number_id, ig_user_id, facebook_page_id, status, created_at
		FROM channels
		WHERE organization_id = $1
		ORDER BY created_at DESC
	`, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var channelsList []*models.Channel
	for rows.Next() {
		channel := &models.Channel{}
		var pnID, igID, fbID sql.NullString
		err := rows.Scan(
			&channel.ID, &channel.OrganizationID, &channel.Type, &channel.Provider, &channel.Name,
			&channel.Config, &pnID, &igID, &fbID, &channel.Status, &channel.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		if pnID.Valid { channel.PhoneNumberID = pnID.String }
		if igID.Valid { channel.IGUserID = igID.String }
		if fbID.Valid { channel.FacebookPageID = fbID.String }
		channelsList = append(channelsList, channel)
	}

	return channelsList, nil
}

// Get retrieves a channel by ID (without org check)
func (s *ChannelService) Get(channelID uuid.UUID) (*models.Channel, error) {
	channel := &models.Channel{}
	var pnID, igID, fbID sql.NullString
	err := s.db.QueryRow(`
		SELECT id, organization_id, type, provider, name, config, access_token, phone_number_id, ig_user_id, facebook_page_id, status, created_at
		FROM channels
		WHERE id = $1
	`, channelID).Scan(
		&channel.ID, &channel.OrganizationID, &channel.Type, &channel.Provider, &channel.Name,
		&channel.Config, &channel.AccessToken, &pnID, &igID, &fbID,
		&channel.Status, &channel.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrChannelNotFound
	}
	if err != nil {
		return nil, err
	}
	if pnID.Valid { channel.PhoneNumberID = pnID.String }
	if igID.Valid { channel.IGUserID = igID.String }
	if fbID.Valid { channel.FacebookPageID = fbID.String }
	return channel, nil
}

// GetByID retrieves a channel by ID
func (s *ChannelService) GetByID(orgID, channelID uuid.UUID) (*models.Channel, error) {
	channel := &models.Channel{}
	var pnID, igID, fbID sql.NullString
	err := s.db.QueryRow(`
		SELECT id, organization_id, type, provider, name, config, access_token, phone_number_id, ig_user_id, facebook_page_id, status, created_at
		FROM channels
		WHERE id = $1 AND organization_id = $2
	`, channelID, orgID).Scan(
		&channel.ID, &channel.OrganizationID, &channel.Type, &channel.Provider, &channel.Name,
		&channel.Config, &channel.AccessToken, &pnID, &igID, &fbID,
		&channel.Status, &channel.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrChannelNotFound
	}
	if err != nil {
		return nil, err
	}
	if pnID.Valid { channel.PhoneNumberID = pnID.String }
	if igID.Valid { channel.IGUserID = igID.String }
	if fbID.Valid { channel.FacebookPageID = fbID.String }
	return channel, nil
}

// Create creates a new channel
func (s *ChannelService) Create(orgID uuid.UUID, input *models.CreateChannelInput) (*models.Channel, error) {
	// Check subscription limits
	var plan string
	err := s.db.QueryRow("SELECT plan FROM organizations WHERE id = $1", orgID).Scan(&plan)
	if err != nil {
		return nil, err
	}

	var channelCount int
	err = s.db.QueryRow("SELECT COUNT(*) FROM channels WHERE organization_id = $1", orgID).Scan(&channelCount)
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
	}

	configJSON, err := json.Marshal(configData)
	if err != nil {
		return nil, err
	}
	channel.Config = configJSON

	_, err = s.db.Exec(`
		INSERT INTO channels (id, organization_id, type, provider, name, config, access_token, phone_number_id, ig_user_id, facebook_page_id, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`, channel.ID, channel.OrganizationID, channel.Type, channel.Provider, channel.Name, channel.Config,
		channel.AccessToken, channel.PhoneNumberID, channel.IGUserID, channel.FacebookPageID, channel.Status)

	if err != nil {
		return nil, err
	}

	return channel, nil
}

// Delete deletes a channel
func (s *ChannelService) Delete(orgID, channelID uuid.UUID) error {
	result, err := s.db.Exec(`
		DELETE FROM channels WHERE id = $1 AND organization_id = $2
	`, channelID, orgID)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrChannelNotFound
	}

	return nil
}

// UpdateStatus updates channel status
func (s *ChannelService) UpdateStatus(channelID uuid.UUID, status models.ChannelStatus) error {
	_, err := s.db.Exec(`
		UPDATE channels SET status = $1, updated_at = NOW() WHERE id = $2
	`, status, channelID)
	return err
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
	channel := &models.Channel{}
	var pnID, igID, fbID sql.NullString
	err := s.db.QueryRow(`
		SELECT id, organization_id, type, provider, name, config, access_token, phone_number_id, ig_user_id, facebook_page_id, status
		FROM channels
		WHERE phone_number_id = $1 AND status = 'active'
	`, phoneNumberID).Scan(
		&channel.ID, &channel.OrganizationID, &channel.Type, &channel.Provider, &channel.Name,
		&channel.Config, &channel.AccessToken, &pnID, &igID, &fbID, &channel.Status,
	)
	if err == sql.ErrNoRows {
		return nil, ErrChannelNotFound
	}
	if err != nil {
		return nil, err
	}
	if pnID.Valid { channel.PhoneNumberID = pnID.String }
	if igID.Valid { channel.IGUserID = igID.String }
	if fbID.Valid { channel.FacebookPageID = fbID.String }
	return channel, nil
}

// GetChannelByIGUserID finds a channel by Instagram user ID
func (s *ChannelService) GetChannelByIGUserID(igUserID string) (*models.Channel, error) {
	channel := &models.Channel{}
	var pnID, igID, fbID sql.NullString
	err := s.db.QueryRow(`
		SELECT id, organization_id, type, provider, name, config, access_token, phone_number_id, ig_user_id, facebook_page_id, status
		FROM channels
		WHERE ig_user_id = $1 AND status = 'active'
	`, igUserID).Scan(
		&channel.ID, &channel.OrganizationID, &channel.Type, &channel.Provider, &channel.Name,
		&channel.Config, &channel.AccessToken, &pnID, &igID, &fbID, &channel.Status,
	)
	if err == sql.ErrNoRows {
		return nil, ErrChannelNotFound
	}
	if err != nil {
		return nil, err
	}
	if pnID.Valid { channel.PhoneNumberID = pnID.String }
	if igID.Valid { channel.IGUserID = igID.String }
	if fbID.Valid { channel.FacebookPageID = fbID.String }
	return channel, nil
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
			_, execErr := s.db.Exec(`UPDATE channels SET access_token=$1, status='active', updated_at=NOW() WHERE id=$2`, pageAccessToken, existing.ID)
			if execErr != nil {
				return nil, execErr
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

	rows, err := s.db.Query(`
		SELECT id, name, type, access_token 
		FROM channels 
		WHERE status = 'active' AND type IN ('instagram', 'whatsapp')
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id uuid.UUID
		var name string
		var chType models.ChannelType
		var token string

		if err := rows.Scan(&id, &name, &chType, &token); err != nil {
			fmt.Printf("[TokenRefresher] Failed to scan channel: %v\n", err)
			continue
		}

		newToken, expiresIn, err := s.metaClient.ExchangeForLongLivedToken(token)
		if err != nil {
			fmt.Printf("[TokenRefresher] Failed to refresh token for channel %s (%s): %v\n", name, id, err)
			continue
		}

		_, err = s.db.Exec(`UPDATE channels SET access_token = $1, updated_at = NOW() WHERE id = $2`, newToken, id)
		if err != nil {
			fmt.Printf("[TokenRefresher] Failed to update token in DB for likely channel %s (%s): %v\n", name, id, err)
			continue
		}

		fmt.Printf("[TokenRefresher] Successfully refreshed token for channel %s. Expires in ~%d days\n", name, expiresIn/86400)
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
				_, execErr := s.db.Exec(`UPDATE channels SET access_token=$1, status='active', updated_at=NOW() WHERE id=$2`, accessToken, existing.ID)
				if execErr != nil {
					return nil, execErr
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
	channel := &models.Channel{}
	var pnID, igID, fbID sql.NullString
	err := s.db.QueryRow(`
		SELECT id, organization_id, type, provider, name, config, access_token, phone_number_id, ig_user_id, facebook_page_id, status
		FROM channels
		WHERE facebook_page_id = $1 AND status = 'active'
	`, pageID).Scan(
		&channel.ID, &channel.OrganizationID, &channel.Type, &channel.Provider, &channel.Name,
		&channel.Config, &channel.AccessToken, &pnID, &igID, &fbID, &channel.Status,
	)
	if err == sql.ErrNoRows {
		return nil, ErrChannelNotFound
	}
	if err != nil {
		return nil, err
	}
	if pnID.Valid { channel.PhoneNumberID = pnID.String }
	if igID.Valid { channel.IGUserID = igID.String }
	if fbID.Valid { channel.FacebookPageID = fbID.String }
	return channel, nil
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
			_, execErr := s.db.Exec(`UPDATE channels SET access_token=$1, status='active', updated_at=NOW() WHERE id=$2`, pageAccessToken, existing.ID)
			if execErr != nil {
				return nil, execErr
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
