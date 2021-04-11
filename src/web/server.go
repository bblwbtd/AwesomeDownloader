package web

import (
	"AwesomeDownloader/src/config"
	"AwesomeDownloader/src/database"
	"AwesomeDownloader/src/web/handlers"
	"AwesomeDownloader/src/web/routers"
	"fmt"
	"github.com/gin-gonic/gin"
)

func StartWebServer() error {
	cfg := config.GetConfig()

	database.InitDB("data.db")

	handlers.StartScheduler()

	server := gin.Default()

	routers.MountAPI(server)

	return server.Run(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
}
