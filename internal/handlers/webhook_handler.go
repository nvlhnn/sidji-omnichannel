package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sidji-omnichannel/internal/config"
	"github.com/sidji-omnichannel/internal/models"
	"github.com/sidji-omnichannel/internal/services"
	"github.com/sidji-omnichannel/internal/websocket"
)

// WebhookHandler handles incoming webhooks from Meta (WhatsApp & Instagram)
type WebhookHandler struct {
	cfg            *config.Config
	channelService *services.ChannelService
	contactService *services.ContactService
	// ...
	convService    *services.ConversationService
	messageService *services.MessageService
	mediaService   *services.MediaService
	aiService      *services.AIService    // Added AIService
	hub            *websocket.Hub
}

// NewWebhookHandler creates a new webhook handler
func NewWebhookHandler(
	cfg *config.Config,
	channelService *services.ChannelService,
	contactService *services.ContactService,
	convService *services.ConversationService,
	messageService *services.MessageService,
	mediaService *services.MediaService,
	aiService *services.AIService, // Added parameter
	hub *websocket.Hub,
) *WebhookHandler {
	return &WebhookHandler{
		cfg:            cfg,
		channelService: channelService,
		contactService: contactService,
		convService:    convService,
		messageService: messageService,
		mediaService:   mediaService,
		aiService:      aiService, // Initialize
		hub:            hub,
	}
}

// MetaWebhookPayload represents the incoming webhook payload from Meta
type MetaWebhookPayload struct {
	Object string `json:"object"`
	Entry  []struct {
		ID      string `json:"id"`
		Time    int64  `json:"time"`
		Changes []struct {
			Value json.RawMessage `json:"value"`
			Field string          `json:"field"`
		} `json:"changes"`
		Messaging []struct {
			Sender    struct{ ID string } `json:"sender"`
			Recipient struct{ ID string } `json:"recipient"`
			Timestamp int64               `json:"timestamp"`
			Message   *struct {
				MID         string `json:"mid"`
				Text        string `json:"text"`
				IsEcho      bool   `json:"is_echo"`
				Attachments []struct {
					Type    string `json:"type"`
					Payload struct {
						URL string `json:"url"`
					} `json:"payload"`
				} `json:"attachments"`
			} `json:"message"`
		} `json:"messaging"`
	} `json:"entry"`
}

// WhatsAppValue represents WhatsApp-specific webhook value
type WhatsAppValue struct {
	MessagingProduct string `json:"messaging_product"`
	Metadata         struct {
		DisplayPhoneNumber string `json:"display_phone_number"`
		PhoneNumberID      string `json:"phone_number_id"`
	} `json:"metadata"`
	Contacts []struct {
		Profile struct {
			Name string `json:"name"`
		} `json:"profile"`
		WaID string `json:"wa_id"`
	} `json:"contacts"`
	Messages []struct {
		ID        string `json:"id"`
		From      string `json:"from"`
		Timestamp string `json:"timestamp"`
		Type      string `json:"type"`
		Text      *struct {
			Body string `json:"body"`
		} `json:"text"`
		Image *struct {
			Caption  string `json:"caption"`
			MimeType string `json:"mime_type"`
			SHA256   string `json:"sha256"`
			ID       string `json:"id"`
		} `json:"image"`
		Document *struct {
			Caption  string `json:"caption"`
			Filename string `json:"filename"`
			MimeType string `json:"mime_type"`
			SHA256   string `json:"sha256"`
			ID       string `json:"id"`
		} `json:"document"`
		Audio *struct {
			MimeType string `json:"mime_type"`
			SHA256   string `json:"sha256"`
			ID       string `json:"id"`
		} `json:"audio"`
		Video *struct {
			Caption  string `json:"caption"`
			MimeType string `json:"mime_type"`
			SHA256   string `json:"sha256"`
			ID       string `json:"id"`
		} `json:"video"`
		Sticker *struct {
			MimeType string `json:"mime_type"`
			SHA256   string `json:"sha256"`
			ID       string `json:"id"`
		} `json:"sticker"`
		Context *struct {
			From string `json:"from"`
			ID   string `json:"id"`
		} `json:"context"`
	} `json:"messages"`
	Statuses []struct {
		ID          string `json:"id"`
		Status      string `json:"status"` // sent, delivered, read
		Timestamp   string `json:"timestamp"`
		RecipientID string `json:"recipient_id"`
	} `json:"statuses"`
}

