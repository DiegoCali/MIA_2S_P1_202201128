package commands

import (
	"backend/structures"
	"backend/utils"
	"fmt"
	"strings"
)

func Rep(id string, route string, name string, pathLs string) (string, error) {
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
	case "bm_inode":
		err := generateBitmap(path, id, route, true)
		if err != nil {
			return "Error creating report", err
		}
	case "bm_block":
		err := generateBitmap(path, id, route, false)
		if err != nil {
			return "Error creating report", err
		}
	case "file":
		err := generateFileReport(path, id, route, pathLs)
		if err != nil {
			return "Error creating report", err
		}
	case "ls":
		err := generateLsReport(path, id, route, pathLs)
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

func generateBitmap(path string, id string, route string, inode bool) error {
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
	// Generate bitmap
	err = spBlock.BitmapInodeTxt(route, path, inode)
	if err != nil {
		return err
	}
	return nil
}

func generateFileReport(path string, id string, route string, pathFile string) error {
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
	// Get path
	pathTokens := strings.Split(pathFile, "/")
	if pathTokens[0] == "" {
		pathTokens = pathTokens[1:]
	}
	// Get file str
	fileStr, err := spBlock.CatInode(pathTokens, path)
	if err != nil {
		return err
	}
	// Generate .txt file
	err = utils.GenerateTxt(route, fileStr)
	if err != nil {
		return err
	}
	return nil
}

func generateLsReport(path string, id string, route string, pathLs string) error {
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
	pathLsTokens := strings.Split(pathLs, "/")
	if pathLsTokens[0] == "" {
		pathLsTokens = pathLsTokens[1:]
	}
	// Generate LS string
	lsHtml, err := spBlock.LsReport(pathLsTokens, path)
	if err != nil {
		return err
	}
	// dot string
	dotStr := "digraph G {\n"
	dotStr += "rankdir=LR;\n"
	dotStr += "node [shape=record, label="
	dotStr += lsHtml
	dotStr += "];\n"
	dotStr += "}\n"
	// Generate file .dot with a node with a label with the lsHtml
	err = utils.GenerateDot(route, dotStr)
	if err != nil {
		return err
	}
	fmt.Println(dotStr)
	return nil
}

func checkIfIDExists(id string) (string, bool) {
	path, exists := utils.GlobalMounts[id]
	if !exists {
		return "", false
	}
	return path, true
}
