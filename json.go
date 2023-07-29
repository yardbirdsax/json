package json

import (
	"bufio"
	"fmt"
	"os"
)

func DecodeFile(filename string, out interface{}) (err error) {
	f, err := os.OpenFile(filename, os.O_RDONLY, 0755)
  if err != nil {
    return fmt.Errorf("error opening file (%q): %w", filename, err)
  }
  defer f.Close()

  l, err := newLexer(bufio.NewReader(f))
  if err != nil {
    return err
  }
  tokens, err := l.Lex()
  if err != nil {
    return err
  }
  decoder, err := newDecoder(tokens)
  if err != nil {
    return err
  }

	return decoder.decode(out)
}

// position is a position in a file
type position struct {
  line int
  column int
  offset int64
}
