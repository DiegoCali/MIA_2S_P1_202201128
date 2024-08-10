package structures

// PBlock Size in bytes: 64
// PBlock is a block for pointers
type PBlock struct {
	Pointers [16]int32 // bytes: 64
}
