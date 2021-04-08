package handlers

import (
	"AwesomeDownloader/src/database"
	"AwesomeDownloader/src/entities"
	"AwesomeDownloader/src/models"
	"os"
	"path"
	"testing"
	"time"
)

var ImageURL = "https://pic.netbian.com/uploads/allimg/170424/104135-14930016950de4.jpg"
var LargeFileURL = "https://npm.taobao.org/mirrors/node/v14.16.0/node-v14.16.0.pkg"

func TestStartScheduler(t *testing.T) {
	filePath := path.Join("temp", "scheduler.jpg")
	AddTask(&models.DownloadRequest{
		URL:  ImageURL,
		Path: filePath,
	})

	time.Sleep(3 * time.Second)

	if _, err := os.Stat(filePath); err != nil {
		t.Error(err)
	}
}

func TestAddTask(t *testing.T) {
	filePath := path.Join("temp", "addTask.jpg")
	task := AddTask(&models.DownloadRequest{
		URL:  ImageURL,
		Path: filePath,
	})

	time.Sleep(3 * time.Second)

	if _, err := os.Stat(filePath); err != nil {
		t.Error(err)
		return
	}

	if task.Status != entities.Finished {
		t.Error("should finished")
	}
}

func TestPauseTask(t *testing.T) {
	filePath := path.Join("temp", "pause.pkg")
	task := AddTask(&models.DownloadRequest{
		URL:  LargeFileURL,
		Path: filePath,
	})

	time.Sleep(3 * time.Second)

	PauseTask(task.ID)

	database.DB.Take(task)
	if task.Status != entities.Paused {
		t.Error("status is not paused")
	}
}

func TestCancelTask(t *testing.T) {
	filePath := path.Join("temp", "cancel.pkg")
	task := AddTask(&models.DownloadRequest{
		URL:  LargeFileURL,
		Path: filePath,
	})

	time.Sleep(3 * time.Second)

	CancelTask(int(task.ID))

	database.DB.Take(task)
	if task.Status != entities.Canceled {
		t.Error("status is not paused")
	}
}

func TestRemoveTask(t *testing.T) {
	filePath := path.Join("temp", "cancel.pkg")
	task := AddTask(&models.DownloadRequest{
		URL:  LargeFileURL,
		Path: filePath,
	})

	time.Sleep(3 * time.Second)

	RemoveTask(task.ID)

	if err := database.DB.Take(task).Error; err == nil {
		t.Error("task is not removed")
	}
}

func TestUnPauseTask(t *testing.T) {
	filePath := path.Join("temp", "pause.pkg")
	task := AddTask(&models.DownloadRequest{
		URL:  LargeFileURL,
		Path: filePath,
	})

	time.Sleep(3 * time.Second)

	PauseTask(task.ID)

	database.DB.Take(task)
	if task.Status != entities.Paused {
		t.Error("status is not paused")
	}

	UnPauseTask(task.ID)

	time.Sleep(1 * time.Second)

	database.DB.Take(task)
	if task.Status != entities.Downloading {
		t.Error("status is not downloading")
	}
}
