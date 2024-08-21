package structures

import (
	"backend/utils"
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
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

func (mbr *MBR) Serialize(path string) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			return
		}
	}(file)
	// Write MBR
	err = binary.Write(file, binary.LittleEndian, mbr)
	if err != nil {
		return err
	}
	return nil
}

func (mbr *MBR) Deserialize(path string) error {
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
	// Read MBR
	buffer := make([]byte, 169)
	_, err = file.Read(buffer)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(buffer)
	err = binary.Read(reader, binary.LittleEndian, mbr)
	if err != nil {
		return err
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

func (mbr *MBR) DotMbr(output string) error {
	file, err := os.Create(output + ".dot")
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			return
		}
	}(file)
	// Write MBR
	strFile := "digraph MBR {\n"
	strFile += "node [shape=record];\n"
	strFile += "MBR [label=<<TABLE>\n"
	strFile += "<TR><TD>mbr_size</TD><TD>" + strconv.Itoa(int(mbr.Size)) + "</TD></TR>\n"
	strFile += "<TR><TD>mbr_date_creation</TD><TD>" + utils.Int64ToDate(mbr.TimeStamp) + "</TD></TR>\n"
	strFile += "<TR><TD>mbr_disk_signature</TD><TD>" + strconv.Itoa(int(mbr.Signature)) + "</TD></TR>\n"
	// Partitions
	for i := 0; i < 4; i++ {
		strFile += "<TR><TD>part_status</TD><TD>" + strconv.Itoa(int(mbr.Partitions[i].Status)) + "</TD></TR>\n"
		// TODO: check if partition is extended
		// check if type and fit is \x00
		if mbr.Partitions[i].Type[0] == '\x00' {
			strFile += "<TR><TD>part_type</TD><TD>none</TD></TR>\n"
		} else {
			strFile += "<TR><TD>part_type</TD><TD>" + string(mbr.Partitions[i].Type[:]) + "</TD></TR>\n"
		}
		if mbr.Partitions[i].Fit[0] == '\x00' {
			strFile += "<TR><TD>part_fit</TD><TD>none</TD></TR>\n"
		} else {
			strFile += "<TR><TD>part_fit</TD><TD>" + string(mbr.Partitions[i].Fit[:]) + "</TD></TR>\n"
		}
		strFile += "<TR><TD>part_start</TD><TD>" + strconv.Itoa(int(mbr.Partitions[i].Start)) + "</TD></TR>\n"
		strFile += "<TR><TD>part_size</TD><TD>" + strconv.Itoa(int(mbr.Partitions[i].Size)) + "</TD></TR>\n"
		// check if name contains \x00
		namestr := string(mbr.Partitions[i].Name[:])
		for j := 0; j < len(namestr); j++ {
			if namestr[j] == '\x00' {
				namestr = namestr[:j]
				break
			}
		}
		if namestr == "" {
			namestr = "none"
		}
		strFile += "<TR><TD>part_name</TD><TD>" + namestr + "</TD></TR>\n"
	}
	strFile += "</TABLE>>];\n"
	strFile += "}\n"
	_, err = file.WriteString(strFile)
	if err != nil {
		return err
	}
	return nil
}

func (mbr *MBR) GetPartitionIndex(name string) (int, error) {
	for i := 0; i < 4; i++ {
		if string(mbr.Partitions[i].Name[:]) == name {
			return i, nil
		}
	}
	return -1, fmt.Errorf("error: partition %s not found", name)
}
