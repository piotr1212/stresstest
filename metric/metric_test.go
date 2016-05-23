package metric

import (
	"fmt"
	"github.com/mlambrichs/stresstest/alphabet"
	"testing"
)

var a alphabet.Nato

func TestNew(*testing.T) {
	m := New(a, 3)
	fmt.Println(m.String())
}