// Verify handles webhook verification (GET request)
// GET /api/webhooks (also accepts /api/webhooks/meta for backward compatibility)
// @Summary      Verify Webhook
// @Description  Meta webhook verification endpoint (Hub Challenge)
// @Tags         webhook
// @Param        hub.mode          query     string  false  "Mode"
// @Param        hub.verify_token  query     string  false  "Verify Token"
// @Param        hub.challenge     query     string  false  "Challenge"
// @Success      200               {string}  string  "Challenge string"
// @Failure      403               {string}  string  "Forbidden"
// @Router       /api/webhooks [get]
func (h *WebhookHandler) Verify(c *gin.Context) {
	mode := c.Query("hub.mode")
	token := c.Query("hub.verify_token")
	challenge := c.Query("hub.challenge")

	if mode == "subscribe" && token == h.cfg.Meta.VerifyToken {
		fmt.Printf("[Webhook] Verification successful. Challenge: %s\n", challenge)
		c.String(http.StatusOK, challenge)
		return
	}

	fmt.Printf("[Webhook] Verification failed. Mode: %s, Token: %s\n", mode, token)

	c.JSON(http.StatusForbidden, gin.H{"error": "Verification failed"})
}

// Handle processes incoming webhooks (POST request)
// POST /api/webhooks (also accepts /api/webhooks/meta for backward compatibility)
// @Summary      Handle Webhook
// @Description  Receives and processes incoming messages from WhatsApp/Instagram
// @Tags         webhook
// @Accept       json
// @Produce      json
// @Param        payload  body      object  true  "Meta Webhook Payload"
// @Success      200      {string}  string  "OK"
// @Router       /api/webhooks [post]
func (h *WebhookHandler) Handle(c *gin.Context) {
	// Read body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	// Verify signature (optional but recommended)
	signature := c.GetHeader("X-Hub-Signature-256")

	// Print raw body for debugging
	fmt.Printf("[Webhook] Received raw body: %s\n", string(body))
	fmt.Printf("[Webhook] Signature Header: %s\n", signature)
	if signature != "" && h.cfg.Meta.AppSecret != "" {
		if !h.verifySignature(body, signature) {
			fmt.Printf("[Webhook] Signature verification failed\n")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature"})
			return
		}
		fmt.Printf("[Webhook] Signature verified successfully\n")
	} else if signature == "" {
		fmt.Printf("[Webhook] No signature header found\n")
	}

	// Parse payload
	var payload MetaWebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		fmt.Printf("[Webhook] Failed to parse payload: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}
	fmt.Printf("[Webhook] Parsed object=%s entries=%d\n", payload.Object, len(payload.Entry))

	// Process based on object type
	switch payload.Object {
	case "whatsapp_business_account":
		h.handleWhatsApp(&payload)
	case "instagram":
		h.handleInstagram(&payload)
	case "page":
		h.handleFacebook(&payload)
	default:
		fmt.Printf("Unhandled webhook object: %s\n", payload.Object)
	}

	// Always respond with 200 to acknowledge receipt
	c.JSON(http.StatusOK, gin.H{"status": "received"})
}

// HandleTwilio processes incoming webhooks from Twilio
func (h *WebhookHandler) HandleTwilio(c *gin.Context) {
	// Twilio sends data as application/x-www-form-urlencoded
	from := c.PostForm("From") // e.g. "whatsapp:+12345678"
	to := c.PostForm("To")     // e.g. "whatsapp:+98765432"
	body := c.PostForm("Body")
	msgSid := c.PostForm("SmsSid")
	numMedia := c.PostForm("NumMedia")

	fmt.Printf("[Twilio Webhook] From: %s, To: %s, Body: %s, Media: %s\n", from, to, body, numMedia)

	if from == "" || to == "" {
		c.JSON(http.StatusOK, gin.H{"status": "ignored"})
		return
	}

	// Remove "whatsapp:" prefix
	cleanFrom := strings.TrimPrefix(from, "whatsapp:")
	cleanTo := strings.TrimPrefix(to, "whatsapp:")

	// Find Channel by To number (the "From" in our provider config, but "To" in Twilio incoming webhook)
	channel, err := h.channelService.GetChannelByPhoneNumberID(cleanTo)
	if err != nil {
		fmt.Printf("[Twilio Webhook] Channel for %s not found: %v\n", cleanTo, err)
		c.JSON(http.StatusOK, gin.H{"status": "ignored"})
		return
	}

	// 1. Find or create contact
	contactName := cleanFrom
	contact, err := h.contactService.FindOrCreateByWhatsAppID(channel.OrganizationID, cleanFrom, contactName)
	if err != nil {
		fmt.Printf("[Twilio Webhook] Error creating contact: %v\n", err)
		c.JSON(http.StatusOK, gin.H{"status": "ignored"})
		return
	}

	// 2. Find or create conversation
	conv, err := h.convService.FindOrCreate(channel.OrganizationID, channel.ID, contact.ID)
	if err != nil {
		fmt.Printf("[Twilio Webhook] Error creating conversation: %v\n", err)
		c.JSON(http.StatusOK, gin.H{"status": "ignored"})
		return
	}

	// 3. Create message record
	message := &models.Message{
		ConversationID: conv.ID,
		SenderType:     models.SenderContact,
		SenderID:       contact.ID,
		Content:        body,
		ExternalID:     msgSid,
		MessageType:    models.MessageTypeText,
		Status:         models.MessageStatusDelivered,
	}

	// Handle media if present
	if numMedia != "" && numMedia != "0" {
		mediaURL := c.PostForm("MediaUrl0")
		contentType := c.PostForm("MediaContentType0")
		if mediaURL != "" {
			message.MediaURL = mediaURL
			message.MediaMimeType = contentType
			if strings.HasPrefix(contentType, "image/") {
				message.MessageType = models.MessageTypeImage
			} else if strings.HasPrefix(contentType, "video/") {
				message.MessageType = models.MessageTypeVideo
			} else if strings.HasPrefix(contentType, "audio/") {
				message.MessageType = models.MessageTypeAudio
			} else {
				message.MessageType = models.MessageTypeDocument
			}
		}
	}

	if err := h.messageService.Create(message); err != nil {
		fmt.Printf("[Twilio Webhook] Error saving message: %v\n", err)
	}

	// 4. Broadcast via WebSocket
	if h.hub != nil {
		h.hub.Broadcast(channel.OrganizationID, "message:new", message)
		updatedConv, err := h.convService.GetByID(channel.OrganizationID, conv.ID)
		if err == nil {
			updatedConv.LastMessage = message
			unreadCount, _ := h.messageService.CountUnread(conv.ID)
			updatedConv.UnreadCount = unreadCount
			h.hub.Broadcast(channel.OrganizationID, "conversation:update", updatedConv)
		}
	}

	// 5. AI Auto-reply
	if message.MessageType == models.MessageTypeText && message.Content != "" {
		h.processAIMessage(channel.ID, contact.ID, conv.ID, message.Content)
	}

	c.JSON(http.StatusOK, gin.H{"status": "received"})
}

// handleWhatsApp processes WhatsApp webhook events
func (h *WebhookHandler) handleWhatsApp(payload *MetaWebhookPayload) {
	for _, entry := range payload.Entry {
		for _, change := range entry.Changes {
			if change.Field != "messages" {
				continue
			}

			var value WhatsAppValue
			if err := json.Unmarshal(change.Value, &value); err != nil {
				continue
			}

			// Ignore if no messages or statuses
			if len(value.Messages) == 0 && len(value.Statuses) == 0 {
				continue
			}

			// Find Channel
			phoneNumberID := value.Metadata.PhoneNumberID
			channel, err := h.channelService.GetChannelByPhoneNumberID(phoneNumberID)
			if err != nil {
				fmt.Printf("Error finding channel for phone_number_id %s: %v\n", phoneNumberID, err)
				continue
			}

			// Process incoming messages
			for _, msg := range value.Messages {
				// 1. Find or create contact by wa_id
				var contactName string
				// Try to find name in contacts array from payload
				for _, c := range value.Contacts {
					if c.WaID == msg.From {
						contactName = c.Profile.Name
						break
					}
				}
				if contactName == "" {
					contactName = msg.From // Fallback to phone number
				}

				contact, err := h.contactService.FindOrCreateByWhatsAppID(channel.OrganizationID, msg.From, contactName)
				if err != nil {
					fmt.Printf("Error finding/creating contact %s: %v\n", msg.From, err)
					continue
				}

				// 2. Find or create conversation
				conv, err := h.convService.FindOrCreate(channel.OrganizationID, channel.ID, contact.ID)
				if err != nil {
					fmt.Printf("Error finding/creating conversation: %v\n", err)
					continue
				}

				// 3. Create message record
				message := &models.Message{
					ConversationID: conv.ID,
					SenderType:     models.SenderContact,
					SenderID:       contact.ID,
					ExternalID:     msg.ID,
					Status:         models.MessageStatusDelivered, // Webhook implies delivered to us
				}

				// Determine message type and content
				switch msg.Type {
				case "text":
					message.MessageType = models.MessageTypeText
					if msg.Text != nil {
						message.Content = msg.Text.Body
					}
				case "image":
					message.MessageType = models.MessageTypeImage
					if msg.Image != nil {
						message.Content = msg.Image.Caption
						message.MediaMimeType = msg.Image.MimeType

						// Download Media
						mediaURL, err := h.mediaService.DownloadMetaMedia(msg.Image.ID, channel.AccessToken)
						if err != nil {
							fmt.Printf("Failed to download image %s: %v\n", msg.Image.ID, err)
							message.Content += " [Download Failed]"
						} else {
							message.MediaURL = mediaURL
						}
					}
				case "document":
					message.MessageType = models.MessageTypeDocument
					if msg.Document != nil {
						message.Content = msg.Document.Caption
						message.MediaFileName = msg.Document.Filename
						message.MediaMimeType = msg.Document.MimeType

						// Download Media
						mediaURL, err := h.mediaService.DownloadMetaMedia(msg.Document.ID, channel.AccessToken)
						if err != nil {
							fmt.Printf("Failed to download document %s: %v\n", msg.Document.ID, err)
							message.Content += " [Download Failed]"
						} else {
							message.MediaURL = mediaURL
						}
					}
				case "audio":
					message.MessageType = models.MessageTypeAudio
					if msg.Audio != nil {
						message.MediaMimeType = msg.Audio.MimeType
						// Download Media
						mediaURL, err := h.mediaService.DownloadMetaMedia(msg.Audio.ID, channel.AccessToken)
						if err != nil {
							fmt.Printf("Failed to download audio %s: %v\n", msg.Audio.ID, err)
							message.Content = "[Audio Download Failed]"
						} else {
							message.MediaURL = mediaURL
							message.Content = "Audio Message"
						}
					}
				case "video":
					message.MessageType = models.MessageTypeVideo
					if msg.Video != nil {
						message.Content = msg.Video.Caption
						message.MediaMimeType = msg.Video.MimeType
						// Download Media
						mediaURL, err := h.mediaService.DownloadMetaMedia(msg.Video.ID, channel.AccessToken)
						if err != nil {
							fmt.Printf("Failed to download video %s: %v\n", msg.Video.ID, err)
							message.Content += " [Video Download Failed]"
						} else {
							message.MediaURL = mediaURL
						}
					}
				default:
					message.MessageType = models.MessageTypeText
					message.Content = "[Unsupported message type: " + msg.Type + "]"
				}

				if err := h.messageService.Create(message); err != nil {
					fmt.Printf("Error creating message: %v\n", err)
					continue
				}

				// 4. Broadcast via WebSocket to connected agents
				if h.hub != nil {
					h.hub.Broadcast(channel.OrganizationID, "message:new", message)

					// Also broadcast conversation update to bring it to top/add to list
					updatedConv, err := h.convService.GetByID(channel.OrganizationID, conv.ID)
					if err == nil {
						updatedConv.LastMessage = message 
						h.hub.Broadcast(channel.OrganizationID, "conversation:update", updatedConv)
					}
				}

				// 5. Trigger AI Auto-Reply (Async)
				if message.MessageType == models.MessageTypeText && message.Content != "" {
					h.processAIMessage(channel.ID, contact.ID, conv.ID, message.Content)
				}
			}

			// Process status updates
			for _, status := range value.Statuses {
				statusEnum := models.MessageStatus(status.Status)
				if statusEnum == "sent" || statusEnum == "delivered" || statusEnum == "read" {
					if err := h.messageService.UpdateStatus(status.ID, statusEnum); err != nil {
						fmt.Printf("Error updating message status %s: %v\n", status.ID, err)
					}
					// Also broadcast status update
					if h.hub != nil {
						h.hub.Broadcast(channel.OrganizationID, "message:status", map[string]interface{}{
							"id":     status.ID,
							"status": statusEnum,
						})
					}
				}
			}
		}
	}
}

// handleInstagram processes Instagram webhook events
func (h *WebhookHandler) handleInstagram(payload *MetaWebhookPayload) {
	for _, entry := range payload.Entry {
		for _, messaging := range entry.Messaging {
			if messaging.Message == nil || messaging.Message.IsEcho {
				continue
			}

			// Find Channel by Recipient ID (Incoming Message)
			igUserID := messaging.Recipient.ID
			channel, err := h.channelService.GetChannelByIGUserID(igUserID)
			
			if err != nil {
				// If not found, check if it's an echo (Sender is the Channel)
				senderChannel, err2 := h.channelService.GetChannelByIGUserID(messaging.Sender.ID)
				if err2 == nil {
					// It's an echo! Sender is our channel.
					fmt.Printf("[Webhook] Detected echo message from channel %s (ID: %s). Ignoring to prevent duplication.\n", senderChannel.Name, messaging.Sender.ID)
					continue
				}

				// Neither recipient nor sender is a known channel
				fmt.Printf("Error finding channel for Recipient=%s (err=%v) OR Sender=%s (err=%v)\n", igUserID, err, messaging.Sender.ID, err2)
				continue
			}

			// 1. Find or create contact through valid channel
			senderID := messaging.Sender.ID

			// Fetch user profile from Instagram Graph API
			name, username, avatarURL, err := h.channelService.GetInstagramUserProfile(senderID, channel.AccessToken)

			contactName := "Instagram User " + senderID
			if err == nil {
				if name != "" {
					contactName = name
				} else if username != "" {
					contactName = username
				}
				// Append username if available and different from name
				if username != "" && name != "" && name != username {
					contactName += " (" + username + ")"
				}
			} else {
				fmt.Printf("Failed to fetch IG profile for %s: %v\n", senderID, err)
			}

			contact, err := h.contactService.FindOrCreateByInstagramID(channel.OrganizationID, senderID, contactName, avatarURL)
			if err != nil {
				fmt.Printf("Error finding/creating contact %s: %v\n", senderID, err)
				continue
			}

			// 2. Find or create conversation
			conv, err := h.convService.FindOrCreate(channel.OrganizationID, channel.ID, contact.ID)
			if err != nil {
				fmt.Printf("Error finding/creating conversation: %v\n", err)
				continue
			}

			// 3. Create message record
			msg := messaging.Message
			message := &models.Message{
				ConversationID: conv.ID,
				SenderType:     models.SenderContact,
				SenderID:       contact.ID,
				ExternalID:     msg.MID,
				Status:         models.MessageStatusDelivered,
				MessageType:    models.MessageTypeText,
			}

			if len(msg.Attachments) > 0 {
				// Handle first attachment
				att := msg.Attachments[0]
				if att.Type == "image" {
					message.MessageType = models.MessageTypeImage
					message.MediaURL = att.Payload.URL
				} else {
					message.Content = "[Attachment: " + att.Type + "]"
				}
			} else {
				message.Content = msg.Text
			}

			if err := h.messageService.Create(message); err != nil {
				fmt.Printf("Error creating message: %v\n", err)
				continue
			}

				// 4. Broadcast via WebSocket to connected agents
				if h.hub != nil {
					h.hub.Broadcast(channel.OrganizationID, "message:new", message)

					// Also broadcast conversation update to bring it to top/add to list
					// We need to reload the conversation to get the latest last_message_at etc
					updatedConv, err := h.convService.GetByID(channel.OrganizationID, conv.ID)
					if err == nil {
						// Manually attach the message we just created so the frontend has the preview text!
						// GetByID does NOT populate LastMessage content by default.
						updatedConv.LastMessage = message
						
						// Calculate unread count to send to client
						unreadCount, _ := h.messageService.CountUnread(conv.ID)
						updatedConv.UnreadCount = unreadCount

						h.hub.Broadcast(channel.OrganizationID, "conversation:update", updatedConv)
					}
				}

				// 5. Trigger AI Auto-Reply (Async)
				if message.MessageType == models.MessageTypeText && message.Content != "" {
				h.processAIMessage(channel.ID, contact.ID, conv.ID, message.Content)
			}
		}
	}
}

// handleFacebook processes Facebook Messenger webhook events
func (h *WebhookHandler) handleFacebook(payload *MetaWebhookPayload) {
	for _, entry := range payload.Entry {
		for _, messaging := range entry.Messaging {
			if messaging.Message == nil || messaging.Message.IsEcho {
				continue
			}

			// Find Channel by Recipient ID (Incoming Message)
			pageID := messaging.Recipient.ID
			channel, err := h.channelService.GetChannelByFacebookPageID(pageID)
			
			if err != nil {
				// If not found, check if it's an echo (Sender is the Page)
				senderChannel, err2 := h.channelService.GetChannelByFacebookPageID(messaging.Sender.ID)
				if err2 == nil {
					// It's an echo! Sender is our channel.
					fmt.Printf("[Webhook] Detected echo message from Facebook Page %s (ID: %s). Ignoring.\n", senderChannel.Name, messaging.Sender.ID)
					continue
				}

				fmt.Printf("Error finding channel for Facebook Recipient=%s (err=%v) OR Sender=%s (err=%v)\n", pageID, err, messaging.Sender.ID, err2)
				continue
			}

			// 1. Find or create contact
			senderID := messaging.Sender.ID

			// Fetch user profile from Facebook Graph API (Messenger)
			name, avatarURL, err := h.channelService.GetFacebookUserProfile(senderID, channel.AccessToken)

			contactName := "Facebook User " + senderID
			if err == nil && name != "" {
				contactName = name
			} else if err != nil {
				fmt.Printf("Failed to fetch Facebook profile for %s: %v\n", senderID, err)
			}

			contact, err := h.contactService.FindOrCreateByFacebookID(channel.OrganizationID, senderID, contactName, avatarURL)
			if err != nil {
				fmt.Printf("Error finding/creating contact %s: %v\n", senderID, err)
				continue
			}

			// 2. Find or create conversation
			conv, err := h.convService.FindOrCreate(channel.OrganizationID, channel.ID, contact.ID)
			if err != nil {
				fmt.Printf("Error finding/creating conversation: %v\n", err)
				continue
			}

			// 3. Create message record
			msg := messaging.Message
			message := &models.Message{
				ConversationID: conv.ID,
				SenderType:     models.SenderContact,
				SenderID:       contact.ID,
				ExternalID:     msg.MID,
				Status:         models.MessageStatusDelivered,
				MessageType:    models.MessageTypeText,
			}

			if len(msg.Attachments) > 0 {
				att := msg.Attachments[0]
				if att.Type == "image" {
					message.MessageType = models.MessageTypeImage
					message.MediaURL = att.Payload.URL
				} else {
					message.Content = "[Attachment: " + att.Type + "]"
				}
			} else {
				message.Content = msg.Text
			}

			if err := h.messageService.Create(message); err != nil {
				fmt.Printf("Error creating message: %v\n", err)
				continue
			}

			// 4. Broadcast via WebSocket
			if h.hub != nil {
				h.hub.Broadcast(channel.OrganizationID, "message:new", message)

				updatedConv, err := h.convService.GetByID(channel.OrganizationID, conv.ID)
				if err == nil {
					updatedConv.LastMessage = message
					unreadCount, _ := h.messageService.CountUnread(conv.ID)
					updatedConv.UnreadCount = unreadCount

					h.hub.Broadcast(channel.OrganizationID, "conversation:update", updatedConv)
				}
			}

			// 5. Trigger AI Auto-Reply
			if message.MessageType == models.MessageTypeText && message.Content != "" {
				h.processAIMessage(channel.ID, contact.ID, conv.ID, message.Content)
			}
		}
	}
}

// verifySignature verifies the X-Hub-Signature-256 header
func (h *WebhookHandler) verifySignature(body []byte, signature string) bool {
	mac := hmac.New(sha256.New, []byte(h.cfg.Meta.AppSecret))
	mac.Write(body)
	expectedSignature := "sha256=" + hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}
