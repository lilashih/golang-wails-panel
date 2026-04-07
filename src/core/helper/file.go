package helper

import (
	"fmt"
	"gbase/src/core/config"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func IsFileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// get current file name without extension
func GetCurrentFileNameNoExt() string {
	// skip 1 level to get this file
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		return ""
	}
	base := filepath.Base(file)
	return strings.TrimSuffix(base, filepath.Ext(base))
}

func GetStorageDirOrCreate(elem ...string) (string, error) {
	path, err := getStoragePath(elem...)

	if err != nil {
		return "", err
	}

	if err = os.MkdirAll(path, os.ModePerm); err != nil {
		return "", err
	}

	return path, nil
}

func getStoragePath(elem ...string) (string, error) {
	basePath, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("Failed to get the project directory: %v", err)
	}

	basePath = filepath.Join(basePath, config.App.BasePath, "storage")

	return filepath.Join(append([]string{basePath}, elem...)...), nil
}
