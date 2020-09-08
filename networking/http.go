package networking

import (
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func Hash(s string) string {
	h := fnv.New32a()
	h.Write([]byte(s))
	return fmt.Sprintf("%d", h.Sum32())
}

func DoGet(route string) string {
	if os.Getenv("EXAMPLE") != "" {
		b, _ := ioutil.ReadFile("example3.html")
		return string(b)
	}
	agent := "https://github.com/andrewarrow/cli-browser"
	urlString := fmt.Sprintf("%s", route)

	hash := Hash(urlString)

	os.Mkdir(".cli-browser-files", 0755)
	b, err := ioutil.ReadFile(".cli-browser-files/" + hash)
	if err == nil {
		return string(b)
	}
	//fmt.Println(url)
	request, _ := http.NewRequest("GET", urlString, nil)
	request.Header.Set("User-Agent", agent)
	//torTransport := &http.Transport{Proxy: http.ProxyURL(torProxyUrl)}
	//client := &http.Client{Transport: torTransport, Timeout: time.Second * 50}
	client := &http.Client{Timeout: time.Second * 50}

	resp, err := client.Do(request)
	if err == nil {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			fmt.Println("resp.StatusCode", resp.StatusCode, len(body))
			if resp.StatusCode == 200 {
				ioutil.WriteFile(".cli-browser-files/"+hash, body, 0775)
				return string(body)
			} else {
				//fmt.Println(id, string(body))
			}
		} else {
			fmt.Println(len(string(body)), err)
		}
	} else {
		fmt.Println(err)
	}
	return ""
}
