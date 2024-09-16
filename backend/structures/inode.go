package structures

import (
	"backend/utils"
	"fmt"
	"strconv"
)

// Inode Size in bytes: 100
type Inode struct {
	UID   int32     // bytes: 4
	GID   int32     // bytes: 4
	Size  int32     // bytes: 4
	Atime int64     // bytes: 8
	Ctime int64     // bytes: 8
	Mtime int64     // bytes: 8
	Block [15]int32 // bytes: 60
	Type  [1]byte   // bytes: 1
	Perm  [3]byte   // bytes: 3
}

func (inode *Inode) GetDotStr(id int) string {
	str := "inode_" + strconv.Itoa(id) + "[label=<<TABLE>\n"
	str += "<TR><TD>INODE</TD></TR>\n"
	str += "<TR><TD>i_uid</TD><TD>" + strconv.Itoa(int(inode.UID)) + "</TD></TR>\n"
	str += "<TR><TD>i_gid</TD><TD>" + strconv.Itoa(int(inode.GID)) + "</TD></TR>\n"
	str += "<TR><TD>i_size</TD><TD>" + strconv.Itoa(int(inode.Size)) + "</TD></TR>\n"
	str += "<TR><TD>i_atime</TD><TD>" + utils.Int64ToDate(inode.Atime) + "</TD></TR>\n"
	str += "<TR><TD>i_ctime</TD><TD>" + utils.Int64ToDate(inode.Ctime) + "</TD></TR>\n"
	str += "<TR><TD>i_mtime</TD><TD>" + utils.Int64ToDate(inode.Mtime) + "</TD></TR>\n"
	// Blocks
	for i := 0; i < 15; i++ {
		str += "<TR><TD>i_block_" + strconv.Itoa(i+1) + "</TD><TD>" + strconv.Itoa(int(inode.Block[i])) + "</TD></TR>\n"
	}
	str += "<TR><TD>i_type</TD><TD>" + string(inode.Type[:]) + "</TD></TR>\n"
	str += "<TR><TD>i_perm</TD><TD>" + string(inode.Perm[:]) + "</TD></TR>\n"
	str += "</TABLE>>];\n"
	return str
}

func (inode *Inode) GetAvailableBlock(offset int32, path string) (*DBlock, int32, int, error) {
	// offset is block start of superblock
	// for in range of blocks
	for i := 0; i < 15; i++ {
		// We're going to check each DBlock and see if they have any available space
		if inode.Block[i] != -1 {
			// Get the block
			block := &DBlock{}
			err := utils.Deserialize(block, path, int(offset+inode.Block[i]*64))
			if err != nil {
				return nil, -1, i, err
			}
			// Check if there's any available space, only need to check 2 and 3 index
			for j := 2; j < 4; j++ {
				content := block.Content[j]
				if content.BInode == -1 {
					return block, inode.Block[i], i, nil
				}
			}
			fmt.Println("Block is full")
			continue
		}
		// If the block is -1, we need to create a new block
		newBlock := &DBlock{
			Content: [4]Content{
				{Name: [12]byte{'.'}, BInode: -1},
				{Name: [12]byte{'.', '.'}, BInode: -1},
				{Name: [12]byte{'-'}, BInode: -1},
				{Name: [12]byte{'-'}, BInode: -1},
			},
		}
		// If we created a new block, we need to send -1 as the block number
		return newBlock, -1, i, nil
	}
	return nil, -1, -1, nil
}
