package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClear(t *testing.T) {
	type testStruct struct {
		a int
		b string
	}
	ts := testStruct{1, "test"}
	Clear(&ts)
	zs := testStruct{}
	assert.Equal(t, ts, zs)
}
