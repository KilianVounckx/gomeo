package main

import (
	"fmt"
)

type Parser struct {
	tokens  []*Token
	index   int
	current *Token
}

func NewParser(tokens []*Token) *Parser {
	res := &Parser{tokens, -1, nil}
	res.Advance()
	return res
}

func (self *Parser) UpdateCurrent() {
	if self.index >= 0 && self.index < len(self.tokens) {
		self.current = self.tokens[self.index]
	}
}

func (self *Parser) Advance() *Token {
	self.index += 1
	self.UpdateCurrent()
	return self.current
}

func (self *Parser) Reverse(amount int) *Token {
	self.index -= amount
	self.UpdateCurrent()
	return self.current
}

func (self *Parser) Parse() *ParseResult {
	res := self.Statements()
	if res.error == nil && self.current.tokenType != EOF {
		return res.Failure(NewInvalidSyntaxError(
			"Expected '+', '-', '*', '/', '^', '==', '!=', '<', '>', '<=', '>=', '&&', or '||'",
			self.current.start, self.current.end,
		))
	}
	return res
}

func (self *Parser) FunctionDefinition() *ParseResult {
	res := NewParseResult()

	if !self.current.Matches(KEYWORD, "function") {
		return res.Failure(NewInvalidSyntaxError(
			"Expected 'function'", self.current.start, self.current.end,
		))
	}

	res.RegisterAdvancement()
	self.Advance()

	if self.current.tokenType != LPAREN {
		return res.Failure(NewInvalidSyntaxError(
			"Expected '('", self.current.start, self.current.end,
		))
	}

	res.RegisterAdvancement()
	self.Advance()

	arguments := make([]*Token, 0)

	if self.current.tokenType == IDENTIFIER {
		arguments = append(arguments, self.current)

		res.RegisterAdvancement()
		self.Advance()

		for self.current.tokenType == COMMA {
			res.RegisterAdvancement()
			self.Advance()

			if self.current.tokenType != IDENTIFIER {
				return res.Failure(NewInvalidSyntaxError(
					"Expected identifier", self.current.start, self.current.end,
				))
			}

			arguments = append(arguments, self.current)

			res.RegisterAdvancement()
			self.Advance()
		}
	}

	if self.current.tokenType != RPAREN {
		var message string
		if len(arguments) == 0 {
			message = "Expected identifier, or ')'"
		} else {
			message = "Expected ',', or ')'"
		}
		return res.Failure(NewInvalidSyntaxError(
			message, self.current.start, self.current.end,
		))
	}

	res.RegisterAdvancement()
	self.Advance()

	if !self.current.Matches(KEYWORD, "do") {
		return res.Failure(NewInvalidSyntaxError(
			"Expected 'do'", self.current.start, self.current.end,
		))
	}

	res.RegisterAdvancement()
	self.Advance()

	var body Node
	if self.current.tokenType == NEWLINE {
		res.RegisterAdvancement()
		self.Advance()
		body = res.Register(self.Statements())
		if res.error != nil {
			return res
		}
	} else {
		body = res.Register(self.Expression())
		if res.error != nil {
			return res
		}
	}

	if !self.current.Matches(KEYWORD, "end") {
		return res.Failure(NewInvalidSyntaxError(
			"Expected 'end'", self.current.start, self.current.end,
		))
	}

	res.RegisterAdvancement()
	self.Advance()

	return res.Success(NewFunctionDefinitionNode(arguments, body))
}

func (self *Parser) WhileExpression() *ParseResult {
	res := NewParseResult()

	if !self.current.Matches(KEYWORD, "while") {
		return res.Failure(NewInvalidSyntaxError(
			"Expected 'while'", self.current.start, self.current.end,
		))
	}

	res.RegisterAdvancement()
	self.Advance()

	condition := res.Register(self.Expression())
	if res.error != nil {
		return res
	}

	if !self.current.Matches(KEYWORD, "do") {
		return res.Failure(NewInvalidSyntaxError(
			"Expected 'do'", self.current.start, self.current.end,
		))
	}

	res.RegisterAdvancement()
	self.Advance()

	var body Node
	if self.current.tokenType == NEWLINE {
		res.RegisterAdvancement()
		self.Advance()
		body = res.Register(self.Statements())
		if res.error != nil {
			return res
		}
	} else {
		body = res.Register(self.Statement())
		if res.error != nil {
			return res
		}
	}

	if !self.current.Matches(KEYWORD, "end") {
		return res.Failure(NewInvalidSyntaxError(
			"Expected 'end'", self.current.start, self.current.end,
		))
	}

	res.RegisterAdvancement()
	self.Advance()

	return res.Success(NewWhileNode(condition, body))
}

