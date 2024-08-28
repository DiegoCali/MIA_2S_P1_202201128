package structures

import (
	"backend/utils"
	"strconv"
)

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

func (inode *Inode) GetDotStr(id int) string {
	str := "inode_" + strconv.Itoa(id) + "[label=<<TABLE>\n"
	str += "<TR><TD>INODE</TD></TR>\n"
	str += "<TR><TD>i_uid</TD><TD>" + strconv.Itoa(int(inode.UID)) + "</TD></TR>\n"
	str += "<TR><TD>i_gid</TD><TD>" + strconv.Itoa(int(inode.GID)) + "</TD></TR>\n"
	str += "<TR><TD>i_size</TD><TD>" + strconv.Itoa(int(inode.Size)) + "</TD></TR>\n"
	str += "<TR><TD>i_atime</TD><TD>" + utils.Int64ToDate(inode.Atime) + "</TD></TR>\n"
	str += "<TR><TD>i_ctime</TD><TD>" + utils.Int64ToDate(inode.Ctime) + "</TD></TR>\n"
	str += "<TR><TD>i_mtime</TD><TD>" + utils.Int64ToDate(inode.Mtime) + "</TD></TR>\n"
	// Blocks
	for i := 0; i < 15; i++ {
		str += "<TR><TD>i_block_" + strconv.Itoa(i+1) + "</TD><TD>" + strconv.Itoa(int(inode.Block[i])) + "</TD></TR>\n"
	}
	str += "<TR><TD>i_type</TD><TD>" + string(inode.Type[:]) + "</TD></TR>\n"
	str += "<TR><TD>i_perm</TD><TD>" + string(inode.Perm[:]) + "</TD></TR>\n"
	str += "</TABLE>>];\n"
	return str
}
