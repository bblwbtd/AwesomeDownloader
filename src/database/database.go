package database

import (
	"AwesomeDownloader/src/database/entities"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("data.db"), &gorm.Config{
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

func RemoveDB() {
	err := os.Remove("data.db")
	if err != nil {
		log.Println(err)
	}
}
