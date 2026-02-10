package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sidji-omnichannel/internal/models"
	"github.com/sidji-omnichannel/internal/services"
)

type AIHandler struct {
	aiService *services.AIService
}

func NewAIHandler(aiService *services.AIService) *AIHandler {
	return &AIHandler{aiService: aiService}
}

// @Summary      Get AI Config
// @Description  Returns the AI configuration for a specific channel
// @Tags         ai
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Channel ID"
// @Success      200  {object}  models.AIConfig
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/channels/{id}/ai [get]
func (h *AIHandler) GetConfig(c *gin.Context) {
	channelID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid channel ID"})
		return
	}

	config, err := h.aiService.GetConfig(channelID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get AI config"})
		return
	}

	c.JSON(http.StatusOK, config)
}

// @Summary      Update AI Config
// @Description  Updates the AI configuration for a specific channel
// @Tags         ai
// @Accept       json
// @Produce      json
// @Param        id    path      string  true  "Channel ID"
// @Param        input body      object  true  "AI Config Input"
// @Success      200   {object}  models.AIConfig
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /api/channels/{id}/ai [put]
func (h *AIHandler) UpdateConfig(c *gin.Context) {
	channelID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid channel ID"})
		return
	}

	var input struct {
		IsEnabled              bool            `json:"is_enabled"`
		Mode                   models.AIMode   `json:"mode"`
		Persona                string          `json:"persona"`
		HandoverTimeoutMinutes int             `json:"handover_timeout_minutes"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config := &models.AIConfig{
		ChannelID:              channelID,
		IsEnabled:              input.IsEnabled,
		Mode:                   input.Mode,
		Persona:                input.Persona,
		HandoverTimeoutMinutes: input.HandoverTimeoutMinutes,
	}

	if err := h.aiService.UpdateConfig(config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update AI config"})
		return
	}

	// Fetch fresh config to return
	updatedConfig, err := h.aiService.GetConfig(channelID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "Config updated", "data": config})
		return
	}

	c.JSON(http.StatusOK, updatedConfig)
}

// @Summary      List Knowledge
// @Description  Returns all items in the knowledge base for a specific channel
// @Tags         ai
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Channel ID"
// @Success      200  {array}   models.KnowledgeBaseItem
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/channels/{id}/ai/knowledge [get]
func (h *AIHandler) ListKnowledge(c *gin.Context) {
	channelID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid channel ID"})
		return
	}

	items, err := h.aiService.ListKnowledge(channelID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list knowledge items"})
		return
	}

	c.JSON(http.StatusOK, items)
}

// @Summary      Add Knowledge
// @Description  Adds a new item to the knowledge base for a specific channel
// @Tags         ai
// @Accept       json
// @Produce      json
// @Param        id    path      string  true  "Channel ID"
// @Param        input body      object  true  "Knowledge Input"
// @Success      201   {object}  map[string]string
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /api/channels/{id}/ai/knowledge [post]
func (h *AIHandler) AddKnowledge(c *gin.Context) {
	channelID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid channel ID"})
		return
	}

	var input struct {
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.aiService.AddKnowledge(channelID, input.Content); err != nil {
		log.Printf("Failed to add knowledge: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add knowledge. Ensure AI service is configured."})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Knowledge added successfully"})
}

// @Summary      Update Knowledge
// @Description  Updates an existing knowledge base item
// @Tags         ai
// @Accept       json
// @Produce      json
// @Param        id    path      string  true  "Channel ID (Ignored but in path)"
// @Param        kid   path      string  true  "Knowledge ID"
// @Param        input body      object  true  "Knowledge Input"
// @Success      200   {object}  map[string]string
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /api/channels/{id}/ai/knowledge/{kid} [put]
func (h *AIHandler) UpdateKnowledge(c *gin.Context) {
	knowledgeID, err := uuid.Parse(c.Param("kid")) // using 'kid' to distinguish from channel 'id'
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid knowledge ID"})
		return
	}

	var input struct {
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.aiService.UpdateKnowledge(knowledgeID, input.Content); err != nil {
		log.Printf("Failed to update knowledge: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update knowledge"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Knowledge updated successfully"})
}

// @Summary      Delete Knowledge
// @Description  Removes an item from the knowledge base
// @Tags         ai
// @Accept       json
// @Produce      json
// @Param        id    path      string  true  "Channel ID (Ignored but in path)"
// @Param        kid   path      string  true  "Knowledge ID"
// @Success      200   {object}  map[string]string
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /api/channels/{id}/ai/knowledge/{kid} [delete]
func (h *AIHandler) DeleteKnowledge(c *gin.Context) {
	knowledgeID, err := uuid.Parse(c.Param("kid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid knowledge ID"})
		return
	}

	if err := h.aiService.DeleteKnowledge(knowledgeID); err != nil {
		log.Printf("Failed to delete knowledge: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete knowledge"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Knowledge deleted successfully"})
}

// @Summary      Test AI Reply
// @Description  Generates an AI reply for a given query using the current config and knowledge base
// @Tags         ai
// @Accept       json
// @Produce      json
// @Param        id    path      string  true  "Channel ID"
// @Param        input body      object  true  "Query Input"
// @Success      200   {object}  map[string]any
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /api/channels/{id}/ai/test [post]
func (h *AIHandler) TestReply(c *gin.Context) {
	channelID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid channel ID"})
		return
	}

	var input struct {
		Query string `json:"query" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. Get embedding for query
	queryVector, err := h.aiService.EmbedText(channelID, input.Query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process query: " + err.Error()})
		return
	}

	// 2. Search knowledge base
	docs, err := h.aiService.SearchKnowledge(channelID, queryVector, 3)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search knowledge base"})
		return
	}

	// 3. Get config for persona
	config, err := h.aiService.GetConfig(channelID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get AI config"})
		return
	}

	// 4. Generate reply
	reply, err := h.aiService.GenerateReply(c.Request.Context(), config, input.Query, docs, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate reply: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"reply":   reply,
		"context": docs,
	})
}
