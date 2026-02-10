package twilio

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// TwilioAccount represents a Twilio account or subaccount response
type TwilioAccount struct {
	Sid          string    `json:"sid"`
	AuthToken    string    `json:"auth_token"`
	FriendlyName string    `json:"friendly_name"`
	Type         string    `json:"type"`
	Status       string    `json:"status"`
	DateCreated  time.Time `json:"date_created"`
}

// CreateSubaccount creates a new subaccount under the master account
func (p *TwilioProvider) CreateSubaccount(friendlyName string) (*TwilioAccount, error) {
	apiURL := "https://api.twilio.com/2010-04-01/Accounts.json"

	data := url.Values{}
	data.Set("FriendlyName", friendlyName)

	req, err := http.NewRequest("POST", apiURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(p.cfg.AccountSID, p.cfg.AuthToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("twilio api error (status %d): %s", resp.StatusCode, string(body))
	}

	var account TwilioAccount
	if err := json.Unmarshal(body, &account); err != nil {
		return nil, err
	}

	return &account, nil
}
