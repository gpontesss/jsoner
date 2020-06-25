package main

import "fmt"

// Token has the information about its kind, its value and its position at input
type Token struct {
	kind   TokenKind
	index  int
	length int
	value  interface{}
}

func (t Token) String() string {
	return fmt.Sprintf("%s char %d:%d, %v", t.kind, t.index, t.index+t.length, t.value)
}

// TokenKind is the token kind
type TokenKind int

// Reference: https://tools.ietf.org/html/rfc7159#section-2
const (
	ErrorToken TokenKind = iota

	// Structural characters
	BeginArrayToken
	EndArrayToken

	BeginObjectToken
	EndObjectToken

	ColonToken
	CommaToken

	// Literal values
	FalseToken
	TrueToken

	NullToken

	NumberToken
	StringToken

	EOFToken
)

func (k TokenKind) String() string {
	switch k {
	case ErrorToken:
		return "error"
	case BeginArrayToken:
		return "["
	case EndArrayToken:
		return "]"
	case BeginObjectToken:
		return "{"
	case EndObjectToken:
		return "}"
	case ColonToken:
		return ":"
	case CommaToken:
		return ","
	case FalseToken:
		return "false"
	case TrueToken:
		return "true"
	case NullToken:
		return "null"
	case NumberToken:
		return "number"
	case StringToken:
		return "string"
	case EOFToken:
		return "EOF"
	default:
		panic("Unexpected token kind")
	}
}
