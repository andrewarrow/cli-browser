package browser

import (
	"cli-browser/files"
	"cli-browser/networking"
	"fmt"
	"os"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

var forms = 0
var tables = 0
var trOn = false
var divId = 0
var curDiv *TopDiv
var prevTag = ""
var prevPrevTag = ""
var prevTagAtts = map[string]string{}
var prevPrevTagAtts = map[string]string{}
var divHolder = TopDiv{}
var divMap = map[int]*TopDiv{}
var childCount = 0
var findTag = ""
var findAtt = ""

type Browser struct {
	Homepage string
}

type TopDiv struct {
	Name     string
	Id       int
	Atts     map[string]string
	Parent   *TopDiv
	Children []*TopDiv
	Text     string
}

func NewBrowser() *Browser {
	b := Browser{}
	b.Homepage = "https://www.amazon.com/"
	return &b
}

func (b *Browser) Start(arg1, arg2 string) {
	// s?i=aps&k=hands with hearts straight borders&ref=nb_sb_noss&url=search-alias=
	s := networking.DoGet(arg1)
	z := html.NewTokenizer(strings.NewReader(s))
	for {
		if handleTag(z) == false {
			break
		}
	}
	requestedDiv := divHolder.Children
	if arg2 != "" {
		if arg2 == "ls" {
			for i, file := range files.OrderOps(arg1) {
				fmt.Printf("%2d. %s\n", i+1, file)
			}
			return
		}
		if strings.HasPrefix(arg2, "push ") {
			tokens := strings.Split(arg2, " ")
			files.Push(arg1, tokens[1])
			return
		}
		if strings.HasPrefix(arg2, "pop") {
			files.Pop(arg1)
			return
		}
	}
	items := files.OrderOps(arg1)
	for _, item := range items {
		tokens := strings.Split(item, ",")
		tagType := tokens[0]
		if tagType == "div" {
			tagId, _ := strconv.Atoi(tokens[1])
			requestedDiv = divMap[tagId].Children
		} else {
			findTag = tagType
			findAtt = tokens[1]
		}
	}
	for _, c := range requestedDiv {
		countAllChildren(0, c)
		fmt.Printf("%d. DIV (%d)\n", c.Id, childCount-1)
		if findTag != "" {
			findTagInChildren(0, c)
		}
	}
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
		if curDiv != nil {
			s := strings.TrimSpace(string(z.Text()))
			if s != "" {
				//combo := prevPrevTag + "," + prevTag
				//if combo == "a,span" {
				//curDiv.Text = curDiv.Text + "|" + prevPrevTagAtts["href"] + " " + combo + "|" + s
				//}
			}
		}
	case html.EndTagToken:
		tn, _ := z.TagName()
		tns := string(tn)
		if tns == "tr" {
			trOn = false
		} else if tns == "div" || tns == "a" {
			curDiv = curDiv.Parent
		}
	case html.StartTagToken:
		tn, _ := z.TagName()
		tns := string(tn)
		prevPrevTag = prevTag
		prevTag = tns
		atts := getAtts(z)
		prevPrevTagAtts = prevTagAtts
		prevTagAtts = atts
		if tns == "form" {
			forms++
			if atts["action"] != "" {
				//fmt.Printf("%d. FORM %s %s\n", forms, atts["method"], atts["action"])
			}
		} else if tns == "input" {
			if atts["name"] != "" {
				//fmt.Printf("          %s\n", atts["name"])
			}
		} else if tns == "table" {
			tables++
			//fmt.Printf("%d. TABLE\n", tables)
		} else if tns == "tr" {
			trOn = true
			//fmt.Printf("\n            %s\n", "tr")
		} else if tns == "div" || tns == "a" {
			divId++
			td := TopDiv{}
			td.Id = divId
			td.Name = tns
			td.Atts = atts
			if curDiv != nil {
				td.Parent = curDiv
			} else {
				td.Parent = &divHolder
			}
			divMap[td.Id] = &td
			td.Parent.Children = append(td.Parent.Children, &td)
			curDiv = &td
		}

	}
	return true
}
func findTagInChildren(start int, td *TopDiv) {
	m := map[int]bool{}
	for {
		if len(td.Children) == 0 {
			return
		}
		if m[td.Id] {
			return
		}
		m[td.Id] = true
		if td.Name == findTag {
			fmt.Println(td.Atts[findAtt])
		}
		for _, c := range td.Children {
			findTagInChildren(start+1, c)
		}
	}
}
func countAllChildren(start int, td *TopDiv) {
	if start == 0 {
		childCount = 0
	}
	m := map[int]bool{}
	for {
		if len(td.Children) == 0 {
			if os.Getenv("VERBOSE") != "" {
				fmt.Println(td.Text)
			}
			childCount++
			return
		}
		if m[td.Id] {
			childCount++
			return
		}
		m[td.Id] = true
		for _, c := range td.Children {
			countAllChildren(start+1, c)
		}
	}
}
func walkDivs(spaces string, td *TopDiv) {
	m := map[int]bool{}
	for {
		if len(td.Children) == 0 {
			fmt.Printf("%s%d\n", spaces, td.Id)
			return
		}
		if m[td.Id] {
			return
		}
		fmt.Printf("%s%d\n", spaces, td.Id)
		m[td.Id] = true
		for _, c := range td.Children {
			walkDivs(spaces+"  ", c)
		}
	}
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
