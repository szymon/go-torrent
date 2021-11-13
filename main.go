package main

import (
	"fmt"
	"github.com/szymon/go-torrent/bencode"
	"log"
)

func main() {

	data, err := bencode.Marshal(int64(64))
	if err != nil {
		log.Fatal("error...")
	}

	fmt.Printf("%s", data)
}
