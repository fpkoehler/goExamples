// A program utilizing an interesting idiom in Go to run
// a state machine.  States are represented by functions
// which returns a function that represents the next state.
//
// Concepts inspired by Rob Pike's talk on Lexical scanning
//   http://cuddle.googlecode.com/hg/talk/lex.html
//   https://www.youtube.com/watch?v=HxaD_trXwRE
//
// The example here is modeling a garage door and its remote.
// The door can be opened or closed, pressing the button moves
// the door to the next position.  Pressing the button while
// it is moving, stops the door.  Pressing the button again
// continues the door in the same direction (though I guess
// normally behavior is for the door to change direction).
//
// Command line interface where each carriage return represents
// a button press.  As the door is moving we see integers representing
// its progress towards the next position.
//
// Concurrent routines (goroutines) are used.  One is for the state
// machine which transitions when it receives input on its channel
// (representing a button press).  The other goroutine is used for 
// the door movement with a channel used to stop the movement.

package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

type garage struct {
	progress   int // value from 0 to 9, used when door is opening or closing
	stopMoving chan int // used to stop the movement of the garage door
	Button     chan int // the channel is used to indicate button press
}

// stateFn represents the state of the garage
// as a function that returns the next state.
type stateFn func(*garage) stateFn

// Represents the movement of the garage door.
// The starting position, either open or closed, is 0 and then
// when it counts up to 10 it has moved to the other position.
// At that point, a 1 is written to the Button channel so that
// the state machine will advance.
// Note also the movement will stop (i.e. this routine will terminate) 
// if it receives a an item on the stopMoving channel.
func (g *garage) moving() {
	for g.progress < 10 {
		select {
		case <-g.stopMoving:  // we received data on the stopMoving channel
			fmt.Println("stopped progress", g.progress)
			return
		default:
			time.Sleep(time.Second)
			fmt.Println(g.progress)
			g.progress++
		}
	}

	// reached the next position.  Send data to the Button channel
	// which the run() method is waiting on.
	g.Button <- 1
}

// These are the various states.  These states
// execute the action that happens when the button
// is pressed.  For example, when in the open state
// the open function gets called when the Button is
// pressed.

func open(g *garage) stateFn {
	g.progress = 0
	go g.moving() // start the goroutine
	fmt.Println("garage closing")
	return closing // Next state.
}

func closing(g *garage) stateFn {
	if g.progress < 10 {
		g.stopMoving <- 1
		return closingStopped
	}
	fmt.Println("garage closed")
	return closed // Next state.
}

func closingStopped(g *garage) stateFn {
	go g.moving()
	fmt.Println("garage closing")
	return closing // Next state.
}

func closed(g *garage) stateFn {
	g.progress = 0
	go g.moving()
	fmt.Println("garage opening")
	return opening // Next state.
}

func opening(g *garage) stateFn {
	if g.progress < 10 {
		g.stopMoving <- 1
		return openingStopped
	}
	fmt.Println("garage open")
	return open // Next state.
}

func openingStopped(g *garage) stateFn {
	go g.moving()
	fmt.Println("garage opening")
	return opening // Next state.
}

// run waits for the Button channel to receive input
// then transitions to the next state
func (g *garage) run() {
	fmt.Println("garage open")
	for state := open; state != nil; {

 		// wait for button press
		<-g.Button

		// state is a function pointer whose function
		// will return another function.
		state = state(g)
	}
}

func main() {
	g := &garage{
		Button:     make(chan int),
		stopMoving: make(chan int),
	}

	//    reader := bufio.NewReader(os.Stdin)
	//    fmt.Print("Enter text: ")
	//    text, _ := reader.ReadString('\n')
	//    fmt.Println(text)

	// Start goroutine which is similar to
	// a thread/coroutine.  Its loop will block
	// on waiting for data for its g.Button channel.
	go g.run()

	// a carriage return simulates pressing the garage
	// door remote.
	reader := bufio.NewReader(os.Stdin)
	for {
		reader.ReadString('\n')
		g.Button <- 1
	}

	close(g.Button) // no more button presses
	close(g.stopMoving)
}
