package api

import (
	"anime-go/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func groupHandler(c *gin.Context) {
	response := []models.Group{}
	models.DB.Find(&response)
	c.JSON(http.StatusOK, response)

}
