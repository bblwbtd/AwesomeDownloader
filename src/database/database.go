package database

import (
	"AwesomeDownloader/src/database/entities"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

var DB *gorm.DB

func InitDB(dsn string) {
	var err error
	DB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Panic(err)
	}
	err = DB.AutoMigrate(&entities.DownloadTask{}, &entities.Batch{})
	if err != nil {
		log.Panic(err)
	}
}
