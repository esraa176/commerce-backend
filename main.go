package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		// return
	}

	// Initialize the Gin router
	router := gin.Default()
	router.GET("/api/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Backend is working!",
		})
	})

	router.Run(":8080") // Start the server on port 8080
}
