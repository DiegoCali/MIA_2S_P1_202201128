package commands

import (
	"backend/structures"
	"backend/utils"
	"fmt"
)

func Rep(id string, path string, name string) (string, error) {
	// Check if ID exists
	route, exists := checkIfIDExists(id)
	if !exists {
		return "Error creating report", fmt.Errorf("ID %s does not exist", id)
	}
	// Generate report
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
	err := utils.Deserialize(mbr, path, 0)
	if err != nil {
		return err
	}
	// Generate report
	mbr.Print()
	// Generate dot file
	err = mbr.DotMbr(route)
	return nil
}

func checkIfIDExists(id string) (string, bool) {
	route := utils.GlobalMounts[id]
	if route == "" {
		return "", false
	}
	return route, true
}
