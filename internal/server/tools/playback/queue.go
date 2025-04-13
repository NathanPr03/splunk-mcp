package playback

import (
	"context"
	"fmt"
	"github.com/zmb3/spotify/v2"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"spotify-mcp/internal/client"
	"spotify-mcp/internal/server/tools"
)

func QueueTools() []tools.ToolEntry {
	return []tools.ToolEntry{
		getQueueTool(),
		queueSongTool(),
	}
}

func getQueueTool() tools.ToolEntry {
	toolDefinition := mcp.NewTool(
		"get_queue",
		mcp.WithDescription("Get the current Spotify playback queue"),
	)

	return tools.ToolEntry{
		ToolDefinition: toolDefinition,
		ToolBehaviour:  getQueueBehaviour,
	}
}

func getQueueBehaviour(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if !client.IsPlaybackAuthenticated() {
		return mcp.NewToolResultText("Not authenticated with Spotify. Please use the spotify_login tool first."), nil
	}

	queue, err := client.AuthenticatedSpotifyClient.GetQueue(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get queue: %w", err)
	}

	currentlyPlaying := formatTrack(&queue.CurrentlyPlaying, "Currently Playing")

	var queueItems []string
	if len(queue.Items) == 0 {
		queueItems = append(queueItems, "No upcoming tracks in the queue.")
	} else {
		for i, track := range queue.Items {
			queueItems = append(queueItems, formatTrack(&track, fmt.Sprintf("Queue #%d", i+1)))
		}
	}

	response := currentlyPlaying + "\n\nUpcoming in Queue:\n" + strings.Join(queueItems, "\n\n")
	return mcp.NewToolResultText(response), nil
}

func queueSongTool() tools.ToolEntry {
	toolDefinition := mcp.NewTool(
		"add_tracks_to_queue",
		mcp.WithDescription("Add tracks to your Spotify queue"),
		mcp.WithString("Track IDs",
			mcp.Description("Comma-separated list of Spotify track IDs"),
			mcp.Required(),
		),
	)

	return tools.ToolEntry{
		ToolDefinition: toolDefinition,
		ToolBehaviour:  queueSongBehaviour,
	}
}

func queueSongBehaviour(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if !client.IsPlaybackAuthenticated() {
		return mcp.NewToolResultText("Not authenticated with Spotify. Please use the spotify_login tool first."), nil
	}

	trackIdsParam, err := tools.GetParamFromRequest(request, "Track IDs")
	if err != nil {
		return nil, fmt.Errorf("failed to get track IDs parameter: %w", err)
	}

	trackIds := strings.Split(trackIdsParam, ",")

	for i := range trackIds {
		trackIds[i] = strings.TrimSpace(trackIds[i])
	}

	var failedTracks []string
	var addedTracks []string

	for _, trackId := range trackIds {
		err := client.AuthenticatedSpotifyClient.QueueSong(ctx, spotify.ID(trackId))
		if err != nil {
			failedTracks = append(failedTracks, trackId)
		} else {
			addedTracks = append(addedTracks, trackId)
		}
	}

	var responseMsg string
	if len(addedTracks) > 0 {
		responseMsg = fmt.Sprintf("Successfully added %d track(s) to your queue.", len(addedTracks))
	}

	if len(failedTracks) > 0 {
		if responseMsg != "" {
			responseMsg += "\n"
		}
		responseMsg += fmt.Sprintf("Failed to add %d track(s): %s", len(failedTracks), strings.Join(failedTracks, ", "))
	}

	return mcp.NewToolResultText(responseMsg), nil
}

func formatTrack(track *spotify.FullTrack, prefix string) string {
	if track == nil {
		return fmt.Sprintf("%s: No track information available", prefix)
	}

	var artists []string
	for _, artist := range track.Artists {
		artists = append(artists, artist.Name)
	}

	return fmt.Sprintf("%s: %s by %s\nAlbum: %s\nDuration: %d ms",
		prefix,
		track.Name,
		strings.Join(artists, ", "),
		track.Album.Name,
		track.Duration,
	)
}
