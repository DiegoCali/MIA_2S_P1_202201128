package structures

// DBlock Size in bytes: 64
// DBlock is a block for directories
type DBlock struct {
	Content [4]Content // bytes: 64
}

// Content Size in bytes: 16
type Content struct {
	Name   [12]byte // bytes: 12
	BInode int32    // bytes: 4
}
