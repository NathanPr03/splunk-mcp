package main

import (
	"context"
	"github.com/joho/godotenv"
	"spotify-mcp/internal/client"
	"spotify-mcp/internal/server"
)

func main() {
	godotenv.Load()

	client.InstantiateSpotifyClient(context.Background())

	server.StartMcpServer()
}
