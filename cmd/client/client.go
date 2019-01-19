package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/fatih/color"
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
	u := url.URL{Scheme: "ws", Host: *server, Path: "/"}

	s := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter your Name: ")
	s.Scan()
	name := s.Text()

	fmt.Printf("\nWelcome %s\n", name)
	fmt.Print("Lets connect to your server.\n\n")
	log.Printf("Connecting to server @ %s", u.String())

	sock, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("Connection error, exiting:", err)
	}

	msg := message{Username: name, Text: "has joined the chat."}
	sock.WriteJSON(msg)

	defer sock.Close()

	go func() {
		for {
			var msg message

			err := sock.ReadJSON(&msg)
			if err != nil {
				log.Println("read:", err)
				return
			}
			color.Red("%s: %s\n", msg.Username, msg.Text)
		}
	}()

	for {
		var msg message
		msg.Username = name
		s.Scan()
		msg.Text = s.Text()
		err := sock.WriteJSON(msg)
		if err != nil {
			log.Fatal("Error sending message, exiting")
		}
		color.Cyan("%s: %s\n", msg.Username, msg.Text)
	}

}
