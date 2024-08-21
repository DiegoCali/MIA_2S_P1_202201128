package structures

// SuperBlock Size in bytes: 68
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
