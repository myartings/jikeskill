package main

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type AppServer struct {
	service   *JikeService
	mcpServer *mcp.Server
	port      string
}

func NewAppServer(tokenPath, port string) *AppServer {
	app := &AppServer{
		service: NewJikeService(tokenPath),
		port:    port,
	}
	app.mcpServer = newMCPServer(app)
	return app
}
