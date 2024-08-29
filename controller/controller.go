package controller

import (
	"anime-go/models"
	"anime-go/qbitorrent"
	"errors"
	"fmt"
)

//	func CheckExist(hash string) bool {
//		ep := models.Episode{}
//		models.DB.Where("hash = ?", hash).First(&ep)
//		if ep.ID == 0 {
//			return false
//		} else {
//			return true
//		}
//	}
func Analize(title string, torrentID int) error {
	// exist := CheckExist(hash)
	// if exist {
	// 	return fmt.Errorf("已经存在")
	// }
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
			if e.Status == "complete" || e.Status == "rename" {
				err = qbitorrent.Delete(e.Torrent.Hash)
				if err != nil {
					return err
				}
			}
			models.DB.Model(&e).Updates(&models.Episode{
				GroupID:    group.ID,
				SubtitleID: subtitle.ID,
				Status:     "pending",
				Resolution: i.Dpi,
				Source:     i.Source,
				TorrentID:  torrentID,
			})
			return nil
		}

	}

	models.DB.Create(&models.Episode{
		Number:     i.Episode,
		Status:     "pending",
		Source:     i.Source,
		Resolution: i.Dpi,
		GroupID:    group.ID,
		SubtitleID: subtitle.ID,
		TorrentID:  torrentID,
		SeasonID:   season.ID,
	})
	return nil
}

// if len(*existEpisodes) > 0 {
// 	best := isTheBest(score, existEpisodes)
// 	if !best {
// 		status = "deleted"
// 	} else {
// 		deleteAll(existEpisodes)
// 	}
// }
// episode := models.Episode{GroupID: group.ID, Episode: i.Episode, SeasonID: season.ID, Resolution: i.Dpi, Source: i.Source, SubtitleID: subtitle.ID, Score: score, Status: status, Hash: hash}
// x, err := FindEpisode(anime.ID, i.Season, i.Episode)
// if err == nil {
// 	episode.AirDate = x.AirDate
// }
// episode.Save()
// return err

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
		groupScore = ep.Subtitle.Score
	}

	score := groupScore * subtitleScore
	return score
}

// func findSameEpisode(seasonID, episode int) *[]models.Episode {
// 	episodes := []models.Episode{}
// 	models.DB.Where("season_id = ? AND episode = ? AND status NOT IN (?)", seasonID, episode, []string{"deleted", "toDelete"}).Find(&episodes)
// 	return &episodes
// }

// func isTheBest(score int, episodes *[]models.Episode) bool {
// 	for _, e := range *episodes {
// 		if score > e.Score {
// 			return true
// 		}
// 	}
// 	return false
// }

// func deleteAll(episodes *[]models.Episode) {
// 	for _, e := range *episodes {
// 		switch e.Status {
// 		case "pending":
// 			e.UpdateStatus("deleted")
// 		case "complete", "rename":
// 			err := qbitorrent.Delete(e.Torrent.Hash)
// 			if err != nil {
// 				e.UpdateStatus("toDelete")
// 			} else {
// 				e.UpdateStatus("deleted")
// 			}
// 		case "toDelete", "deleted":

// 		default:
// 			logger.Log("unknow episode statuts:" + e.Status)
// 		}
// 	}
// }
