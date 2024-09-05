package commands

import "fmt"

func Chgrp(user string, group string) (string, error) {
	// Get superblock
	sb, path, err := getSuperblock()
	if err != nil {
		return "Error getting superblock", err
	}
	// Validate group
	exists, err := sb.ValidateGroup(group, path)
	if err != nil {
		return "Error validating group", err
	}
	if !exists {
		return "Error group does not exist", fmt.Errorf("group does not exist")
	}
	// Change group
	err = sb.ChangeGroup(user, group, path)
	if err != nil {
		return "Error changing group", err
	}
	// Return success
	return "User: " + user + " changed group to: " + group, nil
}
