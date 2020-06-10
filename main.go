package main

import "fmt"

func main() {
	jsonStr := `	[{}, {}		]`
	lexer := NewLexer(jsonStr)
	tokens, err := lexer.Lex()
	if err != nil {
		panic(err)
	}
	for _, token := range tokens {
		fmt.Println(token)
	}
}
