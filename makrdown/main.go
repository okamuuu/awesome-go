package main

import (
	"fmt"
	"github.com/russross/blackfriday"
	"io/ioutil"
)

func main() {
	md, err := ioutil.ReadFile("2015-08-25.something.md")
	if err != nil {

	}
	output := blackfriday.MarkdownCommon([]byte(md))
	fmt.Println(string(output))
}
