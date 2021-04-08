package handlers

import (
	"AwesomeDownloader/src/database"
	"AwesomeDownloader/src/database/entities"
	"AwesomeDownloader/src/web/models"
	"database/sql"
)

func AddBatch(request *models.BatchRequest) *entities.Batch {

	batch := &entities.Batch{
		Name: request.Name,
	}

	database.DB.Save(batch)

	tasks := make([]*entities.DownloadTask, len(request.Tasks))
	for i, t := range request.Tasks {
		tasks[i] = &entities.DownloadTask{
			URL:    t.URL,
			Path:   t.Path,
			Size:   0,
			Status: entities.Pending,
			Batch: sql.NullInt64{
				Int64: int64(batch.ID),
				Valid: true,
			},
		}
	}
	database.DB.Save(tasks)

	for _, task := range tasks {
		taskChannel <- task
	}

	return batch
}

func RemoveBatch(id uint) {
	batch := &entities.Batch{}
	database.DB.Take(batch, id)

	var tasks []entities.DownloadTask
	database.DB.Where("batch = ?", batch.ID).Find(&tasks)

	for _, task := range tasks {
		if task.Status == entities.Downloading {
			cancel(id)
			downloadProgress.Delete(task.ID)
		}
	}

	database.DB.Delete(batch)
	database.DB.Delete(tasks)
}

func PauseBatch(id uint) {
	batch := &entities.Batch{}
	database.DB.Take(batch, id)

	var tasks []entities.DownloadTask
	database.DB.Where("batch = ?", batch.ID).Find(&tasks)

	for _, task := range tasks {
		if task.Status == entities.Downloading {
			cancel(id)
		}

		task.Status = entities.Paused
	}

	database.DB.Save(tasks)
}

func UnPauseBatch(id uint) {
	batch := &entities.Batch{}
	database.DB.Take(batch, id)

	var tasks []entities.DownloadTask
	database.DB.Where("batch = ?", batch.ID).Find(&tasks)

	for _, task := range tasks {
		if task.Status == entities.Paused {
			task.Status = entities.Pending
			taskChannel <- &task
		}
	}
}
