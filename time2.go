package main

import (
	"fmt"
	"time"
)

func main() {

	t1 := time.Now()
	fmt.Println(t1)

	tStr := t1.Format(time.UnixDate)
	fmt.Println("formatting time:", tStr)

	t2, err := time.Parse(time.UnixDate, tStr)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("reading back:", tStr	, "we get:", t2)	
	}

	tStr = t1.String()
	t3, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", tStr)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("reading back:", tStr	, "we get:", t3)
	}

//	t1 = t1.Truncate(time.Duration(time.Second))
//	fmt.Println(t1)
}
