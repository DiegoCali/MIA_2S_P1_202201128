package interpreter

import (
	cmds "backend/commands"
	"backend/utils"
	"fmt"
)

type Instruction struct {
	command string
	options []Option
}

type Stack struct {
	instruction []Instruction
}

type Option struct {
	Name  string
	Value string
}

func Parse(tokens []Token) (Stack, error) {
	var root Stack
	root.instruction = make([]Instruction, 0)
	pos := 0
	for pos < len(tokens) {
		if tokens[pos].kind != "COMMAND" {
			return root, fmt.Errorf("expected COMMAND, got %s", tokens[pos].kind)
		}
		instruction, newPos, err := readCommand(tokens, pos)
		if err != nil {
			return root, err
		}
		root.instruction = append(root.instruction, instruction)
		pos = newPos
	}
	return root, nil
}

func readCommand(tokens []Token, pos int) (Instruction, int, error) {
	var instruction Instruction
	instruction.command = tokens[pos].value
	pos++
	// Check if there are options or not
	if tokens[pos].kind != "OPTION" {
		if tokens[pos].kind == "TERMINATOR" {
			pos++
			return instruction, pos, nil
		}
		return instruction, pos, fmt.Errorf("expected OPTION, got %s", tokens[pos].kind)
	}
	// Read options
	for tokens[pos].kind == "OPTION" {
		option, newPos, err := readOption(tokens, pos)
		if err != nil {
			return instruction, pos, err
		}
		instruction.options = append(instruction.options, option)
		pos = newPos
	}
	if tokens[pos].kind != "TERMINATOR" {
		return instruction, pos, fmt.Errorf("expected TERMINATOR, got %s", tokens[pos].kind)
	}
	pos++
	return instruction, pos, nil
}

func readOption(tokens []Token, pos int) (Option, int, error) {
	var option Option
	option.Name = tokens[pos].value
	pos++
	if tokens[pos].kind != "VALUE" {
		return option, pos, fmt.Errorf("expected VALUE, got %s", tokens[pos].kind)
	}
	option.Value = tokens[pos].value
	pos++
	return option, pos, nil
}

func Execute(root Stack) (string, error) {
	// Execute the instructions
	var output string
	for _, instruction := range root.instruction {
		// Execute the instruction
		if instruction.command == "mkdisk" {
			size, fit, unit, path, err := getDisk(instruction.options)
			if err != nil {
				return "Error reading options", err
			}
			message, err := cmds.MkDisk(size, fit, unit, path)
			if err != nil {
				return output, err
			} else {
				output += message + "\n"
			}
			continue
		}
		if instruction.command == "rmdisk" {
			path, err := getRDisk(instruction.options)
			if err != nil {
				return "Error reading options", err
			}
			message, err := cmds.RmDisk(path)
			if err != nil {
				return output, err
			} else {
				output += message + "\n"
			}
			continue
		}
		if instruction.command == "fdisk" {
			size, unit, path, typeP, fit, name, err := getPartition(instruction.options)
			if err != nil {
				return "Error reading options", err
			}
			message, err := cmds.FDisk(size, unit, path, typeP, fit, name)
			if err != nil {
				return output, err
			} else {
				output += message + "\n"
			}
			continue
		}
		if instruction.command == "mount" {
			path, name, err := getMount(instruction.options)
			if err != nil {
				return "Error reading options", err
			}
			message, err := cmds.Mount(path, name)
			if err != nil {
				return output, err
			} else {
				output += message + "\n"
			}
			continue
		}
		if instruction.command == "mkfs" {
			id, typeF, err := getFileSys(instruction.options)
			if err != nil {
				return "Error reading options", err
			}
			message, err := cmds.MkFS(id, typeF)
			if err != nil {
				return output, err
			} else {
				output += message + "\n"
			}
			continue
		}
		if instruction.command == "cat" {
			files, err := getCat(instruction.options)
			if err != nil {
				return "Error reading options", err
			}
			message, err := cmds.Cat(files)
			if err != nil {
				return output, err
			} else {
				output += message + "\n"
			}
			continue
		}
		if instruction.command == "login" {
			name, pass, id, err := getLogin(instruction.options)
			if err != nil {
				return "Error reading options", err
			}
			message, err := cmds.Login(name, pass, id)
			if err != nil {
				return output, err
			} else {
				output += message + "\n"
			}
			continue
		}
		if instruction.command == "logout" {
			userName := utils.ActualUser.GetName()
			userId := utils.ActualUser.GetId()
			if userName == "" || userId == "" {
				return "Error logging out", fmt.Errorf("no user logged in")
			}
			utils.ActualUser.Set("", "")
			output += "Logged out from: [" + userName + "] successfully\n"
			continue
		}
		if instruction.command == "mkgrp" {
			name, err := getMkGroup(instruction.options)
			if err != nil {
				return "Error reading options", err
			}
			message, err := cmds.MkGroup(name)
			if err != nil {
				return output, err
			} else {
				output += message + "\n"
			}
			continue
		}
		if instruction.command == "rmgrp" {
			name, err := getRmGroup(instruction.options)
			if err != nil {
				return "Error reading options", err
			}
			message, err := cmds.RmGroup(name)
			if err != nil {
				return output, err
			} else {
				output += message + "\n"
			}
			continue
		}
		if instruction.command == "rep" {
			id, route, name, err := getRep(instruction.options)
			if err != nil {
				return "Error reading options", err
			}
			message, err := cmds.Rep(id, route, name)
			if err != nil {
				return output, err
			} else {
				output += message + "\n"
			}
			continue
		}
	}
	return output, nil
}
