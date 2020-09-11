package files

import (
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"os"
	"strings"
)

const DIR = ".cli-browser-files"

func Hash(s string) string {
	h := fnv.New32a()
	h.Write([]byte(s))
	return fmt.Sprintf("%d", h.Sum32())
}

func History() []string {
	d, _ := ioutil.ReadFile(DIR + "/history.txt")
	items := strings.Split(string(d), "\n")
	return items[0 : len(items)-1]
}
func AddToHistory(u string) {
	items := History()
	f, _ := os.OpenFile(DIR+"/history.txt", os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	f.WriteString(u + "\n")
	for _, item := range items {
		f.WriteString(item + "\n")
	}
}

func List(url string) []string {
	files, _ := ioutil.ReadDir(DIR + "/" + Hash(url))
	list := []string{}
	for _, file := range files {
		if file.Name() == "index.html" {
			continue
		}
		list = append(list, file.Name())
	}
	return list
}

func Push(url, payload string) {
	items := List(url)
	ioutil.WriteFile(DIR+"/"+Hash(url)+"/"+
		fmt.Sprintf("%05d.txt", len(items)+1), []byte(payload), 0755)
}
func OrderOps(url string) []string {
	files, _ := ioutil.ReadDir(DIR + "/" + Hash(url))
	list := []string{}
	for _, file := range files {
		if file.Name() == "index.html" {
			continue
		}
		d, _ := ioutil.ReadFile(DIR + "/" + Hash(url) + "/" + file.Name())
		list = append(list, string(d))
	}
	return list
}
func Pop(url string) string {
	items := List(url)
	d, _ := ioutil.ReadFile(DIR + "/" + Hash(url) + "/" + items[len(items)-1])
	os.Remove(DIR + "/" + Hash(url) + "/" + items[len(items)-1])
	return string(d)
}
