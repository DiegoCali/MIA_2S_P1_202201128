package structures

import (
	"encoding/binary"
	"os"
)

// EBR Size in bytes: 33
type EBR struct {
	Mount int32    // bytes: 4
	Fit   [1]byte  // bytes: 1
	Start int32    // bytes: 4
	Size  int32    // bytes: 4
	Next  int32    // bytes: 4
	Name  [16]byte // bytes: 16
}

func (ebr *EBR) Set(mount int32, fit string, start int32, size int32, next int32, name string) error {
	ebr.Mount = mount
	copy(ebr.Fit[:], fit)
	ebr.Start = start
	ebr.Size = size
	ebr.Next = next
	copy(ebr.Name[:], name)
	return nil
}

func (ebr *EBR) Serialize(path string, offset int) error {
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
	// Move pointer to EBR offset
	_, err = file.Seek(int64(offset), 0)
	if err != nil {
		return err
	}
	// Write EBR
	err = binary.Write(file, binary.LittleEndian, ebr)
	if err != nil {
		return err
	}
	return nil
}

func (ebr *EBR) Deserialize(path string, offset int) error {
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
	// Move pointer to EBR offset
	_, err = file.Seek(int64(offset), 0)
	if err != nil {
		return err
	}
	// Read EBR
	err = binary.Read(file, binary.LittleEndian, ebr)
	if err != nil {
		return err
	}
	return nil
}
