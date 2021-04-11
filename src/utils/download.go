package utils

import (
	"AwesomeDownloader/src/config"
	"path"
)

func GetDownloadPath(p string) string {
	cfg := config.GetConfig()

	return path.Join(cfg.DownloadDir, p)
}
