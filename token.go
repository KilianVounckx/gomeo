package main

import (
	"fmt"
)

type TokenType string

const (
	NUMBER TokenType = "NUMBER"
	STRING TokenType = "STRING"

	PLUS  TokenType = "PLUS"
	MINUS TokenType = "MINUS"
	MUL   TokenType = "MUL"
	DIV   TokenType = "DIV"
	POW   TokenType = "POW"
	MOD   TokenType = "MOD"

	LPAREN TokenType = "LPAREN"
	RPAREN TokenType = "RPAREN"

	LBRACKET TokenType = "LBRACKET"
	RBRACKET TokenType = "RBRACKET"

	KEYWORD    TokenType = "KEYWORD"
	IDENTIFIER TokenType = "IDENTIFIER"
	EQ         TokenType = "EQ"

	EE TokenType = "EE"
	NE TokenType = "NE"
	LT TokenType = "LT"
	GT TokenType = "GT"
	LE TokenType = "LE"
	GE TokenType = "GE"

	NOT TokenType = "NOT"
	AND TokenType = "AND"
	OR  TokenType = "OR"

	COMMA TokenType = "COMMA"

	NEWLINE TokenType = "NEWLINE"
	EOF     TokenType = "EOF"
)

func KEYWORDS() []string {
	return []string{
		"var",
		"if", "do", "elseif", "else", "end",
		"while", "for", "from", "to", "step", "continue", "break",
		"function", "return",
	}
}

func InKeywords(s string) bool {
	for _, word := range KEYWORDS() {
		if s == word {
			return true
		}
	}
	return false
}

func (self TokenType) In(tokens []TokenType) bool {
	for _, token := range tokens {
		if self == token {
			return true
		}
	}
	return false
}

type Token struct {
	tokenType  TokenType
	value      interface{}
	start, end *Position
}

func NewToken(tokenType TokenType, value interface{}, start, end *Position) *Token {
	if start != nil {
		start = start.Copy()
		if end == nil {
			end = start.Copy()
			end.Advance(-1)
		}
	}
	if end != nil {
		end = end.Copy()
	}
	return &Token{tokenType, value, start, end}
}

func (self *Token) String() string {
	if self.value != nil {
		return fmt.Sprintf("%s:%v", string(self.tokenType), self.value)
	}
	return string(self.tokenType)
}

func (self *Token) Matches(tokenType TokenType, value interface{}) bool {
	return self.tokenType == tokenType && self.value == value
}