func (self *Parser) ForExpression() *ParseResult {
	res := NewParseResult()

	if !self.current.Matches(KEYWORD, "for") {
		return res.Failure(NewInvalidSyntaxError(
			"Expected 'for'", self.current.start, self.current.end,
		))
	}

	res.RegisterAdvancement()
	self.Advance()

	if self.current.tokenType != IDENTIFIER {
		return res.Failure(NewInvalidSyntaxError(
			"Expected identifier", self.current.start, self.current.end,
		))
	}

	varname := self.current

	res.RegisterAdvancement()
	self.Advance()

	if !self.current.Matches(KEYWORD, "from") {
		return res.Failure(NewInvalidSyntaxError(
			"Expected 'from'", self.current.start, self.current.end,
		))
	}

	res.RegisterAdvancement()
	self.Advance()

	from := res.Register(self.Expression())
	if res.error != nil {
		return res
	}

	if !self.current.Matches(KEYWORD, "to") {
		return res.Failure(NewInvalidSyntaxError(
			"Expected 'to'", self.current.start, self.current.end,
		))
	}

	res.RegisterAdvancement()
	self.Advance()

	to := res.Register(self.Expression())
	if res.error != nil {
		return res
	}

	var step Node
	if self.current.Matches(KEYWORD, "step") {
		res.RegisterAdvancement()
		self.Advance()

		step = res.Register(self.Expression())
		if res.error != nil {
			return res
		}
	}

	if !self.current.Matches(KEYWORD, "do") {
		var message string
		if step == nil {
			message = "Expected 'step', or 'do'"
		} else {
			message = "Expected 'do'"
		}
		return res.Failure(NewInvalidSyntaxError(
			message, self.current.start, self.current.end,
		))
	}

	res.RegisterAdvancement()
	self.Advance()

	var body Node
	if self.current.tokenType == NEWLINE {
		res.RegisterAdvancement()
		self.Advance()

		body = res.Register(self.Statements())
		if res.error != nil {
			return res
		}
	} else {
		body = res.Register(self.Statement())
		if res.error != nil {
			return res
		}
	}

	if !self.current.Matches(KEYWORD, "end") {
		return res.Failure(NewInvalidSyntaxError(
			"Expected 'end'", self.current.start, self.current.end,
		))
	}

	res.RegisterAdvancement()
	self.Advance()

	return res.Success(NewForNode(varname, from, to, step, body))
}

func (self *Parser) ListExpression() *ParseResult {
	res := NewParseResult()
	start := self.current.start.Copy()

	if self.current.tokenType != LBRACKET {
		return res.Failure(NewInvalidSyntaxError(
			"Expected '['", self.current.start, self.current.end,
		))
	}

	res.RegisterAdvancement()
	self.Advance()

	var values []Node

	if self.current.tokenType == RBRACKET {
		res.RegisterAdvancement()
		self.Advance()
	} else {
		values = append(values, res.Register(self.Expression()))
		if res.error != nil {
			return res
		}

		for self.current.tokenType == COMMA {
			res.RegisterAdvancement()
			self.Advance()

			values = append(values, res.Register(self.Expression()))
			if res.error != nil {
				return res.Failure(NewInvalidSyntaxError(
					"Expected ']', 'var', 'if', 'for', 'while', 'function', int, float, "+
						"identifier, '+', '-', '(', or '!'", self.current.start, self.current.end,
				))
			}
		}

		if self.current.tokenType != RBRACKET {
			return res.Failure(NewInvalidSyntaxError(
				"Expected ',', or ']'", self.current.start, self.current.end,
			))
		}

		res.RegisterAdvancement()
		self.Advance()
	}

	return res.Success(NewListNode(values, start, self.current.end.Copy(), false))
}

