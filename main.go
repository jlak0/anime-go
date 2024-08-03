package main

import (
	"anime-go/controller"

	"anime-go/models"
	"anime-go/torrent"
	"fmt"
)

func main() {
	models.Init()
	controller.Test()
}

func CloneMikan() {
	t := []models.Torrent{}
	torrent.GetMikan(&t)
	for _, tx := range t {
		err := tx.Create()
		if err != nil {
			fmt.Println("failed to save")
		}
	}
}

func GetUnreadTorrents() {
	t := models.Torrent{Status: "unread"}
	r := models.DB.Select(&t)
	fmt.Print(r)
}
