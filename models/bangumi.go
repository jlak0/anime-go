package models

type Bangumi struct {
	ID     int    `gorm:"primaryKey" json:"id"`
	Score  string `json:"score"`
	OffSet int    `json:"off_set"`
}
