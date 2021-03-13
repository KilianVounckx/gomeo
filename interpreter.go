package main

import (
	"fmt"
)

func (self WrongNode) Interpret(context *Context) *RuntimeResult {
	panic("No visit method for WrongNode")
}

func (self *NumberNode) Interpret(context *Context) *RuntimeResult {
	return NewRuntimeResult().Success(
		NewNumber(
			self.token.value.(float64),
		).SetContext(context).SetPosition(self.Start(), self.End()),
	)
}

func (self *StringNode) Interpret(context *Context) *RuntimeResult {
	return NewRuntimeResult().Success(
		NewString(
			self.token.value.(string),
		).SetContext(context).SetPosition(self.Start(), self.End()),
	)
}

func (self *ListNode) Interpret(context *Context) *RuntimeResult {
	res := NewRuntimeResult()

	var values []Value
	for _, value := range self.values {
		values = append(values, res.Register(value.Interpret(context)))
		if res.ShouldReturn() {
			return res
		}
	}

	if self.statements {
		if len(values) > 0 {
			return res.Success(values[len(values)-1])
		}
		return res.Success(nil)
	}

	return res.Success(
		NewList(values).SetContext(context).SetPosition(self.Start(), self.End()),
	)
}

func (self *BinaryOperationNode) Interpret(context *Context) *RuntimeResult {
	res := NewRuntimeResult()

	left := res.Register(self.left.Interpret(context))
	if res.ShouldReturn() {
		return res
	}

	right := res.Register(self.right.Interpret(context))
	if res.ShouldReturn() {
		return res
	}

	var value Value
	var err *Error

	switch self.operation.tokenType {
	case PLUS:
		value, err = left.Add(right)
	case MINUS:
		value, err = left.Subtract(right)
	case MUL:
		value, err = left.Multiply(right)
	case DIV:
		value, err = left.Divide(right)
	case MOD:
		value, err = left.Modulo(right)
	case POW:
		value, err = left.Pow(right)
	case AND:
		value, err = left.And(right)
	case OR:
		value, err = left.Or(right)
	case EE:
		value, err = left.Equals(right)
	case NE:
		value, err = left.NotEquals(right)
	case LT:
		value, err = left.LessThan(right)
	case GT:
		value, err = left.GreaterThan(right)
	case LE:
		value, err = left.LessEquals(right)
	case GE:
		value, err = left.GreaterEquals(right)
	}

	if err != nil {
		return res.Failure(err)
	}

	return res.Success(value.SetPosition(self.Start(), self.End()))
}

func (self *UnaryOperationNode) Interpret(context *Context) *RuntimeResult {
	res := NewRuntimeResult()

	left := res.Register(self.node.Interpret(context))
	if res.ShouldReturn() {
		return res
	}

	var value Value
	var err *Error

	switch self.operation.tokenType {
	case PLUS:
		value = left
	case MINUS:
		value, err = left.Multiply(NewNumber(-1))
	case NOT:
		value, err = left.Not()
	}

	if err != nil {
		return res.Failure(err)
	}

	return res.Success(value.SetPosition(self.Start(), self.End()))
}

func (self *VariableAssignmentNode) Interpret(context *Context) *RuntimeResult {
	res := NewRuntimeResult()

	varname := self.name.value.(string)
	value := res.Register(self.node.Interpret(context))
	if res.ShouldReturn() {
		return res
	}

	context.table.Set(varname, value)

	return res.Success(value)
}

func (self *VariableAccessNode) Interpret(context *Context) *RuntimeResult {
	res := NewRuntimeResult()

	varname := self.name.value.(string)
	value := context.table.Get(varname)
	if value == nil {
		return res.Failure(NewRuntimeError(
			fmt.Sprintf("'%s' is not defined", varname),
			self.Start(), self.End(), context,
		))
	}

	value = value.Copy().SetPosition(self.Start(), self.End()).SetContext(context)
	return res.Success(value)
}

func (self *IfNode) Interpret(context *Context) *RuntimeResult {
	res := NewRuntimeResult()

	for _, pair := range self.cases {
		condition := pair[0]
		expression := pair[1]

		conditionValue := res.Register(condition.Interpret(context))
		if res.ShouldReturn() {
			return res
		}

		check := conditionValue.IsTrue()

		if check {
			expressionValue := res.Register(expression.Interpret(context))
			if res.ShouldReturn() {
				return res
			}
			return res.Success(expressionValue)
		}
	}

	if self.elseCase != nil {
		elseValue := res.Register(self.elseCase.Interpret(context))
		if res.ShouldReturn() {
			return res
		}
		return res.Success(elseValue)
	}
	return res.Success(nil)
}