func (self *Parser) IfExpression() *ParseResult {
	res := NewParseResult()
	allCases := res.Register(self.IfElseifExpression("if"))
	if res.error != nil {
		return res
	}
	cases := allCases.(*IfNode).cases
	elseCase := allCases.(*IfNode).elseCase
	return res.Success(NewIfNode(cases, elseCase))
}

func (self *Parser) ElseIfExpression() *ParseResult {
	return self.IfElseifExpression("elseif")
}

func (self *Parser) ElseExpression() *ParseResult {
	res := NewParseResult()
	var elseCase Node

	if self.current.Matches(KEYWORD, "else") {
		res.RegisterAdvancement()
		self.Advance()

		if !self.current.Matches(KEYWORD, "do") {
			return res.Failure(NewInvalidSyntaxError(
				"Expected 'do'", self.current.start, self.current.end,
			))
		}

		res.RegisterAdvancement()
		self.Advance()

		if self.current.tokenType == NEWLINE {
			statements := res.Register(self.Statements())
			if res.error != nil {
				return res
			}
			elseCase = statements
		} else {
			expression := res.Register(self.Statement())
			if res.error != nil {
				return res
			}
			elseCase = expression
		}
	}

	if self.current.Matches(KEYWORD, "end") {
		res.RegisterAdvancement()
		self.Advance()
	} else {
		return res.Failure(NewInvalidSyntaxError(
			"Expected 'end'", self.current.start, self.current.end,
		))
	}

	return res.Success(elseCase)
}

func (self *Parser) Atom() *ParseResult {
	res := NewParseResult()

	if self.current.tokenType == NUMBER {
		token := self.current
		res.RegisterAdvancement()
		self.Advance()
		return res.Success(NewNumberNode(token))

	} else if self.current.tokenType == STRING {
		token := self.current
		res.RegisterAdvancement()
		self.Advance()
		return res.Success(NewStringNode(token))

	} else if self.current.tokenType == IDENTIFIER {
		token := self.current
		res.RegisterAdvancement()
		self.Advance()
		return res.Success(NewVariableAccessNode(token))

	} else if self.current.tokenType == LPAREN {
		res.RegisterAdvancement()
		self.Advance()

		expression := res.Register(self.Expression())
		if res.error != nil {
			return res
		}

		if self.current.tokenType != RPAREN {
			return res.Failure(NewInvalidSyntaxError(
				"Expected ')'",
				self.current.start, self.current.end,
			))
		}

		res.RegisterAdvancement()
		self.Advance()
		return res.Success(expression)

	} else if self.current.Matches(KEYWORD, "if") {
		return self.IfExpression()
	} else if self.current.Matches(KEYWORD, "for") {
		return self.ForExpression()
	} else if self.current.Matches(KEYWORD, "while") {
		return self.WhileExpression()
	} else if self.current.Matches(KEYWORD, "function") {
		return self.FunctionDefinition()
	} else if self.current.tokenType == LBRACKET {
		return self.ListExpression()
	}

	return res.Failure(NewInvalidSyntaxError(
		"Expected 'var', number, identifier, '+', '-', '[', '(', "+
			"'if', 'for', 'while', or 'function'",
		self.current.start, self.current.end,
	))
}

func (self *Parser) Call() *ParseResult {
	res := NewParseResult()

	atom := res.Register(self.Atom())
	if res.error != nil {
		return res
	}

	if self.current.tokenType == LPAREN {
		res.RegisterAdvancement()
		self.Advance()

		var arguments []Node

		if self.current.tokenType == RPAREN {
			res.RegisterAdvancement()
			self.Advance()

		} else {
			arguments = append(arguments, res.Register(self.Expression()))
			if res.error != nil {
				return res.Failure(NewInvalidSyntaxError(
					"Expected ')', 'var', 'if', 'for', 'while', 'function', int, float, "+
						"identifier, '+', '-', '(', or '!'", self.current.start, self.current.end,
				))
			}

			for self.current.tokenType == COMMA {
				res.RegisterAdvancement()
				self.Advance()

				arguments = append(arguments, res.Register(self.Expression()))
				if res.error != nil {
					return res
				}
			}

			if self.current.tokenType != RPAREN {
				return res.Failure(NewInvalidSyntaxError(
					"Expected ',', or ')'", self.current.start, self.current.end,
				))
			}

			res.RegisterAdvancement()
			self.Advance()
		}

		return res.Success(NewFunctionCallNode(atom, arguments))
	}

	return res.Success(atom)
}

