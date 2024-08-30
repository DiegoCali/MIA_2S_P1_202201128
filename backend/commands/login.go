package commands

import (
	"backend/structures"
	"backend/utils"
	"fmt"
)

func Login(name string, pass string, id string) (string, error) {
	// Check if ID exists
	path, exists := checkIfIDExists(id)
	if !exists {
		return "Error logging in", fmt.Errorf("ID %s does not exist", id)
	}
	// Read MBR
	mbr := &structures.MBR{}
	err := utils.Deserialize(mbr, path, 0)
	if err != nil {
		return "Error logging in", err
	}
	// Get Partition Index by ID
	index, err := mbr.GetPartitionId(id)
	if err != nil {
		return "Error logging in", err
	}
	// Read Partition and check if primary
	if string(mbr.Partitions[index].Type[:]) != "P" {
		return "Error logging in", fmt.Errorf("partition is not primary")
	}
	// Read superblock
	sb := &structures.SuperBlock{}
	err = utils.Deserialize(sb, path, int(mbr.Partitions[index].Start))
	if err != nil {
		return "Error logging in", err
	}
	// Check if user exists
	userExists, err := sb.ValidateUser(name, pass, path)
	if err != nil {
		return "Error logging in", err
	}
	if !userExists {
		return "Error logging in", fmt.Errorf("user does not exist or password is incorrect")
	}
	// Save user and id in global variables
	utils.ActualUser.Set(name, id)
	fmt.Println("Actual user is now:", utils.ActualUser)
	return "Logged in as: [" + name + "] successfully", nil
}
