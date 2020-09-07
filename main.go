package main

import (
	"fmt"
	"os"

	"cli-browser/browser"
)

func main() {
	fmt.Println("cli-browser")

	if len(os.Args) > 1 {
		b := browser.NewBrowser()
		b.Start(os.Args[1])
		return
	}

	b := browser.NewBrowser()
	b.Start("")
}
