package models

import (
	"errors"
)

type Torrent struct {
	ID      int    `gorm:"primaryKey;autoIncrement" json:"id,omitempty"`
	Title   string `gorm:"not null" json:"title,omitempty"`
	Link    string `gorm:"not null;unique" json:"link,omitempty"`
	Hash    string `gorm:"not null;unique" json:"hash"`
	PubDate string `gorm:"type:timestamp;not null" json:"pubDate,omitempty"`
	Read    bool   `gorm:"not null" json:"status,omitempty"`
}

func (t *Torrent) Create() error {
	if DB == nil {
		return errors.New("database not initialised")
	}
	r := DB.Create(&t)
	return r.Error
}
