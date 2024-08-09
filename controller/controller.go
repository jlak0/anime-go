package controller

import (
	"anime-go/models"
	"errors"
	"fmt"
)

func CheckExist(hash string) bool {
	ep := models.Episode{}
	models.DB.Where("hash = ?", hash).First(&ep)
	if ep.ID == 0 {
		return false
	} else {
		return true
	}

}
func Analize(title, hash string) error {
	// b := CheckExist("a46c7a3a49054b0823b68cbe5670bf1a374aacff")
	// fmt.Println(b)
	status := "pending"
	i, err := Parse(title)
	if err != nil {
		return err
	}
	anime := GetAnimeInfo(&i)
	if anime.Name == "" {

		fmt.Println("failed to get name")
		fmt.Println(title)
		fmt.Println(i.NameZh)
		fmt.Println(i.NameEn)
		fmt.Println(i.NameJp)
		return errors.New("failed to get name")
	}
	err = anime.ExistOrSave()
	if err != nil {
		return errors.New("create anime error")
	}
	group := models.Group{Score: 1, Group: i.Group}
	err = group.ExistOrSave()
	if err != nil {
		return errors.New("create group error")
	}
	subtitle := models.Subtitle{Score: 1}
	if i.Sub != "" {
		subtitle.Lang = i.Sub
		subtitle.ExistOrSave()
	}
	season := models.Season{AnimeID: anime.ID, SeasonNumber: i.Season}
	season.ExistOrSave()

	episode := models.Episode{GroupID: group.ID, Episode: i.Episode, SeasonID: season.ID, Resolution: i.Dpi, Source: i.Source, SubtitleID: subtitle.ID, Score: group.Score * subtitle.Score, Status: status, Hash: hash}
	x, err := FindEpisode(anime.ID, i.Season, i.Episode)
	if err == nil {
		episode.AirDate = x.AirDate
	}
	episode.Save()
	return err
}
