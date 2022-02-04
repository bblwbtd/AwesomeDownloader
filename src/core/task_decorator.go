package core

import (
	"AwesomeDownloader/src/database"
	"AwesomeDownloader/src/database/entities"
	"context"
	"log"
	"os"
)

type TaskDecorator struct {
	entity         *entities.Task
	cancel         context.CancelFunc
	ctx            context.Context
	downloadedSize uint64
	retryCount     int
}

func NewDecoratedTask(entity *entities.Task) *TaskDecorator {
	ctx, cancel := context.WithCancel(context.TODO())

	return &TaskDecorator{entity: entity, cancel: cancel, ctx: ctx}
}

func (d *TaskDecorator) GetDownloadedSize() uint64 {
	return d.downloadedSize
}

func (d *TaskDecorator) SetDownloadedSize(downloadedSize uint64) {
	d.downloadedSize = downloadedSize
}

func (d *TaskDecorator) GetTaskStatus() (DownloadStatus, error) {
	if tx := database.DB.Take(d.entity); tx.Error != nil {
		return Unknown, tx.Error
	}

	return d.entity.Status, nil
}

func (d *TaskDecorator) SetTaskStatus(status DownloadStatus) error {
	d.entity.Status = status
	if tx := database.DB.Save(d.entity); tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (d *TaskDecorator) SetTaskSize(size uint64) error {
	d.entity.Size = size
	if tx := database.DB.Save(d.entity); tx.Error != nil {
		return tx.Error
	}

	return nil

}

func (d *TaskDecorator) Cancel() {
	if d.cancel != nil {
		d.cancel()
	}

	err := os.Remove(d.entity.Path)
	if err != nil {
		log.Println("Error occur while removing:", d.entity.Path)
		return
	}
}
