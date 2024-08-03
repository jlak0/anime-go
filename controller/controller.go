package controller

import (
	"anime-go/models"
	"fmt"
)

func Test() error {
	i, err := Parse("[LoliHouse] 魔王军最强的魔术师是人类 / Maougun Saikyou no Majutsushi wa Ningen datta - 05 [WebRip 1080p HEVC-10bit AAC][简繁内封字幕]")
	if err != nil {
		return err
	}
	anime := GetAnimeInfo(&i)
	a := models.Anime{}
	models.DB.FirstOrCreate(&a, *anime)
	group := models.Group{}
	subtitle := models.Subtitle{}
	models.DB.FirstOrCreate(&group, models.Group{Score: 1, Group: i.Group})
	if i.Sub != "" {
		models.DB.FirstOrCreate(&subtitle, models.Subtitle{Lang: i.Sub})
	}
	fmt.Println(subtitle)

	return nil
}
