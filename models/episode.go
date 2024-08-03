package models

type Episode struct {
	ID         int      `gorm:"primaryKey;autoIncrement" json:"id"`
	GroupID    int      `gorm:"not null"`
	Group      Group    `gorm:"foreignkey:GroupID" json:"group_id"`
	Episode    int      `gorm:"not null" json:"episode"`
	AirDate    string   `json:"air_date"`
	SeasonID   int      `gorm:"not null"`
	Season     Season   `gorm:"foreignkey:SeasonID" json:"season_id"`
	Hash       string   `gorm:"unique" json:"hash"`
	Status     string   `gorm:"not null" json:"status"`
	Score      int      `gorm:"not null" json:"score"`
	Resolution string   `json:"resolution"`
	Source     string   `json:"source"`
	Subtitle   Subtitle `gorm:"foreignkey:SubtitleID"`
	SubtitleID int      `gorm:"not null" json:"subtitle_id"`
}
