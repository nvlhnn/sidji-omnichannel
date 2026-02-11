package handlers

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/sidji-omnichannel/internal/models"
)

// processAIMessage handles the AI auto-reply logic
func (h *WebhookHandler) processAIMessage(channelID uuid.UUID, contactID uuid.UUID, conversationID uuid.UUID, userMessage string) {
	go func() {
		// Recover from panic to not crash server
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered from AI panic: %v", r)
			}
		}()

		// Fetch fresh channel config 
		channel, err := h.channelService.Get(channelID) // Use Get
		if err != nil {
			log.Printf("Failed to finding channel for AI: %v", err)
			return
		}

		// Fetch conversation to get LastHumanReplyAt
		conv, err := h.convService.GetByID(channel.OrganizationID, conversationID)
		if err != nil {
			log.Printf("Failed to finding conversation for AI: %v", err)
			return
		}

		// Check AI Logic
		shouldAutoReply, shouldDraft, config, err := h.aiService.DetermineAction(channelID, conv.LastHumanReplyAt)
		if err != nil {
			log.Printf("AI Check Error: %v", err)
			return
		}

		log.Printf("[AI processAIMessage] shouldAutoReply=%v shouldDraft=%v mode=%v", shouldAutoReply, shouldDraft, config.Mode)

		if !shouldAutoReply && !shouldDraft {
			return
		}

		// Check Credits
		if err := h.aiService.CheckCredits(channelID); err != nil {
			log.Printf("AI Check Credit Error: %v", err)
			return
		}

		// Simulate "typing..." delay only for auto-reply
		if shouldAutoReply {
			time.Sleep(2 * time.Second)
		}

		// Get Contact External ID needed for sending
		contact, err := h.contactService.GetByID(channel.OrganizationID, contactID)
		if err != nil {
			log.Printf("AI Contact Fetch Error: %v", err)
			return
		}

		// Generate Reply
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Embed User Message
		queryVector, err := h.aiService.EmbedText(channelID, userMessage)
		if err != nil {
			log.Printf("AI Embed Error: %v", err)
			return
		}

		// Search Knowledge Base
		docs, err := h.aiService.SearchKnowledge(channelID, queryVector, 3)
		if err != nil {
			log.Printf("AI Search Error: %v", err)
			docs = []string{} // Fallback to no context
		}

		// Fetch recent history (last 10 messages)
		historyRes, err := h.messageService.List(conversationID, &models.MessageFilter{Limit: 10})
		var history []models.Message
		if err == nil && historyRes != nil {
			for _, m := range historyRes.Messages {
				history = append(history, *m)
			}
		}

		replyContent, err := h.aiService.GenerateReply(ctx, config, userMessage, docs, history)
		if err != nil {
			log.Printf("AI Generate Error: %v", err)
			return
		}

		// Deduct Credit
		if err := h.aiService.DeductCredit(channelID); err != nil {
			log.Printf("AI Deduct Credit Error: %v", err)
			// Continue anyway since reply is generated
		}

		// --- Handle Draft Mode ---
		if shouldDraft {
			if h.hub != nil {
				h.hub.Broadcast(channel.OrganizationID, "conversation:ai_draft", map[string]interface{}{
					"conversation_id": conversationID,
					"content":         replyContent,
					"sender_id":       channel.ID,
					"sender_type":     models.SenderAI,
					"created_at":      time.Now(),
				})
			}
			return
		}

		// --- Handle Auto-Reply Mode ---
		// Create Message
		message := &models.Message{
			ConversationID: conversationID,
			SenderType:     models.SenderAI,
			SenderID:       channel.ID,
			Content:        replyContent,
			MessageType:    models.MessageTypeText,
			Status:         models.MessageStatusPending,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
		
		if err := h.messageService.Create(message); err != nil {
			log.Printf("AI Save Message Error: %v", err)
			return
		}

		// Broadcast to UI immediately (as pending)
		if h.hub != nil {
			// Populate Sender for UI Broadcast
			message.Sender = &models.UserPublic{
				ID:     channel.ID,
				Name:   "AI Assistant",
				Email:  "ai@sidji.internal",
				Role:   models.RoleAgent,
				Status: models.StatusOnline,
			}
			h.hub.Broadcast(channel.OrganizationID, "message:new", message)
			
			// Update conversation list preview
			conv.LastMessage = message
			h.hub.Broadcast(channel.OrganizationID, "conversation:update", conv)
		}

		// Send to Meta
		_, metaErr := h.channelService.SendMessage(channel, contact, message)

		newStatus := models.MessageStatusSent
		if metaErr != nil {
			log.Printf("AI Send Meta Error: %v", metaErr)
			newStatus = models.MessageStatusFailed
		}

		// Update DB
		h.messageService.UpdateStatusByID(message.ID, newStatus)
		message.Status = newStatus

		// Broadcast Status Update
		if h.hub != nil {
			h.hub.Broadcast(channel.OrganizationID, "message:status", map[string]interface{}{
				"id":     message.ID,
				"status": newStatus,
			})
		}
	}()
}
