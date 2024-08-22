package torrent

import (
	"anime-go/models"
	"encoding/xml"
	"io"
	"net/http"
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
		*t = append(*t, models.Torrent{
			Title:   item.Title,
			Link:    item.Link,
			PubDate: item.Torrent.PubDate,
		})

	}
	return nil
}
