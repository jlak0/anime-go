package main

import (
	"anime-go/pkg/logger"
	"anime-go/pkg/torrent"
	"fmt"
)

func main() {
	fmt.Println("开始")
	logger.Log("程序启动")
	defer logger.Close()
	// controller.GetBgmID()
	// cronjobs.StartCronJobs()
	torrent.Test()
	// api.Serve()
}
