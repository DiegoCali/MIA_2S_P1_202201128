package structures

// Inode Size in bytes: 100
type Inode struct {
	UID   int32     // bytes: 4
	GID   int32     // bytes: 4
	Size  int32     // bytes: 4
	Atime int64     // bytes: 8
	Ctime int64     // bytes: 8
	Mtime int64     // bytes: 8
	Block [15]int32 // bytes: 60
	Type  [1]byte   // bytes: 1
	Perm  [3]byte   // bytes: 3
}
