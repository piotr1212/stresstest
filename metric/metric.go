package metric

import (
	"fmt"
	"github.com/mlambrichs/stresstest/alphabet"
	"gopkg.in/fatih/pool.v2"
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

// send a metric. Figures.
func (m Metric) Send(p pool.Pool, timeout int) error {

	start := rand.Intn(59)
	time.Sleep(time.Duration(start) * time.Second)
	for {
		// get timestamp
		tsp := strconv.FormatInt(time.Now().Unix(), 10)
		// get random value
		value := strconv.Itoa(rand.Intn(100))
		metric := strings.Join([]string{strings.Replace(hostname, ".", "_", -1) + "." + m.String(), value, tsp}, " ")

		// Acquire a connection from the pool.
		connection, err := p.Get()
		if err != nil {
			return err
		}

		//		log.Println("Sending", metric)
		fmt.Fprintf(connection, metric+"\n")
		// Release the connection back to the pool.
		connection.Close()
		time.Sleep(time.Duration(timeout) * time.Second)
	}
}
