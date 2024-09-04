package api

import (
	"anime-go/internal/models"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID参数"})
		return
	}

	s := models.Season{ID: id}
	result := models.DB.Model(&s).Update("black_listed", b.BlackListed)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到指定季度"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}