func (self *Parser) Power() *ParseResult {
	return self.BinaryOperation(self.Call, self.Factor, []TokenType{POW})
}

func (self *Parser) Factor() *ParseResult {
	res := NewParseResult()

	if self.current.tokenType == PLUS || self.current.tokenType == MINUS {
		operation := self.current
		res.RegisterAdvancement()
		self.Advance()
		factor := res.Register(self.Factor())
		if res.error != nil {
			return res
		}
		return res.Success(NewUnaryOperationNode(operation, factor))
	}

	return self.Power()
}

func (self *Parser) Term() *ParseResult {
	return self.BinaryOperation(self.Factor, self.Factor, []TokenType{MUL, DIV, MOD})
}

func (self *Parser) ArithmeticExpression() *ParseResult {
	return self.BinaryOperation(self.Term, self.Term, []TokenType{PLUS, MINUS})
}

func (self *Parser) ComparisonExpression() *ParseResult {
	res := NewParseResult()

	if self.current.tokenType == NOT {
		not := self.current

		res.RegisterAdvancement()
		self.Advance()

		comparison := res.Register(self.ComparisonExpression())
		if res.error != nil {
			return res
		}

		return res.Success(NewUnaryOperationNode(not, comparison))
	}

	arithmetic := res.Register(self.BinaryOperation(
		self.ArithmeticExpression, self.ArithmeticExpression,
		[]TokenType{EE, NE, GT, LT, GE, LE},
	))
	if res.error != nil {
		return res.Failure(NewInvalidSyntaxError(
			"Expected number, identifier, '+', '-', '!', or '('",
			self.current.start, self.current.end,
		))
	}

	return res.Success(arithmetic)
}

func (self *Parser) Expression() *ParseResult {
	res := NewParseResult()

	if self.current.Matches(KEYWORD, "var") {
		res.RegisterAdvancement()
		self.Advance()

		if self.current.tokenType != IDENTIFIER {
			return res.Failure(NewInvalidSyntaxError(
				"Expected identifier", self.current.start, self.current.end,
			))
		}

		varname := self.current
		res.RegisterAdvancement()

		self.Advance()

		if self.current.tokenType != EQ {
			return res.Failure(NewInvalidSyntaxError(
				"Expected '='", self.current.start, self.current.end,
			))
		}

		res.RegisterAdvancement()
		self.Advance()

		expression := res.Register(self.Expression())
		if res.error != nil {
			return res
		}

		return res.Success(NewVariableAssignmentNode(varname, expression))
	}

	node := res.Register(self.BinaryOperation(
		self.ComparisonExpression, self.ComparisonExpression, []TokenType{AND, OR},
	))
	if res.error != nil {
		return res.Failure(NewInvalidSyntaxError(
			"Expected 'var', 'if', 'for', 'while', 'function', "+
				"number, identifier, '+', '-', '!', '[', or '('",
			self.current.start, self.current.end,
		))
	}

	return res.Success(node)
}

func (self *Parser) Statement() *ParseResult {
	res := NewParseResult()
	start := self.current.start.Copy()

	if self.current.Matches(KEYWORD, "return") {
		res.RegisterAdvancement()
		self.Advance()

		expression := res.TryRegister(self.Expression())
		if expression == nil {
			self.Reverse(res.reverseCount)
		}
		return res.Success(NewReturnNode(expression, start, self.current.start.Copy()))
	}

	if self.current.Matches(KEYWORD, "continue") {
		res.RegisterAdvancement()
		self.Advance()
		return res.Success(NewContinueNode(start, self.current.start.Copy()))
	}

	if self.current.Matches(KEYWORD, "break") {
		res.RegisterAdvancement()
		self.Advance()
		return res.Success(NewBreakNode(start, self.current.start.Copy()))
	}

	expression := res.Register(self.Expression())
	if res.error != nil {
		return res.Failure(NewInvalidSyntaxError(
			"Expected 'return', 'continue', 'break', 'var', 'if', 'for', 'while', 'function', "+
				"number, identifier, '+', '-', '(', '[', or '!'",
			self.current.start, self.current.end,
		))
	}
	return res.Success(expression)
}

