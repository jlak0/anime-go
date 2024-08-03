package models

type Subtitle struct {
	ID    int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Lang  string `gorm:"unique" json:"lang"`
	Score int    `gorm:"not null" json:"score"`
}
