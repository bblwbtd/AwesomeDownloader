package main

import (
	"AwesomeDownloader/src/web"
	"fmt"
)

func main() {
	err := web.StartWebServer()
	if err != nil {
		fmt.Print(err)
	}
}
