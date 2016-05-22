package metric

import (
	"strings"
)

type Metric []string

func New(a Alphabet, depth int) *Metric {
	var metric Metric = make([]string, depth)
	for i := 0; i < depth; i++ {
		metric[i], _ = a.Get()
	}
	return &metric
}

// make a metric name out of array of strings
func (m Metric) String() string {
	return strings.Join(m, ".")
}
