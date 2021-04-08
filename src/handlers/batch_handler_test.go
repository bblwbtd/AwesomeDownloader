package handlers

import (
	"AwesomeDownloader/src/database"
	"AwesomeDownloader/src/database/entities"
	"AwesomeDownloader/src/models"
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

	if err := database.DB.Take(batch).Error; err != nil {
		t.Error(err)
	}
}

func TestPauseBatch(t *testing.T) {
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
	batch := AddBatch(request)

	PauseBatch(batch.ID)

}

func TestUnPauseBatch(t *testing.T) {

}

func TestRemoveBatch(t *testing.T) {

}
