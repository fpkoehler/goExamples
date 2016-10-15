package main

import (
	"fmt"
)

func insertFloat(l []float64, f float64) []float64 {
	var i int

	if len(l) == 0 {
		return []float64{f}
	}

	if len(l) == 1 {
		if f < l[0] {
			return []float64{f,l[0]}
		} else {
			return append(l, f)
		}
	}

	low := 0
	high := len(l) - 1
	for low+1 != high {
		i = ((high - low)/2) + low
		fmt.Println(low, high, i, f, l[i])

		if f < l[i] {
			high = i
		} else {
			low = i
		}
	}

	fmt.Println(low, high, i, f, l[i])


	switch {
	case f < l[low]:
		fmt.Println(f, l[low])
		l = append([]float64{f}, l...)
	case f > l[high]:
		fmt.Println(l[high], f)
		l = append(l,f)
	default:
		fmt.Println(l[low], f, l[high])
		i = low+1
		l = append(l, 0)      // add a blank element to end of slice
		copy(l[i+1:], l[i:])  // shift from i+1 to the right
		l[i] = f              // insert f
	}

	return l
}

func main() {
	f := make([]float64, 100, 100)
	for i:=0; i<100; i++ {
		f[i] = float64(i)/100.0
	}

	fmt.Println(insertFloat(f, 0.355))
	fmt.Println(insertFloat(f, 10.0))
	fmt.Println(insertFloat(f, -1.0))
	fmt.Println(insertFloat(f, .989))
}