/*
COMS 319 HW01 Client
A simple websocket client written in Go.

Initializes a websocket connection with the server provided with the -server
flag, or localhost:9000 by default.

Author: 	Stephen Onnen
Email: 		onnen@iastate.edu
*/
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"sync"

	"github.com/fatih/color"
	"github.com/gorilla/websocket"
)

type message struct {
	Username string `json:"username"`
	Text     string `json:"text"`
}

var name string       // Name given by user
var s *bufio.Scanner  // Scanner used to read user input
var wg sync.WaitGroup // Waitgroup to force our goroutines to finish

// Main starts an instance of the chat client and connects to the server passed
// in with the --server flag, or 127.0.0.1:8080 by default.
func main() {
	// Create a waitgroup so main doesn't exit prior to threads finishing
	wg.Add(2)
	//Provide the address and port of the server as flag so it isn't hard-coded.
	server := flag.String("server", "localhost:9000", "Server network address")

	flag.Parse()
	u := url.URL{Scheme: "ws", Host: *server, Path: "/"}

	s = bufio.NewScanner(os.Stdin)
	color.Yellow("Enter your Name: ")
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

	go handleIncoming(sock) // Handle incoming messages concurrently
	go handleOutgoing(sock) // Handle outgoing messages concurrently

	wg.Wait() // Wait for handling of incoming/outgoing messages to complete
}

// handleIncoming handles incoming messages on the websocket connection.
// Each message is unmarshalled into a message struct and then printed to the
// console.
func handleIncoming(sock *websocket.Conn) {
	defer wg.Done()
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
func handleOutgoing(sock *websocket.Conn) {
	defer wg.Done()
	for {
		var msg message
		msg.Username = name
		s.Scan()
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
