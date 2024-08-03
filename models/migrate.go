package models

import (
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(
		&Torrent{},
		&Episode{},
		&Season{},
		&Anime{},
		&Bangumi{},
		&Group{},
		&Subtitle{},
	)
}
