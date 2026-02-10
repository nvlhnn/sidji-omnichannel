package twilio

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/sidji-omnichannel/internal/config"
	"github.com/sidji-omnichannel/internal/models"
)

type TwilioProvider struct {
	cfg *config.TwilioConfig
}

func NewTwilioProvider(cfg *config.TwilioConfig) *TwilioProvider {
	return &TwilioProvider{cfg: cfg}
}

func (p *TwilioProvider) GetProviderName() string {
	return "twilio"
}

func (p *TwilioProvider) SendMessage(ctx context.Context, channel *models.Channel, contact *models.Contact, message *models.Message) (string, error) {
	if channel.Type != models.ChannelWhatsApp {
		return "", fmt.Errorf("twilio provider only supports whatsapp")
	}

	// Default credentials from global config
	accountSID := p.cfg.AccountSID
	authToken := p.cfg.AuthToken
	fromNumber := p.cfg.FromNumber

	// Override with per-channel credentials if available (ISV / Subaccount model)
	if channel.Config != nil {
		var waConfig models.WhatsAppConfig
		if err := json.Unmarshal(channel.Config, &waConfig); err == nil {
			if waConfig.TwilioAccountSID != "" {
				accountSID = waConfig.TwilioAccountSID
			}
		}
	}

	// Use channel.AccessToken as AuthToken for Twilio subaccounts if set
	if channel.AccessToken != "" {
		authToken = channel.AccessToken
	}

	// Use channel.PhoneNumberID as the "From" number if set
	if channel.PhoneNumberID != "" {
		fromNumber = channel.PhoneNumberID
	}

	// Twilio WhatsApp API
	apiURL := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", accountSID)

	data := url.Values{}
	data.Set("To", "whatsapp:"+contact.WhatsAppID)
	data.Set("From", "whatsapp:"+fromNumber)

	switch message.MessageType {
	case models.MessageTypeText:
		data.Set("Body", message.Content)
	case models.MessageTypeImage:
		data.Set("MediaUrl", message.MediaURL)
		if message.Content != "" {
			data.Set("Body", message.Content)
		}
	default:
		data.Set("Body", message.Content)
	}

	req, err := http.NewRequest("POST", apiURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(accountSID, authToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("twilio api error (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Sid string `json:"sid"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	return result.Sid, nil
}
