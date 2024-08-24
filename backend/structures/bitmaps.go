package structures

import (
	"encoding/binary"
	"os"
)

func (spBlock *SuperBlock) CreateBitMaps(path string) error {
	file, err := os.OpenFile(path, os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			return
		}
	}(file)
	// Inode bitmap
	// Move pointer to the start of the bitmap
	_, err = file.Seek(int64(spBlock.BMIndoeStart), 0)
	if err != nil {
		return err
	}
	// Create buffer of '0's
	buffer := make([]byte, spBlock.FreeInodesCount)
	for i := 0; i < len(buffer); i++ {
		buffer[i] = '0'
	}
	// Write buffer to file
	err = binary.Write(file, binary.LittleEndian, buffer)
	if err != nil {
		return err
	}
	// Block bitmap
	// Move pointer to the start of the bitmap
	_, err = file.Seek(int64(spBlock.BMBlockStart), 0)
	if err != nil {
		return err
	}
	// Create buffer of 'O's
	buffer = make([]byte, spBlock.FreeBlocksCount)
	for i := 0; i < len(buffer); i++ {
		buffer[i] = 'O'
	}
	// Write buffer to file
	err = binary.Write(file, binary.LittleEndian, buffer)
	if err != nil {
		return err
	}
	return nil
}

func (spBlock *SuperBlock) UpdateBitmapInode(path string) error {
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			return
		}
	}(file)
	// Move pointer to the last inode
	_, err = file.Seek(int64(spBlock.BMIndoeStart)+int64(spBlock.InodesCount), 0)
	if err != nil {
		return err
	}
	// Write '1' to the inode bitmap
	_, err = file.Write([]byte{'1'})
	if err != nil {
		return err
	}
	return nil
}

func (spBlock *SuperBlock) UpdateBitmapBlock(path string) error {
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			return
		}
	}(file)
	// Move pointer to the last block
	_, err = file.Seek(int64(spBlock.BMBlockStart)+int64(spBlock.BlocksCount), 0)
	if err != nil {
		return err
	}
	// Write '1' to the block bitmap
	_, err = file.Write([]byte{'X'})
	if err != nil {
		return err
	}
	return nil
}
