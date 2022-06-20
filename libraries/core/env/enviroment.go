package env

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Config struct {
	DB_NAME string `json:"DB_NAME"`
	DB_USER string `json:"DB_USER"`
	DB_PASS string `json:"DB_PASS"`
}

func HasDigitDir(path string) bool {
	path = filepath.Join(path, ".digit")

	stat, err := os.Stat(path)
	if err != nil {
		return false
	}
	return stat.IsDir()
}

// Get config data in .digit/config.json
func GetConfig(path string) (string, string, string, error) {
	path = filepath.Join(path, ".digit", "config.json")
	jsonFile, _ := ioutil.ReadFile(path)
	var config Config
	err := json.Unmarshal(jsonFile, &config)

	if err != nil {
		return "", "", "", err
	}
	return config.DB_NAME, config.DB_USER, config.DB_PASS, nil
}
