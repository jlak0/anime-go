package main

import (
	"anime-go/controller"
	"anime-go/models"
	"anime-go/qbitorrent"
	"anime-go/torrent"
	"fmt"
	"regexp"
	"strings"

	"github.com/robfig/cron/v3"
)

func containsAny(str string, substrings []string) bool {
	for _, substr := range substrings {
		if strings.Contains(str, substr) {
			return true
		}
	}
	return false
}

func main() {

	c := cron.New()

	ep := models.Find()
	var err error
	for _, v := range *ep {
		err = qbitorrent.AddAndRename(&v)
		if err != nil {
			fmt.Println(err)
		}
	}

	// 添加一个每秒执行一次的任务
	_, err = c.AddFunc("*/2 * * * *", CloneAndAnalizeUnreadTorrent)
	if err != nil {
		fmt.Println("Error adding cron job:", err)
		return
	}

	// 启动cron调度器
	c.Start()

	// 防止主程序退出
	select {}
}

func CloneAndAnalizeUnreadTorrent() {
	CloneMikan()
	substrings := []string{"国漫",
		"整理搬运",
		"外挂",
		"1280x720",
		"720P",
		"ABEMA",
		"CR 1920x1080",
		"BDRip",
		"3840x2160",
		"4K",
		"B-Global",
		"合集",
		"先行版"}

	multiEpisode := regexp.MustCompile(`\d{2}-\d{2}`)

	var items []models.Torrent
	models.DB.Model(&models.Torrent{}).Where("read = ?", false).Select("title", "link", "id").Scan(&items)
	for _, e := range items {
		if containsAny(e.Title, substrings) || multiEpisode.MatchString(e.Title) {
			models.DB.Model(&e).Update("read", true)
			continue
		}

		parts := strings.Split(e.Link, "/")
		hash := parts[len(parts)-1]
		controller.Analize(e.Title, hash)
		models.DB.Model(&e).Update("read", true)

	}
}

func CloneMikan() {
	t := []models.Torrent{}
	torrent.GetMikan(&t)
	for _, tx := range t {
		models.DB.FirstOrCreate(&tx, tx)
	}
}
