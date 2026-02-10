package channels

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sidji-omnichannel/sidji-omnichannel/internal/config"
)

const (
	metaAPIBaseURL = "https://graph.facebook.com/v24.0"
)

// MetaClient handles communication with Meta APIs (WhatsApp & Instagram)
type MetaClient struct {
	httpClient *http.Client
	cfg        *config.MetaConfig
}

// NewMetaClient creates a new Meta API client
func NewMetaClient(cfg *config.MetaConfig) *MetaClient {
	return &MetaClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		cfg: cfg,
	}
}

// ============================================
// WhatsApp Messages
// ============================================

// WhatsAppTextMessage represents a text message to send via WhatsApp
type WhatsAppTextMessage struct {
	To   string `json:"to"`
	Text string `json:"text"`
}

// WhatsAppImageMessage represents an image message to send via WhatsApp
type WhatsAppImageMessage struct {
	To       string `json:"to"`
	ImageURL string `json:"image_url"`
	Caption  string `json:"caption,omitempty"`
}

// WhatsAppDocumentMessage represents a document message to send via WhatsApp
type WhatsAppDocumentMessage struct {
	To          string `json:"to"`
	DocumentURL string `json:"document_url"`
	Filename    string `json:"filename"`
	Caption     string `json:"caption,omitempty"`
}

// WhatsAppTemplateMessage represents a template message to send via WhatsApp
type WhatsAppTemplateMessage struct {
	To           string                 `json:"to"`
	TemplateName string                 `json:"template_name"`
	LanguageCode string                 `json:"language_code"`
	Components   []map[string]interface{} `json:"components,omitempty"`
}

// WhatsAppSendResponse represents the response from sending a message
type WhatsAppSendResponse struct {
	MessagingProduct string `json:"messaging_product"`
	Contacts         []struct {
		Input string `json:"input"`
		WaID  string `json:"wa_id"`
	} `json:"contacts"`
	Messages []struct {
		ID string `json:"id"`
	} `json:"messages"`
}

// SendWhatsAppText sends a text message via WhatsApp
func (c *MetaClient) SendWhatsAppText(phoneNumberID, to, text string, accessToken string) (*WhatsAppSendResponse, error) {
	payload := map[string]interface{}{
		"messaging_product": "whatsapp",
		"recipient_type":    "individual",
		"to":                to,
		"type":              "text",
		"text": map[string]string{
			"preview_url": "false",
			"body":        text,
		},
	}

	return c.sendWhatsAppMessage(phoneNumberID, payload, accessToken)
}

// SendWhatsAppImage sends an image message via WhatsApp
func (c *MetaClient) SendWhatsAppImage(phoneNumberID, to, imageURL, caption string, accessToken string) (*WhatsAppSendResponse, error) {
	imagePayload := map[string]string{
		"link": imageURL,
	}
	if caption != "" {
		imagePayload["caption"] = caption
	}

	payload := map[string]interface{}{
		"messaging_product": "whatsapp",
		"recipient_type":    "individual",
		"to":                to,
		"type":              "image",
		"image":             imagePayload,
	}

	return c.sendWhatsAppMessage(phoneNumberID, payload, accessToken)
}

// SendWhatsAppDocument sends a document message via WhatsApp
func (c *MetaClient) SendWhatsAppDocument(phoneNumberID, to, documentURL, filename, caption string, accessToken string) (*WhatsAppSendResponse, error) {
	docPayload := map[string]string{
		"link":     documentURL,
		"filename": filename,
	}
	if caption != "" {
		docPayload["caption"] = caption
	}

	payload := map[string]interface{}{
		"messaging_product": "whatsapp",
		"recipient_type":    "individual",
		"to":                to,
		"type":              "document",
		"document":          docPayload,
	}

	return c.sendWhatsAppMessage(phoneNumberID, payload, accessToken)
}

