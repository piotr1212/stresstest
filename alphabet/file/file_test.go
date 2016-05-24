package file

import (
	"fmt"
	"testing"
)

var b = NewBuffer("/tmp/test.txt")

func TestGetWithNoIndex(t *testing.T) {
	value, err := b.Get(0)
	if err != nil {
		t.Error(`TestWithGetNoIndex() should return string`)
	}
	fmt.Println("Found value:", value)
}

func TestGetWithOneIndex(t *testing.T) {
	_, err := b.Get()
	if err != nil {
		t.Error(`TestGetWithOneIndex() should return string`)
	}
}

func TestGetWithNegatveIndex(t *testing.T) {
	_, err := b.Get(-1)
	if err == nil {
		t.Error(`TestGetWithNegativeIndex() should return error`)
	}
}

func TestGetWithTooLargeIndex(t *testing.T) {
	_, err := b.Get(27)
	if err == nil {
		t.Error(`TestGetWithTooLargeIndex() should return error`)
	}
}

func TestGetTwoIndexes(t *testing.T) {
	_, err := b.Get(0, 1)
	if err == nil {
		t.Error(`TestGetWithTwoIndexes should fail.`)
	}
}

func TestLen(t *testing.T) {
	if b.Len() != 26 {
		t.Error(`TestLen() should return exact length of alphabet-array.`)
	}
}
