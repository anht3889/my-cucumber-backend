package api

import "github.com/gin-gonic/gin"

// LogoutHandler handles user logout requests
func LogoutHandler(c *gin.Context) {
	// Currently, since we're using JWTs, we can't actually invalidate the token
	// The client should:
	// 1. Delete the token from their storage
	// 2. Stop using the token for future requests

	c.JSON(200, gin.H{
		"message": "Successfully logged out",
	})
}
