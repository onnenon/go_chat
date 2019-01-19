package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type message struct {
	Username string `json:"username"`
	Text     string `json:"text"`
}

var upgrader = websocket.Upgrader{}                // Upgrader instance to upgrade all http connections to a websocket.
var activeClients = make(map[*websocket.Conn]bool) // Map to store currently active client connections.
var chatRoom = make(chan message)                  //Channel to send all messages to.

func main() {
	//Provide the address and port of the server as a flag so it isn't hard-coded.
	addr := flag.String("addr", ":8080", "Server's network address")
	flag.Parse()

	go handleMsg()

	http.HandleFunc("/", handleConn)
	log.Printf("Starting server on %s", *addr)
	err := http.ListenAndServe(*addr, nil)

	if err != nil {
		log.Fatal("Error starting server, exiting.", err)
	}
}

// handleConn handles incomming http connections by adding the connection to a
// global map of current connections and upgrading the connection to a websocket.
func handleConn(w http.ResponseWriter, r *http.Request) {
	sock, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Printf("Error upgrading connection to websocket: %v", err)
	}

	defer sock.Close()

	activeClients[sock] = true

	for {
		var msg message

		err := sock.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error receiving message: %v", err)
			delete(activeClients, sock)
			break
		}
		chatRoom <- msg
	}
}

// handleMsg listens to the chatRoom channel, when a message is read it is sent
// to each client currently in the activeClients map. If a message fails to send
// to an activeClient, the client is removed from the activeClient map.
func handleMsg() {
	for {
		msg := <-chatRoom // Get any messages that are sent to the chatRoom channel
		for client := range activeClients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("Error sending message to client: %v", err)
				client.Close()
				delete(activeClients, client)
			}
		}
	}
}
