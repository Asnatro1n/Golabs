package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	conn *websocket.Conn
	send chan []byte
}

type ChatServer struct {
	clients   map[*Client]bool
	broadcast chan []byte
	mu        sync.Mutex
}

func newChatServer() *ChatServer {
	return &ChatServer{
		clients:   make(map[*Client]bool),
		broadcast: make(chan []byte),
	}
}

func (cs *ChatServer) run() {
	for {
		msg := <-cs.broadcast
		cs.mu.Lock()
		for client := range cs.clients {
			select {
			case client.send <- msg:
			default:
				close(client.send)
				delete(cs.clients, client)
			}
		}
		cs.mu.Unlock()
	}
}

func (c *Client) readMessages(cs *ChatServer) {
	defer func() {
		c.conn.Close()
	}()
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		cs.broadcast <- msg
	}
}

func (c *Client) writeMessages() {
	defer func() {
		c.conn.Close()
	}()
	for msg := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
}

func handleConnection(cs *ChatServer, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error while upgrading connection:", err)
		return
	}
	client := &Client{conn: conn, send: make(chan []byte)}
	cs.mu.Lock()
	cs.clients[client] = true
	cs.mu.Unlock()

	go client.readMessages(cs)
	go client.writeMessages()
}

func main() { //открыть index.html в этой же директории
	chatServer := newChatServer()
	go chatServer.run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleConnection(chatServer, w, r)
	})

	fmt.Println("Chat server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
