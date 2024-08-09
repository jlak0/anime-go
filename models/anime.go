package models

import (
	"errors"
)

type Anime struct {
	ID          int    `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"not null;unique" json:"name"`
	EnglishName string `json:"english_name"`
	ChineseName string `json:"chinese_name"`
	Image       string `json:"image"`
}

func (t *Anime) Save() error {
	if DB == nil {
		return errors.New("database not initialised")
	}
	r := DB.Create(&t)

	return r.Error
}

func (t *Anime) Exist() (bool, error) {
	var anime Anime
	if DB == nil {
		return false, errors.New("database not initialised")
	}
	DB.First(&anime, t.ID)
	if anime.ID != 0 {
		t = &anime
		return true, nil
	} else {
		return false, nil
	}
}

func (g *Anime) ExistOrSave() error {
	exist, err := g.Exist()
	if err == nil && !exist {
		err = g.Save()
	}
	return err
}
