package commands

import (
	"backend/structures"
	"backend/utils"
	"encoding/binary"
	"fmt"
	"time"
)

func MkFS(id string, typeF string) (string, error) {
	// ID structure: '28'  + [numPartition] + [diskLetter]
	// Get path from id
	path, exists := utils.GlobalMounts[id]
	if !exists {
		return "", fmt.Errorf("id %s not found", id)
	}
	// Get mbr from disk
	mbr := &structures.MBR{}
	err := utils.Deserialize(mbr, path, 0)
	if err != nil {
		return "", err
	}
	// Get partition index from id
	partIndex, err := mbr.GetPartitionId(id)
	if err != nil {
		return "", err
	}
	offset := mbr.Partitions[partIndex].Start
	freeSpace := int(mbr.Partitions[partIndex].Size) - binary.Size(structures.SuperBlock{})
	numStructs := freeSpace / (1 + 3 + binary.Size(structures.Inode{}) + 3*binary.Size(structures.FBlock{}))
	// Calculate starts
	// Inode bitmap start
	inodeBitmapStart := offset + int32(binary.Size(structures.SuperBlock{}))
	// Block bitmap start
	blockBitmapStart := inodeBitmapStart + int32(numStructs)
	// Inode start
	inodeStart := blockBitmapStart + int32(3*numStructs)
	// Block start
	blockStart := inodeStart + int32(numStructs*binary.Size(structures.Inode{}))
	// type to int
	typeInt := 2
	// Create superblock
	superBlock := &structures.SuperBlock{
		Type:            int32(typeInt),
		InodesCount:     0,
		BlocksCount:     0,
		FreeBlocksCount: int32(numStructs * 3),
		FreeInodesCount: int32(numStructs),
		MTime:           time.Now().Unix(),
		UMTime:          time.Now().Unix(),
		MCount:          1,
		Magic:           0xEF53,
		InodeSize:       int32(binary.Size(structures.Inode{})),
		BlockSize:       int32(binary.Size(structures.FBlock{})),
		FirstInode:      inodeStart,
		FirstBlock:      blockStart,
		BMIndoeStart:    inodeBitmapStart,
		BMBlockStart:    blockBitmapStart,
		InodeStart:      inodeStart,
		BlockStart:      blockStart,
	}
	// Serialize superblock
	err = utils.Serialize(superBlock, path, int(mbr.Partitions[partIndex].Start))
	if err != nil {
		return "", err
	}
	// Create bitmaps
	err = superBlock.CreateBitMaps(path)
	if err != nil {
		return "", err
	}
	// Init root
	err = superBlock.InitRoot(path)
	if err != nil {
		return "", err
	}
	// Create users
	err = superBlock.CreateUsers(path)
	if err != nil {
		return "", err
	}
	return "Partition formating to ext2 was successful...", nil
}