func (self *Parser) Statements() *ParseResult {
	res := NewParseResult()
	start := self.current.start.Copy()

	for self.current.tokenType == NEWLINE {
		res.RegisterAdvancement()
		self.Advance()
	}

	var statements []Node

	statements = append(statements, res.Register(self.Statement()))
	if res.error != nil {
		return res
	}

	more := true
	for {
		newlines := 0

		for self.current.tokenType == NEWLINE {
			res.RegisterAdvancement()
			self.Advance()
			newlines++
		}

		if newlines == 0 {
			more = false
		}

		if !more {
			break
		}

		statement := res.TryRegister(self.Statement())
		if statement == nil {
			self.Reverse(res.reverseCount)
			more = false
			continue
		}

		statements = append(statements, statement)
	}

	return res.Success(NewListNode(statements, start, self.current.end.Copy(), true))
}

func (self *Parser) BinaryOperation(left_function, right_function func() *ParseResult,
	operations []TokenType) *ParseResult {

	res := NewParseResult()
	left := res.Register(left_function())
	if res.error != nil {
		return res
	}

	for self.current.tokenType.In(operations) {
		operation := self.current
		self.Advance()
		right := res.Register(right_function())
		if res.error != nil {
			return res
		}

		left = NewBinaryOperationNode(left, right, operation)
	}

	return res.Success(left)
}

func (self *Parser) IfElseifExpression(keyword string) *ParseResult {
	res := NewParseResult()

	var cases [][2]Node
	var elseCase Node

	if !self.current.Matches(KEYWORD, keyword) {
		return res.Failure(NewInvalidSyntaxError(
			fmt.Sprintf("Expected '%s'", keyword), self.current.start, self.current.end,
		))
	}

	res.RegisterAdvancement()
	self.Advance()

	condition := res.Register(self.Expression())
	if res.error != nil {
		return res
	}

	if !self.current.Matches(KEYWORD, "do") {
		return res.Failure(NewInvalidSyntaxError(
			"Expected 'do'", self.current.start, self.current.end,
		))
	}

	res.RegisterAdvancement()
	self.Advance()

	if self.current.tokenType == NEWLINE {
		res.RegisterAdvancement()
		self.Advance()

		statements := res.Register(self.Statements())
		if res.error != nil {
			return res
		}
		cases = append(cases, [2]Node{condition, statements})

		if self.current.Matches(KEYWORD, "end") {
			res.RegisterAdvancement()
			self.Advance()
		} else {
			allCases := res.Register(self.ElseifElseExpression())
			if res.error != nil {
				return res
			}
			newCases := allCases.(*IfNode).cases
			elseCase = allCases.(*IfNode).elseCase
			cases = append(cases, newCases...)
		}

	} else {
		expression := res.Register(self.Statement())
		if res.error != nil {
			return res
		}
		cases = append(cases, [2]Node{condition, expression})

		allCases := res.Register(self.ElseifElseExpression())
		if res.error != nil {
			return res
		}
		newCases := allCases.(*IfNode).cases
		elseCase = allCases.(*IfNode).elseCase
		cases = append(cases, newCases...)
	}

	return res.Success(NewIfNode(cases, elseCase))
}

func (self *Parser) ElseifElseExpression() *ParseResult {
	res := NewParseResult()
	var cases [][2]Node
	var elseCase Node

	if self.current.Matches(KEYWORD, "elseif") {
		allCases := res.Register(self.ElseIfExpression())
		if res.error != nil {
			return res
		}
		cases = allCases.(*IfNode).cases
		elseCase = allCases.(*IfNode).elseCase
	} else {
		elseCase = res.Register(self.ElseExpression())
		if res.error != nil {
			return res
		}
	}
	return res.Success(NewIfNode(cases, elseCase))
}
