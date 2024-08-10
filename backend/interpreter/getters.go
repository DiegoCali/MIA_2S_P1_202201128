package interpreter

import (
	"fmt"
	"strconv"
)

func getDisk(opions []Option) (int, string, string, string, error) {
	size := -1
	fit := "FF"
	unit := "M"
	path := ""
	for _, option := range opions {
		switch option.Name {
		case "size":
			size, _ = strconv.Atoi(option.Value)
		case "fit":
			fit = option.Value[0:1]
		case "unit":
			unit = option.Value
		case "path":
			path = option.Value
		}
	}
	if size == -1 {
		return -1, "\x00", "\x00", "\x00", fmt.Errorf("invalid -size option")
	}
	if path == "" {
		return -1, "\x00", "\x00", "\x00", fmt.Errorf("missing -path option")
	}
	return size, fit, unit, path, nil
}

func getRDisk(opions []Option) (string, error) {
	path := ""
	for _, option := range opions {
		if option.Name == "path" {
			path = option.Value
		}
	}
	if path == "" {
		return "\x00", fmt.Errorf("missing -path option")
	}
	return "", nil
}
