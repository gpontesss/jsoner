package main

import "fmt"

func main() {
	jsonStr := `[{"key1": "val1"}, "val2", [1,2,3], {"arr": []}]`

	fmt.Printf("Lexing '%s'\n\n", jsonStr)

	lexer := NewLexer(jsonStr)
	tokens, err := lexer.Lex()

	for _, token := range tokens {
		fmt.Println(token)
	}
	if err != nil {
		panic(err)
	}

	fmt.Println()
	fmt.Println("Parsing tokens")

	parser := NewParser(tokens)
	val, err := parser.Parse()

	fmt.Println(val)
	if err != nil {
		panic(err)
	}
}
