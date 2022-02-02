package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
)

type Proxy struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type Config struct {
	MaxConnections int    `json:"max_connections"`
	MaxRetry       int    `json:"max_retry"`
	DownloadDir    string `json:"download_dir"`
	Port           int    `json:"port"`
	Host           string `json:"host"`
	Proxy          *Proxy `json:"proxy"`
}

var defaultConfig = &Config{
	MaxConnections: 10,
	DownloadDir:    "downloads",
	Port:           1234,
	Host:           "0.0.0.0",
	MaxRetry:       3,
}

var config *Config
var configPath = path.Join("config", "config.json")

func GetConfig() *Config {
	return config
}

func InitConfig() {
	_, err := os.Stat(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			if err = os.MkdirAll(path.Dir(configPath), 0777); err != nil {
				log.Panic(err)
			}
			bytes, _ := json.Marshal(defaultConfig)
			if err = ioutil.WriteFile(configPath, bytes, 0777); err != nil {
				log.Panic(err)
			}
		} else {
			log.Panic(err)
		}
		config = defaultConfig
	} else {
		bytes, err := ioutil.ReadFile(configPath)
		if err != nil {
			log.Panic(err)
		}
		temp := new(Config)
		if err = json.Unmarshal(bytes, temp); err != nil {
			log.Panic(err)
		}
		config = temp
	}
}

func RemoveConfig() {
	_ = os.RemoveAll("config")
}
