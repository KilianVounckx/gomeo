package main

import (
	"math"
	"strconv"
	"strings"
)

type Number struct {
	value      float64
	start, end *Position
	context    *Context
}

func NewNumber(value float64) *Number {
	return &Number{value, nil, nil, nil}
}

var Numbers map[string]*Number = map[string]*Number{
	"NULL":  NewNumber(0),
	"FALSE": NewNumber(0),
	"TRUE":  NewNumber(1),
	"PI":    NewNumber(math.Pi),
}

func (self *Number) Copy() Value {
	res := NewNumber(self.value)
	res.SetPosition(self.start, self.end)
	res.SetContext(self.context)
	return res
}

func (self *Number) String() string {
	return strconv.FormatFloat(self.value, 'f', -1, 64)
}

func (self *Number) SetPosition(start, end *Position) Value {
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

func (self *Number) SetContext(context *Context) Value {
	self.context = context
	return self
}

func (self *Number) Start() *Position {
	return self.start
}

func (self *Number) End() *Position {
	return self.end
}

func (self *Number) Add(value Value) (Value, *Error) {
	switch v := value.(type) {
	case *Number:
		return NewNumber(self.value + v.value).SetContext(self.context), nil
	default:
		return nil, NewRuntimeError(
			"'+' not supported between number and type", self.Start(), value.End(), self.context,
		)
	}
}

func (self *Number) Subtract(value Value) (Value, *Error) {
	switch v := value.(type) {
	case *Number:
		return NewNumber(self.value - v.value).SetContext(self.context), nil
	default:
		return nil, NewRuntimeError(
			"'-' not supported between number and type", self.Start(), value.End(), self.context,
		)
	}
}

func (self *Number) Multiply(value Value) (Value, *Error) {
	switch v := value.(type) {
	case *Number:
		return NewNumber(self.value * v.value).SetContext(self.context), nil
	case *String:
		if math.Floor(self.value) != self.value {
			return nil, NewRuntimeError(
				"'*' not supported for fractional number and string",
				self.Start(), value.End(), self.context,
			)
		}
		return NewString(strings.Repeat(v.value, int(self.value))), nil
	default:
		return nil, NewRuntimeError(
			"'*' not supported between number and type", self.Start(), value.End(), self.context,
		)
	}
}

func (self *Number) Divide(value Value) (Value, *Error) {
	switch v := value.(type) {
	case *Number:
		return NewNumber(self.value / v.value).SetContext(self.context), nil
	default:
		return nil, NewRuntimeError(
			"'/' not supported between number and type", self.Start(), value.End(), self.context,
		)
	}
}

func (self *Number) Modulo(value Value) (Value, *Error) {
	switch v := value.(type) {
	case *Number:
		if math.Floor(self.value) != self.value  || math.Floor(v.value) != v.value {
			return nil, NewRuntimeError(
				"'%' not supported for fractional number", self.Start(), value.End(), self.context,
			)
		}
		return NewNumber(float64(int(self.value) % int(v.value))).SetContext(self.context), nil
	default:
		return nil, NewRuntimeError(
			"'%' not supported between number and type", self.Start(), value.End(), self.context,
		)
	}
}

func (self *Number) Pow(value Value) (Value, *Error) {
	switch v := value.(type) {
	case *Number:
		if self.value < 0 && math.Floor(v.value) != v.value {
			return nil, NewRuntimeError(
				"Negative number cannot be raised to fractional power",
				v.start, v.end, self.context,
			)
		}
		return NewNumber(math.Pow(self.value, v.value)).SetContext(self.context), nil
	default:
		return nil, NewRuntimeError(
			"'^' not supported between number and type", self.Start(), value.End(), self.context,
		)
	}
}

func (self *Number) Equals(value Value) (*Number, *Error) {
	switch v := value.(type) {
	case *Number:
		var res float64
		if self.value == v.value {
			res = 1
		} else {
			res = 0
		}
		return NewNumber(res).SetContext(self.context).(*Number), nil
	default:
		return nil, NewRuntimeError(
			"'==' not supported between number and type", self.Start(), value.End(), self.context,
		)
	}
}

func (self *Number) NotEquals(value Value) (*Number, *Error) {
	switch v := value.(type) {
	case *Number:
		var res float64
		if self.value != v.value {
			res = 1
		} else {
			res = 0
		}
		return NewNumber(res).SetContext(self.context).(*Number), nil
	default:
		return nil, NewRuntimeError(
			"'!=' not supported between number and type", self.Start(), value.End(), self.context,
		)
	}
}

func (self *Number) LessThan(value Value) (*Number, *Error) {
	switch v := value.(type) {
	case *Number:
		var res float64
		if self.value < v.value {
			res = 1
		} else {
			res = 0
		}
		return NewNumber(res).SetContext(self.context).(*Number), nil
	default:
		return nil, NewRuntimeError(
			"'<' not supported between number and type", self.Start(), value.End(), self.context,
		)
	}
}

func (self *Number) GreaterThan(value Value) (*Number, *Error) {
	switch v := value.(type) {
	case *Number:
		var res float64
		if self.value > v.value {
			res = 1
		} else {
			res = 0
		}
		return NewNumber(res).SetContext(self.context).(*Number), nil
	default:
		return nil, NewRuntimeError(
			"'>' not supported between number and type", self.Start(), value.End(), self.context,
		)
	}
}

func (self *Number) LessEquals(value Value) (*Number, *Error) {
	switch v := value.(type) {
	case *Number:
		var res float64
		if self.value <= v.value {
			res = 1
		} else {
			res = 0
		}
		return NewNumber(res).SetContext(self.context).(*Number), nil
	default:
		return nil, NewRuntimeError(
			"'<=' not supported between number and type", self.Start(), value.End(), self.context,
		)
	}
}

func (self *Number) GreaterEquals(value Value) (*Number, *Error) {
	switch v := value.(type) {
	case *Number:
		var res float64
		if self.value >= v.value {
			res = 1
		} else {
			res = 0
		}
		return NewNumber(res).SetContext(self.context).(*Number), nil
	default:
		return nil, NewRuntimeError(
			"'>=' not supported between number and type", self.Start(), value.End(), self.context,
		)
	}
}

func (self *Number) IsTrue() bool {
	return self.value != 0.0
}

func (self *Number) And(value Value) (*Number, *Error) {
	var res float64
	if self.IsTrue() && value.IsTrue() {
		res = 1
	} else {
		res = 0
	}
	return NewNumber(res), nil
}

func (self *Number) Or(value Value) (*Number, *Error) {
	var res float64
	if self.IsTrue() || value.IsTrue() {
		res = 1
	} else {
		res = 0
	}
	return NewNumber(res), nil
}

func (self *Number) Not() (*Number, *Error) {
	var res float64
	if self.IsTrue() {
		res = 0
	} else {
		res = 1
	}
	return NewNumber(res), nil
}
