package models

import "errors"

type Subtitle struct {
	ID    int    `gorm:"primaryKey;autoIncrement" json:"id,omitempty"`
	Lang  string `gorm:"unique" json:"lang,omitempty"`
	Score int    `gorm:"not null" json:"score,omitempty"`
}

func (s *Subtitle) Save() error {

	r := DB.Create(&s)
	return r.Error
}

func (s *Subtitle) Exist() (bool, error) {
	if s.Lang == "" {
		return false, errors.New("lang exist null")
	}
	if DB == nil {
		return false, errors.New("database not initialised")
	}
	DB.Where(s).First(&s)
	if s.ID != 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func (s *Subtitle) ExistOrSave() error {
	exist, err := s.Exist()
	if err == nil && !exist {
		err = s.Save()
	}
	return err
}
