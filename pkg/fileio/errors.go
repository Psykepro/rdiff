package fileio

import (
	"errors"
	"fmt"
)

var (
	ErrReadFile         = fmt.Errorf("error in reading file")
	ErrFileSize         = fmt.Errorf("file is not of appropriate size to generate enough chunks")
	ErrCreateFile       = fmt.Errorf("error in creating file")
	ErrEncodeSignatures = fmt.Errorf("error in encoding signatures")
	ErrDecodeSignatures = fmt.Errorf("error in decoding signatures")
	ErrEncodeDelta      = fmt.Errorf("error in encoding delta")
	ErrDecodeDelta      = fmt.Errorf("error in decoding delta")
)

func NewReadFileError(err error) error {
	return fmt.Errorf("%w. Error Details: %v", ErrReadFile, err)
}

// fmt.Errorf("open " + nonExistentPath + ": no such file or directory")
func NewOpenFileError(fileName string) error {
	return fmt.Errorf("open %v: no such file or directory", errors.New(fileName))
}
