package models

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Season struct {
	ID           int       `gorm:"primaryKey" json:"id,omitempty"`
	SeasonNumber int       `gorm:"not null" json:"season_number,omitempty"`
	AnimeID      int       `gorm:"not null" json:"anime_id,omitempty"`
	Anime        *Anime    `gorm:"foreignKey:AnimeID" json:"anime,omitempty"`
	PosterPath   string    `json:"poster_path,omitempty"`
	AirDate      string    `json:"air_date,omitempty"`
	BlackListed  bool      `json:"black_listed"`
	BgmID        *int      `json:"bgm_id,omitempty"`
	Bangumi      *Bangumi  `gorm:"foreignKey:BgmID" json:"bangumi,omitempty"`
	Episodes     []Episode `gorm:"foreignKey:SeasonID" json:"episodes,omitempty"`
}

type Ep struct {
	AirDate        string  `json:"air_date"`
	EpisodeNumber  int     `json:"episode_number"`
	EpisodeType    string  `json:"episode_type"`
	ID             int     `json:"id"`
	Name           string  `json:"name"`
	Overview       string  `json:"overview"`
	ProductionCode string  `json:"production_code"`
	Runtime        int     `json:"runtime"`
	SeasonNumber   int     `json:"season_number"`
	ShowID         int     `json:"show_id"`
	StillPath      string  `json:"still_path"`
	VoteAverage    float64 `json:"vote_average"`
	VoteCount      int     `json:"vote_count"`
}

// 定义 Main 结构体
type TMDB_SEASON struct {
	ID         int    `json:"id"`
	PosterPath string `json:"poster_path"`
	AirDate    string `json:"air_date"`
	Episodes   []Ep   `json:"episodes"`
}

func (t *Season) Exist() (bool, error) {
	DB.Where(t).First(&t)
	if t.ID != 0 {
		return true, nil
	}
	return false, nil
}

func (s *Season) Find() error {
	url := fmt.Sprintf("https://api.themoviedb.org/3/tv/%d/season/%d?language=ja", s.AnimeID, s.SeasonNumber)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiJ9.eyJhdWQiOiI1Nzg0Y2NmNmU1OTExM2UyNmM0N2EzMDNmYzZiY2EyOSIsIm5iZiI6MTcyMjc5NDEwNy43OTA4NTQsInN1YiI6IjVmMzU1MzE3ZjZmZDE4MDAzNjJiOWFjZiIsInNjb3BlcyI6WyJhcGlfcmVhZCJdLCJ2ZXJzaW9uIjoxfQ.BxT_WUa6DxVbEIZmR567jJDBHjIjw9jDPZOu3yWlDE4")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %v", err)
	}
	defer func() {
		res.Body.Close()
	}()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response status: %s", res.Status)

	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}
	var tmdb_season TMDB_SEASON
	json.Unmarshal(body, &tmdb_season)
	s.ID = tmdb_season.ID
	s.PosterPath = tmdb_season.PosterPath
	s.AirDate = tmdb_season.AirDate
	s.BlackListed = false
	return nil
}

func (s *Season) ExistOrSave() {
	exist, err := s.Exist()
	if err == nil && !exist {
		err = s.Find()
		if err == nil {
			DB.Create(&s)
		}
	}
}
