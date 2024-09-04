package controller

import (
	"anime-go/internal/models"
	"anime-go/pkg/parser"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type SeasonResponse struct {
	Episodes []Episode `json:"episodes"`
}
type Episode struct {
	AirDate        string  `json:"air_date"`
	EpisodeNumber  int     `json:"episode_number"`
	Name           string  `json:"name"`
	Overview       string  `json:"overview"`
	ID             int     `json:"id"`
	ProductionCode string  `json:"production_code"`
	Runtime        int     `json:"runtime"`
	SeasonNumber   int     `json:"season_number"`
	StillPath      string  `json:"still_path"`
	VoteAverage    float64 `json:"vote_average"`
	VoteCount      int     `json:"vote_count"`
}

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

func contains(slice []int, value int) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
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

	var response TVShowResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		fmt.Println(url)
		return data, err
	}
	if len(response.Results) == 0 {
		return data, errors.New("no result")
	}
	var index int
	for i, show := range response.Results {
		if !show.Adult && show.OriginalLanguage == "ja" && contains(show.GenreIDs, 16) {
			index = i
		}
	}
	data = response.Results[index]
	return data, nil
}

func GetAnimeInfo(i *parser.AnimeInfo) *models.Anime {
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

func FindEpisode(animeID, seasonNumber, episodeNumber int) (Episode, error) {
	var ep Episode

	url := fmt.Sprintf("https://api.themoviedb.org/3/tv/%d/season/%d/episode/%d?language=ja", animeID, seasonNumber, episodeNumber)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ep, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiJ9.eyJhdWQiOiI1Nzg0Y2NmNmU1OTExM2UyNmM0N2EzMDNmYzZiY2EyOSIsIm5iZiI6MTcyMjc5NDEwNy43OTA4NTQsInN1YiI6IjVmMzU1MzE3ZjZmZDE4MDAzNjJiOWFjZiIsInNjb3BlcyI6WyJhcGlfcmVhZCJdLCJ2ZXJzaW9uIjoxfQ.BxT_WUa6DxVbEIZmR567jJDBHjIjw9jDPZOu3yWlDE4")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return ep, fmt.Errorf("failed to make request: %v", err)
	}
	defer func() {
		res.Body.Close()
	}()

	if res.StatusCode != http.StatusOK {
		return ep, fmt.Errorf("unexpected response status: %s", res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return ep, fmt.Errorf("failed to read response body: %v", err)
	}
	json.Unmarshal(body, &ep)
	return ep, nil
}

func PreCreateEpisode(animeID, seasonID, seasonNum int) (*[]models.Episode, error) {
	eps := models.FindAllEpisode(seasonID)
	tmdbEps, err := FindAllEpisodes(animeID, seasonNum)
	if err != nil {
		return eps, err
	}
	for _, e := range tmdbEps.Episodes {

		if !episodeCreated(eps, e) {
			initialEpisode(e, seasonID)
		}
	}
	eps = models.FindAllEpisode(seasonID)
	return eps, nil
}

func initialEpisode(e Episode, seasonID int) {
	ep := models.Episode{
		AirDate:  e.AirDate,
		Name:     e.Name,
		SeasonID: seasonID,
		Number:   e.EpisodeNumber,
	}
	models.DB.Create(&ep)
}

func updateDifference(tmdbEp *Episode, ep *models.Episode) {
	if ep.Name != tmdbEp.Name || ep.AirDate != tmdbEp.AirDate {
		ep.Update(tmdbEp.AirDate, tmdbEp.Name)
	}
}

func episodeCreated(eps *[]models.Episode, tmdbEp Episode) bool {
	for _, e := range *eps {
		if e.Number == tmdbEp.EpisodeNumber {
			updateDifference(&tmdbEp, &e)

			return true
		}
	}
	return false
}

func FindAllEpisodes(animeID, seasonNum int) (*SeasonResponse, error) {
	season := &SeasonResponse{}
	url := fmt.Sprintf(`https://api.themoviedb.org/3/tv/%d/season/%d?language=zh-cn`, animeID, seasonNum)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return season, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiJ9.eyJhdWQiOiI1Nzg0Y2NmNmU1OTExM2UyNmM0N2EzMDNmYzZiY2EyOSIsIm5iZiI6MTcyMjc5NDEwNy43OTA4NTQsInN1YiI6IjVmMzU1MzE3ZjZmZDE4MDAzNjJiOWFjZiIsInNjb3BlcyI6WyJhcGlfcmVhZCJdLCJ2ZXJzaW9uIjoxfQ.BxT_WUa6DxVbEIZmR567jJDBHjIjw9jDPZOu3yWlDE4")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return season, fmt.Errorf("failed to make request: %v", err)
	}
	defer func() {
		res.Body.Close()
	}()
	if res.StatusCode != http.StatusOK {
		return season, fmt.Errorf("unexpected response status: %s", res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return season, fmt.Errorf("failed to read response body: %v", err)
	}
	json.Unmarshal(body, &season)
	return season, nil
}
