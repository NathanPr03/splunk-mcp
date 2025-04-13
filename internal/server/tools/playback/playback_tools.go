package playback

import (
	"context"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"spotify-mcp/internal/client"
	"spotify-mcp/internal/server/tools"
)

func PlayerTools() []tools.ToolEntry {
	return []tools.ToolEntry{
		playTool(),
		pauseTool(),
		nextTrackTool(),
		previousTrackTool(),
		shuffleTool(),
	}
}

func playTool() tools.ToolEntry {
	toolDefinition := mcp.NewTool(
		"play",
		mcp.WithDescription("Start or resume playback on your Spotify account"),
	)

	return tools.ToolEntry{
		ToolDefinition: toolDefinition,
		ToolBehaviour:  playBehaviour,
	}
}

func playBehaviour(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	err := client.SpotifyClient.Play(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start playback: %w", err)
	}

	return mcp.NewToolResultText("Playback started"), nil
}

func pauseTool() tools.ToolEntry {
	toolDefinition := mcp.NewTool(
		"pause",
		mcp.WithDescription("Pause playback on your Spotify account"),
	)

	return tools.ToolEntry{
		ToolDefinition: toolDefinition,
		ToolBehaviour:  pauseBehaviour,
	}
}

func pauseBehaviour(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	err := client.SpotifyClient.Pause(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to pause playback: %w", err)
	}

	return mcp.NewToolResultText("Playback paused"), nil
}

func nextTrackTool() tools.ToolEntry {
	toolDefinition := mcp.NewTool(
		"next_track",
		mcp.WithDescription("Skip to the next track in your Spotify queue"),
	)

	return tools.ToolEntry{
		ToolDefinition: toolDefinition,
		ToolBehaviour:  nextTrackBehaviour,
	}
}

func nextTrackBehaviour(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	err := client.SpotifyClient.Next(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to skip to next track: %w", err)
	}

	return mcp.NewToolResultText("Skipped to next track"), nil
}

func previousTrackTool() tools.ToolEntry {
	toolDefinition := mcp.NewTool(
		"previous_track",
		mcp.WithDescription("Skip to the previous track in your Spotify queue"),
	)

	return tools.ToolEntry{
		ToolDefinition: toolDefinition,
		ToolBehaviour:  previousTrackBehaviour,
	}
}

func previousTrackBehaviour(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	err := client.SpotifyClient.Previous(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to skip to previous track: %w", err)
	}

	return mcp.NewToolResultText("Skipped to previous track"), nil
}

func shuffleTool() tools.ToolEntry {
	toolDefinition := mcp.NewTool(
		"shuffle",
		mcp.WithDescription("Toggle shuffle mode on your Spotify account"),
		mcp.WithBoolean("state",
			mcp.Description("Set to true to enable shuffle, false to disable"),
			mcp.Required(),
		),
	)

	return tools.ToolEntry{
		ToolDefinition: toolDefinition,
		ToolBehaviour:  shuffleBehaviour,
	}
}

func shuffleBehaviour(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	shuffleState, err := tools.GetBoolParamFromRequest(request, "state")
	if err != nil {
		return nil, fmt.Errorf("failed to get shuffle state parameter: %w", err)
	}

	err = client.SpotifyClient.Shuffle(ctx, shuffleState)
	if err != nil {
		return nil, fmt.Errorf("failed to set shuffle state: %w", err)
	}

	statusMsg := "Shuffle disabled"
	if shuffleState {
		statusMsg = "Shuffle enabled"
	}

	return mcp.NewToolResultText(statusMsg), nil
}
