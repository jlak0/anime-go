package models

type Season struct {
	ID           int     `gorm:"primaryKey" json:"id"`
	SeasonNumber int     `gorm:"not null" json:"season_number"`
	AnimeID      int     `gorm:"not null" json:"anime_id"` // Anime的ID应该是int类型
	Anime        Anime   `gorm:"foreignKey:AnimeID"`       // 定义关联关系
	PosterPath   string  `json:"poster_path"`
	AirDate      string  `json:"air_date"`
	BlackListed  string  `json:"black_listed"`
	BgmID        *int    `json:"bgm_id"`           // Bangumi的ID应该是int类型
	Bangumi      Bangumi `gorm:"foreignKey:BgmID"` // 定义关联关系
}
