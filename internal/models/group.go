package models

import (
	"errors"
)

type Group struct {
	ID    int    `gorm:"primaryKey;autoIncrement" json:"id,omitempty"`
	Group string `gorm:"not null;unique" json:"group,omitempty"`
	Score int    `gorm:"not null" json:"score,omitempty"`
}

func (g *Group) Save() error {
	if g.Group == "" {
		return errors.New("group save null")
	}
	if DB == nil {
		return errors.New("database not initialised")
	}
	r := DB.Create(&g)
	return r.Error
}

func (g *Group) Exist() (bool, error) {
	if g.Group == "" {
		return false, errors.New("group exist null")
	}
	if DB == nil {
		return false, errors.New("database not initialised")
	}
	DB.Where(`"group" = ?`, g.Group).First(&g)
	if g.ID != 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func (g *Group) ExistOrSave() error {
	exist, err := g.Exist()
	if err == nil && !exist {
		err = g.Save()
	}
	return err
}
