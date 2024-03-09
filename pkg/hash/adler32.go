package hash

import (
	"reflect"
)

// This constant is used as a modulo value when calculating Adler32 hash.
const ADLER_CONSTANT = 65521

type adler32 struct {
	s1     uint
	s2     uint
	window []byte
}

func NewAdler32(chunkSize int) *adler32 {
	return &adler32{
		s1:     1,
		s2:     0,
		window: make([]byte, 0, chunkSize),
	}
}

func (ad *adler32) WindowLength() int {
	return len(ad.window)
}

// Writes chunk byte slice to the window byte slice.
func (ad *adler32) Write(chunk []byte) {
	ad.window = chunk
}

func (ad *adler32) Hash() uint {
	for _, val := range ad.window {
		ad.s1 += uint(val)
		ad.s2 += ad.s1
	}
	ad.s1 = ad.s1 % ADLER_CONSTANT
	ad.s2 = ad.s2 % ADLER_CONSTANT
	return ad.s2<<16 + ad.s1
}

// Appends a single byte c to the window byte slice. Calculates and returns the updated Adler hash.
func (ad *adler32) RollIn(c byte) uint {
	ad.window = append(ad.window, c)
	ad.s1 = (ad.s1 + uint(c)) % ADLER_CONSTANT
	ad.s2 = (ad.s2 + ad.s1) % ADLER_CONSTANT
	return ad.s2<<16 + ad.s1
}

// Removes the first item from the window byte slice. Calculates and returns the updated Adler hash and the removed byte.
func (ad *adler32) RollOut() (uint, byte) {
	removed := ad.window[0]
	//Adding Adler constant so that ad.s1 does not overflow during subtraction since it is an unsigned subtraction
	ad.s1 = (ad.s1 + ADLER_CONSTANT - uint(removed)) % ADLER_CONSTANT
	//Adding Adler constant as many times to make ad.s2 subtraction result positive in order to prevent overflow
	ad.s2 = (ad.s2 + (1+uint(len(ad.window))*uint(removed)/ADLER_CONSTANT)*ADLER_CONSTANT - (uint(len(ad.window)) * uint(removed)) - 1) % ADLER_CONSTANT
	ad.window = ad.window[1:]
	return ad.s2<<16 + ad.s1, removed
}

func (ad *adler32) GetWindowLiterals() []byte {
	return ad.window
}

func (ad *adler32) Reset() {
	p := reflect.ValueOf(ad).Elem()
	p.Set(reflect.Zero(p.Type()))
	ad.s1 = 1
}
