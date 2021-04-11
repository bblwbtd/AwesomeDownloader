package web

import (
	"AwesomeDownloader/src/config"
	"AwesomeDownloader/src/web/routers"
	"fmt"
	"github.com/gin-gonic/gin"
)

func StartWebServer() error {
	cfg := config.GetConfig()
	server := gin.Default()

	routers.MountAPI(server)

	return server.Run(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
}
