package main

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

func (c *Client) Write() {
	defer func() {
		_ = c.Socket.Close()
	}()

	for {
		select {
		case dataToShop, ok := <-c.Send:
			if !ok {
				_ = c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			msg, _ := json.Marshal(dataToShop)
			err := c.Socket.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				return
			}
		}
	}
}
