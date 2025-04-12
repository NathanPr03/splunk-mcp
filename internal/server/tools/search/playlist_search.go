package search

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	_ "github.com/mark3labs/mcp-go/mcp"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"log"
	"spotify-mcp/internal/server/tools"
	"spotify-mcp/internal/token"
)

const playListName = "Playlist Name"

func PlayListSearchTools() []tools.ToolEntry {
	return []tools.ToolEntry{
		simplePlaylistSearch(),
	}
}

func simplePlaylistSearch() tools.ToolEntry {
	toolDefinition := mcp.NewTool(
		"simple_playlist_search",
		mcp.WithDescription("Search for a playlist by name"),
		mcp.WithString(playListName,
			mcp.Required(),
			mcp.Description("Name of the playlist to search for. Extra information: "+SearchQueryInformation),
		),
	)

	toolBehaviour := searchBehaviour

	return tools.ToolEntry{
		ToolDefinition: toolDefinition,
		ToolBehaviour:  toolBehaviour,
	}
}

func searchBehaviour(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	playlistName, err := tools.GetParamFromRequest(request, playListName)
	if err != nil {
		return nil, fmt.Errorf("failed to get playlist name: %w", err)
	}

	accessToken, err := token.GetToken(ctx)
	if err != nil {
		log.Fatalf("couldn't get token: %v", err)
	}

	httpClient := spotifyauth.New().Client(ctx, accessToken)
	client := spotify.New(httpClient)

	results, err := client.Search(ctx, playlistName, spotify.SearchTypePlaylist|spotify.SearchTypeAlbum)
	if err != nil {
		log.Fatal(err)
	}

	jsonPlaylists := ""
	if results.Playlists != nil {
		bytePlaylists, err := json.Marshal(results.Playlists)
		if err != nil {
			log.Fatal(err)
		}

		jsonPlaylists = string(bytePlaylists)
	}

	mcpResult := mcp.NewToolResultText(jsonPlaylists)

	return mcpResult, nil
}
