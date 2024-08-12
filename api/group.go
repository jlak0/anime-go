package api

import (
	"anime-go/models"
	"encoding/json"
	"net/http"
)

func groupHandler(w http.ResponseWriter, r *http.Request) {
	response := []models.Group{}
	models.DB.Find(&response)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
