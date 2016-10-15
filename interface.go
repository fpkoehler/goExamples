package main

import (
	"fmt"
)

//------------------------------------
// Behavior interface and
// tricks() function that utilizes it
//------------------------------------

type Behavior interface {
	Sit() string
	Sleep() string
}

func tricks(b Behavior) {
	fmt.Println(b.Sit())
	fmt.Println(b.Sleep())
}

//------------------------------------
// Dog
//------------------------------------
type Dog struct {
	Name string
}

func (d Dog) Sit() string {
	return d.Name+" sits"
}

func (d Dog) Sleep() string {
	return d.Name+" sleeps"
}

//------------------------------------
// Cat
//------------------------------------
type Cat struct {
	nickName string
}

func (c Cat) Sit() string {
	return "sitting "+c.nickName
}

func (c Cat) Sleep() string {
	return "sleeping "+c.nickName
}

//------------------------------------
// Function Literals and Closures
//------------------------------------

func multX(a int, b int, s string) int {
	fmt.Println(s)
	return a*b
}

func mult(a, b int) int {
	return a*b
}

//func makeMult(fn func(int, int string)) int {
//	return func(a)
//}

//------------------------------------
// main
//------------------------------------
func main() {
	fmt.Println("hello")

	d := Dog{"sonic"}
	c := Cat{"stubs"}

	tricks(d)
	tricks(c)
}