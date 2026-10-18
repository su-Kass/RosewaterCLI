package main

// libra- oh nevermind this is fine
import (
	"fmt"
	"strconv"
)

// the data of the context menu, used for structuring
type Context struct {
	Header       string
	Content	     []string
	Options      []string
	Length       int16
}

// woa........ theme stuff is stored here for some reason............
const (
	reset     = "\033[0m"
	bold      = "\033[1m"
	dim       = "\033[2m"
	italic    = "\033[3m"
	underline = "\033[4m"
	red       = "\033[31m"
	green     = "\033[32m"
	yellow    = "\033[33m"
	blue      = "\033[34m"
	cyan      = "\033[36m"
)

// prompt visuals and display
func Prompt() (response string) {
	fmt.Print(bold+green+"/ᐠ˵- ⩊ -˵マ "+reset)
	fmt.Print(dim+"Choice:"+reset)
	fmt.Print(bold+red)
	var input string
	fmt.Scanln(&input)
	fmt.Print(reset) // yeah im sure this is the proper way to do this
	return input
}

// silly thing that just structures a context
func MakeContext(header string, content []string, options []string) Context {
	return Context{Header: header, Content: content, Options: options} // i should use this more, idk why i dont, pretty cool..
}

// omg CLI slop!!!
// makes the actual visuals of the context, fun enough this makes this whole thing themable lol...
func readContext(readTarget Context) {
	fmt.Println(bold+underline+readTarget.Header+reset)
	fmt.Println(bold+"--------------------------------------------"+reset)
	for _, contentPoint := range readTarget.Content {
		fmt.Println(bold+yellow+"I: "+reset+green+contentPoint+reset)
	}
	fmt.Println(bold+"..OPTIONS.."+reset)
	for option, contentOption := range readTarget.Options {
		fmt.Println(red+strconv.Itoa(option)+": "+reset+yellow+contentOption+reset)
	}
}

//blehh gets the length of the context
func contextLines(cont Context) int {
	return 4 + len(cont.Content) + len(cont.Options)
}

// #ansis slop
func clearContext(cont Context) {
	fmt.Printf("\033[%dA\r\033[J", contextLines(cont))
}

// yay, actually format the context into something usable
func ContextPrompt(cont Context) string {
	readContext(cont)
	answer := Prompt()
	clearContext(cont)
	return answer
}
