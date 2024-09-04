package api

import (
	"anime-go/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func groupHandler(c *gin.Context) {
	response := []models.Group{}
	models.DB.Find(&response)
	c.JSON(http.StatusOK, response)

}

type GroupScore struct {
	ID    int `json:"id"`
	Score int `json:"score"`
}

func groupScore(c *gin.Context) {
	var s GroupScore
	err := c.ShouldBind(&s)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "cant bind"})
		return
	}
	g := models.Group{ID: s.ID}
	o := models.DB.Model(&g).Update("score", s.Score)
	if o.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
