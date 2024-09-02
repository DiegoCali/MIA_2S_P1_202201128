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
	fullString, err := spBlock.getInodeContent(inode, path)
	if err != nil {
		return false, err
	}
	// Get lines
	lines := strings.Split(fullString, "\n")
	// Iterate each line
	for i := 0; i < len(lines); i++ {
		// Get fields
		fields := strings.Split(lines[i], ",")
		// Check if fields are empty
		if len(fields) < 3 {
			break
		}
		// Check if user exists and password is correct
		if fields[1] == "U" && fields[2] == name && fields[3] == pass {
			return true, nil
		}
	}
	return false, nil
}

func (spBlock *SuperBlock) ValidateGroup(name string, path string) (bool, error) {
	inode, err := spBlock.getUsersInode(path)
	if err != nil {
		return false, err
	}
	fullString, err := spBlock.getInodeContent(inode, path)
	if err != nil {
		return false, err
	}
	lines := strings.Split(fullString, "\n")
	for i := 0; i < len(lines); i++ {
		fields := strings.Split(lines[i], ",")
		if len(fields) < 3 {
			break
		}
		if fields[1] == "G" && fields[2] == name {
			return true, nil
		}
	}
	return false, nil
}

func (spBlock *SuperBlock) CreateGroup(name string, path string) error {
	groupText := "1,G," + name + "\n"
	inode, err := spBlock.getUsersInode(path)
	if err != nil {
		return err
	}
	fullString, err := spBlock.getInodeContent(inode, path)
	if err != nil {
		return err
	}
	fullString += groupText
	err = spBlock.writeInode(inode, path, fullString, int(spBlock.InodeStart+100)) // first inode is root, +100 is users.txt
	if err != nil {
		return err
	}
	return nil
}

func (spBlock *SuperBlock) getUsersInode(path string) (*Inode, error) {
	inode := &Inode{}
	err := utils.Deserialize(inode, path, int(spBlock.InodeStart+100)) // first inode is root, +100 is users.txt
	if err != nil {
		return nil, err
	}
	return inode, nil
}

func (spBlock *SuperBlock) getInodeContent(inode *Inode, path string) (string, error) {
	fullString := ""
	for i := 0; i < len(inode.Block); i++ {
		if inode.Block[i] != -1 {
			block := &FBlock{}
			err := utils.Deserialize(block, path, int(spBlock.BlockStart+inode.Block[i]*spBlock.BlockSize))
			if err != nil {
				return "", err
			}
			fullString += utils.CheckNull(block.Content[:])
		}
	}
	return fullString, nil
}

func (spBlock *SuperBlock) writeInode(inode *Inode, path string, content string, offset int) error {
	// String of 64 0s
	zeroes := strings.Repeat("\x00", 64)
	// Get content length
	contentLength := len(content)
	// Write content
	for i := 0; i < len(inode.Block); i++ {
		if inode.Block[i] == -1 {
			var err error
			inode.Block[i], err = spBlock.CreateFBlock(zeroes, path)
			if err != nil {
				return err
			}
		}
		block := &FBlock{}
		err := utils.Deserialize(block, path, int(spBlock.BlockStart+inode.Block[i]*spBlock.BlockSize))
		if err != nil {
			return err
		}
		// Write content
		if contentLength > 64 {
			copy(block.Content[:], content[:64])
			content = content[64:]
			contentLength -= 64
		} else {
			// Write 0s
			copy(block.Content[:], zeroes)
			// Write content
			copy(block.Content[:], content)
			contentLength = 0
		}
		// Serialize block
		err = utils.Serialize(block, path, int(spBlock.BlockStart+inode.Block[i]*spBlock.BlockSize))
		if err != nil {
			return err
		}
		// Check if content is empty
		if contentLength == 0 {
			break
		}
	}
	// Serialize inode
	err := utils.Serialize(inode, path, offset)
	if err != nil {
		return err
	}
	return nil
}
