package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sidji-omnichannel/internal/domain/ports/service"
	"github.com/sidji-omnichannel/internal/middleware"
	"github.com/sidji-omnichannel/internal/models"
	"github.com/sidji-omnichannel/internal/services"
	"github.com/sidji-omnichannel/internal/subscription"
)

// ChannelHandler handles channel management endpoints
type ChannelHandler struct {
	channelService service.ChannelService
}

// NewChannelHandler creates a new channel handler
func NewChannelHandler(channelService service.ChannelService) *ChannelHandler {
	return &ChannelHandler{channelService: channelService}
}

// List returns all channels
// @Summary      List channels
// @Description  Get a list of connected channels (WhatsApp, Instagram)
// @Tags         channels
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string][]models.ChannelPublic
// @Failure      500  {object}  map[string]string
// @Router       /channels [get]
func (h *ChannelHandler) List(c *gin.Context) {
	orgID := middleware.GetOrganizationID(c)

	channels, err := h.channelService.List(orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list channels"})
		return
	}

	// Remove sensitive data
	publicChannels := make([]models.ChannelPublic, len(channels))
	for i, ch := range channels {
		publicChannels[i] = models.ChannelPublic{
			ID:       ch.ID,
			Type:     ch.Type,
			Provider: ch.Provider,
			Name:     ch.Name,
			Status:   ch.Status,
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": publicChannels})
}

// Get returns a channel
// @Summary      Get channel
// @Description  Get details of a specific channel
// @Tags         channels
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Channel ID"
// @Success      200  {object}  models.ChannelPublic
// @Failure      404  {object}  map[string]string
// @Router       /channels/{id} [get]
func (h *ChannelHandler) Get(c *gin.Context) {
	orgID := middleware.GetOrganizationID(c)
	channelID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid channel ID"})
		return
	}

	channel, err := h.channelService.GetByID(orgID, channelID)
	if err != nil {
		if err == services.ErrChannelNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Channel not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get channel"})
		return
	}

	// Return public data only (no access token)
	c.JSON(http.StatusOK, models.ChannelPublic{
		ID:       channel.ID,
		Type:     channel.Type,
		Provider: channel.Provider,
		Name:     channel.Name,
		Status:   channel.Status,
	})
}

// Create creates a new channel
// @Summary      Create channel
// @Description  Connect a new channel (WhatsApp/Instagram) (Admin only)
// @Tags         channels
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        input  body      models.CreateChannelInput  true  "Channel Config"
// @Success      201    {object}  models.ChannelPublic
// @Failure      400    {object}  map[string]string
// @Failure      403    {object}  map[string]string
// @Failure      500    {object}  map[string]string
// @Router       /channels [post]
func (h *ChannelHandler) Create(c *gin.Context) {
	orgID := middleware.GetOrganizationID(c)

	var input models.CreateChannelInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	channel, err := h.channelService.Create(orgID, &input)
	if err != nil {
		if err == subscription.ErrSubscriptionLimitReached {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create channel"})
		return
	}

	c.JSON(http.StatusCreated, models.ChannelPublic{
		ID:     channel.ID,
		Type:   channel.Type,
		Name:   channel.Name,
		Status: channel.Status,
	})
}

// Delete deletes a channel
// @Summary      Delete channel
// @Description  Remove a channel connection (Admin only)
// @Tags         channels
// @Security     BearerAuth
// @Param        id   path      string  true  "Channel ID"
// @Success      200  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /channels/{id} [delete]
func (h *ChannelHandler) Delete(c *gin.Context) {
	orgID := middleware.GetOrganizationID(c)
	channelID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid channel ID"})
		return
	}

	if err := h.channelService.Delete(orgID, channelID); err != nil {
		if err == services.ErrChannelNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Channel not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete channel"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Channel deleted successfully"})
}

// Activate activates a channel
// @Summary      Activate channel
// @Description  Manually activate a channel (Admin only)
// @Tags         channels
// @Security     BearerAuth
// @Param        id   path      string  true  "Channel ID"
// @Success      200  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /channels/{id}/activate [post]
func (h *ChannelHandler) Activate(c *gin.Context) {
	orgID := middleware.GetOrganizationID(c)
	channelID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid channel ID"})
		return
	}

	// Verify channel exists
	_, err = h.channelService.GetByID(orgID, channelID)
	if err != nil {
		if err == services.ErrChannelNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Channel not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get channel"})
		return
	}

	if err := h.channelService.UpdateStatus(channelID, models.ChannelStatusActive); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to activate channel"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Channel activated successfully"})
}

