package models

type Episode struct {
	ID         int    `gorm:"primaryKey;autoIncrement" json:"id,omitempty"`
	Num        int    `gorm:"not null" json:"episode,omitempty"`
	AirDate    string `json:"air_date,omitempty"`
	Name       string
	Status     string `gorm:"not null" json:"status,omitempty"`
	Resolution string `json:"resolution,omitempty"`
	Source     string `json:"source,omitempty"`
	TorrentID  int
	Torrent    *Torrent
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
	DB.Where("season_id = ?", seasonID).Find(&eps)
	return eps
}
