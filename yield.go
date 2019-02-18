package main

import (
	"fmt"
)

func ReadCommands() <-chan string {
	chnl := make(chan string)
	go func() {
		for i := 0; i < 100; i++ {
			chnl <- fmt.Sprintf("command-%d", i)
		}
		close(chnl)
	}()

	return chnl
}

func main() {

	reader := ReadCommands()
	for ok := true; ok; {
		cmd := ""
		commands := []string{}
		for i := 0; i < 7; i++ {
			cmd, ok = <-reader
			if !ok {
				/* channel is closed */
				break
			}
			commands = append(commands, cmd)
		}
		fmt.Println(commands)
	}
}
