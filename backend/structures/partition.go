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
