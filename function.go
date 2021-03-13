package main

import (
	"fmt"
)

type Function struct {
	arguments  []string
	body       Node
	start, end *Position
	context    *Context
}

func NewFunction(arguments []string, body Node) *Function {
	return &Function{arguments, body, nil, nil, nil}
}

func (self *Function) String() string {
	return "<function>"
}

func (self *Function) Copy() Value {
	res := NewFunction(self.arguments, self.body)
	res.SetPosition(self.start, self.end)
	res.SetContext(self.context)
	return res
}

func (self *Function) SetPosition(start, end *Position) Value {
	if start == nil {
		self.start = nil
	} else {
		self.start = start.Copy()
	}
	if end == nil {
		self.end = nil
	} else {
		self.end = end.Copy()
	}
	return self
}

func (self *Function) SetContext(context *Context) Value {
	self.context = context
	return self
}

func (self *Function) Start() *Position {
	return self.start
}

func (self *Function) End() *Position {
	return self.end
}

func (self *Function) Add(value Value) (Value, *Error) {
	return nil, NewRuntimeError(
		"'+' is not supported for functions",
		self.start, self.end, self.context,
	)
}

func (self *Function) Subtract(value Value) (Value, *Error) {
	return nil, NewRuntimeError(
		"'-' is not supported for functions",
		self.start, self.end, self.context,
	)
}

func (self *Function) Multiply(value Value) (Value, *Error) {
	return nil, NewRuntimeError(
		"'*' is not supported for functions",
		self.start, self.end, self.context,
	)
}

func (self *Function) Divide(value Value) (Value, *Error) {
	return nil, NewRuntimeError(
		"'/' is not supported for functions",
		self.start, self.end, self.context,
	)
}

func (self *Function) Modulo(value Value) (Value, *Error) {
	return nil, NewRuntimeError(
		"'%' is not supported for functions",
		self.start, self.end, self.context,
	)
}

func (self *Function) Pow(value Value) (Value, *Error) {
	return nil, NewRuntimeError(
		"'*' is not supported for functions",
		self.start, self.end, self.context,
	)
}

func (self *Function) Equals(value Value) (*Number, *Error) {
	return nil, NewRuntimeError(
		"comparison not supported for functions",
		self.start, self.end, self.context,
	)
}

func (self *Function) NotEquals(value Value) (*Number, *Error) {
	return nil, NewRuntimeError(
		"comparison not supported for functions",
		self.start, self.end, self.context,
	)
}

func (self *Function) LessThan(value Value) (*Number, *Error) {
	return nil, NewRuntimeError(
		"comparison not supported for functions",
		self.start, self.end, self.context,
	)
}

func (self *Function) GreaterThan(value Value) (*Number, *Error) {
	return nil, NewRuntimeError(
		"comparison not supported for functions",
		self.start, self.end, self.context,
	)
}

func (self *Function) LessEquals(value Value) (*Number, *Error) {
	return nil, NewRuntimeError(
		"comparison not supported for functions",
		self.start, self.end, self.context,
	)
}

func (self *Function) GreaterEquals(value Value) (*Number, *Error) {
	return nil, NewRuntimeError(
		"comparison not supported for functions",
		self.start, self.end, self.context,
	)
}

func (self *Function) IsTrue() bool {
	return false
}

func (self *Function) Execute(arguments []Value) *RuntimeResult {
	res := NewRuntimeResult()

	if len(arguments) > len(self.arguments) {
		return res.Failure(NewRuntimeError(
			fmt.Sprintf(
				"%d too many arguments passed into function",
				len(arguments)-len(self.arguments),
			),
			self.start, self.end, self.context,
		))
	}
	if len(arguments) < len(self.arguments) {
		return res.Failure(NewRuntimeError(
			fmt.Sprintf(
				"%d too few arguments passed into",
				len(self.arguments)-len(arguments),
			),
			self.start, self.end, self.context,
		))
	}

	context := NewContext("function", self.context, self.start)
	context.table = NewSymbolTable(context.parent.table)

	for i, argname := range self.arguments {
		argvalue := arguments[i]
		argvalue.SetContext(context)
		context.table.Set(argname, argvalue)
	}

	value := res.Register(self.body.Interpret(context))
	if res.ShouldReturn() && res.returnValue == nil {
		return res
	}

	var returnValue Value
	if res.returnValue != nil {
		returnValue = res.returnValue
	} else {
		returnValue = value
	}
	return res.Success(returnValue)
}

func (self *Function) And(value Value) (*Number, *Error) {
	var res float64
	if self.IsTrue() && value.IsTrue() {
		res = 1
	} else {
		res = 0
	}
	return NewNumber(res), nil
}

func (self *Function) Or(value Value) (*Number, *Error) {
	var res float64
	if self.IsTrue() || value.IsTrue() {
		res = 1
	} else {
		res = 0
	}
	return NewNumber(res), nil
}

func (self *Function) Not() (*Number, *Error) {
	var res float64
	if self.IsTrue() {
		res = 0
	} else {
		res = 1
	}
	return NewNumber(res), nil
}
