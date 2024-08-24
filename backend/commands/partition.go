package commands

import (
	"backend/structures"
	"backend/utils"
	"encoding/binary"
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
	index := 0
	start := int32(0)
	if typeP != "L" {
		// Find if there is an empty partition
		index, err = findEmptyPartition(mbr)
		if err != nil {
			return "Error: Coudn't find empty partition", err
		}
		// Find start of partition
		start = getStartPartition(mbr, index)
	}
	// Create partition
	switch typeP {
	case "P":
		err = createPrimaryPartition(mbr, index, start, int32(sizeBytes), fit, name)
	case "E":
		err = createExtendedPartition(mbr, index, start, int32(sizeBytes), fit, name, path)
	case "L":
		err = createLogicalPartition(mbr, int32(sizeBytes), fit, name, path)
	default:
		return "Error: Partition type not recognized", nil
	}
	if err != nil {
		return "Error: Coudn't create partition", err
	}
	// Serialize MBR
	err = utils.Serialize(mbr, path, 0)
	if err != nil {
		return "Error: Coudn't serialize MBR", err
	}
	mbr.Print()
	return "Partition created succesfully", nil
}

func createPrimaryPartition(mbr *structures.MBR, index int, start int32, size int32, fit string, name string) error {
	mbr.Partitions[index].Status = 1
	copy(mbr.Partitions[index].Type[:], "P")
	copy(mbr.Partitions[index].Fit[:], fit)
	mbr.Partitions[index].Start = start
	mbr.Partitions[index].Size = size
	copy(mbr.Partitions[index].Name[:], name)
	return nil
}

func createExtendedPartition(mbr *structures.MBR, index int, start int32, size int32, fit string, name string, path string) error {
	// Check if there is an extended partition
	_, exists := checkIfExtendedPartitionExists(mbr)
	if exists {
		return fmt.Errorf("error: Extended partition already exists")
	}
	// Create extended partition
	mbr.Partitions[index].Status = 1
	copy(mbr.Partitions[index].Type[:], "E")
	copy(mbr.Partitions[index].Fit[:], fit)
	mbr.Partitions[index].Start = start
	mbr.Partitions[index].Size = size
	copy(mbr.Partitions[index].Name[:], name)
	// Create EBR
	ebr := &structures.EBR{}
	err := ebr.Set(-1, "", start+int32(binary.Size(ebr)), -1, -1, "")
	if err != nil {
		return err
	}
	err = utils.Serialize(ebr, path, int(start))
	if err != nil {
		return err
	}
	return nil
}

func createLogicalPartition(mbr *structures.MBR, size int32, fit string, name string, path string) error {
	// Check if there is an extended partition
	i, exists := checkIfExtendedPartitionExists(mbr)
	if !exists {
		return fmt.Errorf("error: Extended partition does not exist")
	}
	// Get last EBR
	part := &mbr.Partitions[i]
	offset, err := part.GetLastEBR(path)
	if err != nil {
		return err
	}
	// Create logical partition
	ebr := &structures.EBR{}
	err = utils.Deserialize(ebr, path, int(offset))
	if err != nil {
		return err
	}
	nextEbr := offset + int32(binary.Size(ebr)) + size
	// Actualize EBR
	ebr.Mount = 0
	copy(ebr.Fit[:], fit)
	ebr.Size = size
	copy(ebr.Name[:], name)
	ebr.Next = nextEbr
	err = utils.Serialize(ebr, path, int(offset))
	if err != nil {
		return err
	}
	// Create new EBR
	newEbr := &structures.EBR{}
	err = newEbr.Set(-1, "", nextEbr, -1, -1, "")
	if err != nil {
		return err
	}
	err = utils.Serialize(newEbr, path, int(nextEbr))
	if err != nil {
		return err
	}
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
		return 169
	}
	return mbr.Partitions[index-1].Start + mbr.Partitions[index-1].Size
}
