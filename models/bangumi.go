package models

type Bangumi struct {
	ID         int     `gorm:"primaryKey" json:"id,omitempty"`
	SeasonID   int     `json:"season_id,omitempty"`
	Score      float32 `json:"score,omitempty"`
	Eps        int     `json:"eps"`
	AirDate    string  `json:"air_date"`
	AirWeekday int     `json:"air_weekday"`
	Rank       int     `json:"rank"`
	URL        string  `json:"url"`
	OffSet     int     `json:"off_set,omitempty"`
}
