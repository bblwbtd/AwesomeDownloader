package handlers

import (
	"AwesomeDownloader/src/models"
	"os"
	"path"
	"testing"
	"time"
)

func TestStartScheduler(t *testing.T) {
	StartScheduler()
	filePath := path.Join("temp", "scheduler.jpg")
	AddTask(&models.DownloadRequest{
		URL:  "https://pic.netbian.com/uploads/allimg/170424/104135-14930016950de4.jpg",
		Path: filePath,
	})

	time.Sleep(3 * time.Second)

	if _, err := os.Stat(filePath); err != nil {
		t.Error(err)
	}
}
