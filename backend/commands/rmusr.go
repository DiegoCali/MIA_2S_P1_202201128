package commands

func RmUsr(user string) (string, error) {
	// Get Superblock
	sb, path, err := getSuperblock()
	if err != nil {
		return "Error getting Superblock", err
	}
	// Remove user
	err = sb.RemoveUser(user, path)
	if err != nil {
		return "Error removing user", err
	}
	return "User " + user + " removed", nil
}
