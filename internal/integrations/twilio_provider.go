package integrations

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

	// Twilio WhatsApp API
	apiURL := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", p.cfg.AccountSID)

	data := url.Values{}
	data.Set("To", "whatsapp:"+contact.WhatsAppID)
	
	// Determine the "From" number. 
	// If channel has a specific phone number ID or from number in config, use it.
	from := "whatsapp:" + p.cfg.FromNumber
	if channel.PhoneNumberID != "" {
		from = "whatsapp:" + channel.PhoneNumberID
	}
	data.Set("From", from)

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
	req.SetBasicAuth(p.cfg.AccountSID, p.cfg.AuthToken)

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
