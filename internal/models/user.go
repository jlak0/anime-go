package models

type User struct {
	ID       int    `gorm:"primaryKey;autoIncrement"`
	Username string `gorm:"not null;unique"`
	Password string `gorm:"not null"`
}
