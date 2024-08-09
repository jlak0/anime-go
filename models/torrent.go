package models

import (
	"errors"
)

type Torrent struct {
	ID      int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Title   string `gorm:"not null" json:"title"`
	Link    string `gorm:"not null;unique" json:"link"`
	PubDate string `gorm:"type:timestamp;not null" json:"pubDate"`
	Read    bool   `gorm:"not null" json:"status"`
}

func (t *Torrent) Create() error {
	if DB == nil {
		return errors.New("database not initialised")
	}
	r := DB.Create(&t)
	return r.Error
}
