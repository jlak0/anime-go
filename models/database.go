package models

import (
	"anime-go/config"
	"fmt"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	var err error
	c := config.LoadConfig("config.json")
	if c.DB == "sqlite" {

		DB, err = gorm.Open(sqlite.Open("anime.db"), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}
	}
	if c.DB == "postgres" {
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Europe/Paris", c.DB_host, c.DB_user, c.DB_pass, c.DB_name, c.DB_port)
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}
	}
	AutoMigrate(DB)
}

func Find() *[]Episode {
	var ep []Episode
	DB.Joins("JOIN seasons ON episodes.season_id = seasons.id").
		Joins("JOIN animes ON seasons.anime_id = animes.id").
		Where("episodes.status = ? AND seasons.air_date > ? AND seasons.black_listed = ?", "pending", "2024-06-15", false).
		Preload("Season").
		Preload("Season.Anime").
		Preload("Torrent").
		Find(&ep)
	return &ep
}
