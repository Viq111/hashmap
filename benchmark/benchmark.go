package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"sync/atomic"
	"time"
)

var options struct {
	mode string

	getRatio      float64
	nbThreads     int
	rotateSeconds int
	valuesPerKey  int
}

var (
	exit chan bool // Populate when exiting program
)

type Map interface {
	Insert(key int64, value Value)
	Get(key int64) []Value
	Len() int
}

func fatalf(s string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, s, args...)
	os.Exit(2)
}

func init() {
	flag.StringVar(&options.mode, "mode", "hashmap", "What to test. Possible values are: std, hashmap, hashmap_alloc")
	flag.Float64Var(&options.getRatio, "get_ratio", 0.7, "Ratio between set & get. 0.5 means as much of one than the other; 0 means only set; 1 means only get")
	flag.IntVar(&options.nbThreads, "threads", 0, "Number of goroutines to spawn concurrently. default is all")
	flag.IntVar(&options.rotateSeconds, "rotate", 60, "Number of seconds before rotating a bucket")
	flag.IntVar(&options.valuesPerKey, "values", 10, "Number of values per key")
}

func main() {
	exit = make(chan bool, 1)
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, os.Interrupt)
	go func() {
		<-sigterm
		exit <- true
	}()

	flag.Parse()
	if options.nbThreads == 0 {
		options.nbThreads = runtime.NumCPU()
	}
	availableMode := map[string]struct{}{
		"std":           struct{}{},
		"hashmap":       struct{}{},
		"hashmap_alloc": struct{}{},
	}
	if _, exist := availableMode[options.mode]; !exist {
		fatalf("%v is not a known mode, you need to choose a correct one", options.mode)
	}

	rotateTime := time.Duration(options.rotateSeconds) * time.Second
	fmt.Printf("Executing on %v threads, will rotate buckets every %s...\n", options.nbThreads, rotateTime)
	var createFunc (func(int) Map)
	if options.mode == "std" {
		createFunc = NewMapStd
	} else if options.mode == "hashmap" {
		createFunc = NewMapHashmap
	} else if options.mode == "hashmap_alloc" {
		fatalf("not implemented")
	} else {
		fatalf("not implemented")
	}

	// Stats
	statTick, tickers := multiTicker(rotateTime, options.nbThreads)
	var nbGlobalSet int64
	var nbGlobalGet int64
	go func() {
		for range statTick {
			s := atomic.SwapInt64(&nbGlobalSet, 0)
			g := atomic.SwapInt64(&nbGlobalGet, 0)

			fmt.Printf("Did %s INSERTs and %s GETs in %s\n", HumanNb(s), HumanNb(g), rotateTime)
		}
	}()

	// Work

	initialSize := 10
	batchSize := 10000
	every := batchSize / options.valuesPerKey
	for i := 0; i < options.nbThreads; i++ {
		go func(worker int) {
			m := createFunc(initialSize)

			batchNb := 0
			batchOffset := 0
			nbGets := 0

			for {
				select {
				case <-tickers[worker]:
					// Add to stats
					atomic.AddInt64(&nbGlobalSet, int64(batchNb*batchSize+batchOffset))
					atomic.AddInt64(&nbGlobalGet, int64(nbGets))

					// Reset
					m = createFunc(int(float64(m.Len()) * 1.1))
					batchNb = 0
					batchOffset = 0
					nbGets = 0
				default:
				}
				r := rand.Float64()
				if r <= options.getRatio { // Do a get
					nbGets++
					g := m.Get(0)
					runtime.KeepAlive(g) // Simulate using the var
				} else {
					v := Value{
						step:     int32(batchNb*batchSize + batchOffset),
						lastSeen: int64(batchNb*batchSize + batchOffset),
					}
					m.Insert(int64(batchNb*batchSize+(batchOffset%every)), v)
					batchOffset++
					if batchOffset >= batchSize {
						batchNb++
						batchOffset = 0
					}
				}
			}
		}(i)
	}

	RunDebugHTTP(6060)
	<-exit
}
