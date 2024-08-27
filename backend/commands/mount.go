package commands

import (
	"backend/structures"
	"backend/utils"
	"fmt"
	"strconv"
)

func Mount(path string, name string) (string, error) {
	// read mbr in path
	mbr := &structures.MBR{}
	err := utils.Deserialize(mbr, path, 0)
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
	// Get partition correlative from GlobalMounts
	correlative := len(utils.GlobalMounts)
	// Save partition id in mbr, update mount and correlative
	copy(mbr.Partitions[index].Id[:], id)
	mbr.Partitions[index].Status = 1
	mbr.Partitions[index].Correlative = int32(correlative)
	// Serialize mbr
	err = utils.Serialize(mbr, path, 0)
	if err != nil {
		return "Error: Coudn't write MBR", err
	}
	fmt.Println("Gloval mounts: ", utils.GlobalMounts)
	return "Mounted succesfully!!", nil
}

func GenerateId(index int, path string) (string, error) {
	carnet := utils.Carnet
	letter, err := utils.GetLetter(path)
	if err != nil {
		return "", err
	}
	return carnet + strconv.Itoa(index) + letter, nil // "28" + "(index)" + "(diskLetter)"
}
