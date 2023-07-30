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

	d := NewDecoder(f)
	return d.Decode(out)
}

// position is a position in a file
type position struct {
	line   int
	column int
}
