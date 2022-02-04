package utils

import (
	"AwesomeDownloader/src/config"
	"path"
)

func GetDownloadPath(p ...string) string {
	cfg := config.GetConfig()

	paths := append([]string{cfg.DownloadDir}, p...)

	return path.Join(paths...)
}