// SendWhatsAppTemplate sends a template message via WhatsApp
func (c *MetaClient) SendWhatsAppTemplate(phoneNumberID, to, templateName, languageCode string, components []map[string]interface{}, accessToken string) (*WhatsAppSendResponse, error) {
	templatePayload := map[string]interface{}{
		"name": templateName,
		"language": map[string]string{
			"code": languageCode,
		},
	}
	if len(components) > 0 {
		templatePayload["components"] = components
	}

	payload := map[string]interface{}{
		"messaging_product": "whatsapp",
		"recipient_type":    "individual",
		"to":                to,
		"type":              "template",
		"template":          templatePayload,
	}

	return c.sendWhatsAppMessage(phoneNumberID, payload, accessToken)
}

// SendWhatsAppReaction sends a reaction to a message via WhatsApp
func (c *MetaClient) SendWhatsAppReaction(phoneNumberID, to, messageID, emoji string, accessToken string) (*WhatsAppSendResponse, error) {
	payload := map[string]interface{}{
		"messaging_product": "whatsapp",
		"recipient_type":    "individual",
		"to":                to,
		"type":              "reaction",
		"reaction": map[string]string{
			"message_id": messageID,
			"emoji":      emoji,
		},
	}

	return c.sendWhatsAppMessage(phoneNumberID, payload, accessToken)
}

// MarkWhatsAppMessageAsRead marks a message as read
func (c *MetaClient) MarkWhatsAppMessageAsRead(phoneNumberID, messageID string, accessToken string) error {
	payload := map[string]interface{}{
		"messaging_product": "whatsapp",
		"status":            "read",
		"message_id":        messageID,
	}

	_, err := c.sendWhatsAppMessage(phoneNumberID, payload, accessToken)
	return err
}

// sendWhatsAppMessage sends a message to WhatsApp API
func (c *MetaClient) sendWhatsAppMessage(phoneNumberID string, payload map[string]interface{}, accessToken string) (*WhatsAppSendResponse, error) {
	url := fmt.Sprintf("%s/%s/messages", metaAPIBaseURL, phoneNumberID)

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("WhatsApp API error (status %d): %s", resp.StatusCode, string(body))
	}

	var result WhatsAppSendResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	fmt.Println("WhatsApp API response:", result)
	return &result, nil
}

