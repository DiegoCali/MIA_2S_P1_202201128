package commands

import (
	"backend/structures"
	"fmt"
)

func Rep(id string, path string, name string, route string) (string, error) {
	switch name {
	case "mbr":
		err := generateMBRReport(path, route)
		if err != nil {
			return "Error creating report", err
		}
	default:
		return "Error creating report", fmt.Errorf("not implemented yet: %s", name)
	}
	return "Report created succesfully...", nil
}

func generateMBRReport(path string, route string) error {
	// Read MBR
	mbr := &structures.MBR{}
	err := mbr.Deserialize(path)
	if err != nil {
		return err
	}
	// Generate report
	mbr.Print()
	// Generate dot file
	err = mbr.DotMbr(route)
	return nil
}

func checkIfIDExists(mbr *structures.MBR, id string) bool {
	for _, partition := range mbr.Partitions {
		if string(partition.Id[:]) == id {
			return true
		}
	}
	return false
}
