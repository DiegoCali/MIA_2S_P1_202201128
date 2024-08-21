package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
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

func Int64ToDate(timestamp int64) string {
	t := time.Unix(timestamp, 0)
	return t.Format("2006-01-02 15:04:05")
}

const Carnet string = "28"

var alphabet = []string{
	"A", "B", "C", "D", "E", "F", "G", "H", "I", "J",
	"K", "L", "M", "N", "O", "P", "Q", "R", "S", "T",
	"U", "V", "W", "X", "Y", "Z",
}

var pathToLetter = make(map[string]string)

var nextLetter = 0

func GetLetter(path string) (string, error) {
	if _, exists := pathToLetter[path]; exists {
		if nextLetter < len(alphabet) {
			pathToLetter[path] = alphabet[nextLetter]
			nextLetter++
		} else {
			return "", fmt.Errorf("error: maximum number of partitions reached")
		}
	}
	return pathToLetter[path], nil
}

func CheckNull(str []byte) string {
	//
	output := ""
	for i := 0; i < len(str); i++ {
		if str[i] == 0 {
			break
		}
		output += string(str[i])
	}
	return output
}
