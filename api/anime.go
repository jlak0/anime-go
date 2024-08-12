package api

import (
	"anime-go/models"
	"encoding/json"
	"net/http"

	"gorm.io/gorm"
)

func animesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := []models.Season{}
	models.DB.Preload("Episodes", func(db *gorm.DB) *gorm.DB {
		return db.Where("status = ?", "complete").Select("episode,season_id")
	}).Preload("Anime", func(db *gorm.DB) *gorm.DB {
		return db.Select("id,name")
	}).Select("id,season_number,poster_path,air_date").Find(&response, "air_date > ?", "2024-06-15")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
