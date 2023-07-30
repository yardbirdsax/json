package json

import (
	"errors"
	"io"
	"strings"
)

type lexer struct {
	pos    position
	reader io.RuneReader
}

type token struct {
	pos       position
	value     interface{}
	tokenType tokenType
}

type tokenType int

const (
	INVALID tokenType = iota
	SEPARATOR
	IDENT
	EQUAL
	VALUE
)

func newLexer(data io.RuneReader) (lex *lexer, err error) {
	lex = &lexer{}
	lex.pos = position{
		line:   1,
		column: 1,
		offset: 0,
	}
	lex.reader = data
	return lex, err
}

func (l *lexer) Lex() (tokens []*token, err error) {
	currToken := &token{}
	prevToken := &token{}
	for {
		r, o, err := l.reader.ReadRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return tokens, nil
			}
			return tokens, err
		}
		switch r {
		case OBJSTART.Rune(), OBJEND.Rune(), COMMA.Rune():
			currToken = &token{
				pos:       l.pos,
				value:     r,
				tokenType: SEPARATOR,
			}
			tokens = append(tokens, currToken)
			l.pos.column++
		case COLON.Rune():
			currToken = &token{
				pos:       l.pos,
				value:     r,
				tokenType: EQUAL,
			}
			tokens = append(tokens, currToken)
			l.pos.column++
		case '\n', '\r':
			l.pos.column = 1
			l.pos.line++
		case ' ', '\t':
			l.pos.column++
		case QUOTE.Rune():
			s := strings.Builder{}
			var tt tokenType
			switch prevToken.tokenType {
			case EQUAL:
				tt = VALUE
			case SEPARATOR:
				tt = IDENT
			}
			currToken = &token{
				pos:       l.pos,
				tokenType: tt,
			}
			l.pos.column++
		outerloop:
			for {
				rs, os, err := l.reader.ReadRune()
				if err != nil {
					if errors.Is(err, io.EOF) {
						return tokens, io.ErrUnexpectedEOF
					}
					return tokens, err
				}
				l.pos.column++
				l.pos.offset += int64(os)
				switch rs {
				case QUOTE.Rune():
					currToken.value = s.String()
					tokens = append(tokens, currToken)
					break outerloop
				default:
					s.WriteRune(rs)
				}
			}
		}
		l.pos.offset += int64(o)
		*prevToken = *currToken
	}
}
