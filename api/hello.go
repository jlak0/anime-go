package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func helloHandler(c *gin.Context) {
	// 创建一个响应对象

	c.JSON(http.StatusOK, gin.H{
		"message": "Hello, world!",
	})
}
