package commands

import (
	"backend/structures"
	"backend/utils"
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
		return "Error: Coudn't deserialize MBR", err
	}
	// Find [index] of the first partition with status -1
	index := 0
	for i := 0; i < 4; i++ {
		if mbr.Partitions[i].Status == -1 {
			index = i
			break
		}
	}
	// Get start of the partition
	start := int32(166)
	if index != 0 {
		start = mbr.Partitions[index-1].Start + mbr.Partitions[index-1].Size + 1
	}
	// Check if there is space for the partition
	enoughSpace := false
	if index == 0 {
		if int32(sizeBytes) < mbr.Size {
			enoughSpace = true
		}
	} else {
		if start+int32(sizeBytes) < mbr.Size {
			enoughSpace = true
		}
	}
	if !enoughSpace {
		return "Error: Not enough space for the partition", nil
	}
	// Put partition data
	mbr.Partitions[index].Status = 0
	copy(mbr.Partitions[index].Type[:], typeP)
	copy(mbr.Partitions[index].Fit[:], fit)
	if index == 0 {
		mbr.Partitions[index].Start = 166 // MBR size + 1
	} else {
		mbr.Partitions[index].Start = mbr.Partitions[index-1].Start + mbr.Partitions[index-1].Size + 1
	}
	mbr.Partitions[index].Size = int32(sizeBytes)
	copy(mbr.Partitions[index].Name[:], name)
	mbr.Partitions[index].Correlative = int32(index + 1)
	copy(mbr.Partitions[index].Id[:], "----")
	// Write the MBR
	err = mbr.Serialize(path)
	if err != nil {
		return "Error: Coudn't serialize MBR", err
	}
	return "Partition created succesfully", nil
}
