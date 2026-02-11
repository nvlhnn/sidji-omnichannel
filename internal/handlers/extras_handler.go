package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sidji-omnichannel/internal/domain/ports/service"
	"github.com/sidji-omnichannel/internal/middleware"
	"github.com/sidji-omnichannel/internal/models"
	"github.com/sidji-omnichannel/internal/services"
)

// CannedResponseHandler handles canned response endpoints
type CannedResponseHandler struct {
	cannedService service.CannedResponseService
}

// NewCannedResponseHandler creates a new canned response handler
func NewCannedResponseHandler(cannedService service.CannedResponseService) *CannedResponseHandler {
	return &CannedResponseHandler{cannedService: cannedService}
}

// List returns all canned responses
// @Summary      List canned responses
// @Description  Get a list of canned responses
// @Tags         canned-responses
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string][]models.CannedResponse
// @Failure      500  {object}  map[string]string
// @Router       /canned-responses [get]
func (h *CannedResponseHandler) List(c *gin.Context) {
	orgID := middleware.GetOrganizationID(c)

	responses, err := h.cannedService.List(orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list canned responses"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": responses})
}

// Search searches canned responses
// @Summary      Search canned responses
// @Description  Search canned responses by shortcut or content
// @Tags         canned-responses
// @Produce      json
// @Security     BearerAuth
// @Param        q    query     string  true  "Search query"
// @Success      200  {object}  map[string][]models.CannedResponse
// @Failure      500  {object}  map[string]string
// @Router       /canned-responses/search [get]
func (h *CannedResponseHandler) Search(c *gin.Context) {
	orgID := middleware.GetOrganizationID(c)
	query := c.Query("q")

	if query == "" {
		responses, err := h.cannedService.List(orgID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list canned responses"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": responses})
		return
	}

	responses, err := h.cannedService.Search(orgID, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search canned responses"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": responses})
}

// Create creates a new canned response
// @Summary      Create canned response
// @Description  Add a new quick reply template
// @Tags         canned-responses
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        input  body      models.CreateCannedResponseInput  true  "Response Info"
// @Success      201    {object}  models.CannedResponse
// @Failure      400    {object}  map[string]string
// @Failure      500    {object}  map[string]string
// @Router       /canned-responses [post]
func (h *CannedResponseHandler) Create(c *gin.Context) {
	orgID := middleware.GetOrganizationID(c)
	userID := middleware.GetUserID(c)

	var input models.CreateCannedResponseInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.cannedService.Create(orgID, userID, &input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create canned response"})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// Update updates a canned response
// @Summary      Update canned response
// @Description  Update a canned response
// @Tags         canned-responses
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id     path      string                            true  "Response ID"
// @Param        input  body      models.CreateCannedResponseInput  true  "Response Info"
// @Success      200    {object}  models.CannedResponse
// @Failure      404    {object}  map[string]string
// @Router       /canned-responses/{id} [put]
func (h *CannedResponseHandler) Update(c *gin.Context) {
	orgID := middleware.GetOrganizationID(c)
	responseID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid response ID"})
		return
	}

	var input models.CreateCannedResponseInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.cannedService.Update(orgID, responseID, &input)
	if err != nil {
		if err == services.ErrCannedResponseNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Canned response not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update canned response"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// Delete deletes a canned response
// @Summary      Delete canned response
// @Description  Remove a canned response
// @Tags         canned-responses
// @Security     BearerAuth
// @Param        id   path      string  true  "Response ID"
// @Success      200  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /canned-responses/{id} [delete]
func (h *CannedResponseHandler) Delete(c *gin.Context) {
	orgID := middleware.GetOrganizationID(c)
	responseID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid response ID"})
		return
	}

	if err := h.cannedService.Delete(orgID, responseID); err != nil {
		if err == services.ErrCannedResponseNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Canned response not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete canned response"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Canned response deleted successfully"})
}

// ============================================
// Label Handler
// ============================================

// LabelHandler handles label endpoints
type LabelHandler struct {
	labelService service.LabelService
}

// NewLabelHandler creates a new label handler
func NewLabelHandler(labelService service.LabelService) *LabelHandler {
	return &LabelHandler{labelService: labelService}
}

// List returns all labels
// @Summary      List labels
// @Description  Get all conversation labels
// @Tags         labels
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string][]models.Label
// @Failure      500  {object}  map[string]string
// @Router       /labels [get]
func (h *LabelHandler) List(c *gin.Context) {
	orgID := middleware.GetOrganizationID(c)

	labels, err := h.labelService.List(orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list labels"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": labels})
}

// Create creates a new label
// @Summary      Create label
// @Description  Create a new conversation label
// @Tags         labels
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        input  body      models.CreateLabelInput  true  "Label Info"
// @Success      201    {object}  models.Label
// @Failure      400    {object}  map[string]string
// @Failure      500    {object}  map[string]string
// @Router       /labels [post]
func (h *LabelHandler) Create(c *gin.Context) {
	orgID := middleware.GetOrganizationID(c)

	var input models.CreateLabelInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	label, err := h.labelService.Create(orgID, &input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create label"})
		return
	}

	c.JSON(http.StatusCreated, label)
}

// Update updates a label
// @Summary      Update label
// @Description  Update a conversation label
// @Tags         labels
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id     path      string                   true  "Label ID"
// @Param        input  body      models.CreateLabelInput  true  "Label Info"
// @Success      200    {object}  models.Label
// @Failure      404    {object}  map[string]string
// @Router       /labels/{id} [put]
func (h *LabelHandler) Update(c *gin.Context) {
	orgID := middleware.GetOrganizationID(c)
	labelID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid label ID"})
		return
	}

	var input models.CreateLabelInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	label, err := h.labelService.Update(orgID, labelID, &input)
	if err != nil {
		if err == services.ErrLabelNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Label not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update label"})
		return
	}

	c.JSON(http.StatusOK, label)
}

// Delete deletes a label
// @Summary      Delete label
// @Description  Delete a conversation label
// @Tags         labels
// @Security     BearerAuth
// @Param        id   path      string  true  "Label ID"
// @Success      200  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /labels/{id} [delete]
func (h *LabelHandler) Delete(c *gin.Context) {
	orgID := middleware.GetOrganizationID(c)
	labelID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid label ID"})
		return
	}

	if err := h.labelService.Delete(orgID, labelID); err != nil {
		if err == services.ErrLabelNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Label not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete label"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Label deleted successfully"})
}
