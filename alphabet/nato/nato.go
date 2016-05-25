package nato

import (
	"errors"
	"math/rand"
)

type Nato [26]string

var alphabet Nato = [26]string{
	"alfa", "bravo", "charlie", "delta", "echo", "foxtrot", "golf",
	"hotel", "india", "julliett", "kilo", "lima", "mike",
	"november", "oscar", "papa", "quebec", "romeo", "sierra",
	"tango", "uniform", "victor", "whiskey", "x-ray", "yankee", "zulu"}

func NewNato() *Nato {
	a := &alphabet
	return a
}

func (n Nato) Get(i ...int) (metric string, err error) {

	switch len(i) {
	case 0:
		metric = alphabet[rand.Intn(len(alphabet))]
	case 1:
		idx := i[0]
		if idx < 0 || idx > n.Len() {
			err = errors.New("index out of scope")
		} else {
			metric = alphabet[idx]
		}
	default:
		err = errors.New("Get should only have 0 or 1 parameter.")
	}
	return metric, err
}

func (n Nato) Len() int {
	return len(alphabet)
}
