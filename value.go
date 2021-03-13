package main

import (
	"fmt"
)

type Value interface {
	fmt.Stringer
	Copy() Value
	SetPosition(start, end *Position) Value
	SetContext(context *Context) Value

	Start() *Position
	End() *Position

	Add(value Value) (Value, *Error)
	Subtract(value Value) (Value, *Error)
	Multiply(value Value) (Value, *Error)
	Divide(value Value) (Value, *Error)
	Modulo(value Value) (Value, *Error)
	Pow(value Value) (Value, *Error)

	Equals(value Value) (*Number, *Error)
	NotEquals(value Value) (*Number, *Error)
	LessThan(value Value) (*Number, *Error)
	GreaterThan(value Value) (*Number, *Error)
	LessEquals(value Value) (*Number, *Error)
	GreaterEquals(value Value) (*Number, *Error)

	And(value Value) (*Number, *Error)
	Or(value Value) (*Number, *Error)
	Not() (*Number, *Error)

	IsTrue() bool
}
