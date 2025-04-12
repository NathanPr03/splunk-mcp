package server

import (
	mcpServer "github.com/mark3labs/mcp-go/server"
	"spotify-mcp/internal/server/tools/search"
)

func StartMcpServer() {
	s := mcpServer.NewMCPServer(
		"Spotify MCP Server 🚀",
		"1.0.0",
	)

	tools := search.PlayListSearchTools()
	for _, tool := range tools {
		s.AddTool(tool.ToolDefinition, tool.ToolBehaviour)
	}

	baseUrl := "http://localhost:1690"
	option := mcpServer.WithBaseURL(baseUrl)

	sseServer := mcpServer.NewSSEServer(s, option)
	go sseServer.Start(":1690")
}
