// Go_chat Server
// A simple websocket server written in Go.
//
// Creates a persistent webserver using the http library. Listens for incomming
// http connections on the port provided with the -addr flag, or 9000 by default.
//
// Author:		Stephen Onnen
// Email: 		stephen.onnen@gmail.com
package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/fatih/color"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Struct that all incoming messages will be unmarshalled into.
type message struct {
	Username string `json:"username"`
	Text     string `json:"text"`
	ID       uuid.UUID
}

type ConnReadWriter interface {
	WriteJSON(v interface{}) error
	ReadJSON(v interface{}) error
	Close() error
}

type Upgrader interface {
	Upgrade(w http.ResponseWriter, r *http.Request, Header http.Header) (*websocket.Conn, error)
}

type ChatServer struct {
	upgrader      Upgrader
	activeClients map[ConnReadWriter]uuid.UUID
	chatRoom      chan message
}

func (cs *ChatServer) handleConn(w http.ResponseWriter, r *http.Request) {

	sock, err := cs.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading connection to websocket: %v", err)
	}

	defer sock.Close()
	// Generate a UUID for the client and add it to activeClients
	cs.activeClients[sock] = uuid.New()

	for {
		var msg message
		err := sock.ReadJSON(&msg)
		if err != nil {
			log.Printf("Closing connection with ID: %v", cs.activeClients[sock])
			delete(cs.activeClients, sock)
			break
		}
		msg.ID = cs.activeClients[sock]
		cs.chatRoom <- msg
	}
}

func main() {
	// Upgrader instance to upgrade all http connections to a websocket.
	// var upgrader = websocket.Upgrader{}

	var chatServer ChatServer

	chatServer.upgrader = websocket.Upgrader{}
	chatServer.activeClients = make(map[ConnReadWriter]uuid.UUID)
	chatServer.chatRoom = make(chan message)

	//Provide the port of the server as a flag so it isn't hard-coded.
	addr := flag.String("addr", ":9000", "Server's network address")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/", chatServer.handleConn) // We only need one uri, make it root.

	go handleMsg(chatServer.activeClients, chatServer.chatRoom) // Handle incoming messages concurrently.

	log.Printf("Starting se	for {rver on %s", *addr)
	err := http.ListenAndServe(*addr, mux)
	if err != nil {
		log.Fatal("Error starting server, exiting.", err)
	}
}

// handleConn handles incomming http connections by adding the connection to a
// global map of current connections and upgrading connection to a websocket.
// Connections are identified individually by a generated UUID.
func handleConn(w http.ResponseWriter, r *http.Request, u Upgrader, chatRoom chan message, activeClients string) {
	// Upgrade incomming http connections to websocket connections
}

// handleMsg listens to the chatRoom channel, when a message is read it is sent
// to each client currently in the activeClients map. If a message fails to send
// to an activeClient, the client is removed from the activeClient map.
func handleMsg(activeClients map[ConnReadWriter]uuid.UUID, chatRoom chan message) {
	for {
		msg := <-chatRoom // Get messages that are sent to the chatRoom channel

		// Log each message to the server's Stdout
		t := time.Now().Format(time.ANSIC)
		color.Green("%s >> %s: %s\n", t, msg.Username, msg.Text)

		for client, UUID := range activeClients {
			// Check the UUID to prevent sending messages to their origin.
			if msg.ID != UUID {
				err := client.WriteJSON(msg)
				if err != nil {
					log.Printf("Error sending message to client: %v", err)
					client.Close()
					delete(activeClients, client)
				}
			}
		}
	}
}
