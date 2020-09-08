package main

import (
	"cli-browser/browser"
	"fmt"
	"os"
)

func main() {
	fmt.Println("cli-browser")

	arg1 := ""
	if len(os.Args) > 1 {
		arg1 = os.Args[1]
	}

	b := browser.NewBrowser()
	b.Start(arg1)
}
