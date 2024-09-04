package torrent

import (
	"anime-go/internal/models"
	"anime-go/pkg/parser"
	"anime-go/pkg/torrent/adapter"
	"anime-go/pkg/utils"
	"fmt"
)

type TorrentIF interface {
	Get() ([]models.Torrent, error)
}

func Clone(t TorrentIF) (error) {
	torrents, err := t.Get()
	if err != nil {
		return err
	}
	for _, torrent := range torrents {
		models.DB.FirstOrCreate(&torrent, torrent)
	}
	return nil
}

func Search(keyword string) ([]models.Torrent, error) {
	t := adapter.Mikan{}
	return t.Search(keyword)
}

func Test() {

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
		"OVA",
		"剧场版",
		
	}

	t := adapter.Mikan{}
	// Clone(&t)
	data, err := t.Get()
	// data, err := Search("这个世界漏洞百出")
	if err != nil {
		fmt.Println(err)
	}
	for _, v := range data {
		if utils.ContainsAny(v.Title, substrings) {
			continue
		}
		_, err := parser.Parse(v.Title)
		if err != nil {
			fmt.Println(err,v.Title,v.Hash)
			
		}
			// fmt.Printf("%s:%s S%02dE%02d\n",p.Group, p.NameZh, p.Season, p.Episode)
		
	}
}