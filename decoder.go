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
	l, err := newLexer(bufio.NewReaderSize(r, 256))
	if err != nil {
		return nil
	}
	d = &Decoder{
		lexer: l,
	}
	return d
}

func (d *Decoder) Decode(in interface{}) (err error) {
	inVal := reflect.ValueOf(in)
	if inVal.Kind() != reflect.Pointer {
		return fmt.Errorf("in is not a pointer (%s)", inVal.Kind().String())
	}
	inElem := inVal.Elem()
	firstToken, eof, err := d.lexer.Next()
	if eof {
		return nil
	}
	if err != nil {
		return err
	}
	switch firstToken.value {
	case '{':
		switch inElem.Kind() {
		case reflect.Interface:
			inMap := map[string]interface{}{}
			err = d.decodeObject(inMap)
			inElem.Set(reflect.ValueOf(inMap))
		}
	}
	return err
}

func (d *Decoder) decodeObject(m map[string]interface{}) (err error) {
	var key string
	var val interface{}
	var (
		currentToken  *token
		previousToken *token
		eof           bool
	)
outerloop:
	for {
		currentToken, eof, err = d.lexer.Next()
		if eof {
			return nil
		}
		if err != nil {
			return err
		}
		if currentToken == nil {
			continue
		}
		previousToken = d.lexer.Previous()

		switch currentToken.tokenType {
		case IDENT:
			switch previousToken.tokenType {
			case SEPARATOR:
				key = currentToken.value.(string)
			default:
				return fmt.Errorf("error parsing at token (%#v): indentifier was not proceeded by a valid token (%#v)", currentToken, previousToken)
			}
		case VALUE:
			switch previousToken.tokenType {
			case EQUAL:
				val = currentToken.value
			default:
				return fmt.Errorf("error parsing at token (%#v): value was not proceeded by a valid token (%#v)", currentToken, previousToken)
			}
		case SEPARATOR:
			switch currentToken.value {
			case OBJEND.Rune():
				break outerloop
			}
		}
		if key != "" && val != nil {
			m[key] = val
			key = ""
			val = nil
		}
	}
	return err
}
