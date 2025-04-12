package search

import (
	"context"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	_ "github.com/mark3labs/mcp-go/mcp"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2/clientcredentials"
	"log"
	"spotify-mcp/internal/server/tools"
)

const playListName = "Playlist Name"

func PlayListSearchTools() []tools.ToolEntry {
	return []tools.ToolEntry{
		simplePlaylistSearch(),
	}
}

func simplePlaylistSearch() tools.ToolEntry {
	toolDefinition := mcp.NewTool(
		"Simple Playlist Search",
		mcp.WithDescription("Search for a playlist by name"),
		mcp.WithString(playListName,
			mcp.Required(),
			mcp.Description("Name of the playlist to search for"),
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

	config := &clientcredentials.Config{
		ClientID:     "REDACTED",
		ClientSecret: "REDACTED",
		TokenURL:     spotifyauth.TokenURL,
	}
	token, err := config.Token(ctx)
	if err != nil {
		log.Fatalf("couldn't get token: %v", err)
	}

	httpClient := spotifyauth.New().Client(ctx, token)
	client := spotify.New(httpClient)

	results, err := client.Search(ctx, playlistName, spotify.SearchTypePlaylist|spotify.SearchTypeAlbum)
	if err != nil {
		log.Fatal(err)
	}

	// handle album results
	if results.Albums != nil {
		fmt.Println("Albums:")
		for _, item := range results.Albums.Albums {
			fmt.Println("   ", item.Name)
		}
	}

	allPlaylistNames := ""
	if results.Playlists != nil {
		fmt.Println("Playlists:")
		for _, item := range results.Playlists.Playlists {
			fmt.Println("   ", item.Name)
			allPlaylistNames += item.Name + "\n"
		}
	}

	mcpResult := mcp.NewToolResultText(allPlaylistNames)
	return mcpResult, nil
}
