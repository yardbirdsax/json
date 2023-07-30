package json

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecode(t *testing.T) {
	testCases := []struct {
		name          string
		expectedError *string
		obj           interface{}
		want          interface{}
	}{
		{
			name:          "empty",
			expectedError: nil,
			obj:           map[string]interface{}{},
			want:          map[string]interface{}{},
		},
		{
			name:          "simple_object",
			expectedError: nil,
			obj:           map[string]interface{}{},
			want:          map[string]interface{}{"string": "a string", "another_string": "another string"},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			testDataFile := filepath.Join("testdata", tc.name+".json")

			err := DecodeFile(testDataFile, &tc.obj)

			if tc.expectedError == nil {
				assert.NoError(t, err, "error occurred decoding data")
			} else {
				assert.ErrorContains(t, err, *tc.expectedError, "returned error did not contain expected string")
			}
			assert.Equal(t, tc.want, tc.obj)
		})
	}
}
