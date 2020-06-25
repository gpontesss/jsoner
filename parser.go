package main

import "fmt"

// Parser parses tokens from lexer
type Parser struct {
	tokens  []Token
	current int
}

// NewParser creates a mew parser from tokens
func NewParser(tokens []Token) Parser {
	return Parser{tokens: tokens}
}

// Parse parses tokens into an AstValue
func (p *Parser) Parse() (val AstValue, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	return p.ParseValue(), err
}

/*
	JSON-text = ws value ws

	value = false / null / true / object / array / number / string
	false = %x66.61.6c.73.65   ; false
	null  = %x6e.75.6c.6c      ; null
	true  = %x74.72.75.65      ; true

    begin-array     = ws %x5B ws  ; [ left square bracket
    begin-object    = ws %x7B ws  ; { left curly bracket
    end-array       = ws %x5D ws  ; ] right square bracket
    end-object      = ws %x7D ws  ; } right curly bracket
    name-separator  = ws %x3A ws  ; : colon
    value-separator = ws %x2C ws  ; , comma

	object = begin-object [ member *( value-separator member ) ] end-object
	member = string name-separator value

	array = begin-array [ value *( value-separator value ) ] end-array
*/

// ParseValue parses a value
func (p *Parser) ParseValue() AstValue {
	switch p.Current().kind {
	case BeginArrayToken:
		return p.ParseArray()
	case BeginObjectToken:
		return p.ParseObject()
	case FalseToken, TrueToken, NumberToken, NullToken, StringToken:
		return AstValue(p.Next().value)
	default:
		panic(fmt.Errorf("Unexpected token %q", p.Current().kind))
	}
}

// ParseObject parses an object
func (p *Parser) ParseObject() AstValue {
	if !p.Accept(BeginObjectToken) {
		panic(fmt.Errorf("Expected %q, got %q", BeginObjectToken, p.Current().kind))
	}

	var astObj AstObject
	for !p.IsAtEnd() &&
		p.Current().kind != EndObjectToken &&
		(p.Accept(CommaToken) || len(astObj) <= 0) {
		astObj = append(astObj, p.ParseMember())
	}

	if !p.Accept(EndObjectToken) {
		panic(fmt.Errorf("Expected %q, got %q", EndObjectToken, p.Current().kind))
	}
	return AstValue(astObj)
}

// ParseMember parses a member
func (p *Parser) ParseMember() AstMember {
	if p.IsAtEnd() || p.Current().kind != StringToken {
		panic(fmt.Errorf("Expected %q, got %q", StringToken, p.Current().kind))
	}
	keyToken := p.Next()
	if !p.Accept(ColonToken) {
		panic(fmt.Errorf("Expected %q, got %q", ColonToken, p.Current().kind))
	}
	return AstMember{
		key: keyToken.value.(string),
		val: p.ParseValue(),
	}
}

// ParseArray parses an array
func (p *Parser) ParseArray() AstValue {
	if !p.Accept(BeginArrayToken) {
		panic(fmt.Errorf("Expected %q, got %q", BeginArrayToken, p.Current().kind))
	}

	var astArr AstArray
	for !p.IsAtEnd() &&
		p.Current().kind != EndArrayToken &&
		(p.Accept(CommaToken) || len(astArr) <= 0) {
		astArr = append(astArr, p.ParseValue())
	}

	if !p.Accept(EndArrayToken) {
		panic(fmt.Errorf("Expected %q, got %q", EndArrayToken, p.Current().kind))
	}
	return astArr
}

// Accept consumes the next token if it has the desired kind and returns true,
// otherwise returns false and does not advance.
func (p *Parser) Accept(kind TokenKind) bool {
	if p.IsAtEnd() {
		return false
	}
	if p.Next().kind == kind {
		return true
	}
	p.Backup()
	return false
}

// Next advances to next token and returns the last current token.
func (p *Parser) Next() Token {
	p.current++
	return p.tokens[p.current-1]
}

// Backup retreats current position and returns token at new position.
func (p *Parser) Backup() Token {
	p.current--
	return p.tokens[p.current]
}

// IsAtEnd tells if the current token is an EOFToken or if the current position
// exceeded the length of tokens.
func (p *Parser) IsAtEnd() bool {
	if p.current >= len(p.tokens) {
		return true
	}
	return p.Current().kind == EOFToken
}

// Current returns the current token.
func (p *Parser) Current() Token {
	return p.tokens[p.current]
}
