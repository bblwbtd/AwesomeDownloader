package routers

import (
	"AwesomeDownloader/src/utils"
	"AwesomeDownloader/src/web/handlers"
	"AwesomeDownloader/src/web/models"
	"github.com/gin-gonic/gin"
)

func mountDownload(router *gin.Engine) {
	task := router.Group("/task")

	task.POST("/add", addTask)
}

func addTask(ctx *gin.Context) {
	req := new(models.DownloadRequest)
	err := ctx.BindJSON(req)
	if err != nil {
		utils.RespondError(ctx, utils.InvalidBody, "")
		return
	}

	task := handlers.AddTask(req)

	utils.RespondSuccess(ctx, task)
}

func removeTask(ctx *gin.Context) {

}

func pauseTask(ctx *gin.Context) {

}

func unPauseTask(ctx *gin.Context) {

}

func cancelTask(ctx *gin.Context) {

}
