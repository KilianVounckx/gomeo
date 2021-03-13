package main

import (
	"fmt"
	"math"
)

type List struct {
	values     []Value
	start, end *Position
	context    *Context
}

func NewList(values []Value) *List {
	var items []Value
	for _, value := range values {
		if value != nil {
			items = append(items, value)
		}
	}
	return &List{items, nil, nil, nil}
}

func (self *List) String() string {
	res := "["
	if len(self.values) > 1 {
		for _, value := range self.values[0 : len(self.values)-1] {
			res += fmt.Sprintf("%s, ", value.String())
		}
	}
	if len(self.values) > 0 {
		res += self.values[len(self.values)-1].String()
	}
	res += "]"
	return res
}

func (self *List) Copy() Value {
	res := NewList(self.values)
	res.SetPosition(self.start, self.end)
	res.SetContext(self.context)
	return res
}

func (self *List) SetPosition(start, end *Position) Value {
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

func (self *List) SetContext(context *Context) Value {
	self.context = context
	return self
}

func (self *List) Start() *Position {
	return self.start
}

func (self *List) End() *Position {
	return self.end
}

func (self *List) Add(value Value) (Value, *Error) {
	res := self.Copy().(*List)
	res.values = append(res.values, value)
	return res, nil
}

func (self *List) Subtract(value Value) (Value, *Error) {
	switch v := value.(type) {
	case *Number:
		if math.Floor(v.value) != v.value {
			return nil, NewRuntimeError(
				"'-' not supported for list and fractional number",
				self.Start(), value.End(), self.context,
			)
		}
		vint := int(v.value)
		if vint < 0 || vint >= len(self.values) {
			return nil, NewRuntimeError(
				fmt.Sprintf(
					"Index out of range (length %d, index %d)",
					len(self.values), vint,
				), self.Start(), value.End(), self.context,
			)
		}
		res := self.Copy().(*List)
		res.values = remove(res.values, vint)
		return res, nil
	default:
		return nil, NewRuntimeError(
			"'-' not supported for list and type", self.Start(), value.End(), self.context,
		)
	}
}

func (self *List) Multiply(value Value) (Value, *Error) {
	switch v := value.(type) {
	case *List:
		res := self.Copy().(*List)
		res.values = append(res.values, v.values...)
		return res, nil
	default:
		return nil, NewRuntimeError(
			"'*' not supported for list and type", self.Start(), value.End(), self.context,
		)
	}
}

func (self *List) Divide(value Value) (Value, *Error) {
	switch v := value.(type) {
	case *Number:
		if math.Floor(v.value) != v.value {
			return nil, NewRuntimeError(
				"'/' not supported for list and fractional number",
				self.Start(), value.End(), self.context,
			)
		}
		vint := int(v.value)
		if vint < 0 || vint >= len(self.values) {
			return nil, NewRuntimeError(
				fmt.Sprintf(
					"Index out of range (length %d, index %d)",
					len(self.values), vint,
				), self.Start(), value.End(), self.context,
			)
		}
		return self.values[vint], nil
	default:
		return nil, NewRuntimeError(
			"'-' not supported for list and type", self.Start(), value.End(), self.context,
		)
	}
}

func (self *List) Modulo(value Value) (Value, *Error) {
	return nil, NewRuntimeError(
		"'%' not supported for list", self.Start(), value.End(), self.context,
	)
}

func (self *List) Pow(value Value) (Value, *Error) {
	switch v := value.(type) {
	case *Number:
		if math.Floor(v.value) != v.value {
			return nil, NewRuntimeError(
				"'^' not supported for list and fractional number",
				self.Start(), value.End(), self.context,
			)
		}
		vint := int(v.value)
		if vint < 0 {
			return nil, NewRuntimeError(
				"'^' not supported for list and negative number",
				self.Start(), value.End(), self.context,
			)
		}
		res := self.Copy().(*List)
		res.values = repeat(res.values, vint)
		return res, nil
	default:
		return nil, NewRuntimeError(
			"'^' not supported for list and type",
			self.Start(), value.End(), self.context,
		)
	}
}

func (self *List) Equals(value Value) (*Number, *Error) {
	switch v := value.(type) {
	case *List:
		if len(self.values) != len(v.values) {
			return NewNumber(0).SetContext(self.context).(*Number), nil
		}

		for i, selfValue := range self.values {
			otherValue := v.values[i]
			check, err := selfValue.Equals(otherValue)
			if err != nil {
				return nil, err
			}
			if check.value == 0 {
				return NewNumber(0).SetContext(self.context).(*Number), nil
			}
		}
		return NewNumber(1).SetContext(self.context).(*Number), nil
	default:
		return nil, NewRuntimeError(
			"'==' not supported between list and type", self.Start(), value.End(), self.context,
		)
	}
}

func (self *List) NotEquals(value Value) (*Number, *Error) {
	switch v := value.(type) {
	case *List:
		var res float64
		check, err := self.Equals(v)
		if err != nil {
			return nil, err
		}
		if check.value == 0 {
			res = 1
		} else {
			res = 0
		}
		return NewNumber(res).SetContext(self.context).(*Number), nil
	default:
		return nil, NewRuntimeError(
			"'!=' not supported between list and type", self.Start(), value.End(), self.context,
		)
	}
}

func (self *List) LessThan(value Value) (*Number, *Error) {
	return nil, NewRuntimeError(
		"'<' not supported for list", self.Start(), value.End(), self.context,
	)
}

func (self *List) GreaterThan(value Value) (*Number, *Error) {
	return nil, NewRuntimeError(
		"'>' not supported for list", self.Start(), value.End(), self.context,
	)
}

func (self *List) LessEquals(value Value) (*Number, *Error) {
	return nil, NewRuntimeError(
		"'<=' not supported for list", self.Start(), value.End(), self.context,
	)
}

func (self *List) GreaterEquals(value Value) (*Number, *Error) {
	return nil, NewRuntimeError(
		"'>=' not supported for list", self.Start(), value.End(), self.context,
	)
}

func (self *List) And(value Value) (*Number, *Error) {
	var res float64
	if self.IsTrue() && value.IsTrue() {
		res = 1
	} else {
		res = 0
	}
	return NewNumber(res), nil
}

func (self *List) Or(value Value) (*Number, *Error) {
	var res float64
	if self.IsTrue() || value.IsTrue() {
		res = 1
	} else {
		res = 0
	}
	return NewNumber(res), nil
}

func (self *List) Not() (*Number, *Error) {
	var res float64
	if self.IsTrue() {
		res = 0
	} else {
		res = 1
	}
	return NewNumber(res), nil
}

func (self *List) IsTrue() bool {
	return len(self.values) > 0
}
