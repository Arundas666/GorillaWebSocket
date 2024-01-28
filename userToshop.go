package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func UserToShop(c *gin.Context) {
	user_id := c.Query("user_id")
	if user_id == "" {
		return
	}

	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		}}).Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}
	client := &Client{
		ParticipantId: user_id,
		Socket:        conn,
		Send:          make(chan any),
		Usertype:      User,
	}
	Manager.Register <- client
	go client.ReadUser()
	go client.Write()
}

type SendMsg struct {
	Shop_id string `json:"shop_id"`
	Msg     string `json:"msg"`
}

func (c *Client) ReadUser() {
	defer func() {
		Manager.Unregister <- c
		_ = c.Socket.Close()
	}()

	for {
		var sendmsg SendMsg
		if err := c.Socket.ReadJSON(&sendmsg); err != nil {
			fmt.Println(err)
		}

		fmt.Println("Reading message...", sendmsg)
		_, ok := Manager.Clients[User]
		if !ok {
			fmt.Println("User nto there")
		}

		toClient, ok := Manager.Clients[Shop][sendmsg.Shop_id]
		fmt.Println(ok, "OK")
		if ok {
			toClient.Send <- sendmsg.Msg
		}
	}

}
