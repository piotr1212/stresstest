package metric

import (
	"fmt"
	"github.com/mlambrichs/stresstest/alphabet"
	_ "gopkg.in/fatih/pool.v2"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type Metric []string

var hostname string

func init() {
	hostname, _ = os.Hostname()
}

func New(a alphabet.Alphabet, depth int) Metric {
	var metric Metric = make([]string, depth)
	for i := 0; i < depth; i++ {
		metric[i], _ = a.Get()
	}
	return metric
}

// make a metric name out of array of strings
func (m Metric) String() string {
	return strings.Join(m, ".")
}

func (m Metric) Send(timeout int, c chan string) {
	start := rand.Intn(59)
	time.Sleep(time.Duration(start) * time.Second)
	for {
		// get random value
		value := strconv.Itoa(rand.Intn(100))
		// get timestamp
		tsp := strconv.FormatInt(time.Now().Unix(), 10)
		c <- fmt.Sprintf("%s.%s %s %s", strings.Replace(hostname, ".", "_", -1), m.String(), value, tsp)
		time.Sleep(time.Duration(timeout) * time.Second)
	}
}
