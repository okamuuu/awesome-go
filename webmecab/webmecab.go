package main

/*
	http://qiita.com/rerofumi/items/66be3c55405e03dbdcf0
	go run webmecabe.go -p 8001
	curl -i -X POST -H "Content-Type: application/json" -d '{"Token": "01234", "Sentence": "生麦生米生卵"}' http://localhost:8001
*/

import (
	"bitbucket.org/rerofumi/mecab"
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

type MecabOutput struct {
	Status int
	Code   string
	Result []string
}

type MecabInput struct {
	Token    string
	Sentence string
}

func apiRequest(w http.ResponseWriter, r *http.Request) {
	list := make([]string, 0)
	ret := MecabOutput{0, "OK", list}
	request := ""

	// JSON return
	defer func() {
		// result
		outjson, err := json.Marshal(ret)
		if err != nil {
			fmt.Println(err)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(outjson))
	}()

	// type check
	if r.Method != "POST" {
		ret.Status = 1
		ret.Code = "Not POST method"
		return
	}

	rb := bufio.NewReader(r.Body)
	for {
		s, err := rb.ReadString('\n')
		request = request + s
		if err == io.EOF {
			break
		}
	}

	// JSON parse
	var dec MecabInput
	b := []byte(request)
	err := json.Unmarshal(b, &dec)
	if err != nil {
		ret.Status = 2
		ret.Code = "JSON parse error."
		return
	}

	// mecab parse
	// XXX: 解析がうまくできていない??
	result, err := mecab.Parse(dec.Sentence)
	if err == nil {
		for _, n := range result {
			ret.Result = append(ret.Result, n)
		}
	}
}

func main() {
	var (
		portNum int
	)
	flag.IntVar(&portNum, "port", 80, "int flag")
	flag.IntVar(&portNum, "p", 80, "int flag")
	flag.Parse()

	var port string
	port = ":" + strconv.Itoa(portNum)
	fmt.Println("listen port =", port)

	// route handler
	http.HandleFunc("/", apiRequest)

	// do serve
	err := http.ListenAndServe(port, nil)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
