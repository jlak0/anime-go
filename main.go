package main

import (
	"anime-go/api"
	"anime-go/controller"
	"anime-go/logger"
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

	logger.Log("程序启动")
	defer logger.Close()
	startCronJobs()
	api.Serve()
}

func download() {
	ep := models.Find()
	var err error
	for _, v := range *ep {
		err = qbitorrent.AddAndRename(&v)
		if err != nil {
			fmt.Printf("错误%s", err.Error())
		}
	}
}
func startCronJobs() {
	c := cron.New()
	// 添加一个每秒执行一次的任务
	_, err := c.AddFunc("*/10 * * * *", cloneAndAnalizeUnreadTorrent)
	if err != nil {
		e := fmt.Errorf("error adding cron job:%s", err.Error())
		logger.Log(e.Error())
		panic("fatal error starting cron")
	}
	_, err = c.AddFunc("*/10 * * * *", download)
	if err != nil {
		e := fmt.Errorf("error adding cron job:%s", err.Error())
		logger.Log(e.Error())
		panic("fatal error starting cron")
	}
	// 启动cron调度器
	c.Start()
}

func cloneAndAnalizeUnreadTorrent() {
	err := CloneMikan()
	if err != nil {
		logger.Log(err.Error())
		return
	}
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
		err = controller.Analize(e.Title, hash)
		if err != nil {
			logger.Log(err.Error())
		}
		models.DB.Model(&e).Update("read", true)

	}
}

func CloneMikan() error {
	t := []models.Torrent{}
	err := torrent.GetMikan(&t)
	if err != nil {
		return err
	}
	for _, tx := range t {
		models.DB.FirstOrCreate(&tx, tx)
	}
	return nil
}
