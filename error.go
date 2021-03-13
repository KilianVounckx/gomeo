package main

import (
	"fmt"
)

type Error struct {
	name    string
	details string
	start   *Position
	end     *Position
	context *Context
}

func NewError(name, details string, start, end *Position, context *Context) *Error {
	return &Error{name, details, start, end, context}
}

func (self *Error) AsString() string {
	if self.context == nil {
		result := fmt.Sprintf("%s: %s\n", self.name, self.details)
		result += fmt.Sprintf("File %s, line %d", self.start.name, self.start.line)
		result += fmt.Sprintf("\n\n%s", stringWithArrows(self.start.text, self.start, self.end))
		return result
	} else {
		result := self.GenerateTraceback()
		result += fmt.Sprintf("%s: %s\n", self.name, self.details)
		result += fmt.Sprintf("\n\n%s", stringWithArrows(self.start.text, self.start, self.end))
		return result
	}
}

func NewIllegalCharacterError(details string, start, end *Position) *Error {
	return NewError("IllegalCharacterError", details, start, end, nil)
}

func NewExpectedCharacterError(details string, start, end *Position) *Error {
	return NewError("ExpectedCharacterError", details, start, end, nil)
}

func NewInvalidSyntaxError(details string, start, end *Position) *Error {
	return NewError("InvalidSyntaxError", details, start, end, nil)
}

func NewRuntimeError(details string, start, end *Position, context *Context) *Error {
	return NewError("RuntimeError", details, start, end, context)
}

func (self *Error) GenerateTraceback() string {
	res := ""
	position := self.start
	context := self.context

	for context != nil {
		res = fmt.Sprintf(
			"  File %s, line %d, in %s\n", position.name, position.line+1, self.name,
		) + res
		position = context.entry
		context = context.parent
	}

	return "Traceback (most recent call last):\n" + res
}
