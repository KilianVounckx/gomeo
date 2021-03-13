package main

import (
	"fmt"
)

type Node interface {
	fmt.Stringer
	Interpret(context *Context) *RuntimeResult
	Start() *Position
	End() *Position
}

type WrongNode struct{}

func NewWrongNode() WrongNode {
	return WrongNode{}
}

func (self WrongNode) String() string {
	return "(||WRONG NODE||)"
}

func (self WrongNode) Start() *Position {
	return nil
}

func (self WrongNode) End() *Position {
	return nil
}

//--------------------------------------------------------------------------------------------------

type NumberNode struct {
	token *Token
}

func NewNumberNode(token *Token) *NumberNode {
	return &NumberNode{token}
}

func (self *NumberNode) String() string {
	return self.token.String()
}

func (self *NumberNode) Start() *Position {
	return self.token.start
}

func (self *NumberNode) End() *Position {
	return self.token.end
}

//--------------------------------------------------------------------------------------------------

type StringNode struct {
	token *Token
}

func NewStringNode(token *Token) *StringNode {
	return &StringNode{token}
}

func (self *StringNode) String() string {
	return fmt.Sprintf("\"%s\"", self.token.String())
}

func (self *StringNode) Start() *Position {
	return self.token.start
}

func (self *StringNode) End() *Position {
	return self.token.end
}

//--------------------------------------------------------------------------------------------------

type ListNode struct {
	values     []Node
	start, end *Position
	statements bool
}

func NewListNode(values []Node, start, end *Position, statements bool) *ListNode {
	return &ListNode{values, start, end, statements}
}

