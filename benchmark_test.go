package json

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"
)

func BenchmarkDecode(b *testing.B) {
	out := map[string]interface{}{}
	by, err := os.ReadFile("testdata/simple_object.json")
	if err != nil {
		b.Errorf("error opening file: %s", err)
	}
	buf := bytes.NewReader(by)
	for i := 0; i < b.N; i++ {
		buf.Reset(by)
		decoder := NewDecoder(buf)
		err = decoder.Decode(&out)
		if err != nil {
			b.Errorf("error decoding file: %s", err)
		}
	}
}

func BenchmarkDecodeNative(b *testing.B) {
	out := map[string]interface{}{}
	by, err := os.ReadFile("testdata/simple_object.json")
	if err != nil {
		b.Errorf("error opening file: %s", err)
	}
	buf := bytes.NewReader(by)
	for i := 0; i < b.N; i++ {
		buf.Reset(by)
		decoder := json.NewDecoder(buf)
		err = decoder.Decode(&out)
		if err != nil {
			b.Errorf("error decoding file: %s", err)
		}
	}
}
