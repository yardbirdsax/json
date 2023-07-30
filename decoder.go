package json

import (
	"bufio"
	"fmt"
	"io"
	"reflect"
)

type Decoder struct {
	lexer  *lexer
	tokens []*token
	pos    int
}

func NewDecoder(r io.Reader) (d *Decoder) {
	l, err := newLexer(bufio.NewReader(r))
	if err != nil {
		return nil
	}
	d = &Decoder{
		lexer: l,
	}
	return d
}

func (d *Decoder) lex() (err error) {
	tokens, err := d.lexer.Lex()
	d.tokens = tokens
	return err
}

func (p *Decoder) Decode(in interface{}) (err error) {
	err = p.lex()
	if err != nil {
		return err
	}
	inVal := reflect.ValueOf(in)
	if inVal.Kind() != reflect.Pointer {
		return fmt.Errorf("in is not a pointer (%s)", inVal.Kind().String())
	}
	inElem := inVal.Elem()
	firstToken := p.tokens[0]
	switch firstToken.value {
	case '{':
		switch inElem.Kind() {
		case reflect.Interface:
			inMap := map[string]interface{}{}
			_, err = decodeObject(p.tokens, 0, inMap)
			inElem.Set(reflect.ValueOf(inMap))
		}
	}
	return err
}

func decodeObject(tokens []*token, start int, m map[string]interface{}) (newStart int, err error) {
	if tokens[start].value != '{' {
		return newStart, fmt.Errorf("invalid first token for object: %q", tokens[0].value)
	}

	current := start + 1
	var key string
	var val interface{}
outerloop:
	for {
		if current > len(tokens)-1 {
			return 0, fmt.Errorf("current token start (%d) is greater than length of tokens (%d)", current+1, len(tokens))
		}
		switch tokens[current].tokenType {
		case IDENT:
			switch tokens[current-1].tokenType {
			case SEPARATOR:
				key = tokens[current].value.(string)
			default:
				return 0, fmt.Errorf("error parsing at token (%#v): indentifier was not proceeded by a valid token (%#v)", tokens[current], tokens[current-1])
			}
		case VALUE:
			switch tokens[current-1].tokenType {
			case EQUAL:
				val = tokens[current].value
			default:
				return 0, fmt.Errorf("error parsing at token (%#v): value was not proceeded by a valid token (%#v)", tokens[current], tokens[current-1])
			}
		case SEPARATOR:
			switch tokens[current].value {
			case OBJEND.Rune():
				current++
				break outerloop
			}
		}
		if key != "" && val != nil {
			m[key] = val
			key = ""
			val = nil
		}
		current++
	}
	return current, err
}
