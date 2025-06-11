package main

import (
	"fmt"
	"jokylights-backend/handlers"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
		// return
	}
	fmt.Println("✅ .env loaded successfully")

	// Initialize the Gin router
	router := gin.Default()
	router.GET("/api/products", handlers.GetProducts)

	router.Run(":8080") // Start the server on port 8080
}
