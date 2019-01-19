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
	color.Yellow("Enter your Name: ")
	s.Scan()
	name := s.Text()

	color.Green("\nWelcome %s!!\n\n", name)
	color.Green("Connecting to server @ %s\n", *server)
	color.Yellow("Go ahead and send a message, or type quit() to exit.\n")

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
				color.White("Exiting...")
				log.Fatal()
			}
			color.Red("%s: %s\n", msg.Username, msg.Text)
		}
	}()

	defer sock.Close()

	for {
		var msg message
		msg.Username = name
		s.Scan()
		fmt.Printf("\033[1A")
		msg.Text = s.Text()
		if msg.Text == "quit()" {
			sock.WriteJSON(message{Username: name, Text: "has disconnected."})
			break
		}
		color.Cyan("%s: %s\n", msg.Username, msg.Text)
		err := sock.WriteJSON(msg)
		if err != nil {
			log.Fatal("Error sending message, exiting")
		}
	}
}
