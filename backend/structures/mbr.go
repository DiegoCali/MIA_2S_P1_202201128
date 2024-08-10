package structures

import (
	"bytes"
	"encoding/binary"
	"os"
)

// MBR Size in bytes: 165
type MBR struct {
	Size       int32        // bytes: 4
	TimeStamp  float32      // bytes: 4
	Signature  int32        // bytes: 4
	Fit        [1]byte      // bytes: 1
	Partitions [4]Partition // bytes: 38 * 4 = 152
}

func (mbr *MBR) Create(size int32, time float32, sign int32, fit string) error {
	mbr.Size = size
	mbr.TimeStamp = time
	mbr.Signature = sign
	copy(mbr.Fit[:], fit)
	// Initialize partitions
	for i := 0; i < 4; i++ {
		mbr.Partitions[i].Status = -1
		mbr.Partitions[i].Type[0] = '\x00'
		mbr.Partitions[i].Fit[0] = '\x00'
		mbr.Partitions[i].Start = -1
		mbr.Partitions[i].Size = -1
		copy(mbr.Partitions[i].Name[:], "----------------")
		mbr.Partitions[i].Correlative = -1
		copy(mbr.Partitions[i].Id[:], "----")
	}
	return nil
}

func (mbr *MBR) Serialize(path string) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			return
		}
	}(file)
	// Write MBR
	err = binary.Write(file, binary.LittleEndian, mbr)
	if err != nil {
		return err
	}
	return nil
}

func (mbr *MBR) Deserialize(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			return
		}
	}(file)
	// Read MBR
	buffer := make([]byte, 165)
	_, err = file.Read(buffer)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(buffer)
	err = binary.Read(reader, binary.LittleEndian, mbr)
	if err != nil {
		return err
	}
	return nil
}
