package playback

import (
	"context"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"spotify-mcp/internal/client"
	"spotify-mcp/internal/server/tools"
	"sync"
)

var (
	authInProgress = false
	authMutex      sync.Mutex
)

func PlayerTools() []tools.ToolEntry {
	return []tools.ToolEntry{
		loginTool(),
		playTool(),
		pauseTool(),
		nextTrackTool(),
		previousTrackTool(),
		shuffleTool(),
		currentTrackTool(),
	}
}

func loginTool() tools.ToolEntry {
	toolDefinition := mcp.NewTool(
		"spotify_login",
		mcp.WithDescription("Start Spotify authentication process"),
	)

	return tools.ToolEntry{
		ToolDefinition: toolDefinition,
		ToolBehaviour:  loginBehaviour,
	}
}

func loginBehaviour(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	authMutex.Lock()
	defer authMutex.Unlock()

	if client.IsPlaybackAuthenticated() {
		return mcp.NewToolResultText("Already authenticated with Spotify."), nil
	}

	if authInProgress {
		return mcp.NewToolResultText("Authentication already in progress. Please open the auth URL in your browser and complete the process."), nil
	}

	authInProgress = true
	authURL, err := client.InitiateAuth()
	if err != nil {
		authInProgress = false
		return nil, fmt.Errorf("failed to initiate authentication: %w", err)
	}

	return mcp.NewToolResultText(fmt.Sprintf(
		"Please authenticate with Spotify by opening this URL in your browser:\n%s\n\nAfter logging in, you'll be redirected to complete the authentication. Once completed, you can use the other Spotify tools. Please show this link directly to the end user.",
		authURL,
	)), nil
}

func currentTrackTool() tools.ToolEntry {
	toolDefinition := mcp.NewTool(
		"current_track",
		mcp.WithDescription("Get information about the currently playing track"),
	)

	return tools.ToolEntry{
		ToolDefinition: toolDefinition,
		ToolBehaviour:  currentTrackBehaviour,
	}
}

func currentTrackBehaviour(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if !client.IsPlaybackAuthenticated() {
		return mcp.NewToolResultText("Not authenticated with Spotify. Please use the spotify_login tool first."), nil
	}

	currentlyPlaying, err := client.AuthenticatedSpotifyClient.PlayerCurrentlyPlaying(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get currently playing track: %w", err)
	}

	if currentlyPlaying == nil || !currentlyPlaying.Playing {
		return mcp.NewToolResultText("No track is currently playing."), nil
	}

	track := currentlyPlaying.Item
	artists := ""
	for i, artist := range track.Artists {
		if i > 0 {
			artists += ", "
		}
		artists += artist.Name
	}

	response := fmt.Sprintf(
		"Currently playing: %s by %s\nAlbum: %s\nProgress: %d/%d ms\nIs Playing: %t",
		track.Name,
		artists,
		track.Album.Name,
		currentlyPlaying.Progress,
		track.Duration,
		currentlyPlaying.Playing,
	)

	return mcp.NewToolResultText(response), nil
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
	if !client.IsPlaybackAuthenticated() {
		return mcp.NewToolResultText("Not authenticated with Spotify. Please use the spotify_login tool first."), nil
	}

	err := client.AuthenticatedSpotifyClient.Play(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start playback: %w", err)
	}

	return mcp.NewToolResultText("Playback started"), nil
}

// Pause tool
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
	if !client.IsPlaybackAuthenticated() {
		return mcp.NewToolResultText("Not authenticated with Spotify. Please use the spotify_login tool first."), nil
	}

	err := client.AuthenticatedSpotifyClient.Pause(ctx)
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
	if !client.IsPlaybackAuthenticated() {
		return mcp.NewToolResultText("Not authenticated with Spotify. Please use the spotify_login tool first."), nil
	}

	err := client.AuthenticatedSpotifyClient.Next(ctx)
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
	if !client.IsPlaybackAuthenticated() {
		return mcp.NewToolResultText("Not authenticated with Spotify. Please use the spotify_login tool first."), nil
	}

	err := client.AuthenticatedSpotifyClient.Previous(ctx)
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
	if !client.IsPlaybackAuthenticated() {
		return mcp.NewToolResultText("Not authenticated with Spotify. Please use the spotify_login tool first."), nil
	}

	shuffleState, err := tools.GetBoolParamFromRequest(request, "state")
	if err != nil {
		return nil, fmt.Errorf("failed to get shuffle state parameter: %w", err)
	}

	err = client.AuthenticatedSpotifyClient.Shuffle(ctx, shuffleState)
	if err != nil {
		return nil, fmt.Errorf("failed to set shuffle state: %w", err)
	}

	statusMsg := "Shuffle disabled"
	if shuffleState {
		statusMsg = "Shuffle enabled"
	}

	return mcp.NewToolResultText(statusMsg), nil
}
