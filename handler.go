package main

import (
	"fmt"
	"strconv"
)

type Context struct {
	Header       string
	Content	     []string
	Options      []string
	Length       int16
}

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

func Prompt() (response string) {
	fmt.Print(bold+green+"/ᐠ˵- ⩊ -˵マ "+reset)
	fmt.Print(dim+"Choice:"+reset)
	fmt.Print(bold+red)
	var input string
	fmt.Scanln(&input)
	fmt.Print(reset)
	return input
}

func MakeContext(header string, content []string, options []string) Context {
	return Context{Header: header, Content: content, Options: options}
}


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

func contextLines(c Context) int {
	return 4 + len(c.Content) + len(c.Options) // header, separator, "OPTIONS", prompt line
}

func clearContext(c Context) {
	fmt.Printf("\033[%dA\r\033[J", contextLines(c))
}

func ContextPrompt(c Context) string {
	readContext(c)
	answer := Prompt()
	clearContext(c)
	return answer
}
