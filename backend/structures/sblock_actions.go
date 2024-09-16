package structures

import (
	"backend/utils"
	"fmt"
	"strings"
	"time"
)

func (spBlock *SuperBlock) CreatePath(parents []string, isFile bool, inode *Inode, dir int32, dirParent int32, path string, offset int32, create bool) error {
	fmt.Println("Create parents?", create)
	if !create {
		// As we are not creating, we only need search the prev last inode
		nameInode := make([]string, 0)
		nameInode = append(nameInode, parents[len(parents)-1])
		parents = parents[:len(parents)-1]
		fmt.Println("Creating: ", nameInode, ", With parents: ", parents)
		lastInode, inodeRef, err := spBlock.getInode(parents, path)
		if err != nil {
			return err
		}
		fmt.Println("Last Inode: ", lastInode)
		err = spBlock.CreatePath(nameInode, isFile, lastInode, dir, dirParent, path, inodeRef, true)
		if err != nil {
			return err
		}
		return nil
	}
	// First iteration is inode as root
	// Create first parent, we need to get first aviailable block from root inode
	dirBlock, blockRef, index, err := inode.GetAvailableBlock(spBlock.BlockStart, path)
	if err != nil {
		return err
	}
	if blockRef == -1 {
		// means a new block was created, so we need to update inode and superblock
		inode.Block[index] = spBlock.BlocksCount
		blockRef = spBlock.BlocksCount
		inode.Mtime = time.Now().Unix()
		err = spBlock.UpdateBitmapBlock(path)
		if err != nil {
			return err
		}
		// Serialize block
		fmt.Println("Comparing: ", spBlock.FirstBlock, ", And, ", spBlock.BlockStart+blockRef*spBlock.BlockSize)
		err = utils.Serialize(dirBlock, path, int(spBlock.FirstBlock))
		if err != nil {
			return err
		}
		// Update superblock
		spBlock.FreeBlocksCount--
		spBlock.BlocksCount++
		spBlock.FirstBlock += spBlock.BlockSize
	}
	inode.Atime = time.Now().Unix()
	err = utils.Serialize(inode, path, int(spBlock.InodeStart+offset*spBlock.InodeSize))
	if err != nil {
		return err
	}
	// Create new inode for first parent
	newInode, inodeRef, err := spBlock.CreateInode(isFile, path, dir, dirParent)
	if err != nil {
		return err
	}
	// With dirBlock we can create the first parent
	for i := 2; i < 4; i++ {
		if dirBlock.Content[i].BInode == -1 {
			if len(parents) > 0 {
				copy(dirBlock.Content[i].Name[:], parents[0])
				dirBlock.Content[i].BInode = inodeRef
				break
			}
			break
		}
	}
	err = utils.Serialize(dirBlock, path, int(spBlock.BlockStart+inode.Block[index]*spBlock.BlockSize))
	if err != nil {
		return err
	}
	// Create the rest of the parents, recursively
	newParents := parents[1:]
	if len(newParents) > 1 {
		err = spBlock.CreatePath(newParents, false, newInode, inodeRef, dir, path, inodeRef, true)
		if err != nil {
			return err
		}
	} else if len(newParents) == 1 {
		err = spBlock.CreatePath(newParents, isFile, newInode, inodeRef, dir, path, inodeRef, true)
		if err != nil {
			return err
		}
	}
	// Serialize superblock
	err = utils.Serialize(spBlock, path, int(spBlock.BMIndoeStart-76))
	if err != nil {
		return err
	}
	return nil
}
func (spBlock *SuperBlock) FillFile(pathF []string, size int, cont []string, path string) error {
	inode, ref, err := spBlock.getInode(pathF, path)
	if err != nil {
		return err
	}
	fmt.Println("Inode: ", inode.Type)
	contStr := ""
	if len(cont) > 0 {
		inode, _, err = spBlock.getInode(cont, path)
		if err != nil {
			return err
		}
		contStr, err = spBlock.getInodeContent(inode, path)
		if err != nil {
			return err
		}
	} else {
		// Fill str with 0123456789 until size
		strFiller := "0123456789"
		complete := size / 10
		rest := size % 10
		for i := 0; i < complete; i++ {
			contStr += strFiller
		}
		contStr += strFiller[:rest]
	}
	err = spBlock.writeInode(inode, path, contStr, int(spBlock.InodeStart+ref*spBlock.InodeSize))
	if err != nil {
		return err
	}
	return nil
}

