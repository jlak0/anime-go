package api

import (
	"github.com/gin-gonic/gin"
)

func verifyHandler(c *gin.Context) {
	token, err := c.Cookie("token")
	if err != nil {
		c.JSON(401, gin.H{
			"message": "Unauthorized",
		})
		return
	} else {
		if verificationHandler(token) {
			c.JSON(200, gin.H{
				"message": "ok",
			})
		} else {
			c.JSON(401, gin.H{
				"message": "Unauthorized",
			})
		}
	}
}
