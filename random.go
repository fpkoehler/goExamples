package main

import (
	"fmt"
	"time"
	"math/rand"
)

func main() {
	
	const (
		valueLen = 16
	)

	var values [valueLen]int

	for i:=0; i<valueLen; i++ {
		values[i] = i+1
	}

	fmt.Println(values)

	seed := time.Now().UnixNano()
	fmt.Println("seed:", seed)
	rand.Seed(seed)

	for len:=valueLen; len>1; len -- {
		r := rand.Intn(len)
		values[r], values[len-1] = values[len-1], values[r]
	}
	fmt.Println(values)
}