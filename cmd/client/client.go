// Go_chat Client
// A simple websocket client written in Go.
//
// Initializes a websocket connection with the server provided with the -server
// flag, or localhost:9000 by default.
//
// Author: 		Stephen Onnen
// Email: 		stephen.onnen@gmail.com
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

type ConnReader interface {
	ReadJSON(v interface{}) error
}

type ConnWriter interface {
	WriteJSON(v interface{}) error
	Close() error
}

type Scanner interface {
	Scan() bool
	Text() string
}

type message struct {
	Username string `json:"username"`
	Text     string `json:"text"`
}

// Main starts an instance of the chat client and connects to the server passed
// in with the --server flag, or 127.0.0.1:8080 by default.
func main() {
	var name string      // Name given by user
	var s *bufio.Scanner // Scanner used to read user input

	//Provide the address and port of the server as flag so it isn't hard-coded.
	server := flag.String("server", "localhost:9000", "Server network address")
	path := flag.String("path", "/", "Server Path")
	flag.Parse()
	u := url.URL{Scheme: "ws", Host: *server, Path: *path}

	color.Yellow("Enter your Name: ")
	s = bufio.NewScanner(os.Stdin)
	s.Scan()
	name = s.Text()

	color.Green("\nWelcome %s!!\n\n", name)
	color.Green("Connecting to server @ %s\n", *server)
	color.Yellow("Go ahead and send a message, or type quit() to exit.\n")

	sock, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("Connection error, exiting:", err)
	}
	defer sock.Close()

	msg := message{Username: name, Text: "has joined the chat."}
	sock.WriteJSON(msg)

	go handleIncoming(sock)       // Handle incoming messages concurrently.
	handleOutgoing(sock, s, name) // Handle outgoing messages concurrently.
}

// handleIncoming handles incoming messages on the websocket connection.
// Each message is unmarshalled into a message struct and then printed to the
// console.
func handleIncoming(sock ConnReader) {
	for {
		var msg message
		err := sock.ReadJSON(&msg)
		if err != nil {
			color.White("Server closed. Exiting...")
			os.Exit(0)
		}
		color.Red("%s: %s\n", msg.Username, msg.Text)
	}
}

// handleOutgoing scans Stdin and sends each scanned line to the server as a
// message struct marshalled into JSON.
//
// With terminals that do not support escape sequences, the user inputed text
// will not be properly cleared from the screen, and will display twice.
// this should only affect users of Windows.
func handleOutgoing(sock ConnWriter, s Scanner, name string) {
	var msg message
	msg.Username = name

	for {
		if s.Scan() {
			fmt.Printf("\033[A")
			msg.Text = s.Text()

			if msg.Text == "quit()" {
				fmt.Println("Goodbye!")
				sock.WriteJSON(message{Username: name, Text: "has disconnected."})
				sock.Close()
				os.Exit(0)
			}
			color.Cyan("%s: %s\n", msg.Username, msg.Text)

			err := sock.WriteJSON(msg)
			if err != nil {
				log.Fatal("Error sending message, exiting")
			}
		}
	}
}
