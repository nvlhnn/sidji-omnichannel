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

// TeamHandler handles team management endpoints
type TeamHandler struct {
	teamService service.TeamService
}

// NewTeamHandler creates a new team handler
func NewTeamHandler(teamService service.TeamService) *TeamHandler {
	return &TeamHandler{teamService: teamService}
}

// ListMembers returns all team members
// @Summary      List team members
// @Description  Get all users in the organization
// @Tags         team
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string][]models.UserPublic
// @Failure      500  {object}  map[string]string
// @Router       /team/members [get]
func (h *TeamHandler) ListMembers(c *gin.Context) {
	orgID := middleware.GetOrganizationID(c)

	members, err := h.teamService.ListMembers(orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list team members"})
		return
	}

	// Convert to public representation
	publicMembers := make([]*models.UserPublic, len(members))
	for i, m := range members {
		publicMembers[i] = m.ToPublic()
	}

	c.JSON(http.StatusOK, gin.H{"data": publicMembers})
}

// GetMember returns a team member
// @Summary      Get team member
// @Description  Get details of a team member
// @Tags         team
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Member ID"
// @Success      200  {object}  models.UserPublic
// @Failure      404  {object}  map[string]string
// @Router       /team/members/{id} [get]
func (h *TeamHandler) GetMember(c *gin.Context) {
	orgID := middleware.GetOrganizationID(c)
	memberID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid member ID"})
		return
	}

	member, err := h.teamService.GetMember(orgID, memberID)
	if err != nil {
		if err == services.ErrTeamMemberNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Team member not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get team member"})
		return
	}

	c.JSON(http.StatusOK, member.ToPublic())
}

// InviteMember invites a new team member
// @Summary      Invite team member
// @Description  Create a new user/agent (Admin only)
// @Tags         team
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        input  body      models.InviteUserInput  true  "Invite Info"
// @Success      201    {object}  models.UserPublic
// @Failure      400    {object}  map[string]string
// @Failure      403    {object}  map[string]string
// @Failure      409    {object}  map[string]string
// @Router       /team/members [post]
func (h *TeamHandler) InviteMember(c *gin.Context) {
	orgID := middleware.GetOrganizationID(c)

	var input models.InviteUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	member, err := h.teamService.InviteMember(orgID, &input)
	if err != nil {
		if err == services.ErrUserExists {
			c.JSON(http.StatusConflict, gin.H{"error": "User with this email already exists"})
			return
		}
		if err == subscription.ErrSubscriptionLimitReached {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to invite team member"})
		return
	}

	c.JSON(http.StatusCreated, member.ToPublic())
}

// UpdateMember updates a team member
// @Summary      Update team member
// @Description  Update details of a team member (Supervisor/Admin)
// @Tags         team
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id     path      string                  true  "Member ID"
// @Param        input  body      models.UpdateUserInput  true  "Update Info"
// @Success      200    {object}  models.UserPublic
// @Failure      404    {object}  map[string]string
// @Router       /team/members/{id} [patch]
func (h *TeamHandler) UpdateMember(c *gin.Context) {
	orgID := middleware.GetOrganizationID(c)
	memberID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid member ID"})
		return
	}

	// Permission check: Admin/Supervisor can update anyone, agents can only update themselves
	userID := middleware.GetUserID(c)
	role := middleware.GetUserRole(c)
	if memberID != userID && role != "admin" && role != "supervisor" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions to update other members"})
		return
	}

	var input models.UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	member, err := h.teamService.UpdateMember(orgID, memberID, &input)
	if err != nil {
		if err == services.ErrTeamMemberNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Team member not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update team member"})
		return
	}

	c.JSON(http.StatusOK, member.ToPublic())
}

// UpdateMemberRole updates a team member's role
// @Summary      Update member role
// @Description  Promote or demote a team member (Admin only)
// @Tags         team
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id     path      string  true  "Member ID"
// @Param        role   body      string  true  "New Role (admin, supervisor, agent)"
// @Success      200    {object}  map[string]string
// @Failure      404    {object}  map[string]string
// @Router       /team/members/{id}/role [patch]
func (h *TeamHandler) UpdateMemberRole(c *gin.Context) {
	orgID := middleware.GetOrganizationID(c)
	memberID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid member ID"})
		return
	}

	var input struct {
		Role models.UserRole `json:"role" binding:"required,oneof=admin supervisor agent"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.teamService.UpdateMemberRole(orgID, memberID, input.Role); err != nil {
		if err == services.ErrTeamMemberNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Team member not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role updated successfully"})
}

// RemoveMember removes a team member
// @Summary      Remove team member
// @Description  Delete a team member (Admin only)
// @Tags         team
// @Security     BearerAuth
// @Param        id   path      string  true  "Member ID"
// @Success      200  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /team/members/{id} [delete]
func (h *TeamHandler) RemoveMember(c *gin.Context) {
	orgID := middleware.GetOrganizationID(c)
	memberID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid member ID"})
		return
	}

	if err := h.teamService.RemoveMember(orgID, memberID); err != nil {
		if err == services.ErrTeamMemberNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Team member not found"})
			return
		}
		if err.Error() == "cannot remove the last admin" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot remove the last admin"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove team member"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Team member removed successfully"})
}

// GetOrganization returns organization details
// @Summary      Get organization
// @Description  Get organization details
// @Tags         team
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  models.Organization
// @Failure      500  {object}  map[string]string
// @Router       /team/organization [get]
func (h *TeamHandler) GetOrganization(c *gin.Context) {
	orgID := middleware.GetOrganizationID(c)

	org, err := h.teamService.GetOrganization(orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get organization"})
		return
	}

	c.JSON(http.StatusOK, org)
}

// UpdateOrganization updates organization details
// @Summary      Update organization
// @Description  Update organization name or plan (Admin only)
// @Tags         team
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        input  body      object  true  "Update Info"
// @Success      200    {object}  models.Organization
// @Failure      500    {object}  map[string]string
// @Router       /team/organization [patch]
func (h *TeamHandler) UpdateOrganization(c *gin.Context) {
	orgID := middleware.GetOrganizationID(c)

	var input struct {
		Name string `json:"name" binding:"omitempty,min=2,max=100"`
		Plan string `json:"plan" binding:"omitempty,oneof=starter growth scale enterprise"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	org, err := h.teamService.UpdateOrganization(orgID, input.Name, input.Plan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update organization"})
		return
	}

	c.JSON(http.StatusOK, org)
}
