package main

type ParseResult struct {
	error            *Error
	node             Node
	lastAdvanceCount int
	advanceCount     int
	reverseCount     int
}

func NewParseResult() *ParseResult {
	return &ParseResult{nil, NewWrongNode(), 0, 0, 0}
}

func (self *ParseResult) Register(res *ParseResult) Node {
	self.lastAdvanceCount = res.advanceCount
	self.advanceCount += res.advanceCount
	if res.error != nil {
		self.error = res.error
	}
	return res.node
}

func (self *ParseResult) TryRegister(res *ParseResult) Node {
	if res.error != nil {
		self.reverseCount = res.advanceCount
		return nil
	}
	return res.node
}

func (self *ParseResult) RegisterAdvancement() {
	self.advanceCount++
}

func (self *ParseResult) Success(node Node) *ParseResult {
	self.node = node
	return self
}

func (self *ParseResult) Failure(err *Error) *ParseResult {
	if self.error == nil || self.advanceCount == 0 {
		self.error = err
	}
	return self
}
