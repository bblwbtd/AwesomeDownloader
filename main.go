package main

import (
	"AwesomeDownloader/src/config"
	"AwesomeDownloader/src/database"
	"AwesomeDownloader/src/web"
	"fmt"
)

func main() {
	config.InitConfig("config.json")
	database.InitDB("data.db")

	err := web.Start()
	if err != nil {
		fmt.Print(err)
	}
}
