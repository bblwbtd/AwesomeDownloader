package handlers

import (
	"AwesomeDownloader/src/config"
	"AwesomeDownloader/src/core"
	"AwesomeDownloader/src/database"
	"AwesomeDownloader/src/database/entities"
	"AwesomeDownloader/src/utils"
	"AwesomeDownloader/src/web/models"
	"os"
	"path"
	"testing"
	"time"
)

var ImageURL = "https://pic.netbian.com/uploads/allimg/170424/104135-14930016950de4.jpg"
var LargeFileURL = "https://npm.taobao.org/mirrors/node/v14.16.0/node-v14.16.0.pkg"
var downloader *core.Downloader

func addTask(url string, filePath string) []*entities.Task {
	return AddTasks(downloader, []*models.TaskMeta{
		{
			URL:  url,
			Path: filePath,
		},
	})
}

func TestAddTask(t *testing.T) {
	filePath := path.Join("temp", "addTask.jpg")
	tasks := addTask(ImageURL, filePath)

	time.Sleep(3 * time.Second)

	if _, err := os.Stat(utils.GetDownloadPath(filePath)); err != nil {
		t.Error(err)
		return
	}

	if tasks[0].Status != core.Finished {
		t.Error("should finished")
	}
}

func TestPauseTask(t *testing.T) {
	filePath := path.Join("temp", "pause.pkg")
	task := addTask(LargeFileURL, filePath)[0]

	time.Sleep(3 * time.Second)

	err := PauseTask(downloader, []uint{task.ID})
	if err != nil {
		t.Error(err)
		return
	}

	database.DB.Take(task)
	if task.Status != core.Paused {
		t.Error("status is not paused")
	}
}

func TestCancelTask(t *testing.T) {
	filePath := path.Join("temp", "cancel.pkg")
	task := addTask(LargeFileURL, filePath)[0]

	time.Sleep(3 * time.Second)

	err := CancelTasks(downloader, []uint{task.ID})
	if err != nil {
		t.Error(err)
		return
	}

	database.DB.Take(task)
	if task.Status != core.Canceled {
		t.Error("status is not paused")
	}
}

func TestRemoveTask(t *testing.T) {
	filePath := path.Join("temp", "cancel.pkg")
	task := addTask(LargeFileURL, filePath)[0]

	time.Sleep(3 * time.Second)

	err := RemoveTasks(downloader, []uint{task.ID})
	if err != nil {
		return
	}

	if err := database.DB.Take(task).Error; err == nil {
		t.Error("task is not removed")
	}
}

func TestUnpauseTask(t *testing.T) {
	filePath := path.Join("temp", "pause.pkg")
	task := addTask(LargeFileURL, filePath)[0]

	time.Sleep(3 * time.Second)

	err := PauseTask(downloader, []uint{task.ID})
	if err != nil {
		t.Error(err)
		return
	}

	database.DB.Take(task)
	if task.Status != core.Paused {
		t.Error("status is not paused")
	}

	err = UnpauseTask(downloader, []uint{task.ID})
	if err != nil {
		t.Error(err)
		return
	}

	time.Sleep(1 * time.Second)

	database.DB.Take(task)
	if task.Status != core.Downloading {
		t.Error("status is not downloading")
	}
}

func TestMain(m *testing.M) {
	_ = os.RemoveAll("temp")
	_ = os.Mkdir("temp", 0777)

	config.InitConfig(path.Join("temp", "config.json"))
	database.InitDB(path.Join("temp", "test.db"))

	downloader = core.NewDownloader()
	code := m.Run()

	_ = os.RemoveAll("temp")
	_ = os.RemoveAll("downloads")

	os.Exit(code)
}
