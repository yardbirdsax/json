package json

import (
	"fmt"
	"os"
)

func DecodeFile(filename string, out interface{}) (err error) {
	f, err := os.OpenFile(filename, os.O_RDONLY, 0755)
  if err != nil {
    return fmt.Errorf("error opening file (%q): %w", filename, err)
  }
  defer f.Close()

  // lex, err := newLexer(f, tokenSplitFunc)


	return err
}

// position is a position in a file
type position struct {
  line int
  column int
  offset int64
}
