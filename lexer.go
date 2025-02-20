package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// Lexer produces tokens from an input source
type Lexer struct {
	start   int
	current int
	source  string
	err     error
}

// NewLexer creates a Lexer
func NewLexer(source string) Lexer {
	return Lexer{source: source}
}

// Lex tokenizes entire input source
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

// NextToken lexes a token at current position
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
		return l.lexString()
	case '-':
		return l.lexNumber()
	case 'n':
		if l.AcceptAll("null") {
			return Token{
				kind:   NullToken,
				index:  l.start,
				length: l.current - l.start,
			}
		}
		l.err = fmt.Errorf("Bad null token at char %d", l.current)
	case 't':
		if l.AcceptAll("true") {
			return Token{
				kind:   TrueToken,
				index:  l.start,
				length: l.current - l.start,
				value:  true,
			}
		}
		l.err = fmt.Errorf("Bad true token at char %d", l.current)
	case 'f':
		if l.AcceptAll("false") {
			return Token{
				kind:   FalseToken,
				index:  l.start,
				length: l.current - l.start,
				value:  false,
			}
		}
		l.err = fmt.Errorf("Bad false token at char %d", l.current)
	default:
		if unicode.IsSpace(l.Current()) {
			l.Advance()
			goto loop
		}
		if unicode.IsDigit(l.Current()) {
			return l.lexNumber()
		}
		l.err = fmt.Errorf("Unexpected char %q", l.Current())
	}

	l.Advance()
	return token
}

/*
	string 	  	= quotation-mark *char quotation-mark
	unescaped 	= %x20-21 / %x23-5B / %x5D-10FFFF

	char = unescaped /
		escape (
			%x22 /          ; "    quotation mark  U+0022
			%x5C /          ; \    reverse solidus U+005C
			%x2F /          ; /    solidus         U+002F
			%x62 /          ; b    backspace       U+0008
			%x66 /          ; f    form feed       U+000C
			%x6E /          ; n    line feed       U+000A
			%x72 /          ; r    carriage return U+000D
			%x74 /          ; t    tab             U+0009
			%x75 4HEXDIG )  ; uXXXX                U+XXXX
	escape 			= %x5C  ; \
	quotation-mark 	= %x22  ; "
*/
func (l *Lexer) lexString() Token {
	if !l.Accept(`"`) {
		l.err = fmt.Errorf("Bad string")
		return Token{}
	}

	var b strings.Builder
	for !l.IsAtEnd() && !l.Accept(`"`) {
		// if !l.Accept(`"/\bfnrtu`)
		if l.Accept(`\`) {
			switch l.Advance() {
			case '"', '/', '\\':
				b.WriteRune(l.LookBehind())
			case 'b':
				b.WriteRune('\b')
			case 'f':
				b.WriteRune('\f')
			case 'n':
				b.WriteRune('\n')
			case 'r':
				b.WriteRune('\r')
			case 't':
				b.WriteRune('\t')
			case 'u':
				for i := 0; i < 4; i++ {
					if l.IsAtEnd() || !l.Accept("0123456789abcedfABCDEF") {
						l.err = fmt.Errorf("Bad scaped hex")
						return Token{}
					}
				}
				hex, err := strconv.ParseInt(l.source[l.current-4:l.current], 16, 64)
				if err != nil {
					panic(err)
				}
				b.WriteRune(rune(hex))
			default:
				l.err = fmt.Errorf("Bad string scaped char %s", l.source[l.current-1:l.current+1])
				return Token{}
			}
			continue
		}
		b.WriteRune(l.Advance())
	}

	return Token{
		kind:   StringToken,
		index:  l.start,
		length: l.current,
		value:  b.String(),
	}
}

/*
	number = [ minus ] int [ frac ] [ exp ]

	int 			= zero / ( digit1-9 *DIGIT )
	frac 			= decimal-point 1*DIGIT
	exp 			= e [ minus / plus ] 1*DIGIT
	decimal-point 	= %x2E       	; .
	digit1-9 		= %x31-39       ; 1-9
	e 				= %x65 / %x45   ; e E
	minus 			= %x2D          ; -
	plus 			= %x2B          ; +
	zero 			= %x30          ; 0
*/
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

// Current returns rune at current position
func (l *Lexer) Current() rune {
	return rune(l.source[l.current])
}

// Accept consumes the current rune it is contained in a string and returns true;
// otherwise current position remains the same and returns false.
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

// AcceptAll consumes runes if they match the exact sequence of a string and returns
// true; if the string is not entirely matched, current position remains the same
// and returns false.
func (l *Lexer) AcceptAll(s string) bool {
	for i := 0; i < len(s); i++ {
		if l.IsAtEnd() || !l.Accept(s[i:i+1]) {
			l.current -= i
			return false
		}
	}
	return true
}

// LookAhead returns the next rune without advancing.
func (l *Lexer) LookAhead() rune {
	return rune(l.source[l.current+1])
}

// LookBehind returns the last rune read.
func (l *Lexer) LookBehind() rune {
	return rune(l.source[l.current-1])
}

// Advance advances the current position and returns rune at old position.
func (l *Lexer) Advance() rune {
	l.current++
	return rune(l.source[l.current-1])
}

// AdvanceTillSpace consumes runes until it reaches the position of a space, then
// it stops and returns true. If it gets to the end of the input before that
// returns false.
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
