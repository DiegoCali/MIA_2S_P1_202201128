package commands

import (
	"backend/structures"
	"backend/utils"
	"fmt"
)

func MkGroup(name string) (string, error) {
	// Get Superblock
	sb, path, err := getSuperblock()
	if err != nil {
		return "Error getting Superblock", err
	}
	// Check if group exists
	exists, err := sb.ValidateGroup(name, path)
	if err != nil {
		return "Error validating group", err
	}
	if exists {
		return "Error group already exists", fmt.Errorf("group already exists")
	}
	// Create group
	err = sb.CreateGroup(name, path)
	if err != nil {
		return "Error creating group", err
	}
	return "Group " + name + " created successfully", nil
}

func getSuperblock() (*structures.SuperBlock, string, error) {
	// Get Path
	path, exists := utils.GlobalMounts[utils.ActualUser.GetId()]
	if !exists {
		return nil, "", fmt.Errorf("user is not logged in")
	}
	// Read MBR
	mbr := &structures.MBR{}
	err := utils.Deserialize(mbr, path, 0)
	if err != nil {
		return nil, "", err
	}
	// Get Partition Index by ID
	index, err := mbr.GetPartitionId(utils.ActualUser.GetId())
	if err != nil {
		return nil, "", err
	}
	// Read Partition and check if primary
	if string(mbr.Partitions[index].Type[:]) != "P" {
		return nil, "", fmt.Errorf("partition is not primary")
	}
	// Read superblock
	sb := &structures.SuperBlock{}
	err = utils.Deserialize(sb, path, int(mbr.Partitions[index].Start))
	if err != nil {
		return nil, "", err
	}
	return sb, path, nil
}
