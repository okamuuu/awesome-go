package main

import (
	"log"
	"os"
)

func main() {
	f, err := os.OpenFile("/tmp/test.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		log.Fatal("error opening file :", err.Error())
	}

	log.SetOutput(f)

	log.Println("hello log")
}
