package json

import (
	"errors"
	"io"
	"strings"
)

type lexer struct {
	pos           position
	reader        io.RuneReader
	currentToken  *token
	previousToken *token
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
	}
	lex.reader = data
	lex.currentToken = &token{}
	lex.previousToken = &token{}
	return lex, err
}

func (l *lexer) Previous() (t *token) {
	return l.previousToken
}

func (l *lexer) Next() (t *token, eof bool, err error) {
	r, _, err := l.reader.ReadRune()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil, true, nil
		}
		return nil, false, err
	}

	switch r {
	case OBJSTART.Rune(), OBJEND.Rune(), COMMA.Rune():
		*l.previousToken = *l.currentToken
		l.currentToken.pos = l.pos
		l.currentToken.value = r
		l.currentToken.tokenType = SEPARATOR
		l.pos.column++
	case COLON.Rune():
		*l.previousToken = *l.currentToken
		l.currentToken.pos = l.pos
		l.currentToken.value = r
		l.currentToken.tokenType = EQUAL
		l.pos.column++
	case '\n', '\r':
		l.pos.column = 1
		l.pos.line++
		return nil, false, nil
	case ' ', '\t':
		l.pos.column++
		return nil, false, nil
	case QUOTE.Rune():
		*l.previousToken = *l.currentToken
		s := strings.Builder{}
		var tt tokenType
		switch l.previousToken.tokenType {
		case EQUAL:
			tt = VALUE
		case SEPARATOR:
			tt = IDENT
		}
		l.currentToken.pos = l.pos
		l.currentToken.tokenType = tt
		l.pos.column++
	outerloop:
		for {
			rs, _, err := l.reader.ReadRune()
			if err != nil {
				if errors.Is(err, io.EOF) {
					return nil, true, io.ErrUnexpectedEOF
				}
				return nil, false, err
			}
			l.pos.column++
			switch rs {
			case QUOTE.Rune():
				l.currentToken.value = s.String()
				break outerloop
			default:
				s.WriteRune(rs)
			}
		}
	}
	t = &token{}
	*t = *l.currentToken
	return t, false, nil
}

func (l *lexer) Lex() (tokens []*token, err error) {
	for {
		token, eof, err := l.Next()
		if err != nil || (token == nil && eof) {
			return tokens, err
		}
		if token != nil {
			tokens = append(tokens, token)
		}
	}
}
