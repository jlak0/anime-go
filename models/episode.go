package models

type Episode struct {
	ID         int       `gorm:"primaryKey;autoIncrement" json:"id,omitempty"`
	GroupID    int       `gorm:"not null" json:"group_id,omitempty"`
	Group      *Group    `gorm:"foreignkey:GroupID" json:"group,omitempty"`
	Episode    int       `gorm:"not null" json:"episode,omitempty"`
	AirDate    string    `json:"air_date,omitempty"`
	SeasonID   int       `gorm:"not null" json:"season_id,omitempty"`
	Season     *Season   `gorm:"foreignkey:SeasonID" json:"season,omitempty"`
	Hash       string    `gorm:"unique" json:"hash,omitempty"`
	Status     string    `gorm:"not null" json:"status,omitempty"`
	Score      int       `gorm:"not null" json:"score,omitempty"`
	Resolution string    `json:"resolution,omitempty"`
	Source     string    `json:"source,omitempty"`
	Subtitle   *Subtitle `gorm:"foreignkey:SubtitleID" json:"subtitle,omitempty"`
	SubtitleID int       `gorm:"not null" json:"subtitle_id,omitempty"`
}

func (e *Episode) Save() {
	DB.Create(&e)
}

func (e *Episode) UpdateStatus(status string) {
	DB.Model(&e).Update("status", status)
}
