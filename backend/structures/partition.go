package structures

// Partition Size in bytes: 38
type Partition struct {
	Status      int32    // bytes: 4
	Type        [1]byte  // bytes: 1
	Fit         [1]byte  // bytes: 1
	Start       int32    // bytes: 4
	Size        int32    // bytes: 4
	Name        [16]byte // bytes: 16
	Correlative int32    // bytes: 4
	Id          [4]byte  // bytes: 4
}

func (part *Partition) GetLastEBR(path string) (int32, error) {
	offset := part.Start
	// Check if partition is extended
	if part.Type[0] != 'E' {
		return -1, nil
	}
	// Read EBR
	ebr := &EBR{}
	err := ebr.Deserialize(path, int(offset))
	if err != nil {
		return -1, err
	}
	// Get last EBR
	for ebr.Next != -1 {
		offset = ebr.Next
		err = ebr.Deserialize(path, int(offset))
		if err != nil {
			return -1, err
		}
	}
	return offset, nil
}
