package server

import (
	"context"
	mcpServer "github.com/mark3labs/mcp-go/server"
	"os"
	"spotify-mcp/internal/server/tools/search"
)

func StartMcpServer() {
	s := mcpServer.NewMCPServer(
		"Spotify MCP Server 🚀",
		"1.0.0",
	)

	tools := search.GetSearchTools()
	for _, tool := range tools {
		s.AddTool(tool.ToolDefinition, tool.ToolBehaviour)
	}

	sseServer := mcpServer.NewStdioServer(s)
	sseServer.Listen(context.Background(), os.Stdin, os.Stdout)
}
