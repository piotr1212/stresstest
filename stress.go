package main

import (
	"io"
	"log"
	"math/big"
	"math/rand"
	"net"
	"pool"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var nrOfMetrics int = 10000
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
	pooledResources = 50 // The number of connections in our pool
)

type pooledConnection struct {
	ID   int32
	conn net.Conn
}

// createConnection is the factory method that will be called by
// the pool when a new connection is needed.
func createConnection() (io.Closer, error) {
	connection := net.JoinHostPort(server, strconv.Itoa(port))
	id := atomic.AddInt32(&idCounter, 1)

	conn, err := net.Dial("tcp", connection)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Create: New Connection", id)
	return &pooledConnection{id, conn}, nil
}

// Close implements the io.Closer interface so our tcp connection
// can be managed by the pool. Close performs any resource
// release management.
func (connection *pooledConnection) Close() error {
	log.Println("Close: Connection", connection.ID)
	return nil
}

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
func sendMetric(name string, p *pool.Pool) {

	// Acquire a connection from the pool.
	conn, err := p.Acquire()
	if err != nil {
		log.Println(err)
		return
	}

	start := randomStart()
	time.Sleep(time.Duration(start) * time.Second)
	for {
		// get timestamp
		tsp := strconv.FormatInt(time.Now().Unix(), 10)
		// get random value
		value := strconv.Itoa(rand.Intn(100))
		metric := strings.Join([]string{name, value, tsp}, " ")

		log.Println("Sending", metric)
		//		fmt.Fprintf(conn.conn, metric+"\n")
		// Release the connection back to the pool.
		p.Release(conn)
		time.Sleep(1 * time.Minute)
	}
}

func main() {
	numcpu := runtime.NumCPU()
	runtime.GOMAXPROCS(numcpu)

	// create a pool to manage all connections
	p, err := pool.New(createConnection, pooledResources)
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
}
