package controller

import (
	"anime-go/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type TVShowResponse struct {
	Page         int      `json:"page"`
	Results      []TVShow `json:"results"`
	TotalPages   int      `json:"total_pages"`
	TotalResults int      `json:"total_results"`
}

type TVShow struct {
	Adult            bool     `json:"adult"`
	BackdropPath     string   `json:"backdrop_path"`
	GenreIDs         []int    `json:"genre_ids"`
	ID               int      `json:"id"`
	OriginCountry    []string `json:"origin_country"`
	OriginalLanguage string   `json:"original_language"`
	OriginalName     string   `json:"original_name"`
	Overview         string   `json:"overview"`
	Popularity       float64  `json:"popularity"`
	PosterPath       string   `json:"poster_path"`
	FirstAirDate     string   `json:"first_air_date"`
	Name             string   `json:"name"`
	VoteAverage      float64  `json:"vote_average"`
	VoteCount        int      `json:"vote_count"`
}

func getTMDB(name, lang string) (TVShow, error) {
	var data TVShow
	url := fmt.Sprintf(`https://api.themoviedb.org/3/search/tv?query=%s&include_adult=false&language=%s&page=1`, url.QueryEscape(name), lang)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return data, err
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiJ9.eyJhdWQiOiI1Nzg0Y2NmNmU1OTExM2UyNmM0N2EzMDNmYzZiY2EyOSIsIm5iZiI6MTcyMTgyOTgyNi40MDY4MjksInN1YiI6IjVmMzU1MzE3ZjZmZDE4MDAzNjJiOWFjZiIsInNjb3BlcyI6WyJhcGlfcmVhZCJdLCJ2ZXJzaW9uIjoxfQ.hZPFzcmOU4TWJ5QcrQ9eTi0v_j7mVskTNyLsIT5zCSE")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return data, err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return data, err
	}
	bodatString := (string(body))
	var response TVShowResponse
	err = json.Unmarshal([]byte(bodatString), &response)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		fmt.Println(url)
		return data, err
	}

	data = response.Results[0]
	return data, nil
}

func GetAnimeInfo(i *AnimeInfo) *models.Anime {
	var nameList []string
	if i.NameEn != "" {
		nameList = append(nameList, i.NameEn)
	}
	if i.NameJp != "" {
		nameList = append(nameList, i.NameJp)
	}
	if i.NameZh != "" {
		nameList = append(nameList, i.NameZh)
	}
	anime := models.Anime{}
	for _, v := range nameList {
		data, err := getTMDB(v, "en-US")
		if err == nil {
			anime.EnglishName = data.Name
			anime.Name = data.OriginalName
			anime.ID = data.ID
		}
		data, err = getTMDB(v, "zh-CN")
		if err == nil {
			anime.ChineseName = data.Name
			anime.Name = data.OriginalName
			anime.ID = data.ID
			anime.Image = data.PosterPath
		}
	}

	return &anime
}
