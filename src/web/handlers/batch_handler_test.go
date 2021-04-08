package handlers

import (
	"AwesomeDownloader/src/database"
	"AwesomeDownloader/src/database/entities"
	"AwesomeDownloader/src/web/models"
	"path"
	"testing"
)

func NewBatch() *entities.Batch {
	tasks := []*models.DownloadRequest{
		{
			URL:  ImageURL,
			Path: path.Join("temp", "batch1.jpg"),
		},
		{
			URL:  ImageURL,
			Path: path.Join("temp", "batch2.jpg"),
		},
		{
			URL:  ImageURL,
			Path: path.Join("temp", "batch3.jpg"),
		},
	}
	request := &models.BatchRequest{
		Name:  "test",
		Tasks: tasks,
	}
	return AddBatch(request)
}

func TestAddBatch(t *testing.T) {
	batch := NewBatch()

	var tasks []entities.DownloadTask

	err := database.DB.Where("batch = ?", batch.ID).Find(&tasks).Error
	if err != nil {
		t.Error(err)
		return
	}

	if len(tasks) != 3 {
		t.Error("should be equal to 3")
	}
}

func TestPauseBatch(t *testing.T) {
	batch := NewBatch()
	PauseBatch(batch.ID)

	database.DB.Take(batch)

}

func TestUnPauseBatch(t *testing.T) {

}

func TestRemoveBatch(t *testing.T) {

}
