package browser

import (
	"cli-browser/networking"
	"fmt"
)

type Browser struct {
	Homepage string
}

func NewBrowser() *Browser {
	b := Browser{}
	b.Homepage = "https://www.amazon.com/"
	return &b
}

func (b *Browser) Start() {
	s := networking.DoGet(b.Homepage, "")
	fmt.Println(len(s))
}
