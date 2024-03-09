package differ

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiffer(t *testing.T) {
	testCases := []struct {
		name           string
		txt1           string
		txt2           string
		expectedDeltas map[int]PrettyDelta
	}{
		{
			name: "Chunk Change",
			txt1: "This is a Rolling hash file diff algorithm. It should check for changes in file and text",
			txt2: "This is a Rolling hashes file difference algorithm. It should check for changes in file and text",
			expectedDeltas: map[int]PrettyDelta{
				1: {startIndex: 17, endIndex: 32, deleted: true},
				2: {startIndex: 17, endIndex: 32, updatedLiterals: "g hashes file difference"},
			},
		},
		{
			name: "Chunk Deletion",
			txt1: "This is a Rolling hash file diff algorithm. It should check for changes in file and text",
			txt2: "This is a Rolling hash file diff algorithm. It should check for changes",
			expectedDeltas: map[int]PrettyDelta{
				4: {startIndex: 64, endIndex: 80, updatedLiterals: "changes"},
				5: {startIndex: 81, endIndex: 96, deleted: true},
			},
		},
		{
			name: "Chunk Addition",
			txt1: "This is a Rolling hash file diff algorithm. It should check for changes in file and text",
			txt2: "This is a Rolling hash file diff algorithm. It should check for changes in file and text. This is written in a way to detect addition to the text",
			expectedDeltas: map[int]PrettyDelta{
				5: {startIndex: 80, endIndex: 96, updatedLiterals: "and text. This is written in a way to detect addition to the text"},
			},
		},
		{
			name:           "No Changes",
			txt1:           "This is a Rolling hash file diff algorithm. It should check for changes in file and text",
			txt2:           "This is a Rolling hash file diff algorithm. It should check for changes in file and text",
			expectedDeltas: map[int]PrettyDelta{},
		},
		{
			name: "First Chunk Change",
			txt1: "This is a Rolling hash file diff algorithm. It should check for changes in file and text",
			txt2: "The a Rolling hash file diff algorithm. It should check for changes in file and text",
			expectedDeltas: map[int]PrettyDelta{
				0: {startIndex: 1, endIndex: 16, deleted: true},
				1: {startIndex: 1, endIndex: 16, updatedLiterals: "The a Rollin"},
			},
		},
		{
			name: "All Chunks Change",
			txt1: "This is a Rolling hash file diff algorithm. It should check for changes in file and text",
			txt2: "This is a different text and it is different from all the chunks above",
			expectedDeltas: map[int]PrettyDelta{
				0: {startIndex: 0, endIndex: 16, updatedLiterals: "This is a different text and it is different from all the chunks above"},
				1: {startIndex: 17, endIndex: 32, deleted: true},
				2: {startIndex: 33, endIndex: 48, deleted: true},
				3: {startIndex: 49, endIndex: 64, deleted: true},
				4: {startIndex: 65, endIndex: 80, deleted: true},
				5: {startIndex: 81, endIndex: 96, deleted: true},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			differInstance := New(16)
			reader1 := bytes.NewReader([]byte(tc.txt1))
			buffReader1 := bufio.NewReader(reader1)

			reader2 := bytes.NewReader([]byte(tc.txt2))
			buffReader2 := bufio.NewReader(reader2)

			signatures := differInstance.GenerateSignatures(buffReader1)
			deltas := differInstance.GenerateDelta(signatures, buffReader2)
			prettyDelta := PrettifyDelta(deltas)

			assert.Equal(t, tc.expectedDeltas, prettyDelta)
		})
	}
}
