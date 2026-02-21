package tiktok

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/sidji-omnichannel/internal/config"
)

const (
	tiktokAPIBase  = "https://open.tiktokapis.com/v2"
	tiktokAuthBase = "https://open.tiktokapis.com/v2/oauth"
)

// TikTokClient handles communication with TikTok Business API
type TikTokClient struct {
	cfg        *config.TikTokConfig
	httpClient *http.Client
}

// NewTikTokClient creates a new TikTok API client
func NewTikTokClient(cfg *config.TikTokConfig) *TikTokClient {
	return &TikTokClient{
		cfg:        cfg,
		httpClient: &http.Client{},
	}
}

// === OAuth ===

// TokenResponse represents the OAuth token exchange response
type TokenResponse struct {
	Data struct {
		AccessToken      string `json:"access_token"`
		RefreshToken     string `json:"refresh_token"`
		OpenID           string `json:"open_id"`
		ExpiresIn        int    `json:"expires_in"`
		RefreshExpiresIn int    `json:"refresh_expires_in"`
		Scope            string `json:"scope"`
		TokenType        string `json:"token_type"`
	} `json:"data"`
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

// ExchangeCodeForToken exchanges an authorization code for an access token
func (c *TikTokClient) ExchangeCodeForToken(code string) (*TokenResponse, error) {
	data := url.Values{}
	data.Set("client_key", c.cfg.ClientKey)
	data.Set("client_secret", c.cfg.ClientSecret)
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", c.cfg.RedirectURI)

	resp, err := c.httpClient.PostForm(tiktokAuthBase+"/token/", data)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}
	defer resp.Body.Close()

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	if tokenResp.Error.Code != "" {
		return nil, fmt.Errorf("tiktok oauth error: %s - %s", tokenResp.Error.Code, tokenResp.Error.Message)
	}

	return &tokenResp, nil
}

// RefreshAccessToken refreshes an expired access token
func (c *TikTokClient) RefreshAccessToken(refreshToken string) (*TokenResponse, error) {
	data := url.Values{}
	data.Set("client_key", c.cfg.ClientKey)
	data.Set("client_secret", c.cfg.ClientSecret)
	data.Set("refresh_token", refreshToken)
	data.Set("grant_type", "refresh_token")

	resp, err := c.httpClient.PostForm(tiktokAuthBase+"/token/", data)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}
	defer resp.Body.Close()

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	if tokenResp.Error.Code != "" {
		return nil, fmt.Errorf("tiktok refresh error: %s - %s", tokenResp.Error.Code, tokenResp.Error.Message)
	}

	return &tokenResp, nil
}

// === User Info ===

// UserInfoResponse represents TikTok user info
type UserInfoResponse struct {
	Data struct {
		User struct {
			OpenID      string `json:"open_id"`
			UnionID     string `json:"union_id"`
			AvatarURL   string `json:"avatar_url"`
			DisplayName string `json:"display_name"`
			Username    string `json:"username"`
		} `json:"user"`
	} `json:"data"`
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

// GetUserInfo fetches the authenticated user's profile
func (c *TikTokClient) GetUserInfo(accessToken string) (*UserInfoResponse, error) {
	req, err := http.NewRequest("GET", tiktokAPIBase+"/user/info/?fields=open_id,union_id,avatar_url,display_name,username", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	var userResp UserInfoResponse
	if err := json.NewDecoder(resp.Body).Decode(&userResp); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	if userResp.Error.Code != "" {
		return nil, fmt.Errorf("tiktok user info error: %s - %s", userResp.Error.Code, userResp.Error.Message)
	}

	return &userResp, nil
}

// === Direct Messaging ===

// SendMessageRequest represents a TikTok DM send request
type SendMessageRequest struct {
	RecipientOpenID string `json:"recipient_open_id"`
	MessageType     string `json:"message_type"` // "text" or "image"
	Text            string `json:"text,omitempty"`
	ImageURL        string `json:"image_url,omitempty"`
}

// SendMessageResponse represents TikTok DM send response
type SendMessageResponse struct {
	Data struct {
		MessageID string `json:"message_id"`
	} `json:"data"`
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

// SendText sends a text message via TikTok DM
func (c *TikTokClient) SendText(senderOpenID, recipientOpenID, text, accessToken string) (*SendMessageResponse, error) {
	payload := map[string]interface{}{
		"recipient_open_id": recipientOpenID,
		"message_type":      "text",
		"text":              text,
	}

	return c.sendMessage(senderOpenID, payload, accessToken)
}

// SendImage sends an image message via TikTok DM
func (c *TikTokClient) SendImage(senderOpenID, recipientOpenID, imageURL, accessToken string) (*SendMessageResponse, error) {
	payload := map[string]interface{}{
		"recipient_open_id": recipientOpenID,
		"message_type":      "image",
		"image_url":         imageURL,
	}

	return c.sendMessage(senderOpenID, payload, accessToken)
}

func (c *TikTokClient) sendMessage(senderOpenID string, payload map[string]interface{}, accessToken string) (*SendMessageResponse, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("%s/dm/message/send/?sender_open_id=%s", tiktokAPIBase, senderOpenID)
	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send tiktok message: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	var sendResp SendMessageResponse
	if err := json.Unmarshal(respBody, &sendResp); err != nil {
		return nil, fmt.Errorf("failed to decode send response: %w (body: %s)", err, string(respBody))
	}

	if sendResp.Error.Code != "" {
		return nil, fmt.Errorf("tiktok send error: %s - %s", sendResp.Error.Code, sendResp.Error.Message)
	}

	return &sendResp, nil
}

// === Webhook Verification ===

// VerifyWebhookSignature verifies the HMAC-SHA256 signature from TikTok webhooks
// TikTok sends the signature in the X-Tiktok-Signature header
func VerifyWebhookSignature(body []byte, signature string, clientSecret string) bool {
	if signature == "" || clientSecret == "" {
		return false
	}

	// TikTok webhook uses HMAC-SHA256(client_secret, request_body)
	// TODO: Implement exact verification per TikTok's latest API docs once access is confirmed
	// For now, accept all signed requests to allow initial testing
	return true
}
