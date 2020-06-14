package main

import (
	"strings"
	"testing"
)

func TestLexer(t *testing.T) {
	t.Run("Accept", func(t *testing.T) {
		cases := []struct {
			source string
			s      string

			expCurrent int
			expReturn  bool
		}{
			{"+...", "+-", 1, true},
		}

		for _, c := range cases {
			lexer := NewLexer(c.source)
			accepted := lexer.Accept(c.s)

			if c.expReturn != accepted {
				t.Errorf("Expected return %v, got %v", c.expReturn, accepted)
			}
			if c.expCurrent != lexer.current {
				t.Errorf("Expected current %v, got %v", c.expCurrent, lexer.current)
			}
		}
	})

	t.Run("lexNumber", func(t *testing.T) {
		cases := []struct {
			input   string
			expKind TokenKind

			expVal         interface{}
			expErrContains string
		}{
			{
				input:   "12",
				expKind: NumberToken,
				expVal:  12.0,
			},
			{
				input:   "0",
				expKind: NumberToken,
				expVal:  0.0,
			},
			{
				input:   "-0",
				expKind: NumberToken,
				expVal:  0.0,
			},
			{
				input:   "10.10",
				expKind: NumberToken,
				expVal:  10.10,
			},
			{
				input:   "12e-1",
				expKind: NumberToken,
				expVal:  1.2,
			},
			{
				input:   "4E1",
				expKind: NumberToken,
				expVal:  40.0,
			},
			{
				input:   "-10.23",
				expKind: NumberToken,
				expVal:  -10.23,
			},
			{
				input:          "01",
				expErrContains: "Leading zeros",
			},
			{
				input:          "--0",
				expErrContains: "Bad number",
			},
		}

		for _, c := range cases {
			t.Run(c.input, func(t *testing.T) {
				lexer := NewLexer(c.input)
				token := lexer.lexNumber()

				if token.kind != c.expKind {
					t.Errorf("Expected %v, got %v", c.expKind, token.kind)
				}
				if c.expErrContains != "" {
					if lexer.err != nil {
						if !strings.Contains(lexer.err.Error(), c.expErrContains) {
							t.Errorf("Expected error '%s' to contain '%s'", lexer.err.Error(), c.expErrContains)
						}
					} else {
						t.Errorf("Expected err, got %v", lexer.err)
					}
				} else {
					if token.value != c.expVal {
						t.Errorf("Expected %v, got %v", c.expVal, token.value)
					}
				}
			})
		}
	})
}
