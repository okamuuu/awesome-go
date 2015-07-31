package main

/*
	https://jdanger.com/build-a-web-crawler-in-go.html
	go run crawl.go http://google.com
*/

import (
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/jackdanger/collectlinks"
	"net/http"
	"net/url"
	"os"
)

var visited = make(map[string]bool)

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		fmt.Println("Please specify start page")
		os.Exit(1)
	}

	queue := make(chan string)

	go func() { queue <- args[0] }()

	for uri := range queue {
		enqueue(uri, queue)
	}
}

func enqueue(uri string, queue chan string) {
	fmt.Println("fetching", uri)
	visited[uri] = true
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	client := http.Client{Transport: transport}

	resp, err := client.Get(uri)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	links := collectlinks.All(resp.Body)

	for _, link := range links {
		absolute := fixUrl(link, uri)
		if uri != "" && !visited[absolute] {
			go func() { queue <- absolute }()
		}
	}
}

func fixUrl(href, base string) string {
	uri, err := url.Parse(href)
	if err != nil {
		return ""
	}
	baseUrl, err := url.Parse(base)
	if err != nil {
		return ""
	}
	uri = baseUrl.ResolveReference(uri)
	return uri.String()
}