// GetWhatsAppMediaURL retrieves the URL for a media file
func (c *MetaClient) GetWhatsAppMediaURL(mediaID string, accessToken string) (string, error) {
	url := fmt.Sprintf("%s/%s", metaAPIBaseURL, mediaID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var result struct {
		URL string `json:"url"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	return result.URL, nil
}

// ============================================
// Instagram Messages
// ============================================

// InstagramSendResponse represents the response from sending an Instagram message
type InstagramSendResponse struct {
	RecipientID string `json:"recipient_id"`
	MessageID   string `json:"message_id"`
}

// MetaErrorResponse represents a standard error from Meta API
type MetaErrorResponse struct {
	Error struct {
		Message      string `json:"message"`
		Type         string `json:"type"`
		Code         int    `json:"code"`
		ErrorSubcode int    `json:"error_subcode"`
		FBTraceID    string `json:"fbtrace_id"`
	} `json:"error"`
}

// InstagramUserProfile represents an Instagram user's profile
type InstagramUserProfile struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Username       string `json:"username"`
	ProfilePicture string `json:"profile_pic"`
}

// ExchangeForLongLivedToken exchanges a short-lived user access token for a long-lived one (valid for ~60 days)
func (c *MetaClient) ExchangeForLongLivedToken(shortLivedToken string) (string, int, error) {
	url := fmt.Sprintf(
		"https://graph.facebook.com/v24.0/oauth/access_token?grant_type=fb_exchange_token&client_id=%s&client_secret=%s&fb_exchange_token=%s",
		c.cfg.AppID,
		c.cfg.AppSecret,
		shortLivedToken,
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", 0, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", 0, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", 0, fmt.Errorf("Meta API error (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"` // seconds
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", 0, fmt.Errorf("failed to parse response: %w", err)
	}

	return result.AccessToken, result.ExpiresIn, nil
}

// GetInstagramUserProfile ... (existing method)
func (c *MetaClient) GetInstagramUserProfile(userID, accessToken string) (*InstagramUserProfile, error) {
	// Note: 'username' field on IGSID often requires specific permissions or is unavailable.
	// We use 'profile_pic' instead of 'profile_picture_url' for IGSID.
	url := fmt.Sprintf("%s/%s?fields=id,name,username,profile_pic", metaAPIBaseURL, userID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var metaErr MetaErrorResponse
		if jsonErr := json.Unmarshal(body, &metaErr); jsonErr == nil && metaErr.Error.Code != 0 {
			if metaErr.Error.Code == 190 {
				return nil, fmt.Errorf("token expired: %s", metaErr.Error.Message)
			}
			fmt.Printf("Instagram API Error (code %d): %s\nBody: %s\n", metaErr.Error.Code, metaErr.Error.Message, string(body))
		} else {
			fmt.Printf("Instagram API Error Body: %s\n", string(body))
		}
		
		// If username permission fails, try fallback without username
		if resp.StatusCode == 400 || resp.StatusCode == 403 {
			fmt.Println("Retrying fetch without username field...")
			urlRetry := fmt.Sprintf("%s/%s?fields=id,name,profile_pic", metaAPIBaseURL, userID)
			reqRetry, _ := http.NewRequest("GET", urlRetry, nil)
			reqRetry.Header.Set("Authorization", "Bearer "+accessToken)
			respRetry, errRetry := c.httpClient.Do(reqRetry)
			if errRetry == nil {
				defer respRetry.Body.Close()
				bodyRetry, _ := io.ReadAll(respRetry.Body)
				if respRetry.StatusCode == http.StatusOK {
					var user InstagramUserProfile
					if err := json.Unmarshal(bodyRetry, &user); err == nil {
						return &user, nil
					}
				}
			}
		}
		
		return nil, fmt.Errorf("Instagram API error (status %d): %s", resp.StatusCode, string(body))
	}

	var user InstagramUserProfile
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &user, nil
}

// SendInstagramText sends a text message via Instagram
func (c *MetaClient) SendInstagramText(igUserID, recipientID, text string, accessToken string) (*InstagramSendResponse, error) {
	payload := map[string]interface{}{
		"recipient": map[string]string{
			"id": recipientID,
		},
		"message": map[string]string{
			"text": text,
		},
	}

	return c.sendInstagramMessage(igUserID, payload, accessToken)
}

// SendInstagramImage sends an image message via Instagram
func (c *MetaClient) SendInstagramImage(igUserID, recipientID, imageURL string, accessToken string) (*InstagramSendResponse, error) {
	payload := map[string]interface{}{
		"recipient": map[string]string{
			"id": recipientID,
		},
		"message": map[string]interface{}{
			"attachment": map[string]interface{}{
				"type": "image",
				"payload": map[string]string{
					"url": imageURL,
				},
			},
		},
	}

	return c.sendInstagramMessage(igUserID, payload, accessToken)
}

// SendInstagramGenericTemplate sends a generic template via Instagram
func (c *MetaClient) SendInstagramGenericTemplate(igUserID, recipientID string, elements []map[string]interface{}, accessToken string) (*InstagramSendResponse, error) {
	payload := map[string]interface{}{
		"recipient": map[string]string{
			"id": recipientID,
		},
		"message": map[string]interface{}{
			"attachment": map[string]interface{}{
				"type": "template",
				"payload": map[string]interface{}{
					"template_type": "generic",
					"elements":      elements,
				},
			},
		},
	}

	return c.sendInstagramMessage(igUserID, payload, accessToken)
}

// sendInstagramMessage sends a message to Instagram API
func (c *MetaClient) sendInstagramMessage(igUserID string, payload map[string]interface{}, accessToken string) (*InstagramSendResponse, error) {
	// Using 'me' allows the token context to determine the sender, which is often more reliable
	// if the token is already scoped to the correct Page/User.
	url := fmt.Sprintf("%s/me/messages", metaAPIBaseURL)

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		// Log the full error body for debugging
		fmt.Printf("Instagram API Error Body: %s\n", string(body))
		return nil, fmt.Errorf("Instagram API error (status %d): %s", resp.StatusCode, string(body))
	}

	var result InstagramSendResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result, nil
}
// GetMeAccounts fetches the user's Pages
func (c *MetaClient) GetMeAccounts(accessToken string) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("%s/me/accounts?fields=id,name,access_token,instagram_business_account,tasks", metaAPIBaseURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Meta API error (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data []map[string]interface{} `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return result.Data, nil
}
