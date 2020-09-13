package browser

import (
	"cli-browser/files"
	"cli-browser/networking"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

var body bool
var script bool
var style bool
var trOn = false
var tagId = 0
var curTag *ATag
var foundTag *ATag
var foundText []string
var prevTag = ""
var prevPrevTag = ""
var prevTagAtts = map[string]string{}
var prevPrevTagAtts = map[string]string{}
var tagHolder = ATag{}
var tagMap = map[int]*ATag{}
var childCount = 0
var findTag = ""
var findAtt = ""
var mostKidRecord = 0
var mostKidAt *ATag
var allWithKids []*ATag

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
		//tokens := strings.Split(item, ",")
		//tagType := tokens[0]
		//if tagType == "div" {
		tagId, _ := strconv.Atoi(item)
		requestedTag = tagMap[tagId].Children
		//} else {
		//	findTag = tagType
		//	findAtt = tokens[1]
		//}
	}
	mostKids(&tagHolder)
	sort.SliceStable(allWithKids, func(i, j int) bool {
		return len(allWithKids[i].Children) > len(allWithKids[j].Children)
	})
	kidIndex := 0

	for {
		if allWithKids[kidIndex].Children[0].Name != "option" {
			break
		}
		kidIndex++
	}
	requestedTag = allWithKids[kidIndex].Children
	for _, c := range requestedTag {
		countAllChildren(0, "Sponsored", c)
		//out of 5 stars
		fmt.Printf("%d. %10s (%d) (%s)\n", c.Id, c.Name, childCount-1, "")
		if findTag != "" {
			findTagInChildren(0, c)
		}
		foundTag = nil
		foundText = []string{}
		findTextInChildren(c)
		fmt.Println(foundText)
		if c.Id == 480 {
			walkDivs("", c)
		}
	}
	//walkDivs("", &tagHolder)
}

func handleTag(z *html.Tokenizer) bool {

	tt := z.Next()

	switch tt {
	case html.ErrorToken:
		fmt.Println("ERR", z.Err())
		return false
	case html.TextToken:
		if body == false {
			return true
		}
		if script == true {
			return true
		}
		if style == true {
			return true
		}
		t := strings.TrimSpace(string(z.Text()))
		if t == "" {
			return true
		}
		tagId++
		at := ATag{}
		at.Id = tagId
		at.Name = "text"
		at.Text = t
		at.Parent = curTag
		tagMap[at.Id] = &at
		at.Parent.Children = append(at.Parent.Children, &at)
		curTag = &at
	case html.EndTagToken:
		if body == false {
			return true
		}
		if script == true {
			script = false
		}
		if style == true {
			style = false
		}
		//tn, _ := z.TagName()
		//tns := string(tn)
		//fmt.Printf("ending %s -> %s,%s\n", tns, curTag.Name, curTag.Parent.Name)
		if curTag.Name == "text" {
			curTag = curTag.Parent.Parent
		} else {
			curTag = curTag.Parent
		}
	case html.StartTagToken:
		tn, _ := z.TagName()
		tns := string(tn)

		if tns == "body" {
			body = true
		}

		if tns == "script" {
			script = true
		}
		if tns == "style" {
			style = true
		}

		if body == false {
			return true
		}

		prevPrevTag = prevTag
		prevTag = tns
		atts := getAtts(z)
		prevPrevTagAtts = prevTagAtts
		prevTagAtts = atts

		tagId++
		at := ATag{}
		at.Id = tagId
		at.Name = tns
		at.Atts = atts
		if curTag != nil {
			at.Parent = curTag
		} else {
			at.Parent = &tagHolder
		}
		tagMap[at.Id] = &at
		at.Parent.Children = append(at.Parent.Children, &at)
		curTag = &at

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
func findTextInChildren(td *ATag) {
	m := map[int]bool{}
	for {
		if len(td.Children) == 0 {
			foundText = append(foundText, td.Text)
			return
		}
		if m[td.Id] {
			return
		}
		m[td.Id] = true
		for _, c := range td.Children {
			findTextInChildren(c)
		}
	}
}
func countAllChildren(start int, searchText string, td *ATag) {
	if start == 0 {
		childCount = 0
		//foundText = ""
	}
	m := map[int]bool{}
	for {
		if len(td.Children) == 0 {
			if td.Text == searchText {
				//foundText = searchText
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
			countAllChildren(start+1, searchText, c)
		}
	}
}
func mostKids(at *ATag) {
	m := map[int]bool{}
	for {
		if len(at.Children) == 0 {
			return
		}
		if m[at.Id] {
			return
		}
		m[at.Id] = true
		allWithKids = append(allWithKids, at)
		if len(at.Children) > mostKidRecord {
			mostKidRecord = len(at.Children)
			mostKidAt = at
		}
		for _, c := range at.Children {
			mostKids(c)
		}
	}
}
func walkDivs(spaces string, td *ATag) {
	m := map[int]bool{}
	for {
		if len(td.Children) == 0 {
			fmt.Printf("%s%s %s\n", spaces, td.Name, td.Text)
			return
		}
		if m[td.Id] {
			return
		}
		fmt.Printf("%s%s %s\n", spaces, td.Name, td.Text)
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
