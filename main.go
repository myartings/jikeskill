package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	port := flag.String("port", "8080", "Server port")
	tokenPath := flag.String("tokens", "tokens.json", "Path to token storage file")
	flag.Parse()

	app := NewAppServer(*tokenPath, *port)

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	app.setupRoutes(r)

	addr := fmt.Sprintf(":%s", *port)
	log.Printf("Jike MCP Server starting on %s", addr)
	log.Printf("MCP endpoint: http://localhost:%s/mcp", *port)
	log.Printf("REST API: http://localhost:%s/api/v1/", *port)

	if err := r.Run(addr); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
