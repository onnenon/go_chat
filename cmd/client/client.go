package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/gorilla/websocket"
)

type message struct {
	Username string `json:"username"`
	Text     string `json:"text"`
}

var name string

// Main starts an instance of the chat client and connects to the server passed
// in with the --server flag, or 127.0.0.1:8080 by default.
func main() {
	//Provide the address and port of the server as a flag so it isn't hard-coded.
	server := flag.String("server", "localhost:8080", "Server network address")

	flag.Parse()

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your Name: ")
	name, _ := reader.ReadString('\n')

	fmt.Printf("\nWelcome %s\n", name)

	fmt.Print("Lets connect to your server.\n\n")

	u := url.URL{Scheme: "ws", Host: *server, Path: "/"}

	log.Printf("Connecting to server @ %s", u.String())

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("Connection error, exiting:", err)
	}
	msg := message{Username: name, Text: "Hello!"}

	json := json.Marshal(msg)
	for {
		conn.WriteJSON(json)
	}

	// conn.Close()
}
