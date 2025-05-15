package utils

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

func SaveFileToDisk(file *multipart.FileHeader, directory string) (res string, err error) {
	if err := os.MkdirAll(directory, os.ModePerm); err != nil {
		return "", err
	}

	ext := filepath.Ext(file.Filename)
	fileName := fmt.Sprintf("%012x%s", time.Now().UnixNano(), ext)
	destPath := filepath.Join(directory, fileName)

	err = saveFile(file, destPath)
	if err != nil {
		return "", err
	}

	return fileName, nil
}

func saveFile(file *multipart.FileHeader, path string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(path)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = dst.ReadFrom(src)
	return err
}
