package json

import (
	"encoding/json"
	"os"
	"testing"
)

func BenchmarkDecode(b *testing.B) {
	out := map[string]interface{}{}
	f, err := os.Open("testdata/simple_object.json")
	if err != nil {
		b.Errorf("error opening file: %s", err)
	}
	defer f.Close()
	for i := 0; i < b.N; i++ {
		f.Seek(0, 0)
		decoder := NewDecoder(f)
		err = decoder.Decode(&out)
		if err != nil {
			b.Errorf("error decoding file: %s", err)
		}
	}
}

func BenchmarkDecodeNative(b *testing.B) {
	out := map[string]interface{}{}
	f, err := os.Open("testdata/simple_object.json")
	if err != nil {
		b.Errorf("error opening file: %s", err)
	}
	defer f.Close()
	for i := 0; i < b.N; i++ {
		f.Seek(0, 0)
		decoder := json.NewDecoder(f)
		err = decoder.Decode(&out)
		if err != nil {
			b.Errorf("error decoding file: %s", err)
		}
	}
}
