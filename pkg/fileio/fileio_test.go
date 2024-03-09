package fileio

import (
	"encoding/gob"
	"errors"
	"os"
	"testing"

	"github.com/Psykepro/rdiff/pkg/differ"
	"github.com/stretchr/testify/assert"
)

const (
	validFilePath   = "../../test/testdata/original.txt"
	nonExistingPath = "nonExistingPath"
)

func TestNewFileHandler(t *testing.T) {
	fileHandler := NewFileHandler(16)
	assert.Equal(t, 16, fileHandler.ChunkSize())
}

func TestChunkSize(t *testing.T) {
	fileHandler := NewFileHandler(16)
	assert.Equal(t, 16, fileHandler.ChunkSize())
}

func TestOpen(t *testing.T) {
	testCases := []struct {
		name        string
		chunkSize   int
		filePath    string
		expectedErr error
	}{
		{
			name:        "Valid File",
			chunkSize:   16,
			filePath:    validFilePath,
			expectedErr: nil,
		},
		{
			name:        "Not Enough Chunks",
			chunkSize:   1500,
			filePath:    validFilePath,
			expectedErr: ErrFileSize,
		},
		{
			name:        "Error Opening File",
			chunkSize:   16,
			filePath:    nonExistingPath,
			expectedErr: NewReadFileError(errors.New("open " + nonExistingPath + ": no such file or directory")),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fileHandler := NewFileHandler(tc.chunkSize)
			reader, err := fileHandler.Open(tc.filePath)
			if tc.expectedErr != nil {
				assert.NotNil(t, err)
				assert.Nil(t, reader)
				assert.Equal(t, tc.expectedErr.Error(), err.Error())
			} else {
				assert.NotNil(t, reader)
				assert.Nil(t, err)
			}
		})
	}
}

func TestWriteAndReadSignatures(t *testing.T) {
	testCases := []struct {
		name        string
		signatures  map[uint]int
		expectedErr error
	}{
		{
			name: "Valid Signatures",
			signatures: map[uint]int{
				1: 0,
				2: 16,
				3: 32,
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			file, err := os.CreateTemp("", "signatures_test")
			assert.NoError(t, err)
			defer os.Remove(file.Name())

			fileHandler := NewFileHandler(16)
			err = fileHandler.WriteSignatures(tc.signatures, file.Name())
			assert.Equal(t, tc.expectedErr, err)

			readSignatures, err := fileHandler.ReadSignatures(file.Name())
			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.signatures, readSignatures)
		})
	}
}

func TestWriteDelta(t *testing.T) {
	testCases := []struct {
		name        string
		delta       map[int]differ.Delta
		expectedErr error
	}{
		{
			name: "Valid Delta",
			delta: map[int]differ.Delta{
				0: {StartIndex: 0, EndIndex: 16, Deleted: false, UpdatedLiterals: []byte("Updated content")},
				1: {StartIndex: 16, EndIndex: 32, Deleted: true},
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			file, err := os.CreateTemp("", "delta_test")
			assert.NoError(t, err)
			defer os.Remove(file.Name())

			fileHandler := NewFileHandler(16)
			err = fileHandler.WriteDelta(tc.delta, file.Name())
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestReadDelta(t *testing.T) {
	testCases := []struct {
		name        string
		delta       map[int]differ.Delta
		isInvalid   bool
		expectedErr error
	}{
		{
			name: "Valid Delta",
			delta: map[int]differ.Delta{
				0: {StartIndex: 0, EndIndex: 16, Deleted: false, UpdatedLiterals: []byte("Updated content")},
				1: {StartIndex: 16, EndIndex: 32, Deleted: true},
			},
			expectedErr: nil,
		},
		{
			name:        "Non-existent File",
			delta:       nil,
			expectedErr: NewOpenFileError(nonExistingPath),
		},
		{
			name:        "Invalid Delta File",
			delta:       nil,
			isInvalid:   true,
			expectedErr: ErrDecodeDelta,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var deltaFile string
			if tc.delta != nil || tc.isInvalid {
				file, err := os.CreateTemp("", "delta_test")
				assert.NoError(t, err)
				deltaFile = file.Name()
				defer os.Remove(deltaFile)
				fileHandler := NewFileHandler(16)
				err = fileHandler.WriteDelta(tc.delta, deltaFile)
				assert.NoError(t, err)
				if tc.isInvalid {
					// Write invalid data to the file
					file.Seek(0, 0)
					encoder := gob.NewEncoder(file)
					err = encoder.Encode("invalid delta data")
					assert.NoError(t, err)
				}
			} else {
				deltaFile = nonExistingPath
			}

			fileHandler := NewFileHandler(16)
			readDelta, err := fileHandler.ReadDelta(deltaFile)

			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr.Error(), err.Error())
				assert.Nil(t, readDelta)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.delta, readDelta)
			}
		})
	}
}
