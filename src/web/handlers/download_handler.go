package handlers

import (
	"AwesomeDownloader/src/core"
	"AwesomeDownloader/src/database"
	"AwesomeDownloader/src/database/entities"
	"AwesomeDownloader/src/utils"
	"AwesomeDownloader/src/web/models"
	"context"
)

func AddTask(request *models.DownloadRequest) *entities.Task {

	task := &entities.Task{
		URL:    request.URL,
		Path:   utils.GetDownloadPath(request.Path),
		Status: entities.Pending,
	}
	database.DB.Create(task)
	go func() {
		core.TaskChannel <- task
	}()

	return task
}

func RemoveTask(id uint) {
	core.DownloadProgress.Delete(id)
	cancel(id)
	database.DB.Delete(&entities.Task{}, id)
}

func PauseTask(id uint) {
	cancel(id)
	database.DB.Model(&entities.Task{}).Where("id = ?", id).Update("status", entities.Paused)
}

func UnPauseTask(id uint) {
	task := new(entities.Task)
	err := database.DB.Take(task, id).Error
	if err != nil {
		return
	}

	task.Status = entities.Pending
	database.DB.Save(task)

	go func() {
		core.TaskChannel <- task
	}()
}

func CancelTask(id uint) {
	cancel(id)

	database.DB.Model(&entities.Task{}).Where("id = ?", id).Update("status", entities.Canceled)
}

func cancel(id uint) {
	if cancel, loaded := core.Cancellations.LoadAndDelete(id); loaded {
		if cancelFun, ok := cancel.(context.CancelFunc); ok {
			cancelFun()
		}
	}
}
