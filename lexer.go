package main

import (
	"fmt"
	"strconv"
	"unicode"
)

// Lexer docs here
type Lexer struct {
	start   int
	current int
	source  string
	err     error
}

// NewLexer docs here
func NewLexer(source string) Lexer {
	return Lexer{source: source}
}

// Lex docs here
func (l *Lexer) Lex() ([]Token, error) {
	tokens := make([]Token, 0)
	for !l.IsAtEnd() {
		token := l.NextToken()
		if l.err != nil {
			return tokens, l.err
		}
		tokens = append(tokens, token)
	}
	return tokens, l.err
}

// NextToken docs here
func (l *Lexer) NextToken() Token {
	if l.IsAtEnd() {
		return Token{
			kind:  EOFToken,
			index: l.current,
		}
	}

loop:
	l.start = l.current
	token := Token{index: l.start, length: 1}
	switch l.Current() {
	case '[':
		token.kind = BeginArrayToken
	case ']':
		token.kind = EndArrayToken
	case '{':
		token.kind = BeginObjectToken
	case '}':
		token.kind = EndObjectToken
	case ':':
		token.kind = ColonToken
	case ',':
		token.kind = CommaToken
	case '"':
		token = l.lexString()
	case ' ', '\r', '\t', '\n':
		l.Advance()
		goto loop
	case '-':
		l.Advance()
		token = l.lexNumber(true)
	default:
		if unicode.IsDigit(l.Current()) {
			token = l.lexNumber(false)
		} else {
			// Error: unknown token
			l.err = fmt.Errorf("Unexpected char '%s'", l.Current())
			token.kind = UnknownToken
		}
	}

	l.Advance()
	return token
}

func (l *Lexer) lexString() Token {
	return Token{}
}

// number = [ minus ] int [ frac ] [ exp ]
//
// int 				= zero / ( digit1-9 *DIGIT )
// frac 			= decimal-point 1*DIGIT
// exp 				= e [ minus / plus ] 1*DIGIT
// decimal-point 	= %x2E       	; .
// digit1-9 		= %x31-39       ; 1-9
// e 				= %x65 / %x45   ; e E
// minus 			= %x2D          ; -
// plus 			= %x2B          ; +
// zero 			= %x30          ; 0
//
func (l *Lexer) lexNumber(minus bool) Token {
	if l.Current() == '0' && unicode.IsDigit(l.LookAhead()) {
		// "Leading zeros are *not allowed*"
		// Error
		return Token{}
	}

	for unicode.IsDigit(l.Current()) {
		// Consume all integer digits
		l.Advance()
	}

	end := l.current
	if l.Current() == '.' {
		l.Advance()
		if !unicode.IsDigit(l.Current()) {
			// "A fraction part is a decimal point followed by *one or more* digits"
			l.err = fmt.Errorf("Expected integer, got %s", strconv.QuoteRune(l.Current()))
			return Token{}
		}
		for unicode.IsDigit(l.Current()) {
			// Consume all fractional digit characters
			l.Advance()
		}
		end = l.current
	}

	if unicode.ToLower(l.Current()) == 'e' {
		// An exponent part begins with the letter E in upper or lower case,
		// which may be followed by a plus or minus sign. The E and optional
		// sign are followed by one or more digits
		if l.Advance(); l.Current() == '+' || l.Current() == '-' {
			l.Advance()
		}
		if !unicode.IsDigit(l.Current()) {
			l.err = fmt.Errorf("Expected integer, got %s", strconv.QuoteRune(l.Current()))
			return Token{}
		}

		for unicode.IsDigit(l.Advance()) {
			// Consume all fractional digit characters
		}
		end = l.current - 1
	}

	val, err := strconv.ParseFloat(l.source[l.start:end], 64)
	if err != nil {
		// For debugging purposes
		panic(err)
	}

	return Token{
		kind:   NumberToken,
		index:  l.start,
		length: end - l.start,
		value:  val,
	}
}

// Current docs here
func (l *Lexer) Current() rune {
	return rune(l.source[l.current])
}

// LookAhead docs here
func (l *Lexer) LookAhead() rune {
	return rune(l.source[l.current+1])
}

// Advance docs here
func (l *Lexer) Advance() rune {
	l.current++
	return rune(l.source[l.current-1])
}

// IsAtEnd docs here
func (l *Lexer) IsAtEnd() bool {
	return l.current >= len(l.source)
}
