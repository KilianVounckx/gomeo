package main

type BaseFunction interface {
	Value
	Execute(arguments []Value) *RuntimeResult
}
