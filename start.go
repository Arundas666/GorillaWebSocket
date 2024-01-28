package main

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)


func Start() {
	fmt.Println("-----start-----")
	for {
		select {
		case conn := <-Manager.Register:
			fmt.Printf("New connection: %v", conn.ParticipantId)
			Manager.mutex.Lock()
			Manager.Clients[conn.Usertype] = UserTypeMap{}
			Manager.Clients[conn.Usertype][conn.ParticipantId] = conn
			fmt.Println(Manager.Clients)
			Manager.mutex.Unlock()
		case conn := <-Manager.Unregister:
			fmt.Printf("Connection failed: %v", conn.ParticipantId)
			Manager.mutex.Lock()
			if _, ok := Manager.Clients[conn.Usertype][conn.ParticipantId]; ok {
				close(conn.Send)
				delete(Manager.Clients[conn.Usertype], conn.ParticipantId)
			}
			Manager.mutex.Unlock()
		}
	}
}

type Client struct {
	ParticipantId string
	Socket        *websocket.Conn
	Usertype      UserType
	Send          chan any
}

type UserType string

const (
	Agent UserType = "agent"
	Shop  UserType = "shop"
	User  UserType = "user"
)

type UserTypeMap map[string]*Client

type ClientManager struct {
	Clients    map[UserType]UserTypeMap
	Register   chan *Client
	Unregister chan *Client
	mutex      sync.Mutex
}

var Manager = ClientManager{
	Clients:    make(map[UserType]UserTypeMap),
	Register:   make(chan *Client),
	Unregister: make(chan *Client),
}
