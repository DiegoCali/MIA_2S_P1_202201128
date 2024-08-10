package commands

import (
	"backend/structures"
	utils "backend/utils"
	"math/rand"
	"os"
	"time"
)

func MkDisk(size int, fit string, unit string, path string) (string, error) {
	sizeBytes, err := utils.ConvertToBytes(size, unit)
	if err != nil {
		return "Error: Coudn't create disk", err
	}
	// Create file in path
	file, err := os.Create(path)
	if err != nil {
		return "Error: Coudn't create disk", err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			return
		}
	}(file)
	// Today's date
	timeFloat := utils.TimeToFloat(time.Now())
	// Create MBR
	mbr := &structures.MBR{}
	err = mbr.Set(sizeBytes, timeFloat, rand.Int31(), fit)
	if err != nil {
		return "", err
	}
	// Write MBR in file
	return "Disk created succesfully", nil
}

func RmDisk(path string) (string, error) {
	return "Disk removed succesfully", nil
}
