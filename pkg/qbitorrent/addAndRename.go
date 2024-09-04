package qbitorrent

import (
	"anime-go/internal/models"
	"anime-go/pkg/logger"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func AddAndRename(ep *models.Episode) error {
	path := calculatePath(ep)
	err := Add(ep.Torrent.Hash, path)
	if err != nil {
		return fmt.Errorf("无法添加种子: %w", err)
	}
	go func() {
		if err := tryRename(ep); err != nil {
			logger.Log(fmt.Sprintf("重命名错误: %s, Episode: %v", err.Error(), ep))
		}
	}()
	return nil
}

func sanitizeName(name string) string {
	// 定义一个正则表达式，匹配所有文件系统不支持的字符
	// 这里匹配： / \ : * ? " < > | 等
	re := regexp.MustCompile(`[\/\\:*?"<>|]+`)
	// 将这些不支持的字符替换为空字符串
	return re.ReplaceAllString(name, " ")
}

func calculatePath(ep *models.Episode) string {
	splitedDate := strings.Split(ep.Season.AirDate, "-")
	if len(splitedDate) != 3 {
		return "error"
	}
	year := splitedDate[0]
	season := calculateSeason(splitedDate[1])
	path := fmt.Sprintf("/%s/%s月新番/%s/Season %d", year, season, ep.Season.Anime.ChineseName, ep.Season.Number)
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

func tryRename(ep *models.Episode) error {
	var oldName string
	var err error
	for i := 0; i < 60; i++ {
		oldName, err = GetFileName(ep.Torrent.Hash)
		if err == nil {
			break
		} else {
			if err.Error() == "非单一文件" {
				err := Delete(ep.Torrent.Hash)
				if err != nil {
					ep.UpdateStatus("toDelete")
				} else {
					ep.UpdateStatus("deleted")
				}
			}
		}
		time.Sleep(60 * time.Second)
	}
	if err != nil {
		ep.UpdateStatus("rename")
		return fmt.Errorf("重命名错误: %w", err)
	}
	ext := filepath.Ext(oldName)
	newName := fmt.Sprintf(`%s S%02dE%02d%s`, sanitizeName(ep.Season.Anime.ChineseName), ep.Season.Number, ep.Number, ext)
	fmt.Println(newName)
	err = Rename(ep.Torrent.Hash, oldName, newName)
	if err == nil {
		ep.UpdateStatus("complete")
	}
	return nil
}
