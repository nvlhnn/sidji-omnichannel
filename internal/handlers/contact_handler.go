package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sidji-omnichannel/internal/domain/ports/service"
	"github.com/sidji-omnichannel/internal/middleware"
	"github.com/sidji-omnichannel/internal/models"
	"github.com/sidji-omnichannel/internal/services"
)

// ContactHandler handles contact management endpoints
type ContactHandler struct {
	contactService service.ContactService
}

// NewContactHandler creates a new contact handler
func NewContactHandler(contactService service.ContactService) *ContactHandler {
	return &ContactHandler{contactService: contactService}
}

// List returns all contacts with pagination and search
// @Summary      List contacts
// @Description  Get a list of customer contacts
// @Tags         contacts
// @Produce      json
// @Security     BearerAuth
// @Param        search  query     string  false  "Search by name, email, phone"
// @Param        page    query     int     false  "Page number (default 1)"
// @Param        limit   query     int     false  "Page size (default 20)"
// @Success      200     {object}  map[string]interface{}
// @Failure      500     {object}  map[string]string
// @Router       /contacts [get]
func (h *ContactHandler) List(c *gin.Context) {
	orgID := middleware.GetOrganizationID(c)

	page := 1
	limit := 20
	search := c.Query("search")

	if p := c.Query("page"); p != "" {
		if val, err := strconv.Atoi(p); err == nil {
			page = val
		}
	}
	if l := c.Query("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil {
			limit = val
		}
	}

	contacts, total, err := h.contactService.List(orgID, page, limit, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list contacts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  contacts,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// Get returns a contact
// @Summary      Get contact
// @Description  Get details of a specific contact
// @Tags         contacts
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Contact ID"
// @Success      200  {object}  models.Contact
// @Failure      404  {object}  map[string]string
// @Router       /contacts/{id} [get]
func (h *ContactHandler) Get(c *gin.Context) {
	orgID := middleware.GetOrganizationID(c)
	contactID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contact ID"})
		return
	}

	contact, err := h.contactService.GetByID(orgID, contactID)
	if err != nil {
		if err == services.ErrContactNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Contact not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get contact"})
		return
	}

	c.JSON(http.StatusOK, contact)
}

// Create creates a new contact
// @Summary      Create contact
// @Description  Add a new customer contact
// @Tags         contacts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        input  body      models.CreateContactInput  true  "Contact Info"
// @Success      201    {object}  models.Contact
// @Failure      400    {object}  map[string]string
// @Failure      500    {object}  map[string]string
// @Router       /contacts [post]
func (h *ContactHandler) Create(c *gin.Context) {
	orgID := middleware.GetOrganizationID(c)

	var input models.CreateContactInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	contact, err := h.contactService.Create(orgID, &input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create contact"})
		return
	}

	c.JSON(http.StatusCreated, contact)
}

// Update updates a contact
// @Summary      Update contact
// @Description  Update contact details
// @Tags         contacts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id     path      string                     true  "Contact ID"
// @Param        input  body      models.UpdateContactInput  true  "Contact Info"
// @Success      200    {object}  models.Contact
// @Failure      404    {object}  map[string]string
// @Router       /contacts/{id} [patch]
func (h *ContactHandler) Update(c *gin.Context) {
	orgID := middleware.GetOrganizationID(c)
	contactID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contact ID"})
		return
	}

	var input models.UpdateContactInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	contact, err := h.contactService.Update(orgID, contactID, &input)
	if err != nil {
		if err == services.ErrContactNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Contact not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update contact"})
		return
	}

	c.JSON(http.StatusOK, contact)
}

// Delete deletes a contact
// @Summary      Delete contact
// @Description  Remove a contact
// @Tags         contacts
// @Security     BearerAuth
// @Param        id   path      string  true  "Contact ID"
// @Success      200  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /contacts/{id} [delete]
func (h *ContactHandler) Delete(c *gin.Context) {
	orgID := middleware.GetOrganizationID(c)
	contactID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contact ID"})
		return
	}

	if err := h.contactService.Delete(orgID, contactID); err != nil {
		if err == services.ErrContactNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Contact not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete contact"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Contact deleted successfully"})
}

// GetConversations returns all conversations for a contact
// @Summary      Get contact conversations
// @Description  Get history of conversations with a contact
// @Tags         contacts
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Contact ID"
// @Success      200  {object}  map[string][]models.Conversation
// @Failure      404  {object}  map[string]string
// @Router       /contacts/{id}/conversations [get]
func (h *ContactHandler) GetConversations(c *gin.Context) {
	orgID := middleware.GetOrganizationID(c)
	contactID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contact ID"})
		return
	}

	conversations, err := h.contactService.GetConversations(orgID, contactID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get conversations"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": conversations})
}
