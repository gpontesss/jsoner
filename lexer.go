package main

import (
	"fmt"
	"strconv"
	"strings"
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
	case '-':
		return l.lexNumber()
	default:
		if unicode.IsSpace(l.Current()) {
			l.Advance()
			goto loop
		}
		if unicode.IsDigit(l.Current()) {
			return l.lexNumber()
		}
		// Error: unknown token
		l.err = fmt.Errorf("Unexpected char %q", l.Current())
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
func (l *Lexer) lexNumber() Token {
	l.Accept("-")
	if l.IsAtEnd() || !unicode.IsDigit(l.Current()) {
		l.AdvanceTillSpace()
		l.err = fmt.Errorf("Bad number %s", l.source[l.start:l.current])
		return Token{}
	}

	if l.Accept("0") && !l.IsAtEnd() && unicode.IsDigit(l.Current()) {
		// "Leading zeros are *not allowed*"
		l.err = fmt.Errorf("Leading zeros are not allowed")
		return Token{}
	}

	for !l.IsAtEnd() && unicode.IsDigit(l.Current()) {
		// Consume all integer digits
		l.Advance()
	}

	if l.Accept(".") {
		if l.IsAtEnd() || !unicode.IsDigit(l.Current()) {
			// "A fraction part is a decimal point followed by *one or more* digits"
			l.err = fmt.Errorf("Expected integer, got %s", strconv.QuoteRune(l.Current()))
			return Token{}
		}
		for !l.IsAtEnd() && unicode.IsDigit(l.Current()) {
			// Consume all fractional digit characters
			l.Advance()
		}
	}

	if l.Accept("eE") {
		// An exponent part begins with the letter E in upper or lower case,
		// which may be followed by a plus or minus sign. The E and optional
		// sign are followed by one or more digits
		l.Accept("+-")
		if !l.IsAtEnd() && !unicode.IsDigit(l.Current()) {
			l.err = fmt.Errorf("Expected integer, got %s", strconv.QuoteRune(l.Current()))
			return Token{}
		}

		for !l.IsAtEnd() && unicode.IsDigit(l.Current()) {
			// Consume all fractional digit characters
			l.Advance()
		}
	}

	val, err := strconv.ParseFloat(l.source[l.start:l.current], 64)
	if err != nil {
		// For debugging purposes
		panic(err)
	}

	return Token{
		kind:   NumberToken,
		index:  l.start,
		length: l.current - l.start,
		value:  val,
	}
}

// Current docs here
func (l *Lexer) Current() rune {
	return rune(l.source[l.current])
}

// Accept docs here
func (l *Lexer) Accept(s string) bool {
	if l.IsAtEnd() {
		return false
	}

	if strings.ContainsRune(s, l.Advance()) {
		return true
	}
	l.Backup()
	return false
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

// AdvanceTillSpace docs here
func (l *Lexer) AdvanceTillSpace() bool {
	for !l.IsAtEnd() && unicode.IsSpace(l.Current()) {
		l.Advance()
	}
	return l.IsAtEnd()
}

// Backup docs here
func (l *Lexer) Backup() {
	l.current--
}

// IsAtEnd docs here
func (l *Lexer) IsAtEnd() bool {
	return l.current >= len(l.source)
}
