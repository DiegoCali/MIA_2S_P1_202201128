package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func ConvertToBytes(size int, unit string) (int, error) {
	if unit == "b" || unit == "B" {
		return size, nil
	} else if unit == "k" || unit == "K" {
		size = size * 1024
	} else if unit == "m" || unit == "M" {
		size = size * 1024 * 1024
	} else {
		return -1, fmt.Errorf("unit %s not recognized", unit)
	}
	return size, nil
}

func CreateDisk(path string, sizeBytes int) error {
	err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			return
		}
	}(file)

	buffer := make([]byte, 1024*1024)
	for sizeBytes > 0 {
		writeSize := len(buffer)
		if sizeBytes < writeSize {
			writeSize = sizeBytes
		}
		if _, err := file.Write(buffer[:writeSize]); err != nil {
			return err
		}
		sizeBytes -= writeSize
	}
	return nil
}
