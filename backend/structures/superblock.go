package structures

import (
	"backend/utils"
	"time"
)

// SuperBlock Size in bytes: 76
type SuperBlock struct {
	Type            int32 // bytes: 4
	InodesCount     int32 // bytes: 4
	BlocksCount     int32 // bytes: 4
	FreeBlocksCount int32 // bytes: 4
	FreeInodesCount int32 // bytes: 4
	MTime           int64 // bytes: 8
	UMTime          int64 // bytes: 8
	MCount          int32 // bytes: 4
	Magic           int32 // bytes: 4
	InodeSize       int32 // bytes: 4
	BlockSize       int32 // bytes: 4
	FirstInode      int32 // bytes: 4
	FirstBlock      int32 // bytes: 4
	BMIndoeStart    int32 // bytes: 4
	BMBlockStart    int32 // bytes: 4
	InodeStart      int32 // bytes: 4
	BlockStart      int32 // bytes: 4
}

func (spBlock *SuperBlock) InitRoot(path string) error {
	// Create root inode
	rootInode := &Inode{
		UID:   1,
		GID:   1,
		Size:  0,
		Atime: time.Now().Unix(),
		Ctime: time.Now().Unix(),
		Mtime: time.Now().Unix(),
		Block: [15]int32{spBlock.BlocksCount, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		Type:  [1]byte{'0'},
		Perm:  [3]byte{'7', '7', '7'},
	}
	// Serialize inode
	err := utils.Serialize(rootInode, path, int(spBlock.FirstInode))
	if err != nil {
		return err
	}
	// Update bitmap
	err = spBlock.UpdateBitmapInode(path)
	if err != nil {
		return err
	}
	// Update superblock
	spBlock.FreeInodesCount--
	spBlock.InodesCount++
	spBlock.FirstInode += spBlock.InodeSize

	// Create root directory block
	rootBlock := &DBlock{
		Content: [4]Content{
			{Name: [12]byte{'.'}, BInode: 0},
			{Name: [12]byte{'.', '.'}, BInode: 0},
			{Name: [12]byte{'-'}, BInode: -1},
			{Name: [12]byte{'-'}, BInode: -1},
		},
	}
	// Serialize block
	err = utils.Serialize(rootBlock, path, int(spBlock.FirstBlock))
	if err != nil {
		return err
	}
	// Update bitmap
	err = spBlock.UpdateBitmapBlock(path)
	if err != nil {
		return err
	}
	// Update superblock
	spBlock.FreeBlocksCount--
	spBlock.BlocksCount++
	spBlock.FirstBlock += spBlock.BlockSize

	// Print Inode and Block
	err = utils.PrintStruct(rootInode)
	if err != nil {
		return err
	}
	err = utils.PrintStruct(rootBlock)
	if err != nil {
		return err
	}
	return nil
}

func (spBlock *SuperBlock) CreateUsers(path string) error {
	usersText := "1,G,root\n1,U,root,root,123\n"
	// Deserialize rootInode
	rootInode := &Inode{}
	err := utils.Deserialize(rootInode, path, int(spBlock.InodeStart))
	if err != nil {
		return err
	}
	// Update rootInode
	rootInode.Atime = time.Now().Unix()
	// Serialize rootInode
	err = utils.Serialize(rootInode, path, int(spBlock.InodeStart))
	if err != nil {
		return err
	}
	// Deserialize rootBlock
	rootBlock := &DBlock{}
	err = utils.Deserialize(rootBlock, path, int(spBlock.BlockStart))
	if err != nil {
		return err
	}
	// Update rootBlock
	rootBlock.Content[2] = Content{
		Name:   [12]byte{'u', 's', 'e', 'r', 's', '.', 't', 'x', 't'},
		BInode: spBlock.InodesCount,
	}
	// Serialize rootBlock
	err = utils.Serialize(rootBlock, path, int(spBlock.BlockStart))
	if err != nil {
		return err
	}
	// Create users inode
	usersInode := &Inode{
		UID:   1,
		GID:   1,
		Size:  int32(len(usersText)),
		Atime: time.Now().Unix(),
		Ctime: time.Now().Unix(),
		Mtime: time.Now().Unix(),
		Block: [15]int32{spBlock.BlocksCount, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		Type:  [1]byte{'1'},
		Perm:  [3]byte{'7', '7', '7'},
	}
	// Serialize users inode
	err = utils.Serialize(usersInode, path, int(spBlock.FirstInode))
	if err != nil {
		return err
	}
	// Update bitmap
	err = spBlock.UpdateBitmapInode(path)
	if err != nil {
		return err
	}
	// Update superblock
	spBlock.FreeInodesCount--
	spBlock.InodesCount++
	spBlock.FirstInode += spBlock.InodeSize

	// Create users block
	usersBlock := &FBlock{}
	copy(usersBlock.Content[:], usersText)
	// Serialize users block
	err = utils.Serialize(usersBlock, path, int(spBlock.FirstBlock))
	if err != nil {
		return err
	}
	// Update bitmap
	err = spBlock.UpdateBitmapBlock(path)
	if err != nil {
		return err
	}
	// Update superblock
	spBlock.FreeBlocksCount--
	spBlock.BlocksCount++
	spBlock.FirstBlock += spBlock.BlockSize

	// Print Inode and Block
	err = utils.PrintStruct(rootInode)
	if err != nil {
		return err
	}
	err = utils.PrintStruct(rootBlock)
	if err != nil {
		return err
	}
	err = utils.PrintStruct(usersInode)
	if err != nil {
		return err
	}
	err = utils.PrintStruct(usersBlock)
	if err != nil {
		return err
	}
	return nil
}
