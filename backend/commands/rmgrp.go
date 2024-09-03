package commands

import "fmt"

func RmGroup(name string) (string, error) {
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
	if !exists {
		return "Error group does not exists", fmt.Errorf("group does not exists")
	}
	// Remove group
	err = sb.RemoveGroup(name, path)
	if err != nil {
		return "Error deleting group", err
	}
	return "Group " + name + " deleted successfully", nil
}
