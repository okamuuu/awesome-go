package main

import (
	"fmt"
	"runtime"
)

func say(s string) {
	fmt.Printf("say: %s\n", s)
	for i := 0; i < 5; i++ {
		// 下記2行を入れ替えて動作させると分かりやすい
		runtime.Gosched()
		fmt.Println(s)
	}
}

func main() {
	go say("world") // new goroutine
	say("hello")    // current goroutine
}
