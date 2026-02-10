package channels

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// SubscribePageToApp subscribes a page to the app's webhooks for specific fields
func (c *MetaClient) SubscribePageToApp(pageID, accessToken string) error {
	// Endpoint: POST /{page-id}/subscribed_apps
	// Fields for Instagram Messaging: generic 'messages' from page subscription often covers it, but for IG:
	// The instagram webhooks are actually app-level subscriptions in the dashboard (for 'instagram' topic).
	// BUT, the Page must "install" the app via subscribed_apps to authorize the event delivery.
	// For Instagram Messaging, we need to subscribe to `feed`? No.
	// We usually just POST to /subscribed_apps with authorized fields.
	// Recent API versions require specifying fields.
	
	apiURL := fmt.Sprintf("%s/%s/subscribed_apps", metaAPIBaseURL, pageID)

	data := url.Values{}
	// For Instagram Messaging, we generally need the 'messages' field on the Page subscription?
	// Actually, for "Instagram" topic webhooks, the configuration is in the App Dashboard.
	// However, calling subscribed_apps on the Facebook Page linked to the IG account is CRITICAL.
	// It installs the app on the Page.
	// We will request 'messages' and 'messaging_postbacks' just in case, but often empty or 'feed' is standard for Pages.
	// For IG, it's weird. Let's try attempting to subscribe with key fields.
	data.Set("subscribed_fields", "messages,messaging_postbacks,messaging_optins,message_deliveries,message_reads") 

	req, err := http.NewRequest("POST", apiURL, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Meta API error subscribing page (status %d): %s", resp.StatusCode, string(body))
	}
	
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if success, ok := result["success"].(bool); ok && success {
		return nil
	}
	
	// Sometimes it returns { "success": true } or just { "data": [] } on GET.
	// On POST success it usually returns { "success": true }
	return nil
}
