package metric

import (
	"fmt"
	"testing"
)

var a Nato

func TestNew(*testing.T) {
	m := New(a, 3)
	fmt.Println(m.String())
}
