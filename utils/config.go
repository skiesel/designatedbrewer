package config

import (
	"encoding/json"
	"os"
)

var (
	settings = map[string]string{}
)

func init() {
	fileReader, err := os.Open("settings.config")

	decoder := json.NewDecoder(fileReader)
	err = decoder.Decode(&settings)
	if err != nil {
		panic(err)
	}
}

func GetSetting(key string) string {
	return settings[key]
}