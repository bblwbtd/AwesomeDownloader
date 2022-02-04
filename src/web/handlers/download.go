package handlers

import (
	"AwesomeDownloader/src/core"
	"AwesomeDownloader/src/database/entities"
	"AwesomeDownloader/src/utils"
	"AwesomeDownloader/src/web/models"
	"encoding/json"
	"log"
)

func AddTasks(downloader *core.Downloader, taskMetas []*models.TaskMeta) []*entities.Task {

	tasks := make([]*entities.Task, len(taskMetas))


	for index, meta := range taskMetas {
		headersStr, _ := json.Marshal(meta.Headers)

		task := &entities.Task{
			URL:     meta.URL,
			Path:    utils.GetDownloadPath(meta.Path),
			Status:  core.Pending,
			Headers: string(headersStr),
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
