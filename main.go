package main

import (
	"context"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"

	client2 "spotify-mcp/internal/client"
	"spotify-mcp/internal/server"
)

func main() {
	server.StartMcpServer()

	client := client2.NewClient()
	defer client.Close()

	tools, err := client.ListTools(context.Background(), mcp.ListToolsRequest{})
	if err != nil {
		fmt.Printf("List tools error: %v\n", err)
		return
	}

	fmt.Printf("Available tools: %v\n", tools)

	toolCallRequest := mcp.CallToolRequest{
		Request: mcp.Request{
			Method: "tools/call",
		},
	}
	toolCallRequest.Params.Name = "Simple Playlist Search"
	toolCallRequest.Params.Arguments = map[string]interface{}{
		"Playlist Name": "Summer Vibes",
	}

	result, err := client.CallTool(context.Background(), toolCallRequest)
	if err != nil {
		fmt.Printf("Call tool error: %v\n", err)
		return
	}

	fmt.Printf("Tool result: %v\n", result)
}