func (self *ForNode) Interpret(context *Context) *RuntimeResult {
	res := NewRuntimeResult()

	var values []Value

	from := res.Register(self.from.Interpret(context))
	if res.ShouldReturn() {
		return res
	}
	to := res.Register(self.to.Interpret(context))
	if res.ShouldReturn() {
		return res
	}

	var step Value
	if self.step != nil {
		step = res.Register(self.step.Interpret(context))
		if res.ShouldReturn() {
			return res
		}
	} else {
		step = NewNumber(1)
	}

	var condition func() (bool, *Error)
	ascending, err := step.GreaterEquals(NewNumber(0))
	if err != nil {
		return res.Failure(err)
	}
	if ascending.IsTrue() {
		condition = func() (bool, *Error) {
			check, err := from.LessThan(to)
			if err != nil {
				return false, err
			}
			return check.IsTrue(), nil
		}
	} else {
		condition = func() (bool, *Error) {
			check, err := from.GreaterThan(to)
			if err != nil {
				return false, err
			}
			return check.IsTrue(), nil
		}
	}

	for {
		check, err := condition()
		if err != nil {
			return res.Failure(err)
		}
		if !check {
			break
		}

		context.table.Set(self.varname.value.(string), from.Copy())
		from, err = from.Add(step)
		if err != nil {
			return res.Failure(err)
		}

		value := res.Register(self.body.Interpret(context))
		if res.ShouldReturn() && !res.shouldContinue && !res.shouldBreak {
			return res
		}

		if res.shouldContinue {
			continue
		}
		if res.shouldBreak {
			break
		}

		if value != nil {
			values = append(values, value)
		}
	}

	return res.Success(
		NewList(values).SetContext(context).SetPosition(self.Start(), self.End()),
	)
}

func (self *WhileNode) Interpret(context *Context) *RuntimeResult {
	res := NewRuntimeResult()

	var values []Value

	for {
		condition := res.Register(self.condition.Interpret(context))
		if res.ShouldReturn() {
			return res
		}

		if !condition.IsTrue() {
			break
		}

		value := res.Register(self.body.Interpret(context))
		if res.ShouldReturn() && !res.shouldContinue && !res.shouldBreak {
			return res
		}

		if res.shouldContinue {
			continue
		}
		if res.shouldBreak {
			break
		}

		if value != nil {
			values = append(values, value)
		}
	}

	return res.Success(
		NewList(values).SetContext(context).SetPosition(self.Start(), self.End()),
	)
}

func (self *FunctionDefinitionNode) Interpret(context *Context) *RuntimeResult {
	res := NewRuntimeResult()

	body := self.body
	argnames := make([]string, len(self.arguments))
	for i, argument := range self.arguments {
		argnames[i] = argument.value.(string)
	}
	function := NewFunction(argnames, body).
		SetContext(context).
		SetPosition(self.Start(), self.End())

	return res.Success(function)
}

func (self *FunctionCallNode) Interpret(context *Context) *RuntimeResult {
	res := NewRuntimeResult()

	call := res.Register(self.call.Interpret(context))
	if res.ShouldReturn() {
		return res
	}
	call = call.Copy().SetPosition(self.Start(), self.End())

	var arguments []Value

	for _, argument := range self.arguments {
		arguments = append(arguments, res.Register(argument.Interpret(context)))
		if res.ShouldReturn() {
			return res
		}
	}

	value := res.Register(call.(BaseFunction).Execute(arguments))
	if res.ShouldReturn() {
		return res
	}

	if value == nil {
		return res.Success(nil)
	}

	value = value.Copy().SetPosition(self.Start(), self.End()).SetContext(context)
	return res.Success(value)
}

func (self *ReturnNode) Interpret(context *Context) *RuntimeResult {
	res := NewRuntimeResult()

	var value Value
	if self.nodeToReturn != nil {
		value = res.Register(self.nodeToReturn.Interpret(context))
		if res.ShouldReturn() {
			return res
		}
	}
	return res.SuccessReturn(value)
}

func (self *ContinueNode) Interpret(context *Context) *RuntimeResult {
	return NewRuntimeResult().SuccessContinue()
}

func (self *BreakNode) Interpret(context *Context) *RuntimeResult {
	return NewRuntimeResult().SuccessBreak()
}
