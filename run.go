package main

var symbolTable *SymbolTable = NewSymbolTable(nil)

func run(name, text string) (Value, *Error) {
	lexer := NewLexer(name, text)
	tokens, err := lexer.MakeTokens()
	if err != nil {
		return nil, err
	}

	parser := NewParser(tokens)
	syntaxTree := parser.Parse()
	if syntaxTree.error != nil {
		return nil, syntaxTree.error
	}

	context := NewContext("<repl>", nil, nil)

	for name, number := range Numbers {
		symbolTable.Set(name, number)
		number.SetContext(context)
	}

	for name, function := range builtinFunctions {
		symbolTable.Set(name, function)
		function.SetContext(context)
	}

	context.table = symbolTable
	result := syntaxTree.node.Interpret(context)
	if result.error != nil {
		return nil, result.error
	}

	context.table.Set("ans", result.value)
	return result.value, nil
}
