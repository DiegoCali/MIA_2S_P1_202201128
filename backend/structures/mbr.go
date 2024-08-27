package structures

import (
	"backend/utils"
	"fmt"
	"strconv"
)

// MBR Size in bytes: 169
type MBR struct {
	Size       int32        // bytes: 4
	TimeStamp  int64        // bytes: 8
	Signature  int32        // bytes: 4
	Fit        [1]byte      // bytes: 1
	Partitions [4]Partition // bytes: 38 * 4 = 152
}

func (mbr *MBR) Set(size int32, time int64, sign int32, fit string) error {
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
		copy(mbr.Partitions[i].Name[:], "")
		mbr.Partitions[i].Correlative = -1
		copy(mbr.Partitions[i].Id[:], "")
	}
	return nil
}

func (mbr *MBR) Print() {
	println("Size: ", mbr.Size)
	println("TimeStamp: ", utils.Int64ToDate(mbr.TimeStamp))
	println("Signature: ", mbr.Signature)
	println("Fit: ", string(mbr.Fit[:]))
	println("----------------PARTITIONS--------------------")
	for i := 0; i < 4; i++ {
		println("---------------------------------------------")
		println("Partition ", i+1)
		println("Status: ", mbr.Partitions[i].Status)
		println("Type: ", string(mbr.Partitions[i].Type[:]))
		println("Fit: ", string(mbr.Partitions[i].Fit[:]))
		println("Start: ", mbr.Partitions[i].Start)
		println("Size: ", mbr.Partitions[i].Size)
		println("Name: ", string(mbr.Partitions[i].Name[:]))
		println("Correlative: ", mbr.Partitions[i].Correlative)
		println("Id: ", string(mbr.Partitions[i].Id[:]))
		println("---------------------------------------------")
	}
}

func (mbr *MBR) GetPartitionIndex(name string) (int, error) {
	for i := 0; i < 4; i++ {
		// Remove null characters
		partName := utils.CheckNull(mbr.Partitions[i].Name[:])
		if partName == name {
			return i, nil
		}
	}
	return -1, fmt.Errorf("error: partition %s not found", name)
}

func (mbr *MBR) GetPartitionId(id string) (int, error) {
	for i := 0; i < 4; i++ {
		// Remove null characters
		partId := utils.CheckNull(mbr.Partitions[i].Id[:])
		if partId == id {
			return i, nil
		}
	}
	return -1, fmt.Errorf("error: partition %s not found", id)
}

func getFreeSpace(mbr *MBR) int {
	freeSpace := mbr.Size - 169
	for i := 0; i < 4; i++ {
		if mbr.Partitions[i].Status != -1 {
			freeSpace -= mbr.Partitions[i].Size
		}
	}
	return int(freeSpace)
}

func (mbr *MBR) DotMbr(output string, path string) error {
	// Write MBR
	strFile := "digraph MBR {\n"
	strFile += "node [shape=record];\n"
	strFile += "MBR [label=<<TABLE>\n"
	strFile += "<TR><TD>mbr_size</TD><TD>" + strconv.Itoa(int(mbr.Size)) + "</TD></TR>\n"
	strFile += "<TR><TD>mbr_date_creation</TD><TD>" + utils.Int64ToDate(mbr.TimeStamp) + "</TD></TR>\n"
	strFile += "<TR><TD>mbr_disk_signature</TD><TD>" + strconv.Itoa(int(mbr.Signature)) + "</TD></TR>\n"
	// Partitions
	for i := 0; i < 4; i++ {
		strFile += mbr.Partitions[i].GetPartitionStr(path)
	}
	strFile += "</TABLE>>];\n"
	strFile += "}\n"
	err := utils.GenerateDot(output, strFile)
	if err != nil {
		return err
	}
	return nil
}

func (mbr *MBR) DotDisk(output string, path string) error {
	// cleaned path
	cleanPath := utils.CleanPath(path)
	// Disk size
	diskSize := int(mbr.Size)
	// Write Disk
	strFile := "digraph Disk {\n"
	// Label with Disk Name
	strFile += "rankdir=LR;\n"
	strFile += "title [label=\"" + cleanPath + "\", shape=\"plaintext\"];\n"
	strFile += "node [shape=record];\n"
	strFile += "Disk [label=<<TABLE>\n"
	// MBR
	strFile += "<TR>\n<TD>MBR</TD>\n"
	// Partitions
	for i := 0; i < 4; i++ {
		if mbr.Partitions[i].Status != -1 {
			strFile += mbr.Partitions[i].GetDiskStr(path)
		}
	}
	// Free Space
	freeSpace := getFreeSpace(mbr)
	if freeSpace > 0 {
		strFile += "<TD>FREE\n" + strconv.Itoa(freeSpace/diskSize) + "</TD>\n"
	}
	strFile += "</TR></TABLE>>];\n"
	strFile += "}\n"
	err := utils.GenerateDot(output, strFile)
	if err != nil {
		return err
	}
	return nil
}