func (self *ListNode) String() string {
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

func (self *ListNode) Start() *Position {
	return self.start
}

func (self *ListNode) End() *Position {
	return self.end
}

//--------------------------------------------------------------------------------------------------

type BinaryOperationNode struct {
	left, right Node
	operation   *Token
}

func NewBinaryOperationNode(left, right Node, operation *Token) *BinaryOperationNode {
	return &BinaryOperationNode{left, right, operation}
}

func (self *BinaryOperationNode) String() string {
	return fmt.Sprintf(
		"(%s, %s, %s)", self.left.String(), self.operation.String(), self.right.String(),
	)
}

func (self *BinaryOperationNode) Start() *Position {
	return self.left.Start()
}

func (self *BinaryOperationNode) End() *Position {
	return self.right.End()
}

//--------------------------------------------------------------------------------------------------

type UnaryOperationNode struct {
	operation *Token
	node      Node
}

func NewUnaryOperationNode(operation *Token, node Node) *UnaryOperationNode {
	return &UnaryOperationNode{operation, node}
}

func (self *UnaryOperationNode) String() string {
	return fmt.Sprintf("(%s,%s)", self.operation.String(), self.node.String())
}

func (self *UnaryOperationNode) Start() *Position {
	return self.operation.start
}

func (self *UnaryOperationNode) End() *Position {
	return self.node.End()
}

//--------------------------------------------------------------------------------------------------

type VariableAssignmentNode struct {
	name *Token
	node Node
}

func NewVariableAssignmentNode(name *Token, node Node) *VariableAssignmentNode {
	return &VariableAssignmentNode{name, node}
}

func (self *VariableAssignmentNode) String() string {
	return fmt.Sprintf("(var %s = %s)", self.name.String(), self.node.String())
}

func (self *VariableAssignmentNode) Start() *Position {
	return self.name.start
}

func (self *VariableAssignmentNode) End() *Position {
	return self.node.End()
}

//--------------------------------------------------------------------------------------------------

type VariableAccessNode struct {
	name *Token
}

func NewVariableAccessNode(name *Token) *VariableAccessNode {
	return &VariableAccessNode{name}
}

func (self *VariableAccessNode) String() string {
	return fmt.Sprintf("(%s)", self.name)
}

func (self *VariableAccessNode) Start() *Position {
	return self.name.start
}

func (self *VariableAccessNode) End() *Position {
	return self.name.end
}

//--------------------------------------------------------------------------------------------------

type IfNode struct {
	cases    [][2]Node
	elseCase Node
}

func NewIfNode(cases [][2]Node, elseCase Node) *IfNode {
	return &IfNode{cases, elseCase}
}

func (self *IfNode) String() string {
	res := fmt.Sprintf("(if %s do %s ", self.cases[0][0].String(), self.cases[0][1].String())
	if len(self.cases) > 0 {
		for _, pair := range self.cases[1:len(self.cases)] {
			res += fmt.Sprintf("elseif %s do %s ", pair[0], pair[1])
		}
	}
	if self.elseCase != nil {
		res += fmt.Sprintf("else %s ", self.elseCase.String())
	}
	res += "end)"
	return res
}

func (self *IfNode) Start() *Position {
	return self.cases[0][0].Start()
}

func (self *IfNode) End() *Position {
	if self.elseCase == nil {
		return self.cases[len(self.cases)-1][1].End()
	}
	return self.elseCase.End()
}

//--------------------------------------------------------------------------------------------------

type ForNode struct {
	varname        *Token
	from, to, step Node
	body           Node
}

func NewForNode(varname *Token, from, to, step, body Node) *ForNode {
	return &ForNode{varname, from, to, step, body}
}

func (self *ForNode) String() string {
	res := fmt.Sprintf(
		"(for %s from %s to %s ", self.varname.String(), self.from.String(), self.to.String(),
	)
	if self.step != nil {
		res += fmt.Sprintf("step %s ", self.step.String())
	}
	res += fmt.Sprintf("do %s end)", self.body.String())
	return res
}

func (self *ForNode) Start() *Position {
	return self.varname.start
}

func (self *ForNode) End() *Position {
	return self.body.End()
}

//--------------------------------------------------------------------------------------------------

type WhileNode struct {
	condition, body Node
}

func NewWhileNode(condition, body Node) *WhileNode {
	return &WhileNode{condition, body}
}

func (self *WhileNode) String() string {
	return fmt.Sprintf("(while %s do %s end)", self.condition.String(), self.body.String())
}

func (self *WhileNode) Start() *Position {
	return self.condition.Start()
}

func (self *WhileNode) End() *Position {
	return self.body.End()
}

//--------------------------------------------------------------------------------------------------

type FunctionDefinitionNode struct {
	arguments []*Token
	body      Node
}

func NewFunctionDefinitionNode(arguments []*Token, body Node) *FunctionDefinitionNode {
	return &FunctionDefinitionNode{arguments, body}
}

func (self *FunctionDefinitionNode) String() string {
	res := "(function ("
	if len(self.arguments) > 1 {
		for _, argument := range self.arguments[0 : len(self.arguments)-1] {
			res += fmt.Sprintf("%s, ", argument.String())
		}
	}
	if len(self.arguments) > 0 {
		res += self.arguments[len(self.arguments)-1].String()
	}
	res += "))"
	return res
}

func (self *FunctionDefinitionNode) Start() *Position {
	if len(self.arguments) == 0 {
		return self.body.Start()
	}
	return self.arguments[0].start
}

func (self *FunctionDefinitionNode) End() *Position {
	return self.body.End()
}

//--------------------------------------------------------------------------------------------------

type FunctionCallNode struct {
	call      Node
	arguments []Node
}

func NewFunctionCallNode(call Node, arguments []Node) *FunctionCallNode {
	return &FunctionCallNode{call, arguments}
}

func (self *FunctionCallNode) String() string {
	res := fmt.Sprintf("(%s (", self.call.String())
	if len(self.arguments) > 1 {
		for _, argument := range self.arguments[0 : len(self.arguments)-1] {
			res += fmt.Sprintf("%s, ", argument.String())
		}
	}
	if len(self.arguments) > 0 {
		res += self.arguments[len(self.arguments)-1].String()
	}
	res += "))"
	return res
}

func (self *FunctionCallNode) Start() *Position {
	return self.call.Start()
}

func (self *FunctionCallNode) End() *Position {
	if len(self.arguments) == 0 {
		return self.call.End()
	}
	return self.arguments[len(self.arguments)-1].End()
}

//--------------------------------------------------------------------------------------------------

type ReturnNode struct {
	nodeToReturn Node
	start, end   *Position
}

func NewReturnNode(nodeToReturn Node, start, end *Position) *ReturnNode {
	return &ReturnNode{nodeToReturn, start, end}
}

func (self *ReturnNode) String() string {
	return "return"
}

func (self *ReturnNode) Start() *Position {
	return self.start
}

func (self *ReturnNode) End() *Position {
	return self.end
}

//--------------------------------------------------------------------------------------------------

type ContinueNode struct {
	start, end *Position
}

func NewContinueNode(start, end *Position) *ContinueNode {
	return &ContinueNode{start, end}
}

func (self *ContinueNode) String() string {
	return "continue"
}

func (self *ContinueNode) Start() *Position {
	return self.start
}

func (self *ContinueNode) End() *Position {
	return self.end
}

//--------------------------------------------------------------------------------------------------

type BreakNode struct {
	start, end *Position
}

func NewBreakNode(start, end *Position) *BreakNode {
	return &BreakNode{start, end}
}

func (self *BreakNode) String() string {
	return "break"
}

func (self *BreakNode) Start() *Position {
	return self.start
}

func (self *BreakNode) End() *Position {
	return self.end
}
