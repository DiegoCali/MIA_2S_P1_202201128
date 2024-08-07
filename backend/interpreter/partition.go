package interpreter

import (
	"encoding/binary"
	"fmt"
	"sort"
)

func createPartition(size int, unit string, path string, typePartition string, fit string, name string) error {
	var partition Partition
	// Look for MBR in disk path
	diskMBR, err := readMBR(path)
	// If not found, throw error
	if err != nil {
		return err
	}
	// Find the first partition with N Type
	var partitionIndex int
	for i, partition := range diskMBR.Partitions {
		if string(partition.Type[:]) == "N" {
			partitionIndex = i
			break
		}
	}
	// Check if name doesn't repeat and if it's the right length
	if len(name) > 16 {
		return fmt.Errorf("name is too long")
	}
	notLegalName := false
	for _, partition := range diskMBR.Partitions {
		if string(partition.Name[:]) == name {
			notLegalName = true
			break
		}
	}
	if notLegalName {
		return fmt.Errorf("name already exists")
	}
	// Check partition type
	if typePartition == "L" {
		// TODO: LOGIC FOR LOGICAL PARTITION
		// - Check if there is an extended partition
		// - Check if there is a free space
		// - Create logical partition
	} else if !(typePartition == "P" || typePartition == "E") {
		return fmt.Errorf("invalid partition type")
	}
	// Get only first char
	fit = fit[:1]
	// Get the size in bytes
	sizeBytes, err := convertToBytes(size, unit)
	if err != nil {
		return err
	}
	// Search free spaceand start for partition
	diskSize := binary.LittleEndian.Uint32(diskMBR.Size[:])
	start, err := searchFreeSpace(diskMBR.Partitions, fit, int(diskSize), sizeBytes)
	if err != nil {
		return err
	}
	// Update info in MBR partition
	partition = diskMBR.Partitions[partitionIndex]
	copy(partition.Type[:], typePartition)
	copy(partition.Fit[:], fit)
	binary.LittleEndian.PutUint32(partition.Start[:], uint32(start))
	binary.LittleEndian.PutUint32(partition.Size[:], uint32(sizeBytes))
	copy(partition.Name[:], name)
	binary.LittleEndian.PutUint32(partition.Correlative[:], uint32(partitionIndex+1))
	return nil
}

func searchFreeSpace(partitions [4]Partition, fit string, sizeD int, sizeP int) (int, error) {
	start := 154 // Default start, discarding MBR
	filledSpaces := make([]Space, 0)
	for _, partition := range partitions {
		// Check if size is "NULL"
		if string(partition.Size[:]) == "NULL" {
			continue
		}
		partitionSize := binary.LittleEndian.Uint32(partition.Size[:])
		partitionStart := binary.LittleEndian.Uint32(partition.Start[:])
		partitionEnd := int(partitionStart) + int(partitionSize)
		// Create space
		space := Space{
			Start: int(partitionStart),
			End:   partitionEnd,
			Size:  int(partitionSize),
		}
		// Add space to filled spaces
		filledSpaces = append(filledSpaces, space)
	}
	// Sort by start
	sort.Slice(filledSpaces, func(i, j int) bool {
		return filledSpaces[i].Start < filledSpaces[j].Start
	})
	// Create new list of free spaces
	freeSpaces := make([]Space, 0)
	for i := 0; i < len(filledSpaces)-1; i++ {
		// Calculate space between two partitions
		space := Space{
			Start: filledSpaces[i].End + 1,
			End:   filledSpaces[i+1].Start,
			Size:  filledSpaces[i+1].Start - filledSpaces[i].End,
		}
		// Add space to free spaces
		freeSpaces = append(freeSpaces, space)
	}
	// Add the free space after the last partition
	space := Space{
		Start: filledSpaces[len(filledSpaces)-1].End + 1,
		End:   sizeD,
		Size:  sizeD - filledSpaces[len(filledSpaces)-1].End,
	}
	freeSpaces = append(freeSpaces, space)
	// Dependind of fit, sort the slice then search for the right space
	if fit == "F" {
		// Sort by start, ascending
		sort.Slice(freeSpaces, func(i, j int) bool {
			return freeSpaces[i].Start < freeSpaces[j].Start
		})
	} else if fit == "B" {
		// Sort by size, ascending
		sort.Slice(freeSpaces, func(i, j int) bool {
			return freeSpaces[i].Size < freeSpaces[j].Size
		})
	} else if fit == "W" {
		// Sort by size, descending
		sort.Slice(freeSpaces, func(i, j int) bool {
			return freeSpaces[i].Size > freeSpaces[j].Size
		})
	} else {
		return -1, fmt.Errorf("fit %s not recognized", fit)
	}
	// Search for the right space, take the first that fits the sizeP
	for _, space := range freeSpaces {
		if space.Size >= sizeP {
			start = space.Start
			break
		}
	}
	return start, nil
}