// ConnectInstagram endpoints
// @Summary      Connect Instagram
// @Description  Automatically connect an Instagram channel using Facebook Login token
// @Tags         channels
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        input  body      models.ConnectInstagramInput  true  "Access Token"
// @Success      200    {object}  models.ChannelPublic
// @Failure      400    {object}  map[string]string
// @Failure      403    {object}  map[string]string
// @Failure      500    {object}  map[string]string
// @Router       /channels/connect/instagram [post]
func (h *ChannelHandler) ConnectInstagram(c *gin.Context) {
	orgID := middleware.GetOrganizationID(c)

	var input models.ConnectInstagramInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	channel, err := h.channelService.ConnectInstagram(orgID, input.AccessToken, input.SelectedID)
	if err != nil {
		if err == subscription.ErrSubscriptionLimitReached {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.ChannelPublic{
		ID:     channel.ID,
		Type:   channel.Type,
		Name:   channel.Name,
		Status: channel.Status,
	})
}
// ConnectWhatsApp godoc
// @Summary      Connect WhatsApp
// @Description  Automatically discover and connect WhatsApp accounts using a user access token
// @Tags         channels
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        input  body      models.ConnectWhatsAppInput  true  "Access Token"
// @Success      200    {object}  models.ChannelPublic
// @Failure      400    {object}  map[string]string
// @Failure      403    {object}  map[string]string
// @Failure      500    {object}  map[string]string
// @Router       /channels/whatsapp/connect [post]
func (h *ChannelHandler) ConnectWhatsApp(c *gin.Context) {
	orgID := middleware.GetOrganizationID(c)

	var input models.ConnectWhatsAppInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	channel, err := h.channelService.ConnectWhatsApp(orgID, input.AccessToken, input.SelectedID)
	if err != nil {
		if err == subscription.ErrSubscriptionLimitReached {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.ChannelPublic{
		ID:     channel.ID,
		Type:   channel.Type,
		Name:   channel.Name,
		Status: channel.Status,
	})
}

// ConnectFacebook godoc
// @Summary      Connect Facebook
// @Description  Automatically discover and connect Facebook pages using a user access token
// @Tags         channels
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        input  body      models.ConnectFacebookInput  true  "Access Token"
// @Success      200    {object}  models.ChannelPublic
// @Failure      400    {object}  map[string]string
// @Failure      403    {object}  map[string]string
// @Failure      500    {object}  map[string]string
// @Router       /channels/facebook/connect [post]
func (h *ChannelHandler) ConnectFacebook(c *gin.Context) {
	orgID := middleware.GetOrganizationID(c)

	var input models.ConnectFacebookInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	channels, err := h.channelService.ConnectFacebook(orgID, input.AccessToken, input.SelectedID)
	if err != nil {
		if err == subscription.ErrSubscriptionLimitReached {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(channels) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No pages connected", "data": []interface{}{}})
		return
	}

	// Just return the first one as successful connection indicator for the automated flow
	// or return the list if needed.
	c.JSON(http.StatusOK, models.ChannelPublic{
		ID:     channels[0].ID,
		Type:   channels[0].Type,
		Name:   channels[0].Name,
		Status: channels[0].Status,
	})
}
// DiscoverMeta godoc
// @Summary      Discover Meta accounts
// @Description  Get a list of available Facebook Pages, Instagram accounts, and WhatsApp numbers for a token
// @Tags         channels
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        input  body      models.ConnectFacebookInput  true  "Access Token"
// @Success      200    {object}  map[string]interface{}
// @Failure      400    {object}  map[string]string
// @Router       /channels/discover/meta [post]
func (h *ChannelHandler) DiscoverMeta(c *gin.Context) {
	var input struct {
		AccessToken string `json:"access_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accounts, err := h.channelService.DiscoverMetaAccounts(input.AccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, accounts)
}
