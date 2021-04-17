package handlers

import (
	"AwesomeDownloader/src/database"
	"AwesomeDownloader/src/database/entities"
	"AwesomeDownloader/src/downloader"
	"AwesomeDownloader/src/utils"
	"AwesomeDownloader/src/web/models"
	"database/sql"
	"log"
)

func QueryBatchByName(name string) *entities.Batch {
	batch := new(entities.Batch)

	if err := database.DB.First(batch, "name = ?", name).Error; err != nil {
		return nil
	}

	return batch
}

func AddBatch(request *models.BatchRequest) *entities.Batch {

	batch := &entities.Batch{
		Name: request.Name,
	}

	database.DB.Save(batch)

	tasks := make([]*entities.DownloadTask, len(request.Tasks))
	for i, t := range request.Tasks {
		tasks[i] = &entities.DownloadTask{
			URL:    t.URL,
			Path:   utils.GetDownloadPath(t.Path),
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
		downloader.TaskChannel <- task
	}

	return batch
}

func RemoveBatch(id uint) {
	batch := &entities.Batch{}
	database.DB.Take(batch, id)

	var tasks []*entities.DownloadTask
	database.DB.Where("batch = ?", batch.ID).Find(&tasks)

	for _, task := range tasks {
		if task.Status == entities.Downloading {
			cancel(id)
			downloader.DownloadProgress.Delete(task.ID)
		}
	}

	database.DB.Delete(batch)
	database.DB.Delete(tasks)
}

func PauseBatch(id uint) {
	batch := &entities.Batch{}
	database.DB.Take(batch, id)

	var tasks []*entities.DownloadTask
	database.DB.Where("batch = ?", batch.ID).Find(&tasks)
	taskID := make([]uint, len(tasks))

	for index, task := range tasks {
		if task.Status == entities.Downloading {
			cancel(id)
		}
		taskID[index] = task.ID
		task.Status = entities.Paused
		if err := database.DB.Save(task).Error; err != nil {
			log.Println(err)
		}
	}

}

func UnPauseBatch(id uint) {
	batch := &entities.Batch{}
	database.DB.Take(batch, id)

	var tasks []*entities.DownloadTask
	database.DB.Where("batch = ?", batch.ID).Find(&tasks)

	for _, task := range tasks {
		if task.Status == entities.Paused {
			task.Status = entities.Pending
			database.DB.Save(task)
			downloader.TaskChannel <- task
		}
	}
}

func CancelBatch(id uint) {

	var tasks []*entities.DownloadTask
	database.DB.Where("batch = ? and status = 'Downloading'", id).Find(&tasks)

	for _, task := range tasks {
		cancel(task.ID)
	}

	database.DB.Where("batch = ? and status = 'Downloading", id).Update("status", entities.Canceled)
}
