package json

import (
	"errors"
	"fmt"
	"reflect"
)

type decoder struct {
	tokens []*token
	pos    int
}

func newDecoder(t []*token) (p *decoder, err error) {
	if len(t) == 0 {
		return nil, errors.New("cannot initialize decoder with zero length tokens slice")
	}
	p = &decoder{
		tokens: t,
	}
	return p, nil
}

func (p *decoder) decode(in interface{}) (err error) {
	inVal := reflect.ValueOf(in)
	firstToken := p.tokens[0]
	switch firstToken.value {
	case '{':
		if inVal.Kind() != reflect.Map {
			return fmt.Errorf("JSON object cannot be decoded into type other than map (type is %q)", inVal.Kind())
		}
		if inVal.Type().Elem().Kind() != reflect.Interface {
			return fmt.Errorf("JSON object cannot be decoded into map of type other than interface (type is %q)", inVal.Elem().Kind())
		}
		inMap := in.(map[string]interface{})
		_, err = decodeObject(p.tokens, 0, inMap)
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
