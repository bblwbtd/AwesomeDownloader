package core

import (
	"AwesomeDownloader/src/database"
	"AwesomeDownloader/src/database/entities"
	"context"
)

type TaskDecorator struct {
	entity *entities.Task
	Cancel context.CancelFunc
}

func NewDecoratedTask(entity *entities.Task) *TaskDecorator {
	return &TaskDecorator{entity: entity}
}

func (d *TaskDecorator) getTaskStatus() (DownloadStatus, error) {
	if tx := database.DB.Take(d.entity); tx.Error != nil {
		return Unknown, tx.Error
	}

	return d.entity.Status, nil
}

func (d *TaskDecorator) setTaskStatus(status DownloadStatus) error {
	d.entity.Status = status
	if tx := database.DB.Save(d.entity.Status); tx.Error != nil {
		return tx.Error
	}

	return nil
}
