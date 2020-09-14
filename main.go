package main

import (
	"cli-browser/browser"
	"fmt"
	"os"
)

func main() {
	//fmt.Println("cli-browser")

	if len(os.Args) == 1 {
		fmt.Println("first arg should be url")
		return
	}
	arg1 := os.Args[1]
	arg2 := ""
	if len(os.Args) > 2 {
		arg2 = os.Args[2]
	}

	b := browser.NewBrowser()
	b.Start(arg1, arg2)
}
