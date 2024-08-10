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
	// Today's date
	timeFloat := float32(time.Now().Unix())
	// Create MBR
	mbr := &structures.MBR{}
	err = mbr.Set(int32(sizeBytes), timeFloat, rand.Int31(), fit)
	if err != nil {
		return "Error: Coudn't create disk", err
	}
	// Create disk
	err = utils.CreateDisk(path, sizeBytes)
	if err != nil {
		return "Error: Coudn't create disk", err
	}
	// Serialize MBR
	err = mbr.Serialize(path)
	if err != nil {
		return "Error: Coudn't serialize MBR", err
	}
	return "Disk created succesfully", nil
}

func RmDisk(path string) (string, error) {
	err := os.Remove(path)
	if err != nil {
		return "Error: Coudn't remove disk", err
	}
	return "Disk removed succesfully", nil
}

func FDisk(size int, unit string, path string, typeP string, fit string, name string) (string, error) {

	return "Partition createed succesfully", nil
}
