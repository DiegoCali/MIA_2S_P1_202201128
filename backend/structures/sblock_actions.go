package structures

import (
	"backend/utils"
	"fmt"
	"strings"
)

func (spBlock *SuperBlock) CatInode(pathFile []string, path string) (string, error) {
	inode, err := spBlock.getInode(pathFile, path)
	if err != nil {
		return "File not found", err
	}
	output, err := spBlock.getInodeContent(inode, path)
	if err != nil {
		return "Error getting inode content", err
	}
	return output, nil
}

func (spBlock *SuperBlock) getInode(pathFile []string, path string) (*Inode, error) {
	inode := &Inode{}
	// Get root inode
	err := utils.Deserialize(inode, path, int(spBlock.InodeStart))
	if err != nil {
		return nil, err
	}
	// Get inode
	for i := 0; i < len(pathFile); i++ {
		fmt.Println("Searching for: ", pathFile[i])
		namePath := pathFile[i]
		// Check if inode is directory, if not is because it is a file
		if inode.Type[0] != '0' {
			break
		}
		// Loop through inode blocks
		for j := 0; j < len(inode.Block); j++ {
			// Get DBlock
			if inode.Block[j] != -1 {
				block := &DBlock{}
				err := utils.Deserialize(block, path, int(spBlock.BlockStart+inode.Block[j]*spBlock.BlockSize))
				if err != nil {
					return nil, err
				}
				// Ass block has 4 contents and the first two are parents
				// You only need to check the last 2
				firstContent := block.Content[2]
				if namePath == utils.CheckNull(firstContent.Name[:]) {
					// Replace inode
					err := utils.Deserialize(inode, path, int(spBlock.InodeStart+firstContent.BInode*spBlock.InodeSize))
					if err != nil {
						return nil, err
					}
					break
				}
				secondContent := block.Content[3]
				if namePath == utils.CheckNull(secondContent.Name[:]) {
					// Replace inode
					err := utils.Deserialize(inode, path, int(spBlock.InodeStart+firstContent.BInode*spBlock.InodeSize))
					if err != nil {
						return nil, err
					}
					break
				}
			}
		}
	}
	if inode.Type[0] == '0' {
		return nil, fmt.Errorf("file: " + strings.Join(pathFile, "/") + "not found")
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
