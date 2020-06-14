package main

import "fmt"

func main() {
	jsonStr := `[ {1.23} , {1e-1}, "this is my \"problematic\" string \u1F00", true, false, null ]`
	fmt.Printf("Lexing '%s'\n\n", jsonStr)
	lexer := NewLexer(jsonStr)
	tokens, err := lexer.Lex()
	for _, token := range tokens {
		fmt.Println(token)
	}
	if err != nil {
		panic(err)
	}
}
