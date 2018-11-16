package net_test

import (
	"encoding/binary"
	"fmt"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/maprost/should"
	"github.com/maprost/testbox/must"
	"github.com/maprost/toolbox/net"
)

var clients = make(map[*websocket.Conn]bool) // connected clients
var broadcast = make(chan Message)           // broadcast channel

// Configure the upgrader
var upgrader = websocket.Upgrader{}

// Define our message object
type Message struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

func TestWebSockets(t *testing.T) {
	// Create a simple file server
	fs := http.FileServer(http.Dir("web"))
	http.Handle("/", fs)

	// Configure websocket route
	http.HandleFunc("/ws", handleConnections)

	// Start listening for incoming chat messages
	go handleMessages()
	go runClient()

	// Start the server on localhost port 8000 and log any errors
	log.Println("http server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Make sure we close the connection when the function returns
	defer ws.Close()

	// Register our new client
	clients[ws] = true

	for {
		// Send the newly received message to the broadcast channel
		broadcast <- convertMessage(ws)
	}
}

func handleMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-broadcast
		// Send it out to every client that is currently connected
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func runClient() {
	time.Sleep(2 * time.Second)

	log.Println("Connect client:")
	c, _, err := websocket.DefaultDialer.Dial("ws://localhost:8000/ws", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	for {
		msg := convertMessage(c)
		if msg.Username != "echo" && strings.HasPrefix(msg.Message, ">>echo:") {
			msg = Message{
				Username: "echo",
				Email:    "support@chat.de",
				Message:  strings.TrimPrefix(msg.Message, ">>echo:"),
			}
			c.WriteJSON(&msg)
		}

		if msg.Message == ">>time:" {
			msg = Message{
				Username: "timer",
				Email:    "support@chat.de",
				Message:  time.Now().String(),
			}
			c.WriteJSON(&msg)
		}

		if msg.Message == ">>help:" {
			msg = Message{
				Username: "help",
				Email:    "support@chat.de",
				Message:  "Command are: '>>help:', '>>time:', '>>echo:'",
			}
			c.WriteJSON(&msg)
		}
	}
}

func convertMessage(ws *websocket.Conn) Message {
	var msg Message
	// Read in a new message as JSON and map it to a Message object
	err := ws.ReadJSON(&msg)
	if err != nil {
		log.Printf("error: %v", err)
		delete(clients, ws)
	}

	return msg
}

func TestWebSocketChannel(t *testing.T) {
	openWS := false
	closedWithError := false
	// create echo web socket server
	server := NewGetServer(map[string]net.HandlerFunc{
		"ws": func(con *net.Connection) {
			ws, err := con.WebSocketChannel()
			if err != nil {
				con.SendResponse(nil, err)
				return
			}

			defer ws.Close()

			oldMsg := ""
			for {
				openWS = true
				msgData, open, err := ws.Read()
				if err != nil {
					fmt.Println("Error!", err)
					openWS = false
					closedWithError = true
					con.SendResponse(nil, err)
					return
				}
				if !open {
					fmt.Println("Close!")
					openWS = false
					con.SendResponse(nil, nil)
					return
				}

				msg := string(msgData)
				if msg != oldMsg {
					oldMsg = msg
					ws.WriteText(msg)
				}
			}
		},
	})

	// create client
	client := net.Client{}
	ws, err := client.WebSocketChannel(server.URL + "/ws")
	must.BeNoError(t, err)

	t.Run("simple echo", func(t *testing.T) {
		// send message
		err = ws.WriteText("Hello")
		must.BeNoError(t, err)

		// get message
		data, open, err := ws.Read()
		must.BeNoError(t, err)
		must.BeTrue(t, open)
		should.BeEqual(t, string(data), "Hello")
	})

	t.Run("double echo", func(t *testing.T) {
		// send message
		err = ws.WriteText("World")
		must.BeNoError(t, err)

		err = ws.WriteText("Peace")
		must.BeNoError(t, err)

		// get message
		data, open, err := ws.Read()
		must.BeNoError(t, err)
		must.BeTrue(t, open)
		should.BeEqual(t, string(data), "World")

		data, open, err = ws.Read()
		must.BeNoError(t, err)
		must.BeTrue(t, open)
		should.BeEqual(t, string(data), "Peace")
	})

	fmt.Println("------------------------------------")

	// send message
	err = ws.Close()
	must.BeNoError(t, err)

	// close web socket
	//err = ws.Close()
	//must.BeNoError(t, err)

	time.Sleep(time.Millisecond)
	must.BeFalse(t, openWS)
	must.BeFalse(t, closedWithError)
}

func intToByteSlice(id uint16) []byte {
	var buf [2]byte
	binary.BigEndian.PutUint16(buf[:], id)
	return buf[:]
}
