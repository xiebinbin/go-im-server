package filehelper

import (
	"os"
	"path"
)

func GetAppPath(dir string) string {
	APP_PATH, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return path.Join(APP_PATH, dir)
}
