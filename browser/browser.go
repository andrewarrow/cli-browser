package browser

import (
	"cli-browser/networking"
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

var forms = 0
var tables = 0
var trOn = false
var divId = 0
var curDiv *TopDiv
var divHolder = TopDiv{}

type Browser struct {
	Homepage string
}

type TopDiv struct {
	Id       int
	Name     string
	Parent   *TopDiv
	Children []*TopDiv
}

func NewBrowser() *Browser {
	b := Browser{}
	b.Homepage = "https://www.amazon.com/"
	return &b
}

func (b *Browser) Start(u string) {
	// s?i=aps&k=hands with hearts straight borders&ref=nb_sb_noss&url=search-alias=
	s := networking.DoGet(u, "")
	z := html.NewTokenizer(strings.NewReader(s))
	for {
		if handleTag(z) == false {
			break
		}
	}
	fmt.Printf("%d. DIV (%s)\n", divId, walkDivs(&divHolder))
}

func handleTag(z *html.Tokenizer) bool {

	tt := z.Next()

	switch tt {
	case html.ErrorToken:
		fmt.Println("ERR", z.Err())
		return false
	case html.TextToken:
		if trOn {
			//fmt.Printf("|%s|\n", string(z.Text()))
		}
	case html.EndTagToken:
		tn, _ := z.TagName()
		tns := string(tn)
		if tns == "tr" {
			trOn = false
		} else if tns == "div" {
			curDiv = curDiv.Parent
		}
	case html.StartTagToken:
		tn, _ := z.TagName()
		tns := string(tn)
		if tns == "form" {
			forms++
			atts := getAtts(z)
			if atts["action"] != "" {
				fmt.Printf("%d. FORM %s %s\n", forms, atts["method"], atts["action"])
			}
		} else if tns == "input" {
			atts := getAtts(z)
			if atts["name"] != "" {
				fmt.Printf("          %s\n", atts["name"])
			}
		} else if tns == "table" {
			tables++
			fmt.Printf("%d. TABLE\n", tables)
		} else if tns == "tr" {
			trOn = true
			//fmt.Printf("\n            %s\n", "tr")
		} else if tns == "div" {
			m := getAtts(z)
			divId++
			td := TopDiv{}
			td.Id = divId
			td.Name = m["hi"]
			if curDiv != nil {
				td.Parent = curDiv
			} else {
				td.Parent = &divHolder
			}
			td.Parent.Children = append(td.Parent.Children, &td)
			curDiv = &td
		}

	}
	return true
}
func walkDivs(td *TopDiv) string {
	m := map[int]bool{}
	for {
		fmt.Printf("%+v\n", *td)
		if len(td.Children) == 0 || m[td.Id] {
			return ""
		}
		m[td.Id] = true
		for _, c := range td.Children {
			walkDivs(c)
		}
	}
	return ""
}
func getAtts(z *html.Tokenizer) map[string]string {
	atts := map[string]string{}
	k, v, b := z.TagAttr()
	for {
		atts[string(k)] = string(v)
		if b == false {
			break
		}
		k, v, b = z.TagAttr()
	}
	return atts
}
