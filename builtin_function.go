package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"os/exec"
	"strconv"
)

type BuiltinFunction struct {
	arguments  []string
	body       func(context *Context) *RuntimeResult
	context    *Context
	start, end *Position
}

func NewBuiltinFunction(arguments []string,
	body func(context *Context) *RuntimeResult) *BuiltinFunction {
	return &BuiltinFunction{arguments, body, nil, nil, nil}
}

func (self *BuiltinFunction) String() string {
	return "<built-in function>"
}

func (self *BuiltinFunction) Copy() Value {
	res := NewBuiltinFunction(self.arguments, self.body)
	res.SetPosition(self.start, self.end)
	res.SetContext(self.context)
	return res
}

func (self *BuiltinFunction) SetPosition(start, end *Position) Value {
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

func (self *BuiltinFunction) SetContext(context *Context) Value {
	self.context = context
	return self
}

func (self *BuiltinFunction) Start() *Position {
	return self.start
}

func (self *BuiltinFunction) End() *Position {
	return self.end
}

func (self *BuiltinFunction) Add(value Value) (Value, *Error) {
	return nil, NewRuntimeError(
		"'+' is not supported for functions",
		self.start, self.end, self.context,
	)
}

func (self *BuiltinFunction) Subtract(value Value) (Value, *Error) {
	return nil, NewRuntimeError(
		"'-' is not supported for functions",
		self.start, self.end, self.context,
	)
}

func (self *BuiltinFunction) Multiply(value Value) (Value, *Error) {
	return nil, NewRuntimeError(
		"'*' is not supported for functions",
		self.start, self.end, self.context,
	)
}

func (self *BuiltinFunction) Divide(value Value) (Value, *Error) {
	return nil, NewRuntimeError(
		"'/' is not supported for functions",
		self.start, self.end, self.context,
	)
}

func (self *BuiltinFunction) Modulo(value Value) (Value, *Error) {
	return nil, NewRuntimeError(
		"'%' is not supported for functions",
		self.start, self.end, self.context,
	)
}

func (self *BuiltinFunction) Pow(value Value) (Value, *Error) {
	return nil, NewRuntimeError(
		"'*' is not supported for functions",
		self.start, self.end, self.context,
	)
}

func (self *BuiltinFunction) Equals(value Value) (*Number, *Error) {
	return nil, NewRuntimeError(
		"comparison not supported for functions",
		self.start, self.end, self.context,
	)
}

func (self *BuiltinFunction) NotEquals(value Value) (*Number, *Error) {
	return nil, NewRuntimeError(
		"comparison not supported for functions",
		self.start, self.end, self.context,
	)
}

func (self *BuiltinFunction) LessThan(value Value) (*Number, *Error) {
	return nil, NewRuntimeError(
		"comparison not supported for functions",
		self.start, self.end, self.context,
	)
}

func (self *BuiltinFunction) GreaterThan(value Value) (*Number, *Error) {
	return nil, NewRuntimeError(
		"comparison not supported for functions",
		self.start, self.end, self.context,
	)
}

func (self *BuiltinFunction) LessEquals(value Value) (*Number, *Error) {
	return nil, NewRuntimeError(
		"comparison not supported for functions",
		self.start, self.end, self.context,
	)
}

func (self *BuiltinFunction) GreaterEquals(value Value) (*Number, *Error) {
	return nil, NewRuntimeError(
		"comparison not supported for functions",
		self.start, self.end, self.context,
	)
}

func (self *BuiltinFunction) IsTrue() bool {
	return false
}

func (self *BuiltinFunction) And(value Value) (*Number, *Error) {
	var res float64
	if self.IsTrue() && value.IsTrue() {
		res = 1
	} else {
		res = 0
	}
	return NewNumber(res), nil
}

func (self *BuiltinFunction) Or(value Value) (*Number, *Error) {
	var res float64
	if self.IsTrue() || value.IsTrue() {
		res = 1
	} else {
		res = 0
	}
	return NewNumber(res), nil
}

func (self *BuiltinFunction) Not() (*Number, *Error) {
	var res float64
	if self.IsTrue() {
		res = 0
	} else {
		res = 1
	}
	return NewNumber(res), nil
}

func (self *BuiltinFunction) Execute(arguments []Value) *RuntimeResult {
	res := NewRuntimeResult()

	context := NewContext("function", self.context, self.start)
	context.table = NewSymbolTable(context.parent.table)

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

	for i, argname := range self.arguments {
		argvalue := arguments[i]
		argvalue.SetContext(context)
		context.table.Set(argname, argvalue)
	}

	value := res.Register(self.body(context))
	if res.error != nil {
		return res
	}
	return res.Success(value)
}

//--------------------------------------------------------------------------------------------------

var builtinFunctions map[string]*BuiltinFunction = map[string]*BuiltinFunction{
	"print": NewBuiltinFunction([]string{"value"}, func(context *Context) *RuntimeResult {
		fmt.Print(context.table.Get("value"))
		return NewRuntimeResult().Success(nil)
	}),

	"println": NewBuiltinFunction([]string{"value"}, func(context *Context) *RuntimeResult {
		fmt.Println(context.table.Get("value"))
		return NewRuntimeResult().Success(nil)
	}),

	"printReturn": NewBuiltinFunction([]string{"value"}, func(context *Context) *RuntimeResult {
		value := context.table.Get("value")
		fmt.Println(value)
		return NewRuntimeResult().Success(value)
	}),

	"input": NewBuiltinFunction([]string{"prompt"}, func(context *Context) *RuntimeResult {
		prompt := context.table.Get("prompt")
		fmt.Print(prompt)

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		text := scanner.Text()
		res := NewRuntimeResult().Success(NewString(text))
		return res
	}),

	"inputNumber": NewBuiltinFunction([]string{"prompt"}, func(context *Context) *RuntimeResult {
		prompt := context.table.Get("prompt")
		fmt.Print(prompt)

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		text := scanner.Text()
		res, err := strconv.ParseFloat(text, 64)
		if err != nil {
			return NewRuntimeResult().Failure(NewRuntimeError(
				fmt.Sprintf("Invalid number '%s'", text), context.entry, context.entry, context,
			))
		}
		return NewRuntimeResult().Success(NewNumber(res))
	}),

	"isNumber": NewBuiltinFunction([]string{"number"}, func(context *Context) *RuntimeResult {
		number := context.table.Get("number")
		switch number.(type) {
		case *Number:
			return NewRuntimeResult().Success(NewNumber(1))
		default:
			return NewRuntimeResult().Success(NewNumber(0))
		}
	}),

	"isString": NewBuiltinFunction([]string{"string"}, func(context *Context) *RuntimeResult {
		number := context.table.Get("string")
		switch number.(type) {
		case *String:
			return NewRuntimeResult().Success(NewNumber(1))
		default:
			return NewRuntimeResult().Success(NewNumber(0))
		}
	}),

	"isList": NewBuiltinFunction([]string{"list"}, func(context *Context) *RuntimeResult {
		number := context.table.Get("list")
		switch number.(type) {
		case *List:
			return NewRuntimeResult().Success(NewNumber(1))
		default:
			return NewRuntimeResult().Success(NewNumber(0))
		}
	}),

	"isFunction": NewBuiltinFunction([]string{"function"}, func(context *Context) *RuntimeResult {
		number := context.table.Get("function")
		switch number.(type) {
		case BaseFunction:
			return NewRuntimeResult().Success(NewNumber(1))
		default:
			return NewRuntimeResult().Success(NewNumber(0))
		}
	}),

	"clear": NewBuiltinFunction([]string{}, func(context *Context) *RuntimeResult {
		command := exec.Command("clear")
		command.Stdout = os.Stdout
		_ = command.Run()
		return NewRuntimeResult().Success(nil)
	}),

	"exit": NewBuiltinFunction([]string{}, func(context *Context) *RuntimeResult {
		os.Exit(0)
		return nil
	}),

	"len": NewBuiltinFunction([]string{"value"}, func(context *Context) *RuntimeResult {
		value := context.table.Get("value")
		switch v := value.(type) {
		case *Number:
			return NewRuntimeResult().Success(NewNumber(1))
		case BaseFunction:
			return NewRuntimeResult().Success(NewNumber(1))
		case *String:
			return NewRuntimeResult().Success(NewNumber(float64(len(v.value))))
		case *List:
			return NewRuntimeResult().Success(NewNumber(float64(len(v.values))))
		default:
			return NewRuntimeResult().Failure(NewRuntimeError(
				"'len' not supported for type", context.entry, context.entry, context,
			))
		}
	}),

	"change": NewBuiltinFunction([]string{"list", "index", "value"},
			func(context *Context) *RuntimeResult {
		v := context.table.Get("list")
		switch v.(type) {
		case *List:
		default:
			return NewRuntimeResult().Failure(NewRuntimeError(
				"'change!' first parameter must be a list", context.entry, context.entry, context,
			))
		}
		list := v.(*List)
		v = context.table.Get("index")
		switch v.(type) {
		case *Number:
		default:
			return NewRuntimeResult().Failure(NewRuntimeError(
				"'change!' second parameter must be a number",
				context.entry, context.entry, context,
			))
		}
		number := v.(*Number)
		if math.Floor(number.value) != number.value {
			return NewRuntimeResult().Failure(NewRuntimeError(
				"'change!' second parameter must not be fractional",
				context.entry, context.entry, context,
			))
		}
		index := int(number.value)

		if index < 0 || index >= len(list.values) {
			return NewRuntimeResult().Failure(NewRuntimeError(
				fmt.Sprintf("Index out of range (length %d, index %d)", len(list.values), index),
				context.entry, context.entry, context,
			))
		}

		value := context.table.Get("value")
		list.values[index] = value
		return NewRuntimeResult().Success(list)
	}),
}

var RUN *BuiltinFunction = NewBuiltinFunction([]string{"file"},
	func(context *Context) *RuntimeResult {
		res := NewRuntimeResult()
		file := context.table.Get("file")
		switch file.(type) {
		case *String:
		default:
			return res.Failure(NewRuntimeError(
				"Filename must be a string", context.entry, context.entry, context,
			))
		}

		fileString := file.(*String)
		content, err := ioutil.ReadFile(fileString.value)
		if err != nil {
			return res.Failure(NewRuntimeError(
				fmt.Sprintf("File '%s' not found", fileString.value),
				context.entry, context.entry, context,
			))
		}

		value, err2 := run(fileString.value, string(content))
		if err2 != nil {
			return res.Failure(err2)
		}

		return res.Success(value)
	})
