package models

type Episode struct {
	ID         int       `gorm:"primaryKey;autoIncrement" json:"id,omitempty"`
	Number     int       `gorm:"not null" json:"number"`
	AirDate    string    `json:"air_date"`
	Name       string    `json:"name"`
	Status     string    `gorm:"not null" json:"status,omitempty"`
	Resolution string    `json:"resolution,omitempty"`
	Source     string    `json:"source,omitempty"`
	TorrentID  int       `json:"torrent_id,omitempty"`
	Torrent    *Torrent  `json:"torrent,omitempty"`
	SeasonID   int       `gorm:"not null" json:"season_id,omitempty"`
	Season     *Season   `gorm:"foreignkey:SeasonID" json:"season,omitempty"`
	GroupID    int       `gorm:"not null" json:"group_id,omitempty"`
	Group      *Group    `gorm:"foreignkey:GroupID" json:"group,omitempty"`
	SubtitleID int       `gorm:"not null" json:"subtitle_id,omitempty"`
	Subtitle   *Subtitle `gorm:"foreignkey:SubtitleID" json:"subtitle,omitempty"`
}

func (e *Episode) Save() {
	DB.Create(&e)
}

func (e *Episode) UpdateStatus(status string) {
	DB.Model(&e).Update("status", status)
}

func (e *Episode) Update(airDate, name string) {
	DB.Model(&e).Updates(&Episode{Name: name, AirDate: airDate})
}

func FindAllEpisode(seasonID int) *[]Episode {
	eps := &[]Episode{}
	DB.Preload("Torrent").Where("season_id = ?", seasonID).Find(&eps)
	return eps
}
