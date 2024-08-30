package structures

import (
	"backend/utils"
	"strings"
)

func (spBlock *SuperBlock) ValidateUser(name string, pass string, path string) (bool, error) {
	// Look for /users.txt
	inode, err := spBlock.getUsersInode(path)
	if err != nil {
		return false, err
	}
	// Read Inode FBlocks
	fullString := ""
	for i := 0; i < len(inode.Block); i++ {
		if inode.Block[i] != -1 {
			block := &FBlock{}
			err = utils.Deserialize(block, path, int(spBlock.BlockStart+inode.Block[i]*spBlock.BlockSize))
			if err != nil {
				return false, err
			}
			fullString += string(block.Content[:])
		}
	}
	// Get lines
	lines := strings.Split(fullString, "\n")
	// Iterate each line
	for i := 0; i < len(lines); i++ {
		// Get fields
		fields := strings.Split(lines[i], ",")
		// Check if user exists
		if fields[1] == "U" && fields[2] == name && fields[3] == pass {
			return true, nil
		}
	}
	return false, nil
}

func (spBlock *SuperBlock) getUsersInode(path string) (*Inode, error) {
	inode := &Inode{}
	err := utils.Deserialize(inode, path, int(spBlock.InodeStart+100)) // first inode is root, +100 is users.txt
	if err != nil {
		return nil, err
	}
	return inode, nil
}
