package files

import (
	"fmt"
	"hash/fnv"
	"io/ioutil"
)

const DIR = ".cli-browser-files"

func Hash(s string) string {
	h := fnv.New32a()
	h.Write([]byte(s))
	return fmt.Sprintf("%d", h.Sum32())
}

func List(url string) []string {
	files, _ := ioutil.ReadDir(DIR + "/" + Hash(url))
	list := []string{}
	for _, file := range files {
		list = append(list, file.Name())
	}
	return list
}

func Push(url, payload string) {
	items := List(url)
	ioutil.WriteFile(DIR+"/"+Hash(url)+"/"+
		fmt.Sprintf("%05d.txt", len(items)), []byte(payload), 0755)
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
