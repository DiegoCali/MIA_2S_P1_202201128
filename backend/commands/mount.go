package commands

import (
	"backend/structures"
	"backend/utils"
	"strconv"
)

func Mount(path string, name string) (string, error) {
	// read mbr in path
	mbr := &structures.MBR{}
	err := mbr.Deserialize(path)
	if err != nil {
		return "Error: Coudn't read MBR", err
	}
	index, err := mbr.GetPartitionIndex(name)
	if err != nil {
		return "Error: Partition not found", err
	}
	// Generate partition id
	id, err := GenerateId(index, path)
	if err != nil {
		return "Error: Coudn't generate id", err
	}
	// Save partition id in memory
	utils.GlobalMounts[id] = path
	return "Mounted succesfully!!", nil
}

func GenerateId(index int, path string) (string, error) {
	carnet := utils.Carnet
	letter, err := utils.GetLetter(path)
	if err != nil {
		return "", err
	}
	return carnet + strconv.Itoa(index) + letter, nil
}
