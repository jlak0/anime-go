package cronjobs

import (
	"anime-go/internal/controller"
	"anime-go/internal/models"
	"anime-go/pkg/logger"
	"anime-go/pkg/qbitorrent"
	"anime-go/pkg/torrent"
	"anime-go/pkg/torrent/adapter"
	"anime-go/pkg/utils"
	"fmt"
	"regexp"

	"github.com/robfig/cron/v3"
)

func StartCronJobs() {
	c := cron.New()
	addCronJob(c, "*/10 * * * *", cloneAndAnalizeUnreadTorrent)
	addCronJob(c, "*/10 * * * *", download)
	c.Start()
}

func addCronJob(c *cron.Cron, spec string, cmd func()) {
	_, err := c.AddFunc(spec, cmd)
	if err != nil {
		logger.Log(fmt.Sprintf("添加定时任务错误: %s, Spec: %s", err.Error(), spec))
		panic(fmt.Sprintf("启动定时任务时发生致命错误: %s", err.Error()))
	}
}

func download() {
	ep := models.Find()
	for _, v := range *ep {
		err := qbitorrent.AddAndRename(&v)
		if err != nil {
			logger.Log(fmt.Sprintf("下载错误: %s, Episode: %v", err.Error(), v))
		}
	}
}

func cloneAndAnalizeUnreadTorrent() {
	t := adapter.Mikan{}
	err := torrent.Clone(&t)
	if err != nil {
		logger.Log(err.Error())
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
		"先行版",
		"粤语",
	}

	multiEpisode := regexp.MustCompile(`\d{2}-\d{2}`)

	items := &[]models.Torrent{}
	models.DB.Model(&models.Torrent{}).Where("read = ?", false).Select("title", "hash", "id").Scan(&items)
	for _, e := range *items {
		if utils.ContainsAny(e.Title, substrings) || multiEpisode.MatchString(e.Title) {
			models.DB.Model(&e).Update("read", true)
			continue
		}

		err = controller.Analize(e.Title, e.ID)
		if err != nil {
			logger.Log(err.Error())
		}
		models.DB.Model(&e).Update("read", true)

	}
}
