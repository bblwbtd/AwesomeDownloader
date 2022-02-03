package handlers

import (
	"AwesomeDownloader/src/core"
	"AwesomeDownloader/src/database/entities"
	"AwesomeDownloader/src/utils"
	"AwesomeDownloader/src/web/models"
)

func AddTask(request *models.DownloadRequest) *entities.Task {

	task := &entities.Task{
		URL:    request.URL,
		Path:   utils.GetDownloadPath(request.Path),
		Status: core.Pending,
	}

	downloader.Enqueue(task)

	return task
}

func RemoveTask(id uint) error {
	return downloader.DeleteTasks([]uint{id})
}

func PauseTask(id uint) error {
	return downloader.PauseTasks([]uint{id})
}

func UnPauseTask(id uint) error {
	return downloader.UnPauseTasks([]uint{id})
}

func CancelTask(id uint) error {
	return downloader.CancelTasks([]uint{id})
}
