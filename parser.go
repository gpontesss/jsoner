package main

import "fmt"

// Parser docs here
type Parser struct {
	tokens  []Token
	current int
}

// NewParser docs here
func NewParser(tokens []Token) Parser {
	return Parser{tokens: tokens}
}

// Parse docs here
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

// ParseValue docs here
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

// ParseObject docs here
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

// ParseMember docs here
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

// ParseArray docs here
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

// Accept docs here
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

// Next docs here
func (p *Parser) Next() Token {
	p.current++
	return p.tokens[p.current-1]
}

// Backup docs here
func (p *Parser) Backup() Token {
	p.current--
	return p.tokens[p.current]
}

// IsAtEnd docs here
func (p *Parser) IsAtEnd() bool {
	if p.current >= len(p.tokens) {
		return true
	}
	return p.Current().kind == EOFToken
}

// Current docs here
func (p *Parser) Current() Token {
	return p.tokens[p.current]
}
