package interpreter

import (
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
	if tokens[pos].kind != "OPTION" {
		return instruction, pos, fmt.Errorf("expected OPTION, got %s", tokens[pos].kind)
	}
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
			message, err := MkDisk(instruction.options)
			if err != nil {
				return output, err
			} else {
				output += message + "\n"
			}
		}
	}
	return output, nil
}
