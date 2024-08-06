package interpreter

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type MBR struct {
	Size       [4]byte
	TimeStamp  [4]byte
	Signature  [4]byte
	Fit        [1]byte
	Partitions [4]Partition
}

type Partition struct {
	Status      [1]byte
	Type        [1]byte
	Fit         [1]byte
	Start       [4]byte
	Size        [4]byte
	Name        [16]byte
	Correlative [4]byte
	Id          [4]byte
}

func MkDisk(options []Option) (string, error) {
	var message string
	size := -1
	fit := "none"
	unit := "none"
	path := "none"
	for _, option := range options {
		if option.Name == "size" {
			//Parse string to int
			size, _ = strconv.Atoi(option.Value)
			continue
		}
		if option.Name == "fit" {
			fit = option.Value
			continue
		}
		if option.Name == "unit" {
			unit = option.Value
			continue
		}
		if option.Name == "path" {
			path = option.Value
			continue
		}
	}
	if size != -1 {
		err := createDisk(size, fit, unit, path)
		if err != nil {
			return "ERROR: Disk not created", err
		}
		message = "Disk created successfully, size: " + strconv.Itoa(size) + ", fit: " + fit + ", unit: " + unit + ", path: " + path
	} else {
		message = "ERROR: Disk not created"
		return message, fmt.Errorf("-size is required")
	}
	return message, nil
}

func createDisk(size int, fit string, unit string, path string) error {
	var mbr MBR
	if fit == "none" {
		fit = "F"
	} else {
		// Get only first char
		fit = fit[:1]
	}
	if unit == "none" {
		unit = "M"
	}
	if path == "none" {
		// Create file in /home/diego/Documents/Archivos_2024/MIA_2S_P1_202201128/backend/disks
		path = "/home/diego/Documents/Archivos_2024/MIA_2S_P1_202201128/backend/disks/disk" + strconv.Itoa(rand.Intn(1000)) + ".mia"
	}
	sizeBytes, err := convertToBytes(size, unit)
	if err != nil {
		return err
	}
	// Get current time
	timeNow := time.Now()
	// Parse time to float
	timeFloat := timeToFloat(timeNow)
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Error closing file")
		}
	}(file)
	copy(mbr.Size[:], strconv.Itoa(size))
	copy(mbr.TimeStamp[:], strconv.Itoa(int(timeFloat)))
	copy(mbr.Signature[:], strconv.Itoa(rand.Intn(1000)))
	copy(mbr.Fit[:], fit)
	// Write MBR to file
	err = writeFile(file, mbr, sizeBytes)
	if err != nil {
		return err
	}
	return nil
}

func convertToBytes(size int, unit string) (int, error) {
	if unit == "k" || unit == "K" {
		size = size * 1024
	} else if unit == "m" || unit == "M" {
		size = size * 1024 * 1024
	} else {
		return -1, fmt.Errorf("unit %s not recognized", unit)
	}
	return size, nil
}

func writeFile(file *os.File, mbr MBR, size int) error {
	// Fill file with 0s
	buffer := make([]byte, 1024*1024)
	for size > 0 {
		writeSize := len(buffer)
		if size < writeSize {
			writeSize = size
		}
		if _, err := file.Write(buffer[:writeSize]); err != nil {
			return err
		}
		size -= writeSize
	}
	fmt.Println("[File created successfully] in " + file.Name())
	// Write MBR to file
	_, err := file.Seek(0, 0)
	if err != nil {
		return err
	}
	err = binary.Write(file, binary.LittleEndian, mbr)
	if err != nil {
		return err
	}
	return nil
}

func timeToFloat(time time.Time) float64 {
	year := time.Year()
	month := monthToInt(time.Month())
	day := time.Day()
	hour := time.Hour()
	minute := time.Minute()
	second := time.Second()
	return float64(year*10000000000 + month*100000000 + day*1000000 + hour*10000 + minute*100 + second)
}

func monthToInt(month time.Month) int {
	switch month {
	case time.January:
		return 1
	case time.February:
		return 2
	case time.March:
		return 3
	case time.April:
		return 4
	case time.May:
		return 5
	case time.June:
		return 6
	case time.July:
		return 7
	case time.August:
		return 8
	case time.September:
		return 9
	case time.October:
		return 10
	case time.November:
		return 11
	case time.December:
		return 12
	}
	return 0
}
