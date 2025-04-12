package search

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/zmb3/spotify/v2"
	"log"
	"spotify-mcp/internal/client"
	"spotify-mcp/internal/server/tools"
)

const songNameParameter = "Song Name"

func SongSearchTools() []tools.ToolEntry {
	return []tools.ToolEntry{
		simpleSongSearch(),
	}
}

func simpleSongSearch() tools.ToolEntry {
	toolDefinition := mcp.NewTool(
		"simple_song_search",
		mcp.WithDescription("Search for a song by name"),
		mcp.WithString(songNameParameter,
			mcp.Required(),
			mcp.Description("Name of the playlist to search for. Extra information: "+SearchQueryInformation+FilterInformation),
		),
	)

	toolBehaviour := songSearchBehaviour

	return tools.ToolEntry{
		ToolDefinition: toolDefinition,
		ToolBehaviour:  toolBehaviour,
	}
}

func songSearchBehaviour(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	songName, err := tools.GetParamFromRequest(request, songNameParameter)
	if err != nil {
		return nil, fmt.Errorf("failed to get playlist name: %w", err)
	}

	// We can set this limit quite low as songs generally don't clash names.
	// The client can also specify an album/artist to narrow down the search.
	songLimit := spotify.Limit(5)
	results, err := client.SpotifyClient.Search(ctx, songName, spotify.SearchTypeTrack, songLimit)
	if err != nil {
		log.Fatal(err)
	}

	jsonPlaylists := ""
	if results.Tracks != nil {
		bytePlaylists, err := json.Marshal(results.Tracks.Tracks)
		if err != nil {
			log.Fatal(err)
		}

		jsonPlaylists = string(bytePlaylists)
	}

	mcpResult := mcp.NewToolResultText(jsonPlaylists)

	return mcpResult, nil
}
