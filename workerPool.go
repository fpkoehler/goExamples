package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const NumItems = 5
const NumWorkers = 2

var doneReading chan bool

func doWork(goRoutine int, v string) string {
	time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
	return fmt.Sprintf("go%d-%s", goRoutine, v)
}

func readResults(r <-chan string) {
	// Not needed, but used just so that I can comment out other time.Sleep
	// lines and not have to reedit imports to remove "time" and "math/rand"
	time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)

	for i := 0; i < NumItems; i++ {
		v := <-r
		fmt.Println(v)
	}
	doneReading <- true
}

func main() {
	doneReading = make(chan bool)

	c1 := make(chan string)
	c2 := make(chan string)

	// Worker pool of Go routines.  Each Go routine will read from
	// c1, then call a doWork() function and put the results into c2
	var wg sync.WaitGroup
	wg.Add(NumWorkers)
	for i := 0; i < NumWorkers; i++ {
		go func(goRoutine int, ci <-chan string, co chan<- string) {
			// Decrement the wait group counter when the goroutine completes.
			defer wg.Done()

			// For the ci channel, each go routine is a reader
			// for the co channel, each go routine is a writer
			for v := range ci {
				result := doWork(goRoutine, v)
				co <- result
			}
			// the for loop ends when there are no more items
			// in the channel and the channel is closed.
		}(i, c1, c2)
	}

	// Reads and prints from c2
	go readResults(c2)

	// Add the items to c1
	for i := 0; i < NumItems; i++ {
		v := fmt.Sprintf("foo%d", i)
		c1 <- v
	}
	close(c1)

	// Waits for the worker pool to finish
	wg.Wait()

	// Wait for readResults() to finish
	<-doneReading

	fmt.Println("done")
}
