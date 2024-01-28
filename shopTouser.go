package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func ShopToUser(c *gin.Context) {
	shop_id := c.Query("shop_id")
	if shop_id == "" {
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
		ParticipantId: shop_id,
		Socket:        conn,
		Send:          make(chan any),
		Usertype:      Shop,
	}
	Manager.Register <- client
	go client.ReadShop()
	go client.Write()
}

type SendMsgFromShop struct {
	UserId string `json:"user_id"`
	Msg     string `json:"msg"`
}

func (c *Client) ReadShop() {
	defer func() {
		Manager.Unregister <- c
		_ = c.Socket.Close()
	}()

	for {
		var sendmsg SendMsgFromShop
		if err := c.Socket.ReadJSON(&sendmsg); err != nil {
			fmt.Println(err)
		}

		fmt.Println("Reading message...", sendmsg)
		_, ok := Manager.Clients[User]
		if !ok {
			fmt.Println("Shop is not there")
		}

		toClient, ok := Manager.Clients[Shop][sendmsg.UserId]
		fmt.Println(ok, "OK")
		if ok {
			toClient.Send <- sendmsg.Msg
		}
	}

}
