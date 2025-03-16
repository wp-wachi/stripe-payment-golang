package main

import (
	"log"

	"github.com/wp-wachi/stripe-payment-golang/middlewares"
	"github.com/wp-wachi/stripe-payment-golang/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// config.LoadEnv()
	// database.ConnectDB()

	r := gin.Default()

	// Apply CORS Middleware
	r.Use(middlewares.CORSMiddleware())

	routes.RegisterRoutes(r)

	log.Println("Server running on port 8080")
	r.Run(":8080")
}