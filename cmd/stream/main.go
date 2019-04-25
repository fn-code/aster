package main

import (
	"log"

	ff "github.com/fn-code/aster"
)

func main() {

	log.Printf("Audio server running on port: %d", 9092)
	conn, err := ff.Listen(":9092")
	if err != nil {
		log.Printf("error open connection: %v\n", err)
	}
	defer conn.Close()
	conn.ReadDataStream()

}
