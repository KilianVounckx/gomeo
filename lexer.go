package main

import (
	"strconv"
	"strings"
)

const DIGITS = "0123456789"
const LETTERS = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const LETTER_DIGITS = LETTERS + DIGITS

type Lexer struct {
	name     string
	text     string
	position *Position
	current  rune
}

func NewLexer(name, text string) *Lexer {
	res := &Lexer{name, text, NewPosition(-1, 0, -1, name, text), -1}
	res.Advance()
	return res
}

func (self *Lexer) Advance() {
	self.position.Advance(self.current)
	if self.position.index < len(self.text) {
		self.current = rune(self.text[self.position.index])
	} else {
		self.current = -1
	}
}

func (self *Lexer) MakeTokens() ([]*Token, *Error) {
	var tokens []*Token

	for self.current != -1 {
		switch self.current {
		case ' ', '\t':
			self.Advance()
		case '#':
			self.SkipComment()
		case '+':
			tokens = append(tokens, NewToken(PLUS, nil, self.position, nil))
			self.Advance()
		case '-':
			tokens = append(tokens, NewToken(MINUS, nil, self.position, nil))
			self.Advance()
		case '*':
			tokens = append(tokens, NewToken(MUL, nil, self.position, nil))
			self.Advance()
		case '/':
			tokens = append(tokens, NewToken(DIV, nil, self.position, nil))
			self.Advance()
		case '%':
			tokens = append(tokens, NewToken(MOD, nil, self.position, nil))
			self.Advance()
		case '^':
			tokens = append(tokens, NewToken(POW, nil, self.position, nil))
			self.Advance()
		case '(':
			tokens = append(tokens, NewToken(LPAREN, nil, self.position, nil))
			self.Advance()
		case ')':
			tokens = append(tokens, NewToken(RPAREN, nil, self.position, nil))
			self.Advance()
		case '[':
			tokens = append(tokens, NewToken(LBRACKET, nil, self.position, nil))
			self.Advance()
		case ']':
			tokens = append(tokens, NewToken(RBRACKET, nil, self.position, nil))
			self.Advance()
		case ',':
			tokens = append(tokens, NewToken(COMMA, nil, self.position, nil))
			self.Advance()
		case '\n', ';':
			tokens = append(tokens, NewToken(NEWLINE, nil, self.position, nil))
			self.Advance()
		case '!':
			tokens = append(tokens, self.MakeNotEqual())
		case '=':
			tokens = append(tokens, self.MakeEqual())
		case '<':
			tokens = append(tokens, self.MakeLess())
		case '>':
			tokens = append(tokens, self.MakeGreater())
		case '"':
			tokens = append(tokens, self.MakeString())
		case '|':
			token, err := self.MakeOr()
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, token)
		case '&':
			token, err := self.MakeAnd()
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, token)
		default:
			if strings.ContainsRune(DIGITS, self.current) {
				tokens = append(tokens, self.MakeNumber())
			} else if strings.ContainsRune(LETTERS, self.current) {
				tokens = append(tokens, self.MakeIdentifier())
			} else {
				start := self.position.Copy()
				character := self.current
				self.Advance()
				return nil, NewIllegalCharacterError(
					"'"+string(character)+"'", start, self.position.Copy(),
				)
			}
		}
	}

	tokens = append(tokens, NewToken(EOF, nil, self.position, nil))
	return tokens, nil
}

func (self *Lexer) MakeNumber() *Token {
	start := self.position.Copy()
	number := ""
	dotCount := 0

	for self.current != -1 && strings.ContainsRune(DIGITS+".", self.current) {
		if self.current == '.' {
			if dotCount >= 1 {
				break
			}
			dotCount++
			number += "."
		} else {
			number += string(self.current)
		}
		self.Advance()
	}

	res, _ := strconv.ParseFloat(number, 64)
	return NewToken(NUMBER, res, start, self.position)
}

func (self *Lexer) MakeIdentifier() *Token {
	start := self.position.Copy()
	word := ""

	for self.current != -1 && strings.ContainsRune(LETTER_DIGITS+"_", self.current) {
		word += string(self.current)
		self.Advance()
	}

	if InKeywords(word) {
		return NewToken(KEYWORD, word, start, self.position)
	}
	return NewToken(IDENTIFIER, word, start, self.position)
}

func (self *Lexer) MakeEqual() *Token {
	start := self.position.Copy()

	self.Advance()
	if self.current == '=' {
		self.Advance()
		return NewToken(EE, nil, start, self.position)
	}
	return NewToken(EQ, nil, start, self.position)
}

func (self *Lexer) MakeNotEqual() *Token {
	start := self.position.Copy()

	self.Advance()
	if self.current == '=' {
		self.Advance()
		return NewToken(NE, nil, start, self.position)
	}
	return NewToken(NOT, nil, start, self.position)
}

func (self *Lexer) MakeLess() *Token {
	start := self.position.Copy()

	self.Advance()
	if self.current == '=' {
		self.Advance()
		return NewToken(LE, nil, start, self.position)
	}
	return NewToken(LT, nil, start, self.position)
}

func (self *Lexer) MakeGreater() *Token {
	start := self.position.Copy()

	self.Advance()
	if self.current == '=' {
		self.Advance()
		return NewToken(GE, nil, start, self.position)
	}
	return NewToken(GT, nil, start, self.position)
}

func (self *Lexer) MakeOr() (*Token, *Error) {
	start := self.position.Copy()

	self.Advance()
	if self.current != '|' {
		return nil, NewExpectedCharacterError("'|'", start, self.position)
	}

	self.Advance()
	return NewToken(OR, nil, start, self.position), nil
}

func (self *Lexer) MakeAnd() (*Token, *Error) {
	start := self.position.Copy()

	self.Advance()
	if self.current != '&' {
		return nil, NewExpectedCharacterError("'&'", start, self.position)
	}

	self.Advance()
	return NewToken(AND, nil, start, self.position), nil
}
func (self *Lexer) MakeString() *Token {
	res := ""
	start := self.position.Copy()
	escape := false

	escape_characters := map[rune]string{
		'n':  "\n",
		't':  "\t",
		'r':  "\r",
		'\\': "\\",
		'"':  "\"",
	}

	self.Advance()
	for self.current != -1 && (self.current != '"' || escape) {
		if escape {
			next := escape_characters[self.current]
			if next == "" {
				res += string(self.current)
			} else {
				res += next
			}
			escape = false
		} else {
			if self.current == '\\' {
				escape = true
			} else {
				res += string(self.current)
			}
		}
		self.Advance()
	}

	self.Advance()
	return NewToken(STRING, res, start, self.position)
}

func (self *Lexer) SkipComment() {
	self.Advance()
	block := false

	if self.current == '=' {
		self.Advance()
		block = true
	}

	for self.current != -1 {
		if block {
			if self.current == '=' {
				self.Advance()
				if self.current == '#' {
					self.Advance()
					break
				}
			}
		} else {
			if self.current == '\n' || self.current == ';' {
				self.Advance()
				break
			}
		}
		self.Advance()
	}
}
