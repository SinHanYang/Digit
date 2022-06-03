package env

import (
	"os"
	"path/filepath"
)

func HasDigitDir(path string) bool {
	path = filepath.Join(path, ".digit")

	stat, err := os.Stat(path)
	if err != nil {
		return false
	}
	return stat.IsDir()
}
