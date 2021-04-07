package database

import (
	"AwesomeDownloader/src/entities"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("data.DB"))
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
