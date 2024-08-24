package commands

import (
	"backend/structures"
	"backend/utils"
	"encoding/binary"
	"fmt"
	"strconv"
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
	partIndex := int(id[2]) - 48 // Convert ASCII to int
	// Get num of structures
	freeSpace := int(mbr.Partitions[partIndex].Size) - binary.Size(structures.SuperBlock{})
	numStructs := freeSpace / (1 + 3 + binary.Size(structures.Inode{}) + 3*binary.Size(structures.FBlock{}))
	// type to int
	typeInt, err := strconv.Atoi(typeF)
	if err != nil {
		return "", err
	}
	// Create superblock
	superBlock := &structures.SuperBlock{
		Type:            int32(typeInt),
		InodesCount:     int32(numStructs),
		BlocksCount:     int32(numStructs * 3),
		FreeBlocksCount: int32(numStructs * 3),
		FreeInodesCount: int32(numStructs),
		MTime:           time.Now().Unix(),
		UMTime:          time.Now().Unix(),
		MCount:          0,
		Magic:           0xEF53,
		InodeSize:       int32(binary.Size(structures.Inode{})),
		BlockSize:       int32(binary.Size(structures.FBlock{})),
		FirstInode:      int32(binary.Size(structures.SuperBlock{})),
		FirstBlock:      int32(binary.Size(structures.SuperBlock{})) + int32(numStructs),
		BMIndoeStart:    int32(binary.Size(structures.SuperBlock{})),
		BMBlockStart:    int32(binary.Size(structures.SuperBlock{})) + int32(numStructs),
		InodeStart:      int32(binary.Size(structures.SuperBlock{})) + int32(numStructs) + 3*int32(numStructs),
		BlockStart:      int32(binary.Size(structures.SuperBlock{})) + int32(numStructs) + 3*int32(numStructs) + int32(numStructs*binary.Size(structures.Inode{})),
	}
	// Serialize superblock
	err = utils.Serialize(superBlock, path, int(mbr.Partitions[partIndex].Start))
	if err != nil {
		return "", err
	}
	// TODO: Create root folder and users.txt file
	return "Partition formating to ext2 was successful...", nil
}
