package main

import (
	"AwesomeDownloader/src/config"
	"AwesomeDownloader/src/web"
	"fmt"
)

func main() {
	config.InitConfig()

	err := web.StartWebServer()
	if err != nil {
		fmt.Print(err)
	}
}
