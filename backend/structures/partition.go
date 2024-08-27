package structures

import (
	"backend/utils"
	"fmt"
	"strconv"
)

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
	err := utils.Deserialize(ebr, path, int(offset))
	if err != nil {
		return -1, err
	}
	// Get last EBR
	for ebr.Next != -1 {
		offset = ebr.Next
		err = utils.Deserialize(ebr, path, int(offset))
		if err != nil {
			return -1, err
		}
	}
	return offset, nil
}

func (part *Partition) GetPartitionStr(path string) string {
	str := "<TR><TD>PARTITION</TD></TR>\n"
	str += "<TR><TD>part_status</TD><TD>" + strconv.Itoa(int(part.Status)) + "</TD></TR>\n"
	str += "<TR><TD>part_type</TD><TD>" + utils.CheckNull(part.Type[:]) + "</TD></TR>\n"
	str += "<TR><TD>part_fit</TD><TD>" + utils.CheckNull(part.Fit[:]) + "</TD></TR>\n"
	str += "<TR><TD>part_start</TD><TD>" + strconv.Itoa(int(part.Start)) + "</TD></TR>\n"
	str += "<TR><TD>part_size</TD><TD>" + strconv.Itoa(int(part.Size)) + "</TD></TR>\n"
	str += "<TR><TD>part_name</TD><TD>" + utils.CheckNull(part.Name[:]) + "</TD></TR>\n"
	// Check if partition is extended
	if part.Type[0] == 'E' {
		str += part.GetEBRStr(path)
	}
	return str
}

func (part *Partition) GetEBRStr(path string) string {
	strEBR := ""
	// Get first EBR
	offset := part.Start
	ebr := &EBR{}
	for {
		err := utils.Deserialize(ebr, path, int(offset))
		if err != nil {
			strEBR += "<TR><TD>Somethin went wrong...</TD></TR>\n"
			break
		}
		fmt.Println(ebr)
		strEBR += "<TR><TD>EBR, Logic Partition</TD></TR>\n"
		strEBR += "<TR><TD>part_status</TD><TD>" + strconv.Itoa(int(ebr.Mount)) + "</TD></TR>\n"
		strEBR += "<TR><TD>part_fit</TD><TD>" + utils.CheckNull(ebr.Fit[:]) + "</TD></TR>\n"
		strEBR += "<TR><TD>part_start</TD><TD>" + strconv.Itoa(int(ebr.Start)) + "</TD></TR>\n"
		strEBR += "<TR><TD>part_size</TD><TD>" + strconv.Itoa(int(ebr.Size)) + "</TD></TR>\n"
		strEBR += "<TR><TD>part_name</TD><TD>" + utils.CheckNull(ebr.Name[:]) + "</TD></TR>\n"
		// Check if there is another EBR
		if ebr.Next == -1 {
			break
		}
		offset = ebr.Next
	}
	return strEBR
}

func (part *Partition) GetDiskStr(path string) string {
	str := "<TD>" + utils.CheckNull(part.Name[:]) + "</TD>\n"
	if part.Type[0] == 'E' {
		nested, numParts := part.GetEBRDiskStr(path)
		str = "<TD>\n"
		str += "<TABLE>\n"
		str += "<TR><TD colspan='" + strconv.Itoa(numParts) + "'>" + utils.CheckNull(part.Name[:]) + "</TD></TR>\n"
		str += nested
		str += "</TABLE>\n"
		str += "</TD>\n"
	}
	return str
}

func (part *Partition) GetEBRDiskStr(path string) (string, int) {
	numParts := 0
	str := "<TR>"
	// Get first EBR
	offset := part.Start
	ebr := &EBR{}
	for {
		err := utils.Deserialize(ebr, path, int(offset))
		if err != nil {
			str = "<TR><TD>Somethin went wrong...</TD></TR>\n"
			break
		}
		if ebr.Mount != -1 {
			str += "<TD>" + utils.CheckNull(ebr.Name[:]) + "</TD>\n"
		}
		numParts++
		// Check if there is another EBR
		if ebr.Next == -1 {
			str += "<TD>FREE</TD>\n</TR>\n"
			break
		}
		offset = ebr.Next
	}
	return str, numParts
}
