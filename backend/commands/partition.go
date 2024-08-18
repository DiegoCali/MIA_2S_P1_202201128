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
		err = createExtendedPartition(mbr, sizeBytes, fit, name, path)
	case "L":
		err = createLogicalPartition(mbr, sizeBytes, fit, name, path)
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
	// TODO: Create SuperBlocks and Inodes
	return nil
}

func createExtendedPartition(mbr *structures.MBR, sizeBytes int, fit string, name string, path string) error {
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
	// Create EBR
	ebr := &structures.EBR{}
	err = ebr.Set(-1, "", mbr.Partitions[index].Start+33, -1, -1, "")
	if err != nil {
		return err
	}
	err = ebr.Serialize(path, int(mbr.Partitions[index].Start))
	if err != nil {
		return err
	}
	return nil
}

func createLogicalPartition(mbr *structures.MBR, sizeBytes int, fit string, name string, path string) error {
	index, exists := checkIfExtendedPartitionExists(mbr)
	if !exists {
		return fmt.Errorf("error: Extended partition doesn't exist")
	}
	// Read EBR
	offset, err := mbr.Partitions[index].GetLastEBR(path)
	if err != nil {
		return err
	}
	ebr := &structures.EBR{}
	err = ebr.Deserialize(path, int(offset))
	if err != nil {
		return err
	}
	// Check if there is enough space
	if ebr.Start+int32(sizeBytes)+33 > mbr.Partitions[index].Start+mbr.Partitions[index].Size {
		return fmt.Errorf("error: Not enough space in extended partition")
	}
	// Set new data to EBR
	err = ebr.Set(0, fit, ebr.Start, int32(sizeBytes), ebr.Start+int32(sizeBytes)+1, name)
	if err != nil {
		return err
	}
	err = ebr.Serialize(path, int(offset))
	if err != nil {
		return err
	}
	// Write next EBR
	nextEbr := &structures.EBR{}
	err = nextEbr.Set(-1, "", ebr.Start+int32(sizeBytes)+1, -1, -1, "")
	if err != nil {
		return err
	}
	err = nextEbr.Serialize(path, int(nextEbr.Start))
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
		return 170 // MBR size
	}
	return mbr.Partitions[index-1].Start + mbr.Partitions[index-1].Size + 1
}
