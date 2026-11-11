package main

import (
	"log"

	"github.com/vitor-franz/book-log/infrastructure"
)

func main() {
	err := infrastructure.Start()
	if err != nil {
		log.Fatal(err)
	}
}
