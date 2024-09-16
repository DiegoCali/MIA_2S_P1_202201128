package commands

import (
	"backend/structures"
	"backend/utils"
	"strings"
)

func MkFile(path string, parent bool, size int, cont string) (string, error) {
	// Get superblock
	sb, pathF, err := getSuperblock()
	if err != nil {
		return "", err
	}
	pathList := strings.Split(path, "/")
	if pathList[0] == "" {
		pathList = pathList[1:]
	}
	content := strings.Split(cont, "/")
	if content[0] == "" {
		content = content[1:]
	}
	// Get root inode
	inode := &structures.Inode{}
	err = utils.Deserialize(inode, pathF, int(sb.InodeStart))
	if err != nil {
		return "", err
	}
	err = sb.CreatePath(pathList, true, inode, 0, 0, pathF, 0)
	if err != nil {
		return "", err
	}
	err = sb.FillFile(pathList, size, content, pathF)
	if err != nil {
		return "", err
	}
	return "File Created Successfully", nil
}