func (spBlock *SuperBlock) CreateInode(isFile bool, path string, dir int32, parentDir int32) (*Inode, int32, error) {
	// Type 0 is directory, 1 is file by default
	typeInode := [1]byte{'0'}
	if isFile {
		typeInode = [1]byte{'1'}
	}
	newInode := &Inode{
		UID:   1,
		GID:   1,
		Size:  0,
		Atime: time.Now().Unix(),
		Ctime: time.Now().Unix(),
		Mtime: time.Now().Unix(),
		Block: [15]int32{spBlock.BlocksCount, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		Type:  typeInode,
		Perm:  [3]byte{'6', '6', '4'},
	}
	// Serialize inode
	err := utils.Serialize(newInode, path, int(spBlock.FirstInode))
	if err != nil {
		return nil, -1, err
	}
	// Update bitmap
	err = spBlock.UpdateBitmapInode(path)
	if err != nil {
		return nil, -1, err
	}
	// Save inode count as a reference
	inodeRef := spBlock.InodesCount
	// Update superblock
	spBlock.FreeInodesCount--
	spBlock.InodesCount++
	spBlock.FirstInode += spBlock.InodeSize

	// Create block for inode, depending on the type
	if isFile {
		newBlock := &FBlock{}
		err = utils.Serialize(newBlock, path, int(spBlock.FirstBlock))
		if err != nil {
			return nil, -1, err
		}
	} else {
		newBlock := &DBlock{
			Content: [4]Content{
				{Name: [12]byte{'.'}, BInode: dir},
				{Name: [12]byte{'.', '.'}, BInode: parentDir},
				{Name: [12]byte{'-'}, BInode: -1},
				{Name: [12]byte{'-'}, BInode: -1},
			},
		}
		err = utils.Serialize(newBlock, path, int(spBlock.FirstBlock))
		if err != nil {
			return nil, -1, err
		}
	}
	// Update bitmap
	err = spBlock.UpdateBitmapBlock(path)
	if err != nil {
		return nil, -1, err
	}
	// Update superblock
	spBlock.FreeBlocksCount--
	spBlock.BlocksCount++
	spBlock.FirstBlock += spBlock.BlockSize
	// Serialize superblock
	err = utils.Serialize(spBlock, path, int(spBlock.BMIndoeStart-76))
	if err != nil {
		return nil, -1, err
	}
	return newInode, inodeRef, err
}

func (spBlock *SuperBlock) CatInode(pathFile []string, path string) (string, error) {
	inode, _, err := spBlock.getInode(pathFile, path)
	if err != nil {
		return "File not found", err
	}
	output, err := spBlock.getInodeContent(inode, path)
	if err != nil {
		return "Error getting inode content", err
	}
	return output, nil
}

func (spBlock *SuperBlock) getInode(pathFile []string, path string) (*Inode, int32, error) {
	inode := &Inode{}
	// Get root inode
	err := utils.Deserialize(inode, path, int(spBlock.InodeStart))
	if err != nil {
		return nil, -1, err
	}
	inodeRef := int32(-1)
	// Get inode
	for i := 0; i < len(pathFile); i++ {
		fmt.Println("Searching for: ", pathFile[i])
		namePath := pathFile[i]
		// Check if inode is directory, if not is because it is a file
		// Loop through inode blocks
		for j := 0; j < len(inode.Block); j++ {
			// Get DBlock
			if inode.Block[j] != -1 {
				block := &DBlock{}
				err := utils.Deserialize(block, path, int(spBlock.BlockStart+inode.Block[j]*spBlock.BlockSize))
				if err != nil {
					return nil, -1, err
				}
				// Ass block has 4 contents and the first two are parents
				// You only need to check the last 2
				firstContent := block.Content[2]
				fmt.Println("Comparing to:", utils.CheckNull(firstContent.Name[:]))
				if namePath == utils.CheckNull(firstContent.Name[:]) {
					fmt.Println("Found")
					// Replace inode
					err := utils.Deserialize(inode, path, int(spBlock.InodeStart+firstContent.BInode*spBlock.InodeSize))
					if err != nil {
						return nil, -1, err
					}
					inodeRef = firstContent.BInode
					break
				}
				secondContent := block.Content[3]
				fmt.Println("Comparing to:", utils.CheckNull(secondContent.Name[:]))
				if namePath == utils.CheckNull(secondContent.Name[:]) {
					fmt.Println("Found")
					// Replace inode
					err := utils.Deserialize(inode, path, int(spBlock.InodeStart+secondContent.BInode*spBlock.InodeSize))
					if err != nil {
						return nil, -1, err
					}
					inodeRef = secondContent.BInode
					break
				}
			}
		}
	}
	return inode, inodeRef, nil
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

func (spBlock *SuperBlock) getInodeDirs(inode *Inode, path string) ([]string, error) {
	dirs := make([]string, 0)
	for i := 0; i < len(inode.Block); i++ {
		if inode.Block[i] != -1 {
			fmt.Println("Block: ", inode.Block[i])
			block := &DBlock{}
			err := utils.Deserialize(block, path, int(spBlock.BlockStart+inode.Block[i]*spBlock.BlockSize))
			if err != nil {
				return nil, err
			}
			for j := 2; j < 4; j++ {
				if block.Content[j].BInode != -1 {
					fmt.Println("Appending: ", utils.CheckNull(block.Content[j].Name[:]))
					dirs = append(dirs, utils.CheckNull(block.Content[j].Name[:]))
				}
			}
		}
	}
	return dirs, nil
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
	// Change inode size
	inode.Size = int32(len(content))
	// Change inode mtime
	inode.Mtime = time.Now().Unix()
	// Serialize inode
	err := utils.Serialize(inode, path, offset)
	if err != nil {
		return err
	}
	return nil
}

func (spBlock *SuperBlock) LsReport(pathFile []string, path string) (string, error) {
	inode, _, err := spBlock.getInode(pathFile, path)
	if err != nil {
		return "File not found", err
	}
	fullString := "<<TABLE>\n"
	dirs, err := spBlock.getInodeDirs(inode, path)
	if err != nil {
		return "Error getting inode content", err
	}
	fmt.Println(dirs)
	fullString += "<TR><TH>NAME</TH><TH>TYPE</TH></TR>\n"
	for i := 0; i < len(dirs); i++ {
		if strings.ContainsRune(dirs[i], '.') {
			fullString += "<TR><TD>" + dirs[i] + "</TD><TD>File</TD></TR>\n"
			continue
		}
		fullString += "<TR><TD>" + dirs[i] + "</TD><TD>Directory</TD></TR>\n"
	}
	return fullString, nil
}
