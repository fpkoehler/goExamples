package main

import (
	"fmt"
)

var nextIndex int
var nextIndexChan chan int

func setupNextChannel() {
	if nextIndexChan != nil {
		for {
			_, ok := <- nextIndexChan
			if !ok {
				break
			}
		}
	}
	nextIndexChan = make(chan int)

	// feeds the channel
	go func() {
		for i:=0; i<10; i++ {
			nextIndexChan <- i
		}
		close(nextIndexChan)
	}()
}

func main() {
	setupNextChannel()

	for i:=0; i<4; i++ {
		x, ok := <- nextIndexChan
		if !ok {
			break
		}
		fmt.Println(x)
	}

	setupNextChannel()

	for {
		x, ok := <- nextIndexChan
		if !ok {
			break
		}
		fmt.Println(x)
	}

}
