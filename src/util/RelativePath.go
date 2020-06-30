package util

import "os"

func RelativePath() string {
	path, _ := os.Getwd()
	return path
}
