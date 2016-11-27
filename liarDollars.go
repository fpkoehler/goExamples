package main

import (
	"fmt"
	"time"
	"math/rand"
)

func main() {
	
	const (
		valueLen = 8
		numValues = 5
	)

	var values [valueLen]int
	var stats [10]int

	seed := time.Now().UnixNano()
//	fmt.Println("seed:", seed)
	rand.Seed(seed)

	for j:=0; j<numValues; j++ {
		for i:=0; i<valueLen; i++ {
			values[i] = rand.Intn(10)
			fmt.Print(values[i])
			stats[values[i]]++
		}
		fmt.Println()
	}

	for i:=0; i<10; i++ {
		fmt.Println(i, stats[i])
	}

	randIntN := 1;
	for i:=0; i<valueLen; i++ {
		randIntN *= 10
	}

	for j:=0; j<numValues; j++ {
		x := rand.Intn(randIntN)
		fmt.Println(x)
	}
}