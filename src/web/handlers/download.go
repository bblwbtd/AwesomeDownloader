package handlers

import (
	"AwesomeDownloader/src/core"
	"AwesomeDownloader/src/database/entities"
	"AwesomeDownloader/src/web/models"
	"log"
)

func AddTasks(downloader *core.Downloader, taskMetas []*models.TaskMeta) []*entities.Task {

	tasks := make([]*entities.Task, len(taskMetas))

	for index, meta := range taskMetas {
		task := &entities.Task{
			URL:    meta.URL,
			Path:   meta.Path,
			Status: core.Pending,
		}
		tasks[index] = task
		err := downloader.CreateAndEnqueue(task)
		if err != nil {
			log.Println("Error occur while add tasks:", err)
		}
	}

	return tasks
}

func RemoveTasks(downloader *core.Downloader, id []uint) error {
	return downloader.DeleteTasks(id)
}

func PauseTask(downloader *core.Downloader, id []uint) error {
	return downloader.PauseTasks(id)
}

func UnpauseTask(downloader *core.Downloader, id []uint) error {
	return downloader.UnpauseTasks(id)
}

func CancelTasks(downloader *core.Downloader, id []uint) error {
	return downloader.CancelTasks(id)
}
