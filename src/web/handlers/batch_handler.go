package handlers

import (
	"AwesomeDownloader/src/database"
	"AwesomeDownloader/src/database/entities"
	"AwesomeDownloader/src/web/models"
)

func AddBatch(request *models.BatchRequest) *entities.Batch {

	tasks := make([]entities.DownloadTask, len(request.Tasks))
	for i, t := range request.Tasks {
		tasks[i] = entities.DownloadTask{
			URL:    t.URL,
			Path:   t.Path,
			Size:   0,
			Status: entities.Pending,
		}
	}

	batch := &entities.Batch{
		Name:  request.Name,
		Tasks: tasks,
	}

	database.DB.Create(batch)

	for _, task := range tasks {
		taskChannel <- &task
	}

	return batch
}

func RemoveBatch(id uint) {
	batch := &entities.Batch{}
	database.DB.Take(batch, id)

	for _, task := range batch.Tasks {
		if task.Status == entities.Downloading {
			cancel(id)
			downloadProgress.Delete(task.ID)
		}
	}

	database.DB.Delete(batch)
}

func PauseBatch(id uint) {
	batch := &entities.Batch{}
	database.DB.Take(batch, id)

	for _, task := range batch.Tasks {
		if task.Status == entities.Downloading {
			cancel(id)
		}

		task.Status = entities.Paused
		database.DB.Save(task)
	}
}

func UnPauseBatch(id uint) {
	batch := &entities.Batch{}
	database.DB.Take(batch, id)

	for _, task := range batch.Tasks {
		if task.Status == entities.Paused {
			task.Status = entities.Pending
			taskChannel <- &task
		}
	}
}
