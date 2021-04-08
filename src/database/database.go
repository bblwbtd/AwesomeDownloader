package database

import (
	"AwesomeDownloader/src/entities"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("data.DB"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Panic(err)
	}
	err = DB.AutoMigrate(&entities.DownloadTask{})
	if err != nil {
		log.Panic(err)
	}
}

func RemoveDB() {
	_ = os.Remove("data.DB")
}
