package files

import (
	"fmt"
	"hash/fnv"
	"io/ioutil"
)

func Hash(s string) string {
	h := fnv.New32a()
	h.Write([]byte(s))
	return fmt.Sprintf("%d", h.Sum32())
}

func List(url string) []string {
	files, _ := ioutil.ReadDir(".cli-browser-files/" + Hash(url))
	list := []string{}
	for _, file := range files {
		list = append(list, file.Name())
	}
	return list
}
