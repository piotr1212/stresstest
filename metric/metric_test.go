package metric

import (
	"fmt"
	"github.com/mlambrichs/stresstest/alphabet/nato"
	"testing"
)

var a nato.Nato

func TestNew(*testing.T) {
	m := New(a, 3)
	fmt.Println(m.String())
}
