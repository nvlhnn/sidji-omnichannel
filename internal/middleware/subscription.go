package middleware

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sidji-omnichannel/internal/subscription"
)

// Subscription check middleware prevents actions if organization is over limits
func Subscription(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		orgID := GetOrganizationID(c)
		if orgID == uuid.Nil {
			c.Next()
			return
		}

		// Fetch organization plan and counts
		var plan string
		var msgUsed, msgLimit int
		err := db.QueryRow("SELECT plan, message_usage_used, message_usage_limit FROM organizations WHERE id = $1", orgID).Scan(&plan, &msgUsed, &msgLimit)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify subscription status"})
			return
		}

		var userCount, channelCount int
		_ = db.QueryRow("SELECT COUNT(*) FROM users WHERE organization_id = $1", orgID).Scan(&userCount)
		_ = db.QueryRow("SELECT COUNT(*) FROM channels WHERE organization_id = $1", orgID).Scan(&channelCount)

		// Check if organization is compliant (Basic limits)
		if !subscription.IsCompliance(plan, userCount, channelCount) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Subscription resource limit exceeded. Please upgrade your plan.",
				"code":  "subscription_limit_reached",
			})
			return
		}

		// Check message usage limit
		if msgLimit != -1 && msgUsed >= msgLimit {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Monthly message volume quota reached. Please upgrade your plan to continue sending messages.",
				"code":  "message_limit_reached",
			})
			return
		}

		c.Next()
	}
}
