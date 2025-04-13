package playback

import (
	"github.com/mark3labs/mcp-go/mcp"
)

func PauseDefinition() mcp.Tool {
	pauseToolDef := mcp.NewTool(
		"pause_current_song",
		mcp.WithDescription("Pause the current song"),
	)

	return pauseToolDef
}

//func pauseBehaviour(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
//
//}
