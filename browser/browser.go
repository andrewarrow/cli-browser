package browser

import (
	"cli-browser/networking"
	"fmt"
	"strings"

	"golang.org/x/net/html"
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
	z := html.NewTokenizer(strings.NewReader(s))
	for {
		tt := z.Next()

		switch tt {
		case html.ErrorToken:
			fmt.Println("ERR", z.Err())
			return
		case html.TextToken:
			fmt.Println(string(z.Text()))
		case html.StartTagToken, html.EndTagToken:
			tn, _ := z.TagName()
			fmt.Println(string(tn))
		}
	}
}
