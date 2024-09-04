package controller

import (
	"anime-go/internal/models"
	"anime-go/pkg/qbitorrent"
	"anime-go/pkg/parser"
	"errors"
)

const (
	StatusComplete = "complete"
	StatusRename   = "rename"
	StatusPending  = "pending"
)

func Analize(title string, torrentID int) error {
	
	i, err := parser.Parse(title)
	if err != nil {
		return err
	}
	anime := GetAnimeInfo(&i)
	if anime.Name == "" {

		return errors.New("failed to get name" + title)
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
	season := models.Season{AnimeID: anime.ID, Number: i.Season}
	season.ExistOrSave()
	if season.Bangumi.ID == 0 {
		saveBgmId(anime.Name, season.AirDate, season.ID)
	}
	// existEpisodes := findSameEpisode(season.ID, i.Episode)
	score := group.Score * subtitle.Score
	currentScore := findCurrentScore(season.ID, i.Episode)
	if score <= currentScore {
		return nil
	}

	eps, _ := PreCreateEpisode(anime.ID, season.ID, season.Number)
	for _, e := range *eps {
		if e.Number == i.Episode {
			if e.Status == StatusComplete || e.Status == StatusRename {
				err = qbitorrent.Delete(e.Torrent.Hash)
				if err != nil {
					return err
				}
			}
			models.DB.Model(&e).Updates(&models.Episode{
				GroupID:    group.ID,
				SubtitleID: subtitle.ID,
				Status:     StatusPending,
				Resolution: i.Dpi,
				Source:     i.Source,
				TorrentID:  torrentID,
			})
			return nil
		}

	}

	models.DB.Create(&models.Episode{
		Number:     i.Episode,
		Status:     StatusPending,
		Source:     i.Source,
		Resolution: i.Dpi,
		GroupID:    group.ID,
		SubtitleID: subtitle.ID,
		TorrentID:  torrentID,
		 SeasonID:   season.ID,
	})
	return nil
}


func findCurrentScore(seasonID, episodeNum int) int {
	ep := &models.Episode{}
	models.DB.Preload("Group").Preload("Subtitle").Where("season_id = ? AND number = ?", seasonID, episodeNum).First(&ep)
	if (ep.ID) == 0 {
		return 0
	}

	groupScore, subtitleScore := 1, 1
	if ep.Group != nil {
		groupScore = ep.Group.Score
	}
	if ep.Subtitle != nil {
		subtitleScore = ep.Subtitle.Score
	}

	score := groupScore * subtitleScore
	return score
}
