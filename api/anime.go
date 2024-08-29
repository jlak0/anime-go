package api

import (
	"anime-go/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func animesHandler(c *gin.Context) {
	response := []models.Season{}
	models.DB.Preload("Episodes").Preload("Episodes.Group").Preload("Episodes.Subtitle").Preload("Anime", func(db *gorm.DB) *gorm.DB {
		return db.Select("id,chinese_name,image")
	}).Preload("Bangumi").Find(&response, "air_date > ?", "2024-06-15")
	c.JSON(http.StatusOK, response)
}

type BlackList struct {
	BlackListed *bool `json:"black_listed" binding:"required" `
}

func blackAnime(c *gin.Context) {
	var b BlackList
	if err := c.ShouldBindJSON(&b); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	s := models.Season{ID: id}
	models.DB.Model(&s).Update("black_listed", b.BlackListed)
	c.JSON(http.StatusOK, s)
}
