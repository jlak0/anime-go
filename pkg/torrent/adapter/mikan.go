package adapter

import (
	"anime-go/internal/models"
	"encoding/xml"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// RSS represents the top-level RSS XML structure.
type RSS struct {
	Channel Channel `xml:"channel"`
}

// Channel represents the <channel> element in the RSS feed.
type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Items       []Item `xml:"item"`
}

// Item represents the <item> element in the RSS feed.
type Item struct {
	GUID        string  `xml:"guid"`
	Link        string  `xml:"link"`
	Title       string  `xml:"title"`
	Description string  `xml:"description"`
	Torrent     Torrent `xml:"torrent"`
}

// Torrent represents the custom <torrent> element in the RSS feed.
type Torrent struct {
	Link          string `xml:"link"`
	ContentLength int    `xml:"contentLength"`
	PubDate       string `xml:"pubDate"`
}

type Mikan struct {
	Title string
	Hash  string
	Link  string
	PubDate string
}

func GetMikan(t *[]models.Torrent) error {
	resp, err := http.Get("https://mikanproxy.flak.workers.dev/") // 替换为你的RSS URL
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return (err)
	}

	var rss RSS
	err = xml.Unmarshal(body, &rss)
	if err != nil {
		return (err)
	}

	for _, item := range rss.Channel.Items {
		parts := strings.Split(item.Link, "/")
		hash := parts[len(parts)-1]
		*t = append(*t, models.Torrent{
			Title: item.Title,
			Hash:  hash,
			Link:  item.Link,

			PubDate: item.Torrent.PubDate,
		})

	}
	return nil
}


func (t *Mikan) Get() ([]models.Torrent, error) {
	torrents := []models.Torrent{}
	err := GetMikan(&torrents)
	if err != nil {
		return nil, err
	}
	return torrents, nil
}

func (t *Mikan) Search(keyword string) ([]models.Torrent, error) {
	torrents := []models.Torrent{}
	resp, err := http.Get("https://mikanani.me/RSS/Search?searchstr=" + url.QueryEscape(keyword)) // 替换为你的RSS URL
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var rss RSS
	err = xml.Unmarshal(body, &rss)
	if err != nil {
		return nil, err
	}
	for _, item := range rss.Channel.Items {
		parts := strings.Split(item.Link, "/")
		hash := parts[len(parts)-1]
		torrents = append(torrents, models.Torrent{
			Title: item.Title,
			Hash:  hash,
			Link:  item.Link,

			PubDate: item.Torrent.PubDate,
		})

	}
	return torrents, nil
}