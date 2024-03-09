package hash

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdler32(t *testing.T) {
	testCases := []struct {
		name         string
		input        string
		expectedHash uint
	}{
		{
			name:         "Test Hash",
			input:        "This is a test",
			expectedHash: 611517686,
		},
		{
			name:         "Test RollIn",
			input:        "This is a test",
			expectedHash: 611517686,
		},
		{
			name:         "Test RollOut",
			input:        "aThis is a test",
			expectedHash: 611517686,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			adler32 := NewAdler32(4)
			switch tc.name {
			case "Test Hash":
				adler32.Write([]byte(tc.input))
				hash := adler32.Hash()
				assert.Equal(t, tc.expectedHash, hash)
			case "Test RollIn":
				stringBytes := []byte(tc.input)
				var hash uint
				for len(stringBytes) > 0 {
					hash = adler32.RollIn(stringBytes[0])
					stringBytes = stringBytes[1:]
				}
				assert.Equal(t, tc.expectedHash, hash)
			case "Test RollOut":
				stringBytes := []byte(tc.input)
				var hash uint
				adler32.Write(stringBytes)
				adler32.Hash()
				hash, removed := adler32.RollOut()
				assert.Equal(t, hash, tc.expectedHash)
				assert.Equal(t, uint8([]byte("a")[0]), removed)
			}
		})
	}
}

func TestGetWindowLiterals(t *testing.T) {
	adler32 := NewAdler32(4)
	adler32.Write([]byte("Test"))
	assert.Equal(t, []byte("Test"), adler32.GetWindowLiterals())
}

func TestGetWindowLength(t *testing.T) {
	adler32 := NewAdler32(4)
	adler32.RollIn([]byte("A")[0])
	assert.Equal(t, 1, adler32.WindowLength())
}
