package differ

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrettifyDelta(t *testing.T) {
	testCases := []struct {
		name           string
		input          map[int]Delta
		expectedOutput map[int]PrettyDelta
	}{
		{
			name: "Regular case",
			input: map[int]Delta{
				0: {StartIndex: 0, EndIndex: 16, Deleted: false, UpdatedLiterals: []byte("This is a test")},
				1: {StartIndex: 17, EndIndex: 32, Deleted: true},
			},
			expectedOutput: map[int]PrettyDelta{
				0: {startIndex: 0, endIndex: 16, deleted: false, updatedLiterals: "This is a test"},
				1: {startIndex: 17, endIndex: 32, deleted: true, updatedLiterals: ""},
			},
		},
		{
			name:           "Empty input",
			input:          map[int]Delta{},
			expectedOutput: map[int]PrettyDelta{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualOutput := PrettifyDelta(tc.input)
			assert.Equal(t, tc.expectedOutput, actualOutput)
		})
	}
}
