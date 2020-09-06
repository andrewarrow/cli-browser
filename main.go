package main

import (
	"fmt"
	"os"

	"github.com/andrewarrow/cli-browser/browser"
)

func main() {
	fmt.Println("cli-browser")

	if len(os.Args) > 1 {
		fmt.Println("hi")
	}

	b := browser.NewBrowser()
	fmt.Println(b)
}
