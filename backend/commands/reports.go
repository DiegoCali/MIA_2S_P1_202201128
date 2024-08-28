package commands

import (
	"backend/structures"
	"backend/utils"
	"fmt"
)

func Rep(id string, route string, name string) (string, error) {
	// Check if ID exists
	path, exists := checkIfIDExists(id)
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
	case "disk":
		err := generateDiskReport(path, route)
		if err != nil {
			return "Error creating report", err
		}
	case "sb":
		err := generateSuperBlockReport(path, id, route)
		if err != nil {
			return "Error creating report", err
		}
	case "inode":
		err := generateInodesReport(path, id, route)
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
	err = mbr.DotMbr(route, path)
	if err != nil {
		return err
	}
	return nil
}

func generateDiskReport(path string, route string) error {
	// Read MBR
	mbr := &structures.MBR{}
	err := utils.Deserialize(mbr, path, 0)
	if err != nil {
		return err
	}
	// Generate report
	mbr.Print()
	// Generate dot file
	err = mbr.DotDisk(route, path)
	if err != nil {
		return err
	}
	return nil
}

func generateSuperBlockReport(path string, id string, route string) error {
	// Read MBR
	mbr := &structures.MBR{}
	err := utils.Deserialize(mbr, path, 0)
	if err != nil {
		return err
	}
	partIndex, err := mbr.GetPartitionId(id)
	if err != nil {
		return err
	}
	offset := mbr.Partitions[partIndex].Start
	// Read SuperBlock
	spBlock := &structures.SuperBlock{}
	err = utils.Deserialize(spBlock, path, int(offset))
	if err != nil {
		return err
	}
	// Generate report
	err = utils.PrintStruct(spBlock)
	if err != nil {
		return err
	}
	// Generate dot file
	err = spBlock.SuperBlockDot(route, path)
	if err != nil {
		return err
	}
	return nil
}

func generateInodesReport(path string, id string, route string) error {
	// Read MBR
	mbr := &structures.MBR{}
	err := utils.Deserialize(mbr, path, 0)
	if err != nil {
		return err
	}
	partIndex, err := mbr.GetPartitionId(id)
	if err != nil {
		return err
	}
	offset := mbr.Partitions[partIndex].Start
	// Read SuperBlock
	spBlock := &structures.SuperBlock{}
	err = utils.Deserialize(spBlock, path, int(offset))
	if err != nil {
		return err
	}
	// Generate dot file
	err = spBlock.InodesDot(route, path)
	if err != nil {
		return err
	}
	return nil
}

func checkIfIDExists(id string) (string, bool) {
	path, exists := utils.GlobalMounts[id]
	if !exists {
		return "", false
	}
	return path, true
}
