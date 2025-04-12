package main

import (
	"github.com/joho/godotenv"
	"spotify-mcp/internal/server"
)

func main() {
	godotenv.Load()

	server.StartMcpServer()
}
