package json

import (
	"bufio"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLex(t *testing.T) {
	testCases := []struct {
		name          string
		expectedError *string
		want          []*token
	}{
		{
			name:          "empty",
			expectedError: nil,
			want: []*token{
				{
					pos: position{
						line:   1,
						column: 1,
					},
					value:     '{',
					tokenType: SEPARATOR,
				},
				{
					pos: position{
						line:   1,
						column: 2,
					},
					value:     '}',
					tokenType: SEPARATOR,
				},
			},
		},
		{
			name:          "simple_object",
			expectedError: nil,
			want: []*token{
				{
					pos: position{
						line:   1,
						column: 1,
					},
					value:     '{',
					tokenType: SEPARATOR,
				},
				{
					pos: position{
						line:   2,
						column: 3,
					},
					value:     "string",
					tokenType: IDENT,
				},
				{
					pos: position{
						line:   2,
						column: 11,
					},
					value:     ':',
					tokenType: EQUAL,
				},
				{
					pos: position{
						line:   2,
						column: 13,
					},
					value:     "a string",
					tokenType: VALUE,
				},
				{
					pos: position{
						line:   2,
						column: 23,
					},
					value:     ',',
					tokenType: SEPARATOR,
				},
				{
					pos: position{
						line:   3,
						column: 3,
					},
					value:     "another_string",
					tokenType: IDENT,
				},
				{
					pos: position{
						line:   3,
						column: 19,
					},
					value:     ':',
					tokenType: EQUAL,
				},
				{
					pos: position{
						line:   3,
						column: 21,
					},
					value:     "another string",
					tokenType: VALUE,
				},
				{
					pos: position{
						line:   4,
						column: 1,
					},
					value:     '}',
					tokenType: SEPARATOR,
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			testDataFile := filepath.Join("testdata", tc.name+".json")
			f, err := os.Open(testDataFile)
			require.NoError(t, err, "error opening file")
			r := bufio.NewReader(f)
			l, err := newLexer(r)
			require.NoError(t, err, "error creating lexer")

			got, err := l.Lex()

			assert.NoError(t, err, "error lexing")
			assert.EqualValues(t, tc.want, got)
		})
	}
}
