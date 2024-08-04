package interpreter

import "unicode"

type Token struct {
	kind  string
	value string
}

func Lex(input string) ([]Token, error) {
	var tokens []Token
	pos := 0
	for pos < len(input) {
		character := input[pos]
		// skip whitespace
		if character == ' ' || character == '\n' || character == '\t' {
			pos++
			continue
		}
		// number is a sequence of digits
		if unicode.IsDigit(rune(character)) {
			start := pos
			for pos < len(input) && unicode.IsDigit(rune(input[pos])) {
				pos++
			}
			tokens = append(tokens, Token{"NUMBER", input[start:pos]})
			continue
		}
		// options are words with a dash before them
		if character == '-' {
			start := pos
			pos++
			for pos < len(input) && unicode.IsLetter(rune(input[pos])) {
				pos++
			}
			tokens = append(tokens, Token{"OPTION", input[start:pos]})
			continue
		}
		// command is just a word
		if unicode.IsLetter(rune(character)) {
			start := pos
			for pos < len(input) && (unicode.IsLetter(rune(input[pos])) || unicode.IsDigit(rune(input[pos]))) {
				pos++
			}
			tokens = append(tokens, Token{"COMMAND", input[start:pos]})
			continue
		}
		// value is any character after an equal sign and before a space
		if character == '=' {
			pos++
			start := pos
			// stop until whitespace
			for pos < len(input) && input[pos] != ' ' && input[pos] != '\n' && input[pos] != '\t' {
				pos++
			}
			tokens = append(tokens, Token{"VALUE", input[start:pos]})
			continue
		}
		// any other character is invalid
		tokens = append(tokens, Token{"INVALID", string(character)})
	}
	return tokens, nil
}
