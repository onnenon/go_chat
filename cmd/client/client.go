package main

import (
	"flag"
	"log"
)

func main() {
	//Provide the address and port of the server as a flag so it isn't hard-coded.
	addr := flag.String("addr", ":8080", "Server's network address")

	flag.Parse()

	log.Printf("Connecting to server @ %s", *addr)

}
