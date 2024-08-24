package commands

import (
	"backend/structures"
	"backend/utils"
	"fmt"
)

func FDisk(size int, unit string, path string, typeP string, fit string, name string) (string, error) {
	sizeBytes, err := utils.ConvertToBytes(size, unit)
	if err != nil {
		return "Error: Coudn't transform size to bytes", err
	}
	// Read the MBR
	mbr := &structures.MBR{}
	err = utils.Deserialize(mbr, path, 0)
	if err != nil {
		return "Error: Coudn't read MBR", err
	}
	// Check if there is enough space
	if !checkIfEnoughSpace(mbr, sizeBytes) {
		return "Error: Not enough space in disk", fmt.Errorf("error: Not enough space in disk")
	}
	// Create partition
	switch typeP {
	case "P":
		// TODO: create primary partition
	case "E":
		// TODO: create extended partition
	case "L":
		// TODO: create logical partition
	default:
		return "Error: Partition type not recognized", nil
	}
	// Serialize MBR
	err = utils.Serialize(mbr, path, 0)
	if err != nil {
		return "Error: Coudn't serialize MBR", err
	}
	mbr.Print()
	return "Partition created succesfully", nil
}

func findEmptyPartition(mbr *structures.MBR) (int, error) {
	for i, partition := range mbr.Partitions {
		if partition.Status == -1 {
			return i, nil
		}
	}
	return -1, nil
}

func checkIfExtendedPartitionExists(mbr *structures.MBR) (int, bool) {
	for i, partition := range mbr.Partitions {
		if partition.Type[0] == 'E' {
			return i, true
		}
	}
	return -1, false
}

func checkIfEnoughSpace(mbr *structures.MBR, sizeBytes int) bool {
	var totalSize int
	for _, partition := range mbr.Partitions {
		if partition.Status != -1 {
			totalSize += int(partition.Size)
		}
	}
	return totalSize+sizeBytes <= int(mbr.Size)
}

func getStartPartition(mbr *structures.MBR, index int) int32 {
	if index == 0 {
		return 170 // MBR size
	}
	return mbr.Partitions[index-1].Start + mbr.Partitions[index-1].Size + 1
}
