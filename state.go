package main

import (
	"fmt"
)

// lexer holds the state of the scanner.
type lexer struct {
    name  string    // used only for error reports.
    input string    // the string being scanned.
    start int       // start position of this item.
    pos   int       // current position in the input.
    width int       // width of last rune read from input.
//    items chan item // channel of scanned items.
}

// stateFn represents the state of the scanner
// as a function that returns the next state.
type stateFn func(*lexer) stateFn

func lexText(l *lexer) stateFn {
	fmt.Println("lexText")
    return lexA    // Next state.
}

func lexA(l *lexer) stateFn {
	fmt.Println("lexA")
    return lexB    // Next state.
}

func lexB(l *lexer) stateFn {
	fmt.Println("lexB")
    return nil    // Next state.
}

// run lexes the input by executing state functions until
// the state is nil.
func (l *lexer) run() {
    for state := lexText; state != nil; {
        state = state(l)
    }
//    close(l.items) // No more tokens will be delivered.
}

func main() {
	l := &lexer{
        name:  "test",
        input: "abcd",
    }

    l.run()


}
