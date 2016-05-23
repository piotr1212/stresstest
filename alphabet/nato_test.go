package alphabet

import (
	"testing"
)

var n Nato

func TestGetWithNoIndex(t *testing.T) {
	_, err := n.Get()
	if err != nil {
		t.Error(`TestWithGetNoIndex() should return string`)
	}
}

func TestGetWithOneIndex(t *testing.T) {
	_, err := n.Get()
	if err != nil {
		t.Error(`TestGetWithOneIndex() should return string`)
	}
}

func TestGetWithNegatveIndex(t *testing.T) {
	_, err := n.Get(-1)
	if err == nil {
		t.Error(`TestGetWithNegativeIndex() should return error`)
	}
}

func TestGetWithTooLargeIndex(t *testing.T) {
	_, err := n.Get(27)
	if err == nil {
		t.Error(`TestGetWithTooLargeIndex() should return error`)
	}
}

func TestGetTwoIndexes(t *testing.T) {
	_, err := n.Get(0, 1)
	if err == nil {
		t.Error(`TestGetWithTwoIndexes should fail.`)
	}
}

func TestLen(t *testing.T) {
	if n.Len() != 26 {
		t.Error(`TestLen() should return exact length of alphabet-array.`)
	}
}
