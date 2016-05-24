package main

import (
	"flag"
	"fmt"
	"github.com/mlambrichs/stresstest/alphabet"
	"github.com/mlambrichs/stresstest/alphabet/file"
	"gopkg.in/fatih/pool.v2"
	"log"
	"math/big"
	"math/rand"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var nrOfMetrics int
var path string
var poolCapacity int
var port int
var server string
var timeout int
var idCounter int32
var waitGrp sync.WaitGroup

func init() {
	const (
		defaultNrOfMetrics  = 20000
		defaultPoolCapacity = 20
		defaultPort         = 2003
		defaultServer       = "127.0.0.1"
		defaultTimeout      = 60

		usageNrOfMetrics  = "the number of metrics being sent"
		usagePath         = "the path of your inputfile containing alphabet"
		usagePoolCapacity = "the size of the pool of connections"
		usagePort         = "the port to connect to"
		usageServer       = "the server to connect to"
		usageTimeout      = "timeout between sent messages of same metric"
	)
	// define flag for nr of metrics
	flag.IntVar(&nrOfMetrics, "nr_of_metrics", defaultNrOfMetrics, usageNrOfMetrics)
	flag.IntVar(&nrOfMetrics, "n", defaultNrOfMetrics, usageNrOfMetrics+" (shorthand)")

	// define flag for path
	flag.StringVar(&path, "path", "", usagePath)
	flag.StringVar(&path, "p", "", usagePath+" (shorthand)")

	// define flag for pool capacity
	flag.IntVar(&poolCapacity, "pool_capacity", defaultPoolCapacity, usagePoolCapacity)
	flag.IntVar(&poolCapacity, "pc", defaultPoolCapacity, usagePoolCapacity+" (shorthand)")

	// define flag for port
	flag.IntVar(&port, "port", defaultPort, usagePort)
	flag.IntVar(&port, "po", defaultPort, usagePort+" (shorthand)")

	// define flag for server
	flag.StringVar(&server, "server", defaultServer, usageServer)
	flag.StringVar(&server, "s", defaultServer, usageServer+"(shorthand)")

	flag.IntVar(&timeout, "timeout", defaultTimeout, usageTimeout)
	flag.IntVar(&timeout, "t", defaultTimeout, usageTimeout)
}

// create new metric based on a 'depth' random selections out of our alphabet
func createMetricParts(alphabet alphabet.Alphabet, depth int64) []string {
	metric := make([]string, depth)
	for i := 0; int64(i) < depth; i++ {
		metric[i], _ = alphabet.Get(rand.Intn(alphabet.Len()))
		//		metric[i] = alphabet[rand.Intn(len(alphabet))]
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
func calculateDepth(nrOfItems int) int64 {
	log.Printf("calculateDepth with nr of items %d", nrOfItems)
	d := big.NewInt(int64(nrOfItems))
	k := int64(1)
	n := int64(nrOfItems)

	result := new(big.Int)

	for {
		if result.Binomial(n, k); result.Cmp(d) > 0 {
			break
		}
		log.Printf("result = %s", result.Binomial(n, k).String())
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
		time.Sleep(time.Duration(timeout) * time.Second)
	}
}

func main() {

	flag.Parse()

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

	// get alphabet
	var alphabet alphabet.Alphabet

	if path != "" {
		alphabet = file.NewBuffer(path)
		depth := calculateDepth(alphabet.Len())
		log.Printf("Starting with %d metrics and depth %d", nrOfMetrics, depth)
		for i := 0; i < nrOfMetrics; i++ {
			// create new metric
			var newMetric string
			for {
				newMetric = joinParts(createMetricParts(alphabet, depth))
				if !isNotNewMetric(newMetric, metrics) {
					break
				}
			}
			// ...and save new metric into map
			metrics[newMetric]++
		}

	} else {
		depth := calculateDepth(nrOfMetrics)

		for i := 0; i < nrOfMetrics; i++ {
			// create new metric
			var newMetric string
			for {
				newMetric = joinParts(createMetricParts(alphabet, depth))
				if !isNotNewMetric(newMetric, metrics) {
					break
				}
			}
			// ...and save new metric into map
			metrics[newMetric]++
		}
	}
	log.Println("metrics", metrics)

	// start sending metrics to graphite
	for key, _ := range metrics {
		log.Println("Starting client with metric", key)
		// add 1 to waitGroup
		waitGrp.Add(1)
		go func(k string) {
			hostname, _ := os.Hostname()
			sendMetric(strings.Join([]string{hostname, key}, "."), p)
			waitGrp.Done()
		}(key)
	}

	waitGrp.Wait()
	// Close pool. This means closing all connedctions in pool.
	p.Close()
}
