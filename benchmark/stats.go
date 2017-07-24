package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"time"
)

// multiTicker creates multiple out tickers
// It also gives a ticker that tick one second after the others
func multiTicker(tick time.Duration, nb int) (chan time.Time, []chan time.Time) {
	srcTicker := time.NewTicker(tick)
	dst := make([]chan time.Time, nb)
	afterTicker := make(chan time.Time, 1)

	// Create the tickers
	for i := 0; i < nb; i++ {
		dst[i] = make(chan time.Time, 1)
	}

	// Send the ticks
	go func() {
		for t := range srcTicker.C {
			for i := 0; i < nb; i++ {
				select { // Push in a non-blocking way
				case dst[i] <- t:
				default:
				}
			}
			time.Sleep(time.Second)
			select {
			case afterTicker <- t:
			default:
			}
		}
	}()

	return afterTicker, dst
}

func HumanNb(nb int64) string {
	exp := []string{"", "K", "M", "G", "T", "P"}
	n := float64(nb)
	index := 0
	for n >= 1000 {
		n = n / 1000
		index++
	}
	return fmt.Sprintf("%.2f%s", n, exp[index])
}

func RunDebugHTTP(port int) {
	go func() {
		for i := 0; i < 20; i++ {
			address := fmt.Sprintf("localhost:%d", port)
			fmt.Printf("Starting debug server on %s\n", address)
			err := http.ListenAndServe(address, nil)
			if err != nil {
				fmt.Printf("Error running debug server: %s\n", err)
				port++
				continue
			}
			return
		}
	}()
}
