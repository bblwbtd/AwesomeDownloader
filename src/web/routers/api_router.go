package routers

import (
	"AwesomeDownloader/src/core"
	"github.com/gin-gonic/gin"
)

var downloader *core.Downloader

func MountAPI(router *gin.Engine) {
	api := router.Group("/api")

	downloader = core.NewDownloader()

	mountTaskRouter(api)
}
