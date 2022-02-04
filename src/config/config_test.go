package config

import (
	"os"
	"path"
	"testing"
)

func TestGetConfig(t *testing.T) {
	_ = os.RemoveAll("temp")
	_ = os.Mkdir("temp", 0777)

	InitConfig(path.Join("temp", "config.json"))
	config := GetConfig()
	if config == nil {
		t.Error("no config")
	}

	_ = os.RemoveAll("temp")
}
