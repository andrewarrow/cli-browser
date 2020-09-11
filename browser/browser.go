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
var tagId = 0
var curTag *ATag
var prevTag = ""
var prevPrevTag = ""
var prevTagAtts = map[string]string{}
var prevPrevTagAtts = map[string]string{}
var tagHolder = ATag{}
var tagMap = map[int]*ATag{}
var childCount = 0
var findTag = ""
var findAtt = ""

type Browser struct {
	Homepage string
}

type ATag struct {
	Name     string
	Id       int
	Atts     map[string]string
	Parent   *ATag
	Children []*ATag
	Text     string
}

func NewBrowser() *Browser {
	b := Browser{}
	b.Homepage = "https://www.amazon.com/"
	return &b
}

func (b *Browser) Start(arg1, arg2 string) {
	// s?i=aps&k=hands with hearts straight borders&ref=nb_sb_noss&url=search-alias=
	if arg1 == "ls" {
		for i, url := range files.History() {
			fmt.Printf("%2d. %s\n", i+1, url)
		}
		return
	}
	s := networking.DoGet(arg1)
	z := html.NewTokenizer(strings.NewReader(s))
	for {
		if handleTag(z) == false {
			break
		}
	}
	requestedTag := tagHolder.Children
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
			requestedTag = tagMap[tagId].Children
		} else {
			findTag = tagType
			findAtt = tokens[1]
		}
	}
	for _, c := range requestedTag {
		countAllChildren(0, c)
		fmt.Printf("%d. %10s (%d)\n", c.Id, c.Name, childCount-1)
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
		if curTag != nil {
			s := strings.TrimSpace(string(z.Text()))
			if s != "" {
				//combo := prevPrevTag + "," + prevTag
				//if combo == "a,span" {
				//curTag.Text = curTag.Text + "|" + prevPrevTagAtts["href"] + " " + combo + "|" + s
				//}
			}
		}
	case html.EndTagToken:
		//tn, _ := z.TagName()
		//tns := string(tn)
		curTag = curTag.Parent
	case html.StartTagToken:
		tn, _ := z.TagName()
		tns := string(tn)
		prevPrevTag = prevTag
		prevTag = tns
		atts := getAtts(z)
		prevPrevTagAtts = prevTagAtts
		prevTagAtts = atts

		tagId++
		td := ATag{}
		td.Id = tagId
		td.Name = tns
		td.Atts = atts
		if curTag != nil {
			td.Parent = curTag
		} else {
			td.Parent = &tagHolder
		}
		tagMap[td.Id] = &td
		td.Parent.Children = append(td.Parent.Children, &td)
		curTag = &td

	}
	return true
}
func findTagInChildren(start int, td *ATag) {
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
func countAllChildren(start int, td *ATag) {
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
func walkDivs(spaces string, td *ATag) {
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
