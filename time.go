package main

import (
	"fmt"
	"time"
	"os"
	"bufio"
	)

func main() {
//    const t1Str = "2016-01-13T17:00:00.000000-05:00"
//    t1, err := time.Parse(time.RFC3339, t1Str)

// See this to understand why using EST below we
// actually get UTC time:
// https://code.google.com/p/go/issues/detail?id=3604
    const t1Str = "Tue Jan 13 17:00:00 EST 2016"
    t1, err := time.Parse(time.UnixDate, t1Str)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
    fmt.Println("t1", t1, t1.Location().String())
    fmt.Println("t1", t1.Local())

    const t2Str = "2016-01-13T10:30:00.000000-08:00"
//    t2, err := time.Parse(time.RFC3339, t2Str)
	t2 := time.Now()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
    fmt.Println("t2", t2)
	fmt.Println("t2-t1", t2.Sub(t1))


    const xStreamTimeStr = "Tue Jan 12 10:30:32 EST 2016"
    tXstream, err := time.Parse(time.UnixDate, xStreamTimeStr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
    fmt.Println("xStream time:", tXstream)

	tNow := time.Now()
    fmt.Println("Current time:", tNow)

	delta := tNow.Sub(tXstream)
	fmt.Println("Now-xStreamTime =", delta)

	fmt.Println("syslog time mapped to real time:", tNow.Add(delta))

	var text string
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter syslog time: ")
	text, _ = reader.ReadString('\n')
	text = text[:len(text)-1] // chop off \n
    t, err := time.Parse(time.RFC3339, text)
	t = t.Local()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("syslog time mapped to real time:", t.Add(delta))

}
