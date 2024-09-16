package commands

import (
	"backend/structures"
	"backend/utils"
	"strings"
)

func MkDir(path string, parent bool) (string, error) {
	// Get superblock
	sb, pathF, err := getSuperblock()
	if err != nil {
		return "", err
	}
	pathList := strings.Split(path, "/")
	if pathList[0] == "" {
		pathList = pathList[1:]
	}
	// Get root inode
	inode := &structures.Inode{}
	err = utils.Deserialize(inode, pathF, int(sb.InodeStart))
	if err != nil {
		return "", err
	}
	err = sb.CreatePath(pathList, false, inode, 0, 0, pathF, 0, parent)
	if err != nil {
		return "", err
	}
	return "Path created Succesfully", nil
}
