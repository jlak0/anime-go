package models

type Bangumi struct {
	ID     int    `gorm:"primaryKey" json:"id,omitempty"`
	Score  string `json:"score,omitempty"`
	OffSet int    `json:"off_set,omitempty"`
}
