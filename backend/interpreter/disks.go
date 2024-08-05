package interpreter

import (
	"fmt"
	"strconv"
)

type Disk struct {
	Size       [4]byte
	TimeStamp  [4]byte
	Signature  [4]byte
	Fit        [1]byte
	Partitions []Partition
}

type Partition struct {
	Status      [1]byte
	Type        [1]byte
	Fit         [1]byte
	Start       [4]byte
	Size        [4]byte
	Name        [16]byte
	Correlative [4]byte
	Id          [4]byte
}

func MkDisk(options []Option) (string, error) {
	var message string
	size := -1
	fit := "none"
	unit := "none"
	path := "none"
	for _, option := range options {
		if option.Name == "size" {
			//Parse string to int
			size, _ = strconv.Atoi(option.Value)
			continue
		}
		if option.Name == "fit" {
			fit = option.Value
			continue
		}
		if option.Name == "unit" {
			unit = option.Value
			continue
		}
		if option.Name == "path" {
			path = option.Value
			continue
		}
	}
	if size != -1 {
		message = "Disk created successfully, size: " + strconv.Itoa(size) + ", fit: " + fit + ", unit: " + unit + ", path: " + path
	} else {
		message = "Disk not created"
		return message, fmt.Errorf("-size is required")
	}
	return message, nil
}
