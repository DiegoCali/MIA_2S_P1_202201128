package interpreter

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func MkDisk(options []Option) (string, error) {
	var message string
	size := -1
	fit := "FF"
	unit := "M"
	path := "/home/diego/Documents/Archivos_2024/MIA_2S_P1_202201128/backend/disks/disk" + strconv.Itoa(rand.Intn(1000)) + ".mia"
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

func RmDisk(options []Option) (string, error) {
	var message string
	path := "none"
	for _, option := range options {
		if option.Name == "path" {
			path = option.Value
			break
		}
	}
	if path == "none" {
		return "ERROR: Disk not removed", fmt.Errorf("-path is required")
	}
	err := os.Remove(path)
	if err != nil {
		return "ERROR: Disk not removed", err
	}
	message = "Disk removed successfully, path: " + path
	return message, nil
}
func FDisk(options []Option) (string, error) {
	var message string
	size := -1
	unit := "K"
	path := "none"
	typePartition := "P"
	fit := "WF"
	name := "none"
	for _, option := range options {
		if option.Name == "size" {
			//Parse string to int
			size, _ = strconv.Atoi(option.Value)
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
		if option.Name == "type" {
			typePartition = option.Value
			continue
		}
		if option.Name == "fit" {
			fit = option.Value
			continue
		}
		if option.Name == "name" {
			name = option.Value
			continue
		}
	}
	if size == -1 {
		return "ERROR: Partition not created", fmt.Errorf("-size is required")
	}
	if path == "none" {
		return "ERROR: Partition not created", fmt.Errorf("-path is required")
	}
	if name == "none" {
		return "ERROR: Partition not created", fmt.Errorf("-name is required")
	}
	err := createPartition(size, unit, path, typePartition, fit, name)
	if err != nil {
		return "ERROR: Partition not created", err
	}
	return message, nil
}

func createDisk(size int, fit string, unit string, path string) error {
	var mbr MBR
	// Get only first char
	fit = fit[:1]
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
	// Fill MBR
	binary.LittleEndian.PutUint32(mbr.Size[:], uint32(sizeBytes))
	binary.LittleEndian.PutUint32(mbr.TimeStamp[:], uint32(timeFloat))
	binary.LittleEndian.PutUint32(mbr.Signature[:], rand.Uint32())
	copy(mbr.Fit[:], fit)
	// Fill partitions
	for i := range mbr.Partitions {
		copy(mbr.Partitions[i].Status[:], "0")
		copy(mbr.Partitions[i].Type[:], "N")
		copy(mbr.Partitions[i].Fit[:], "N")
		copy(mbr.Partitions[i].Start[:], "NULL")
		copy(mbr.Partitions[i].Size[:], "NULL")
		copy(mbr.Partitions[i].Name[:], "----------------")
		copy(mbr.Partitions[i].Correlative[:], "NULL")
		copy(mbr.Partitions[i].Id[:], "NULL")
	}
	err = writeFile(file, mbr, sizeBytes)
	if err != nil {
		return err
	}
	return nil
}

func convertToBytes(size int, unit string) (int, error) {
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

func readMBR(path string) (MBR, error) {
	var mbr MBR
	file, err := os.Open(path)
	if err != nil {
		return mbr, err
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		return mbr, err
	}
	err = binary.Read(file, binary.LittleEndian, &mbr)
	if err != nil {
		return mbr, err
	}
	fmt.Println("MBR read successfully")
	// Read MBR data
	fmt.Println("Size: ", binary.LittleEndian.Uint32(mbr.Size[:]), " bytes")
	fmt.Println("TimeStamp: ", binary.LittleEndian.Uint32(mbr.TimeStamp[:]))
	fmt.Println("Signature: ", binary.LittleEndian.Uint32(mbr.Signature[:]))
	fmt.Println("Fit: ", string(mbr.Fit[:]))
	// Read partitions
	for i, partition := range mbr.Partitions {
		fmt.Println("Partition ", i)
		fmt.Println("Status: ", string(partition.Status[:]))
		fmt.Println("Type: ", string(partition.Type[:]))
		fmt.Println("Fit: ", string(partition.Fit[:]))
		fmt.Println("Start: ", string(partition.Start[:]))
		fmt.Println("Size: ", string(partition.Size[:]))
		fmt.Println("Name: ", string(partition.Name[:]))
		fmt.Println("Correlative: ", string(partition.Correlative[:]))
		fmt.Println("Id: ", string(partition.Id[:]))
	}
	return mbr, nil
}

func timeToFloat(time time.Time) float64 {
	year := time.Year()
	month := monthToInt(time.Month())
	day := time.Day()
	hour := time.Hour()
	minute := time.Minute()
	second := time.Second()
	dateFloat := float64(year*10000 + month*100 + day + hour/100 + minute/10000 + second/1000000)
	return dateFloat
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
