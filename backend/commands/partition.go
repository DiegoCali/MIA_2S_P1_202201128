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
	err = mbr.Deserialize(path)
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
		err = createPrimaryPartition(mbr, sizeBytes, fit, name)
	case "E":
		err = createExtendedPartition(mbr, sizeBytes, fit, name)
	case "L":
		err = createLogicalPartition(mbr, sizeBytes, fit, name)
	default:
		return "Error: Partition type not recognized", nil
	}
	if err != nil {
		return "Error: Coudn't create partition", err
	}
	// Serialize MBR
	err = mbr.Serialize(path)
	if err != nil {
		return "Error: Coudn't serialize MBR", err
	}
	mbr.Print()
	return "Partition created succesfully", nil
}

func createPrimaryPartition(mbr *structures.MBR, sizeBytes int, fit string, name string) error {
	index, err := findEmptyPartition(mbr)
	if err != nil {
		return err
	}
	// Copy data to partition
	mbr.Partitions[index].Status = 0
	copy(mbr.Partitions[index].Type[:], "P")
	copy(mbr.Partitions[index].Fit[:], fit)
	mbr.Partitions[index].Start = getStartPartition(mbr, index)
	mbr.Partitions[index].Size = int32(sizeBytes)
	copy(mbr.Partitions[index].Name[:], name)
	mbr.Partitions[index].Correlative = int32(index)
	copy(mbr.Partitions[index].Id[:], "----")
	return nil
}

func createExtendedPartition(mbr *structures.MBR, sizeBytes int, fit string, name string) error {
	_, exists := checkIfExtendedPartitionExists(mbr)
	if exists {
		return fmt.Errorf("error: Extended partition already exists")
	}
	index, err := findEmptyPartition(mbr)
	if err != nil {
		return err
	}
	// Copy data to partition
	mbr.Partitions[index].Status = 0
	copy(mbr.Partitions[index].Type[:], "E")
	copy(mbr.Partitions[index].Fit[:], fit)
	mbr.Partitions[index].Start = getStartPartition(mbr, index)
	mbr.Partitions[index].Size = int32(sizeBytes)
	copy(mbr.Partitions[index].Name[:], name)
	mbr.Partitions[index].Correlative = int32(index)
	copy(mbr.Partitions[index].Id[:], "----")
	// TODO: Create EBR
	return nil
}

func createLogicalPartition(mbr *structures.MBR, sizeBytes int, fit string, name string) error {
	index, exists := checkIfExtendedPartitionExists(mbr)
	if !exists {
		return fmt.Errorf("error: Extended partition doesn't exist")
	}
	// Get EBR offset
	offset := mbr.Partitions[index].Start
	// Read EBR
	ebr := &structures.EBR{}
	fmt.Println("EBR offset: ", offset, ebr)
	return nil
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
		return 166 // MBR size
	}
	return mbr.Partitions[index-1].Start + mbr.Partitions[index-1].Size + 1
}
