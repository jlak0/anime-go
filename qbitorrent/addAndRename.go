package qbitorrent

import (
	"anime-go/models"
	"fmt"
	"strings"
	"time"
)

func AddAndRename(ep *models.Episode) error {
	path := calculatePath(ep)
	err := Add(ep.Hash, path)
	if err != nil {
		return fmt.Errorf("无法添加种子%s", err)
	}
	go tryRename(ep)
	return nil
}

func calculatePath(ep *models.Episode) string {
	splitedDate := strings.Split(ep.Season.AirDate, "-")
	if len(splitedDate) != 3 {
		return "error"
	}
	year := splitedDate[0]
	season := calculateSeason(splitedDate[1])
	path := fmt.Sprintf("/%s/%s月新番/%s/Season %d", year, season, ep.Season.Anime.ChineseName, ep.Season.SeasonNumber)
	return path
}

func calculateSeason(month string) string {
	switch month {
	case "12", "01", "02":
		return "1"
	case "03", "04", "05":
		return "4"
	case "06", "07", "08":
		return "7"
	case "09", "10", "11":
		return "10"
	}

	return ""
}
func tryRename(ep *models.Episode) {
	var oldName string
	var err error
	for i := 0; i < 10; i++ {
		oldName, err = GetFileName(ep.Hash)
		if err == nil {
			break
		} else {
			if err.Error() == "非单一文件" {
				ep.UpdateStatus("toDelete")
			}
		}
		time.Sleep(60 * time.Second)
	}
	if err != nil {
		ep.UpdateStatus("rename")
		return
	}
	err = Rename(ep.Hash, oldName, "newName")
	if err == nil {
		ep.UpdateStatus("complete")
	}
}
