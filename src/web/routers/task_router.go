package routers

import (
	"AwesomeDownloader/src/utils"
	"AwesomeDownloader/src/web/handlers"
	"AwesomeDownloader/src/web/models"
	"github.com/gin-gonic/gin"
)

func mountDownload(router *gin.RouterGroup) {
	task := router.Group("/task")

	task.POST("/add", addTask)
	task.POST("/remove/:id", removeTask)
	task.POST("/pause/:id", pauseTask)
	task.POST("/unpause/:id", unpauseTask)
	task.POST("/cancel/:id", cancelTask)
}

func addTask(ctx *gin.Context) {
	req := new(models.DownloadRequest)
	err := ctx.BindJSON(req)
	if err != nil {
		utils.RespondError(ctx, utils.InvalidBody)
		return
	}

	task := handlers.AddTask(req)

	utils.RespondSuccess(ctx, task)
}

func removeTask(ctx *gin.Context) {
	id, err := utils.ExtractID(ctx)
	if err != nil {
		utils.RespondError(ctx, utils.InvalidID)
		return
	}

	handlers.RemoveTask(id)

	utils.RespondSuccess(ctx, "")
}

func pauseTask(ctx *gin.Context) {
	id, err := utils.ExtractID(ctx)
	if err != nil {
		utils.RespondError(ctx, utils.InvalidID)
		return
	}

	handlers.PauseTask(id)

	utils.RespondSuccess(ctx, "")
}

func unpauseTask(ctx *gin.Context) {
	id, err := utils.ExtractID(ctx)
	if err != nil {
		utils.RespondError(ctx, utils.InvalidID)
		return
	}

	handlers.UnPauseBatch(id)

	utils.RespondSuccess(ctx, "")
}

func cancelTask(ctx *gin.Context) {
	id, err := utils.ExtractID(ctx)
	if err != nil {
		utils.RespondError(ctx, utils.InvalidID)
		return
	}

	handlers.CancelTask(id)

	utils.RespondSuccess(ctx, "")
}
