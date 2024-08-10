package commands

func MkDisk(size int, fit string, unit string, path string) (string, error) {
	return "Disk created succesfully", nil
}

func RmDisk(path string) (string, error) {
	return "Disk removed succesfully", nil
}
