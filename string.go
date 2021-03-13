package main

import (
	"math"
	"strings"
)

type String struct {
	value      string
	start, end *Position
	context    *Context
}

func NewString(value string) *String {
	return &String{value, nil, nil, nil}
}

func (self *String) Copy() Value {
	res := NewString(self.value)
	res.SetPosition(self.start, self.end)
	res.SetContext(self.context)
	return res
}

func (self *String) String() string {
	return self.value
}

func (self *String) SetPosition(start, end *Position) Value {
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

func (self *String) SetContext(context *Context) Value {
	self.context = context
	return self
}

func (self *String) Start() *Position {
	return self.start
}

func (self *String) End() *Position {
	return self.end
}

func (self *String) Add(value Value) (Value, *Error) {
	switch v := value.(type) {
	case *String:
		return NewString(self.value + v.value).SetContext(self.context), nil
	default:
		return nil, NewRuntimeError(
			"'+' not supported for string and type", self.Start(), value.End(), self.context,
		)
	}
}

func (self *String) Subtract(value Value) (Value, *Error) {
	return nil, NewRuntimeError(
		"'-' not supported for string", self.Start(), value.End(), self.context,
	)
}

func (self *String) Multiply(value Value) (Value, *Error) {
	switch v := value.(type) {
	case *Number:
		if math.Floor(v.value) != v.value {
			return nil, NewRuntimeError(
				"'*' not supported for string and fractional number",
				self.Start(), value.End(), self.context,
			)
		}
		return NewString(strings.Repeat(self.value, int(v.value))), nil
	default:
		return nil, NewRuntimeError(
			"'*' not supported for string and type", self.Start(), value.End(), self.context,
		)
	}
}

func (self *String) Divide(value Value) (Value, *Error) {
	return nil, NewRuntimeError(
		"'/' not supported for string", self.Start(), value.End(), self.context,
	)
}

func (self *String) Modulo(value Value) (Value, *Error) {
	return nil, NewRuntimeError(
		"'%' not supported for string", self.Start(), value.End(), self.context,
	)
}

func (self *String) Pow(value Value) (Value, *Error) {
	return nil, NewRuntimeError(
		"'^' not supported for string", self.Start(), value.End(), self.context,
	)
}

func (self *String) Equals(value Value) (*Number, *Error) {
	switch v := value.(type) {
	case *String:
		var res float64
		if self.value == v.value {
			res = 1
		} else {
			res = 0
		}
		return NewNumber(res).SetContext(self.context).(*Number), nil
	default:
		return nil, NewRuntimeError(
			"'==' not supported between string and type", self.Start(), value.End(), self.context,
		)
	}
}

func (self *String) NotEquals(value Value) (*Number, *Error) {
	switch v := value.(type) {
	case *String:
		var res float64
		if self.value != v.value {
			res = 1
		} else {
			res = 0
		}
		return NewNumber(res).SetContext(self.context).(*Number), nil
	default:
		return nil, NewRuntimeError(
			"'!=' not supported between string and type", self.Start(), value.End(), self.context,
		)
	}
}

func (self *String) LessThan(value Value) (*Number, *Error) {
	switch v := value.(type) {
	case *String:
		var res float64
		if self.value < v.value {
			res = 1
		} else {
			res = 0
		}
		return NewNumber(res).SetContext(self.context).(*Number), nil
	default:
		return nil, NewRuntimeError(
			"'<' not supported between string and type", self.Start(), value.End(), self.context,
		)
	}
}

func (self *String) GreaterThan(value Value) (*Number, *Error) {
	switch v := value.(type) {
	case *String:
		var res float64
		if self.value > v.value {
			res = 1
		} else {
			res = 0
		}
		return NewNumber(res).SetContext(self.context).(*Number), nil
	default:
		return nil, NewRuntimeError(
			"'>' not supported between string and type", self.Start(), value.End(), self.context,
		)
	}
}

func (self *String) LessEquals(value Value) (*Number, *Error) {
	switch v := value.(type) {
	case *String:
		var res float64
		if self.value <= v.value {
			res = 1
		} else {
			res = 0
		}
		return NewNumber(res).SetContext(self.context).(*Number), nil
	default:
		return nil, NewRuntimeError(
			"'<=' not supported between string and type", self.Start(), value.End(), self.context,
		)
	}
}

func (self *String) GreaterEquals(value Value) (*Number, *Error) {
	switch v := value.(type) {
	case *String:
		var res float64
		if self.value <= v.value {
			res = 1
		} else {
			res = 0
		}
		return NewNumber(res).SetContext(self.context).(*Number), nil
	default:
		return nil, NewRuntimeError(
			"'<=' not supported between string and type", self.Start(), value.End(), self.context,
		)
	}
}

func (self *String) IsTrue() bool {
	return self.value != ""
}

func (self *String) And(value Value) (*Number, *Error) {
	var res float64
	if self.IsTrue() && value.IsTrue() {
		res = 1
	} else {
		res = 0
	}
	return NewNumber(res), nil
}

func (self *String) Or(value Value) (*Number, *Error) {
	var res float64
	if self.IsTrue() || value.IsTrue() {
		res = 1
	} else {
		res = 0
	}
	return NewNumber(res), nil
}

func (self *String) Not() (*Number, *Error) {
	var res float64
	if self.IsTrue() {
		res = 0
	} else {
		res = 1
	}
	return NewNumber(res), nil
}
