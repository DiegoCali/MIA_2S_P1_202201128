package interpreter

import (
	"fmt"
	"strconv"
	"strings"
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
			fit = strings.ToUpper(option.Value[0:1])
		case "unit":
			unit = strings.ToUpper(option.Value)
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
			unit = strings.ToUpper(option.Value)
		case "path":
			path = option.Value
		case "type":
			typeP = strings.ToUpper(option.Value)
		case "fit":
			fit = strings.ToUpper(option.Value[0:1])
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

func getMount(options []Option) (string, string, error) {
	path := ""
	name := ""
	for _, option := range options {
		switch option.Name {
		case "path":
			path = option.Value
		case "name":
			name = option.Value
		default:
			return "\x00", "\x00", fmt.Errorf("invalid option %s", option.Name)
		}
	}
	if path == "" {
		return "\x00", "\x00", fmt.Errorf("missing -path option")
	}
	if name == "" {
		return "\x00", "\x00", fmt.Errorf("missing -name option")
	}
	return path, name, nil
}

func getRep(options []Option) (string, string, string, string, error) {
	id := ""
	path := ""
	name := ""
	route := ""
	for _, option := range options {
		switch option.Name {
		case "id":
			id = option.Value
		case "path":
			path = option.Value
		case "name":
			name = option.Value
		case "route":
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

func getFileSys(options []Option) (string, string, error) {
	id := ""
	typeF := "FULL"
	for _, option := range options {
		switch option.Name {
		case "id":
			id = option.Value
		case "type":
			typeF = strings.ToUpper(option.Value)
		default:
			return "\x00", "\x00", fmt.Errorf("invalid option %s", option.Name)
		}
	}
	if id == "" {
		return "\x00", "\x00", fmt.Errorf("missing -id option")
	}
	return id, typeF, nil
}

func getCat(options []Option) ([]string, error) {
	files := make([]string, 0)
	for _, option := range options {
		if option.Name == "file" {
			files = append(files, option.Value)
		} else {
			return nil, fmt.Errorf("invalid option %s", option.Name)
		}
	}
	return files, nil
}

func getLogin(option []Option) (string, string, string, error) {
	user := ""
	password := ""
	id := ""
	for _, opt := range option {
		switch opt.Name {
		case "user":
			user = opt.Value
		case "pass":
			password = opt.Value
		case "id":
			id = opt.Value
		default:
			return "\x00", "\x00", "\x00", fmt.Errorf("invalid option %s", opt.Name)
		}
	}
	if user == "" {
		return "\x00", "\x00", "\x00", fmt.Errorf("missing -user option")
	}
	if password == "" {
		return "\x00", "\x00", "\x00", fmt.Errorf("missing -pass option")
	}
	if id == "" {
		return "\x00", "\x00", "\x00", fmt.Errorf("missing -id option")
	}
	return user, password, id, nil
}

func getMkGroup(options []Option) (string, error) {
	name := ""
	for _, option := range options {
		if option.Name == "name" {
			name = option.Value
		} else {
			return "\x00", fmt.Errorf("invalid option %s", option.Name)
		}
	}
	if name == "" {
		return "\x00", fmt.Errorf("missing -name option")
	}
	return name, nil
}

func getRmGroup(option []Option) (string, error) {
	name := ""
	for _, opt := range option {
		if opt.Name == "name" {
			name = opt.Value
		} else {
			return "\x00", fmt.Errorf("invalid option %s", opt.Name)
		}
	}
	if name == "" {
		return "\x00", fmt.Errorf("missing -name option")
	}
	return name, nil
}

func getMkUser(option []Option) (string, string, string, error) {
	user := ""
	password := ""
	group := ""
	for _, opt := range option {
		switch opt.Name {
		case "user":
			user = opt.Value
		case "pass":
			password = opt.Value
		case "grp":
			group = opt.Value
		default:
			return "\x00", "\x00", "\x00", fmt.Errorf("invalid option %s", opt.Name)
		}
	}
	if user == "" {
		return "\x00", "\x00", "\x00", fmt.Errorf("missing -user option")
	}
	if password == "" {
		return "\x00", "\x00", "\x00", fmt.Errorf("missing -pass option")
	}
	if group == "" {
		return "\x00", "\x00", "\x00", fmt.Errorf("missing -group option")
	}
	return user, password, group, nil
}

func getRmUser(option []Option) (string, error) {
	user := ""
	for _, opt := range option {
		if opt.Name == "user" {
			user = opt.Value
		} else {
			return "\x00", fmt.Errorf("invalid option %s", opt.Name)
		}
	}
	if user == "" {
		return "\x00", fmt.Errorf("missing -user option")
	}
	return user, nil
}

func getChgrp(option []Option) (string, string, error) {
	user := ""
	group := ""
	for _, opt := range option {
		switch opt.Name {
		case "user":
			user = opt.Value
		case "grp":
			group = opt.Value
		default:
			return "\x00", "\x00", fmt.Errorf("invalid option %s", opt.Name)
		}
	}
	if user == "" {
		return "\x00", "\x00", fmt.Errorf("missing -usr option")
	}
	if group == "" {
		return "\x00", "\x00", fmt.Errorf("missing -grp option")
	}
	return user, group, nil
}

func getMkFile(option []Option) (string, bool, int, string, error) {
	path := ""
	createParents := false
	size := 0
	cont := ""
	for _, opt := range option {
		switch opt.Name {
		case "path":
			path = opt.Value
		case "r":
			createParents = opt.Value == "true"
		case "size":
			size, _ = strconv.Atoi(opt.Value)
		case "cont":
			cont = opt.Value
		default:
			return "\x00", false, -1, "\x00", fmt.Errorf("invalid option %s", opt.Name)
		}
	}
	if path == "" {
		return "\x00", false, -1, "\x00", fmt.Errorf("missing -path option")
	}
	if size < 0 {
		return "\x00", false, -1, "\x00", fmt.Errorf("missing -size option")
	}
	return path, createParents, size, cont, nil
}

func getMkDir(option []Option) (string, bool, error) {
	createParents := false
	newDir := ""
	for _, opt := range option {
		switch opt.Name {
		case "path":
			newDir = opt.Value
		case "p":
			createParents = opt.Value == "true"
		default:
			return "\x00", false, fmt.Errorf("invalid option %s", opt.Name)
		}
	}
	if newDir == "" {
		return "\x00", false, fmt.Errorf("missing -path option")
	}
	return newDir, createParents, nil
}
