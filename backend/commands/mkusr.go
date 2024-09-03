package commands

import "fmt"

func MkUsr(user string, pass string, group string) (string, error) {
	// Get Superblock
	sb, path, err := getSuperblock()
	if err != nil {
		return "Error getting Superblock", err
	}
	// Check if user exists
	exists, err := sb.ValidateUser(user, pass, path)
	if err != nil {
		return "Error validating user", err
	}
	if exists {
		return "Error user already exists", fmt.Errorf("user already exists")
	}
	// Check if group exists
	exists, err = sb.ValidateGroup(group, path)
	if err != nil {
		return "Error validating group", err
	}
	if !exists {
		return "Error group does not exists", fmt.Errorf("group does not exists")
	}
	// Create user
	err = sb.CreateUser(user, pass, group, path)
	if err != nil {
		return "Error creating user", err
	}
	return "User " + user + " created", nil
}
