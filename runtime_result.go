package main

type RuntimeResult struct {
	value          Value
	error          *Error
	returnValue    Value
	shouldContinue bool
	shouldBreak    bool
}

func NewRuntimeResult() *RuntimeResult {
	return &RuntimeResult{nil, nil, nil, false, false}
}

func (self *RuntimeResult) Reset() {
	self.value = nil
	self.error = nil
	self.returnValue = nil
	self.shouldContinue = false
	self.shouldBreak = false
}

func (self *RuntimeResult) Register(res *RuntimeResult) Value {
	self.error = res.error
	self.returnValue = res.returnValue
	self.shouldContinue = res.shouldContinue
	self.shouldBreak = res.shouldBreak
	return res.value
}

func (self *RuntimeResult) Success(value Value) *RuntimeResult {
	self.Reset()
	self.value = value
	return self
}

func (self *RuntimeResult) SuccessReturn(returnValue Value) *RuntimeResult {
	self.Reset()
	self.returnValue = returnValue
	return self
}

func (self *RuntimeResult) SuccessContinue() *RuntimeResult {
	self.Reset()
	self.shouldContinue = true
	return self
}

func (self *RuntimeResult) SuccessBreak() *RuntimeResult {
	self.Reset()
	self.shouldBreak = true
	return self
}

func (self *RuntimeResult) Failure(err *Error) *RuntimeResult {
	self.Reset()
	self.error = err
	return self
}

func (self *RuntimeResult) ShouldReturn() bool {
	return self.error != nil || self.returnValue != nil || self.shouldContinue || self.shouldBreak
}
