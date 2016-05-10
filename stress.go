package main

import (
	"fmt"
	"gopkg.in/fatih/pool.v2"
	"log"
	"math/big"
	"math/rand"
	"net"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var nrOfMetrics int = 100000
var server string = "127.0.0.1"
var port int = 2003

var idCounter int32
var waitGrp sync.WaitGroup

// create an array to hold some random strings
var alphabet [26]string = [26]string{
	"alfa", "bravo", "charlie", "delta", "echo", "foxtrot", "golf",
	"hotel", "india", "julliett", "kilo", "lima", "mike",
	"november", "oscar", "papa", "quebec", "romeo", "sierra",
	"tango", "uniform", "victor", "whiskey", "x-ray", "yankee", "zulu"}

const (
	poolCapacity = 20 // The number of connections in our pool
)

// create new metric based on a 'depth' random selections out of our alphabet
func createMetricParts(depth int64) []string {
	metric := make([]string, depth)
	for i := 0; int64(i) < depth; i++ {
		metric[i] = alphabet[rand.Intn(len(alphabet))]
	}
	return metric
}

// check if meetric doesn't already exist
func isNotNewMetric(metric string, metrics map[string]int) bool {
	_, ok := metrics[metric]
	return ok
}

// make a metric name out of array of strings
func joinParts(parts []string) string {
	return strings.Join(parts, ".")
}

// use this for letting all sends not starting at the same time
func randomStart() int {
	return rand.Intn(59)
}

// calculate depth
func calculateDepth(depth int) int64 {
	d := big.NewInt(int64(depth))
	k := int64(1)
	n := int64(26)

	result := new(big.Int)

	for {
		if result.Binomial(n, k); result.Cmp(d) > 0 {
			break
		}
		k++
	}
	return k
}

// send a metric. Figures.
func sendMetric(name string, p pool.Pool) {

	start := randomStart()
	time.Sleep(time.Duration(start) * time.Second)
	for {
		// get timestamp
		tsp := strconv.FormatInt(time.Now().Unix(), 10)
		// get random value
		value := strconv.Itoa(rand.Intn(100))
		metric := strings.Join([]string{name, value, tsp}, " ")

		// Acquire a connection from the pool.
		connection, err := p.Get()
		if err != nil {
			log.Println(err)
			return
		}

		//		log.Println("Sending", metric)
		fmt.Fprintf(connection, metric+"\n")
		// Release the connection back to the pool.
		connection.Close()
		time.Sleep(1 * time.Minute)
	}
}

func main() {
	numcpu := runtime.NumCPU()
	runtime.GOMAXPROCS(numcpu)

	// create factory function to be used with channel based pool
	connection := net.JoinHostPort(server, strconv.Itoa(port))
	factory := func() (net.Conn, error) { return net.Dial("tcp", connection) }

	// create a channel based pool to manage all connections
	p, err := pool.NewChannelPool(5, poolCapacity, factory)
	if err != nil {
		log.Println(err)
	}

	// create a map for containing al metrics
	metrics := make(map[string]int)

	// calculate the necessary depth
	depth := calculateDepth(nrOfMetrics)

	for i := 0; i < nrOfMetrics; i++ {
		// create new metric
		var newMetric string
		for {
			newMetric = joinParts(createMetricParts(depth))
			if !isNotNewMetric(newMetric, metrics) {
				break
			}
		}
		// ...and save new metric into map
		metrics[newMetric]++
	}

	// start sending metrics to graphite
	for key, _ := range metrics {
		log.Println("Starting client with metric", key)
		// add 1 to waitGroup
		waitGrp.Add(1)
		go func(k string) {
			sendMetric(key, p)
			waitGrp.Done()
		}(key)
	}

	waitGrp.Wait()
	// Close pool. This means closing all connedctions in pool.
	p.Close()
}
