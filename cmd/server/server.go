package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/fatih/color"
	"github.com/google/uuid"
	"github.com/gorilla/websocket" // Reccomended by Golang over it's STD Library
)

type message struct {
	Username string `json:"username"`
	Text     string `json:"text"`
	ID       uuid.UUID
}

var upgrader = websocket.Upgrader{}                     // Upgrader instance to upgrade all http connections to a websocket.
var activeClients = make(map[*websocket.Conn]uuid.UUID) // Map to store currently active client connections.
var chatRoom = make(chan message)                       //Channel to send all messages to.

func main() {
	//Provide the port of the server as a flag so it isn't hard-coded.
	addr := flag.String("addr", ":8080", "Server's network address")
	flag.Parse()

	http.HandleFunc("/", handleConn) // Since we only need one endpoint, make it root.

	go handleMsg() // Create thread to handle messages

	log.Printf("Starting server on %s", *addr)
	err := http.ListenAndServe(*addr, nil)

	if err != nil {
		log.Fatal("Error starting server, exiting.", err)
	}
}

// handleConn handles incomming http connections by adding the connection to a
// global map of current connections and upgrading the connection to a websocket.
// Connections are identified individually by a UUID
func handleConn(w http.ResponseWriter, r *http.Request) {
	sock, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Printf("Error upgrading connection to websocket: %v", err)
	}

	defer sock.Close()

	activeClients[sock] = uuid.New() // Generate a UUID for the client and add it to activeClients

	for {
		var msg message

		err := sock.ReadJSON(&msg)
		if err != nil {
			log.Printf("Closing connection with ID: %v", activeClients[sock])
			delete(activeClients, sock)
			break
		}

		msg.ID = activeClients[sock]
		chatRoom <- msg
	}
}

// handleMsg listens to the chatRoom channel, when a message is read it is sent
// to each client currently in the activeClients map. If a message fails to send
// to an activeClient, the client is removed from the activeClient map.
func handleMsg() {
	for {
		msg := <-chatRoom // Get any messages that are sent to the chatRoom channel
		color.Green("%s >> %s: %s\n", time.Now().Format(time.ANSIC), msg.Username, msg.Text)
		for client, UUID := range activeClients {
			if msg.ID != UUID { // Check the UUID to prevent sending messages to their origin.
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
