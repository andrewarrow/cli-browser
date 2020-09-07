package networking

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func DoGet(route, params string) string {
	if os.Getenv("EXAMPLE") != "" {
		b, _ := ioutil.ReadFile("example.html")
		return string(b)
	}
	agent := "https://github.com/andrewarrow/cli-browser"
	urlString := fmt.Sprintf("%s", route)
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
