package models

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	dsn := "host=192.168.111.79 user=anime-go password=bJP8ZcE6FpFrcFzZ dbname=anime-go port=5432 sslmode=disable TimeZone=Europe/Paris"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	AutoMigrate(DB)

}
