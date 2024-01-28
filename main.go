package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	// Create a Gin router
	router := gin.Default()

	// Initialize WebSocket handler

	router.GET("/usertoshop", UserToShop)
	router.GET("/shoptouser", ShopToUser)

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Hello, this is the root route!",
		})
	})

	go Start()
	// Run the server on port 8080
	err := router.Run("localhost:8080")
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
