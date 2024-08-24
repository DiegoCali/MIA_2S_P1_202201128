package utils

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
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
	if _, exists := pathToLetter[path]; !exists {
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

func Serialize[T any](data *T, path string, offset int) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			return
		}
	}(file)
	// Seek to offset
	_, err = file.Seek(int64(offset), 0)
	if err != nil {
		return err
	}
	// Write data
	err = binary.Write(file, binary.LittleEndian, data)
	if err != nil {
		return err
	}
	return nil
}

func Deserialize[T any](data *T, path string, offset int) error {
	file, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			return
		}
	}(file)
	// Seek to offset
	_, err = file.Seek(int64(offset), 0)
	if err != nil {
		return err
	}
	// Get size of data
	size := binary.Size(data)
	// Create buffer
	buffer := make([]byte, size)
	// Read data
	_, err = file.Read(buffer)
	if err != nil {
		return err
	}
	// Create buffer reader
	reader := bytes.NewReader(buffer)
	// Deserialize data
	err = binary.Read(reader, binary.LittleEndian, data)
	if err != nil {
		return err
	}
	return nil
}

func PrintStruct[T any](data *T) error {
	// Convert struct to json
	jsonData, err := json.MarshalIndent(data, "", "   ")
	if err != nil {
		return err
	}
	// Print json
	fmt.Println(string(jsonData))
	return nil
}
