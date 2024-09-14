package commands

import (
	"fmt"
	"strings"
)

func Cat(files []string) (string, error) {
	// First get superblock
	sb, path, err := getSuperblock()
	if err != nil {
		return "Error getting superblock", err
	}
	// Initialize output
	output := ""
	// Get content of each file, file is a path
	for _, file := range files {
		pathFile := strings.Split(file, "/")
		if pathFile[0] == "" {
			pathFile = pathFile[1:]
		}
		// Get inode
		inodeContent, err := sb.CatInode(pathFile, path)
		if err != nil {
			output += "Error getting inode content: " + file + "\n"
		} else {
			output += inodeContent + "\n"
		}
	}
	fmt.Println("output: ", output)
	return output, nil
}
