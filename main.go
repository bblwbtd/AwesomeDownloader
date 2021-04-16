package main

import (
	"AwesomeDownloader/src/config"
	"AwesomeDownloader/src/database"
	"AwesomeDownloader/src/downloader"
	"AwesomeDownloader/src/web"
	"fmt"
)

func main() {
	config.InitConfig()
	database.InitDB("data.db")

	downloader.StartScheduler()

	err := web.StartWebServer()
	if err != nil {
		fmt.Print(err)
	}
}
