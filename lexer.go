package main

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
		tokens = append(tokens, l.NextToken())
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
	token := Token{index: l.current}
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
		token = LexString()
	case ' ', '\r', '\t', '\n':
		l.Advance()
		goto loop
	default:
		token.kind = UnknownToken
	}

	l.Advance()
	return token
}

// LexString docs here
func LexString() Token {
	return Token{}
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
