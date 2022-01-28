package handlers

import (
	"AwesomeDownloader/src/database"
	"AwesomeDownloader/src/database/entities"
	"AwesomeDownloader/src/web/models"
	"path"
	"testing"
	"time"
)

func NewBatch() *entities.Batch {
	tasks := []*models.DownloadRequest{
		{
			URL:  LargeFileURL,
			Path: path.Join("temp", "batch1"),
		},
		{
			URL:  LargeFileURL,
			Path: path.Join("temp", "batch2"),
		},
		{
			URL:  LargeFileURL,
			Path: path.Join("temp", "batch3"),
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

	var tasks []entities.Task

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

	time.Sleep(1 * time.Second)

	PauseBatch(batch.ID)

	var tasks []entities.Task
	database.DB.Where("batch = ?", batch.ID).Find(&tasks)

	for _, task := range tasks {
		if task.Status != entities.Paused {
			t.Error("should be paused")
		}
	}
}

func TestUnPauseBatch(t *testing.T) {
	batch := NewBatch()

	time.Sleep(1 * time.Second)

	PauseBatch(batch.ID)

	time.Sleep(1 * time.Second)

	UnPauseBatch(batch.ID)

	var tasks []entities.Task
	database.DB.Where("batch = ?", batch.ID).Find(&tasks)

	for _, task := range tasks {
		if task.Status != entities.Pending {
			t.Error("should be downloading")
		}
	}
}

func TestRemoveBatch(t *testing.T) {
	batch := NewBatch()

	RemoveBatch(batch.ID)

	err := database.DB.Take(batch).Error
	if err == nil {
		t.Error("should be removed")
	}

}
