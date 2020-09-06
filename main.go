package main

import (
	"fmt"
	"os"
	"time"

	"cli-browser/browser"
)

func main() {
	fmt.Println("cli-browser")

	if len(os.Args) > 1 {
		fmt.Println("hi")
	}

	b := browser.NewBrowser()
	go b.Start()

	for {
		time.Sleep(time.Second)
	}
}
