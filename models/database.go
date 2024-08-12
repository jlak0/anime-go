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
	if config.AppConfig.DB == "sqlite" {

		DB, err = gorm.Open(sqlite.Open("anime.db"), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}
	}
	if config.AppConfig.DB == "postgres" {
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Europe/Paris", config.AppConfig.DB_host, config.AppConfig.DB_user, config.AppConfig.DB_pass, config.AppConfig.DB_name, config.AppConfig.DB_port)
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
		Find(&ep)
	return &ep
}
