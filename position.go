package main

type Position struct {
	index        int
	line, column int
	name, text   string
}

func NewPosition(index, line, column int, name, text string) *Position {
	return &Position{index, line, column, name, text}
}

func (self *Position) Advance(current rune) {
	self.index++
	self.column++
	if current == '\n' {
		self.column = 0
		self.line++
	}
}

func (self *Position) Copy() *Position {
	return &Position{self.index, self.line, self.column, self.name, self.text}
}
