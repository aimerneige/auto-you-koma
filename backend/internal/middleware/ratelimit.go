package middleware

import (
	"net/http"

	"github.com/aimerneige/auto-you-koma/internal/repository"
	"github.com/gin-gonic/gin"
)

// QuotaMiddleware checks if user has enough quota and increments usage.
// Note: In production, this should ideally be handled asynchronously or within a transaction
// around the actual generation logic.
func QuotaMiddleware(userRepo repository.UserRepository, cost int) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		user, err := userRepo.GetByID(c.Request.Context(), userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user info"})
			c.Abort()
			return
		}

		if user.QuotaUsed+cost > user.QuotaLimit {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "quota exceeded"})
			c.Abort()
			return
		}

		// Simple optimistic update
		user.QuotaUsed += cost
		if err := userRepo.Update(c.Request.Context(), user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update quota"})
			c.Abort()
			return
		}

		c.Next()
	}
}
