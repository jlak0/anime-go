package controller

import (
	"anime-go/logger"
	"anime-go/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type APIResponse struct {
	Results int    `json:"results"`
	List    []Item `json:"list"`
}

type Item struct {
	ID         int    `json:"id"`
	URL        string `json:"url"`
	Type       int    `json:"type"`
	Name       string `json:"name"`
	NameCN     string `json:"name_cn"`
	Eps        int    `json:"eps"`
	AirDate    string `json:"air_date"`
	AirWeekday int    `json:"air_weekday"`
	Rank       int    `json:"rank"`
	Rating     Rating `json:"rating"`
}

type Rating struct {
	Score float32 `json:"score"`
}

func DateDifference(date1Str, date2Str string) (int, error) {
	// 定义日期格式
	const layout = "2006-01-02"

	// 解析日期字符串
	date1, err := time.Parse(layout, date1Str)
	if err != nil {
		return 999, fmt.Errorf("error parsing date1: %w", err)
	}

	date2, err := time.Parse(layout, date2Str)
	if err != nil {
		return 999, fmt.Errorf("error parsing date2: %w", err)
	}

	// 计算两个日期之间的差值
	diff := date2.Sub(date1)

	// 将时间差转换为天数

	days := int(diff.Hours() / 24)
	if days < 0 {
		days *= -1
	}

	return days, nil
}

func searchBgm(name string) (APIResponse, error) {
	var apiResponse APIResponse
	resp, err := http.Get(fmt.Sprintf(`https://api.bgm.tv/search/subject/%s?type=2&responseGroup=large`, name))
	if err != nil {
		return apiResponse, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return apiResponse, err
	}
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return apiResponse, err
	}
	return apiResponse, nil
}

func GetBgmID() {
	animes := []models.Season{}
	models.DB.Where("id NOT IN (?)", models.DB.Model(&models.Bangumi{}).Select("season_id")).Preload("Anime").Find(&animes)
	for _, a := range animes {
		s, err := searchBgm(a.Anime.Name)
		if err != nil {
			logger.Log("错误")
		}
		for _, item := range s.List {
			diff, _ := DateDifference(item.AirDate, a.AirDate)
			if diff < 2 {
				bgm := models.Bangumi{
					ID:         item.ID,
					SeasonID:   a.ID,
					Score:      item.Rating.Score,
					Eps:        item.Eps,
					AirDate:    item.AirDate,
					AirWeekday: item.AirWeekday,
					Rank:       item.Rank,
					URL:        item.URL,
				}
				models.DB.Create(&bgm)
				fmt.Println(item.NameCN)
				break
			}
		}
	}
}

func Test() {
	var result []models.Season
	models.DB.Where("id NOT IN (?)", models.DB.Model(&models.Bangumi{}).Select("season_id")).Preload("Anime").Find(&result)
	for _, r := range result {
		fmt.Println(r.ID)
	}
}
