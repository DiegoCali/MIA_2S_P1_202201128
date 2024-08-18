package interpreter

import (
	"fmt"
	"strconv"
)

func getDisk(options []Option) (int, string, string, string, error) {
	size := -1
	fit := "FF"
	unit := "M"
	path := ""
	for _, option := range options {
		switch option.Name {
		case "size":
			size, _ = strconv.Atoi(option.Value)
		case "fit":
			fit = option.Value[0:1]
		case "unit":
			unit = option.Value
		case "path":
			path = option.Value
		default:
			return -1, "\x00", "\x00", "\x00", fmt.Errorf("invalid option %s", option.Name)
		}
	}
	if size <= 0 {
		return -1, "\x00", "\x00", "\x00", fmt.Errorf("invalid -size option")
	}
	if path == "" {
		return -1, "\x00", "\x00", "\x00", fmt.Errorf("missing -path option")
	}
	return size, fit, unit, path, nil
}

func getRDisk(options []Option) (string, error) {
	path := ""
	for _, option := range options {
		if option.Name == "path" {
			path = option.Value
		} else {
			return "\x00", fmt.Errorf("invalid option %s", option.Name)
		}
	}
	if path == "" {
		return "\x00", fmt.Errorf("missing -path option")
	}
	return path, nil
}

func getPartition(options []Option) (int, string, string, string, string, string, error) {
	size := -1
	unit := "K"
	path := ""
	typeP := "P"
	fit := "WF"
	name := ""
	for _, option := range options {
		switch option.Name {
		case "size":
			size, _ = strconv.Atoi(option.Value)
		case "unit":
			unit = option.Value
		case "path":
			path = option.Value
		case "type":
			typeP = option.Value
		case "fit":
			fit = option.Value[0:1]
		case "name":
			name = option.Value
		default:
			return -1, "\x00", "\x00", "\x00", "\x00", "\x00", fmt.Errorf("invalid option %s", option.Name)
		}
	}
	if size <= 0 {
		return -1, "\x00", "\x00", "\x00", "\x00", "\x00", fmt.Errorf("invalid -size option")
	}
	if path == "" {
		return -1, "\x00", "\x00", "\x00", "\x00", "\x00", fmt.Errorf("missing -path option")
	}
	if name == "" {
		return -1, "\x00", "\x00", "\x00", "\x00", "\x00", fmt.Errorf("missing -name option")
	}
	return size, unit, path, typeP, fit, name, nil
}

func getRep(options []Option) (string, string, string, string, error) {
	id := ""
	path := ""
	name := ""
	route := "/home/diego/Documents/report"
	for _, option := range options {
		switch option.Name {
		case "id":
			id = option.Value
		case "path":
			path = option.Value
		case "name":
			name = option.Value
		case "ruta":
			route = option.Value
		default:
			return "\x00", "\x00", "\x00", "\x00", fmt.Errorf("invalid option %s", option.Name)
		}
	}
	if id == "" {
		return "\x00", "\x00", "\x00", "\x00", fmt.Errorf("missing -id option")
	}
	if path == "" {
		return "\x00", "\x00", "\x00", "\x00", fmt.Errorf("missing -path option")
	}
	if name == "" {
		return "\x00", "\x00", "\x00", "\x00", fmt.Errorf("missing -name option")
	}
	return id, path, name, route, nil
}
