package fileio

import (
	"bufio"
	"encoding/gob"
	"os"

	"github.com/Psykepro/rdiff/pkg/differ"
)

type FileHandler struct {
	chunkSize int
}

func NewFileHandler(chunkSize int) FileHandler {
	return FileHandler{
		chunkSize: chunkSize,
	}
}

func (f FileHandler) ChunkSize() int {
	return f.chunkSize
}

func (f FileHandler) Open(filepath string) (*bufio.Reader, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, NewReadFileError(err)
	}

	fileInfo, _ := file.Stat()
	size := fileInfo.Size()

	chunks := size / int64(f.chunkSize)
	if chunks < 2 {
		return nil, ErrFileSize
	}
	return bufio.NewReader(file), nil
}

func (f FileHandler) WriteSignatures(signatures map[uint]int, output string) error {
	file, err := os.Create(output)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	err = encoder.Encode(signatures)
	if err != nil {
		return err
	}

	return nil
}

func (f FileHandler) ReadSignatures(filePath string) (map[uint]int, error) {
	signatures := make(map[uint]int)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&signatures)
	if err != nil {
		return nil, err
	}

	return signatures, nil
}

func (f FileHandler) WriteDelta(delta map[int]differ.Delta, output string) error {
	file, err := os.Create(output)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	err = encoder.Encode(delta)
	if err != nil {
		return err
	}

	return nil
}

func (f FileHandler) ReadDelta(filePath string) (map[int]differ.Delta, error) {
	delta := make(map[int]differ.Delta)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&delta)
	if err != nil {
		return nil, ErrDecodeDelta
	}

	return delta, nil
}
