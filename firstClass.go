/* https://dave.cheney.net/2016/11/13/do-not-fear-first-class-functions */

package main

import (
	"fmt"
	"math"
)

type Calculator struct {
        acc float64
}

func Add(n float64) func(float64) float64 {
        return func(acc float64) float64 {
                return acc + n
        }
}

func (c *Calculator) Do(op func(float64) float64) float64 {
        c.acc = op(c.acc)
        return c.acc
}

func main() {
    var c Calculator
    fmt.Println(c.Do(Add(2)))
    fmt.Println(c.Do(math.Sqrt))   // 1.41421356237
    fmt.Println(c.Do(math.Cos))    // 0.99969539804
}

