package structures

import (
	"backend/utils"
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

func (part *Partition) GetPartitionStr() string {
	str := "<TR><TD>PARTITION</TD></TR>\n"
	str += "<TR><TD>part_status</TD><TD>" + strconv.Itoa(int(part.Status)) + "</TD></TR>\n"
	str += "<TR><TD>part_type</TD><TD>" + utils.CheckNull(part.Type[:]) + "</TD></TR>\n"
	str += "<TR><TD>part_fit</TD><TD>" + utils.CheckNull(part.Fit[:]) + "</TD></TR>\n"
	str += "<TR><TD>part_start</TD><TD>" + strconv.Itoa(int(part.Start)) + "</TD></TR>\n"
	str += "<TR><TD>part_size</TD><TD>" + strconv.Itoa(int(part.Size)) + "</TD></TR>\n"
	str += "<TR><TD>part_name</TD><TD>" + utils.CheckNull(part.Name[:]) + "</TD></TR>\n"
	return str
}
