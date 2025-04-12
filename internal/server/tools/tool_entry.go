package tools

import (
	"github.com/mark3labs/mcp-go/mcp"
	mcpServer "github.com/mark3labs/mcp-go/server"
)

type ToolEntry struct {
	ToolDefinition mcp.Tool
	ToolBehaviour  mcpServer.ToolHandlerFunc
}
