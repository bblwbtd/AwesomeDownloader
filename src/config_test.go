package src

import "testing"

func TestGetConfig(t *testing.T) {
	InitConfig()
	config := GetConfig()
	if config == nil {
		t.Error("no config")
	}
}
