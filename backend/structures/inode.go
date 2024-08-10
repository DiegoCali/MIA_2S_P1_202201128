package structures

// Inode Size in bytes: 88
type Inode struct {
	UID   int32     // bytes: 4
	GID   int32     // bytes: 4
	Size  int32     // bytes: 4
	Atime float32   // bytes: 4
	Ctime float32   // bytes: 4
	Mtime float32   // bytes: 4
	Block [15]int32 // bytes: 60
	Type  [1]byte   // bytes: 1
	Perm  [3]byte   // bytes: 3
}
