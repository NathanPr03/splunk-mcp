package playlist

import (
	"context"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/zmb3/spotify/v2"
	"spotify-mcp/internal/client"
	"spotify-mcp/internal/server/tools"
	"strings"
)

func PlaylistTools() []tools.ToolEntry {
	return []tools.ToolEntry{
		getPlaylistTool(),
		getPlaylistTracksTool(),
		createPlaylistTool(),
		addTracksToPlaylistTool(),
		removeTracksFromPlaylistTool(),
		getUserPlaylistsTool(),
	}
}

func getPlaylistTool() tools.ToolEntry {
	toolDefinition := mcp.NewTool(
		"get_playlist",
		mcp.WithDescription("Get detailed information about a specific playlist"),
		mcp.WithString("Playlist ID",
			mcp.Required(),
			mcp.Description("Spotify ID of the playlist"),
		),
	)

	return tools.ToolEntry{
		ToolDefinition: toolDefinition,
		ToolBehaviour:  getPlaylistBehaviour,
	}
}

func getPlaylistBehaviour(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	playlistID, err := tools.GetParamFromRequest(request, "Playlist ID")
	if err != nil {
		return nil, fmt.Errorf("failed to get playlist ID: %w", err)
	}

	spotifyClient := client.SpotifyClient
	if spotifyClient == nil {
		return mcp.NewToolResultText("Spotify client not initialized. Please use the spotify_login tool first."), nil
	}

	playlist, err := spotifyClient.GetPlaylist(ctx, spotify.ID(playlistID))
	if err != nil {
		return nil, fmt.Errorf("failed to get playlist: %w", err)
	}

	owner := playlist.Owner.DisplayName
	if owner == "" {
		owner = playlist.Owner.ID
	}

	response := fmt.Sprintf("Playlist: %s (by %s)\n", playlist.Name, owner)
	response += fmt.Sprintf("ID: %s\n", playlist.ID)
	response += fmt.Sprintf("Tracks: %d\n", playlist.Tracks.Total)
	response += fmt.Sprintf("Public: %t\n", playlist.IsPublic)
	response += fmt.Sprintf("Collaborative: %t\n", playlist.Collaborative)

	if playlist.Description != "" {
		response += fmt.Sprintf("Description: %s\n", playlist.Description)
	}

	if len(playlist.Images) > 0 {
		response += fmt.Sprintf("Image URL: %s\n", playlist.Images[0].URL)
	}

	response += fmt.Sprintf("Followers: %d\n", playlist.Followers.Count)
	response += fmt.Sprintf("URL: %s\n", playlist.ExternalURLs["spotify"])

	return mcp.NewToolResultText(response), nil
}

func getPlaylistTracksTool() tools.ToolEntry {
	toolDefinition := mcp.NewTool(
		"get_playlist_tracks",
		mcp.WithDescription("Get the tracks in a playlist"),
		mcp.WithString("Playlist ID",
			mcp.Required(),
			mcp.Description("Spotify ID of the playlist"),
		),
		mcp.WithNumber("Limit",
			mcp.Description("Maximum number of tracks to return (default: 20, max: 100)"),
		),
		mcp.WithNumber("Offset",
			mcp.Description("The index of the first track to return (default: 0)"),
		),
	)

	return tools.ToolEntry{
		ToolDefinition: toolDefinition,
		ToolBehaviour:  getPlaylistTracksBehaviour,
	}
}

func getPlaylistTracksBehaviour(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	playlistID, err := tools.GetParamFromRequest(request, "Playlist ID")
	if err != nil {
		return nil, fmt.Errorf("failed to get playlist ID: %w", err)
	}

	limit, _ := tools.GetIntParamFromRequest(request, "Limit")

	offset, _ := tools.GetIntParamFromRequest(request, "Offset")

	if limit <= 0 || limit > 100 {
		limit = 20
	}

	spotifyClient := client.SpotifyClient
	if spotifyClient == nil {
		return mcp.NewToolResultText("Spotify client not initialized. Please use the spotify_login tool first."), nil
	}

	opts := []spotify.RequestOption{
		spotify.Limit(limit),
		spotify.Offset(offset),
	}

	playlistItems, err := spotifyClient.GetPlaylistItems(ctx, spotify.ID(playlistID), opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to get playlist tracks: %w", err)
	}

	response := fmt.Sprintf("Tracks in playlist (showing %d of %d total):\n\n",
		len(playlistItems.Items), playlistItems.Total)

	for i, item := range playlistItems.Items {
		if item.Track.Track == nil {
			continue
		}

		track := item.Track.Track

		artists := ""
		for j, artist := range track.Artists {
			if j > 0 {
				artists += ", "
			}
			artists += artist.Name
		}

		trackNum := i + offset + 1

		response += fmt.Sprintf("%d. %s - %s\n", trackNum, track.Name, artists)
		response += fmt.Sprintf("   Album: %s\n", track.Album.Name)
		response += fmt.Sprintf("   Duration: %d ms\n", track.Duration)
		response += fmt.Sprintf("   Track ID: %s\n", track.ID)

		if item.AddedBy.ID != "" {
			addedBy := item.AddedBy.DisplayName
			if addedBy == "" {
				addedBy = item.AddedBy.ID
			}
			response += fmt.Sprintf("   Added by: %s\n", addedBy)
		}

		if item.AddedAt != "" {
			response += fmt.Sprintf("   Added at: %s\n", item.AddedAt)
		}

		response += "\n"
	}

	if int(playlistItems.Total) > limit {
		response += fmt.Sprintf("\nShowing tracks %d-%d of %d. Use the Offset parameter to see more tracks.",
			offset+1, offset+len(playlistItems.Items), playlistItems.Total)
	}

	return mcp.NewToolResultText(response), nil
}

func createPlaylistTool() tools.ToolEntry {
	toolDefinition := mcp.NewTool(
		"create_playlist",
		mcp.WithDescription("Create a new Spotify playlist"),
		mcp.WithString("Name",
			mcp.Required(),
			mcp.Description("Name of the playlist"),
		),
		mcp.WithString("Description",
			mcp.Description("Description of the playlist"),
		),
		mcp.WithBoolean("Public",
			mcp.Description("Whether the playlist should be public (default: false)"),
		),
		mcp.WithBoolean("Collaborative",
			mcp.Description("Whether the playlist should be collaborative (default: false)"),
		),
	)

	return tools.ToolEntry{
		ToolDefinition: toolDefinition,
		ToolBehaviour:  createPlaylistBehaviour,
	}
}

func createPlaylistBehaviour(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := tools.GetParamFromRequest(request, "Name")
	if err != nil {
		return nil, fmt.Errorf("failed to get playlist name: %w", err)
	}

	if !client.IsPlaybackAuthenticated() {
		return mcp.NewToolResultText("Not authenticated with Spotify for playlist creation. Please use the spotify_login tool first."), nil
	}

	description, err := tools.GetParamFromRequest(request, "Description")
	if err != nil {
		return nil, fmt.Errorf("failed to get playlist description: %w", err)
	}

	isPublic, err := tools.GetBoolParamFromRequest(request, "Public")
	if err != nil {
		return nil, fmt.Errorf("failed to get public parameter: %w", err)
	}

	isCollaborative, err := tools.GetBoolParamFromRequest(request, "Collaborative")
	if err != nil {
		return nil, fmt.Errorf("failed to get collaborative parameter: %w", err)
	}

	user, err := client.AuthenticatedSpotifyClient.CurrentUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get current user: %w", err)
	}

	playlist, err := client.AuthenticatedSpotifyClient.CreatePlaylistForUser(
		ctx,
		user.ID,
		name,
		description,
		isPublic,
		isCollaborative,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create playlist: %w", err)
	}

	response := fmt.Sprintf("Successfully created playlist!\n\n")
	response += fmt.Sprintf("Name: %s\n", playlist.Name)
	response += fmt.Sprintf("ID: %s\n", playlist.ID)
	response += fmt.Sprintf("Public: %t\n", playlist.IsPublic)
	response += fmt.Sprintf("Collaborative: %t\n", playlist.Collaborative)

	if description != "" {
		response += fmt.Sprintf("Description: %s\n", description)
	}

	response += fmt.Sprintf("URL: %s\n", playlist.ExternalURLs["spotify"])

	return mcp.NewToolResultText(response), nil
}

func addTracksToPlaylistTool() tools.ToolEntry {
	toolDefinition := mcp.NewTool(
		"add_tracks_to_playlist",
		mcp.WithDescription("Add tracks to a playlist"),
		mcp.WithString("Playlist ID",
			mcp.Required(),
			mcp.Description("Spotify ID of the playlist"),
		),
		mcp.WithString("Track IDs",
			mcp.Required(),
			mcp.Description("Comma-separated list of Spotify track IDs"),
		),
	)

	return tools.ToolEntry{
		ToolDefinition: toolDefinition,
		ToolBehaviour:  addTracksToPlaylistBehaviour,
	}
}

func addTracksToPlaylistBehaviour(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	playlistID, err := tools.GetParamFromRequest(request, "Playlist ID")
	if err != nil {
		return nil, fmt.Errorf("failed to get playlist ID: %w", err)
	}

	trackIDsStr, err := tools.GetParamFromRequest(request, "Track IDs")
	if err != nil {
		return nil, fmt.Errorf("failed to get track IDs: %w", err)
	}

	if !client.IsPlaybackAuthenticated() {
		return mcp.NewToolResultText("Not authenticated with Spotify for playlist modification. Please use the spotify_login tool first."), nil
	}

	trackIDsList := strings.Split(trackIDsStr, ",")
	for i := range trackIDsList {
		trackIDsList[i] = strings.TrimSpace(trackIDsList[i])
	}

	trackIDs := make([]spotify.ID, 0, len(trackIDsList))
	for _, id := range trackIDsList {
		if id != "" {
			trackIDs = append(trackIDs, spotify.ID(id))
		}
	}

	if len(trackIDs) == 0 {
		return mcp.NewToolResultText("No valid track IDs provided."), nil
	}

	if len(trackIDs) > 100 {
		return mcp.NewToolResultText("Too many track IDs provided. Maximum is 100 tracks per request."), nil
	}

	snapshotID, err := client.AuthenticatedSpotifyClient.AddTracksToPlaylist(ctx, spotify.ID(playlistID), trackIDs...)
	if err != nil {
		return nil, fmt.Errorf("failed to add tracks to playlist: %w", err)
	}

	response := fmt.Sprintf("Successfully added %d tracks to the playlist!\n", len(trackIDs))
	response += fmt.Sprintf("Playlist ID: %s\n", playlistID)
	response += fmt.Sprintf("New snapshot ID: %s\n", snapshotID)

	return mcp.NewToolResultText(response), nil
}

func removeTracksFromPlaylistTool() tools.ToolEntry {
	toolDefinition := mcp.NewTool(
		"remove_tracks_from_playlist",
		mcp.WithDescription("Remove tracks from a playlist"),
		mcp.WithString("Playlist ID",
			mcp.Required(),
			mcp.Description("Spotify ID of the playlist"),
		),
		mcp.WithString("Track IDs",
			mcp.Required(),
			mcp.Description("Comma-separated list of Spotify track IDs to remove"),
		),
	)

	return tools.ToolEntry{
		ToolDefinition: toolDefinition,
		ToolBehaviour:  removeTracksFromPlaylistBehaviour,
	}
}

func removeTracksFromPlaylistBehaviour(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	playlistID, err := tools.GetParamFromRequest(request, "Playlist ID")
	if err != nil {
		return nil, fmt.Errorf("failed to get playlist ID: %w", err)
	}

	trackIDsStr, err := tools.GetParamFromRequest(request, "Track IDs")
	if err != nil {
		return nil, fmt.Errorf("failed to get track IDs: %w", err)
	}

	if !client.IsPlaybackAuthenticated() {
		return mcp.NewToolResultText("Not authenticated with Spotify for playlist modification. Please use the spotify_login tool first."), nil
	}

	trackIDsList := strings.Split(trackIDsStr, ",")
	for i := range trackIDsList {
		trackIDsList[i] = strings.TrimSpace(trackIDsList[i])
	}

	trackIDs := make([]spotify.ID, 0, len(trackIDsList))
	for _, id := range trackIDsList {
		if id != "" {
			trackIDs = append(trackIDs, spotify.ID(id))
		}
	}

	if len(trackIDs) == 0 {
		return mcp.NewToolResultText("No valid track IDs provided."), nil
	}

	snapshotID, err := client.AuthenticatedSpotifyClient.RemoveTracksFromPlaylist(ctx, spotify.ID(playlistID), trackIDs...)
	if err != nil {
		return nil, fmt.Errorf("failed to remove tracks from playlist: %w", err)
	}

	response := fmt.Sprintf("Successfully removed %d tracks from the playlist!\n", len(trackIDs))
	response += fmt.Sprintf("Playlist ID: %s\n", playlistID)
	response += fmt.Sprintf("New snapshot ID: %s\n", snapshotID)

	return mcp.NewToolResultText(response), nil
}

func getUserPlaylistsTool() tools.ToolEntry {
	toolDefinition := mcp.NewTool(
		"get_user_playlists",
		mcp.WithDescription("Get playlists for a Spotify user"),
		mcp.WithString("User ID",
			mcp.Description("Spotify user ID (leave empty for current user)"),
		),
		mcp.WithNumber("Limit",
			mcp.Description("Maximum number of playlists to return (default: 20, max: 50)"),
		),
		mcp.WithNumber("Offset",
			mcp.Description("The index of the first playlist to return (default: 0)"),
		),
	)

	return tools.ToolEntry{
		ToolDefinition: toolDefinition,
		ToolBehaviour:  getUserPlaylistsBehaviour,
	}
}

func getUserPlaylistsBehaviour(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if !client.IsPlaybackAuthenticated() {
		return mcp.NewToolResultText("Not authenticated with Spotify. Please use the spotify_login tool first."), nil
	}

	userID, err := tools.GetParamFromRequest(request, "User ID")
	if err != nil {
		return nil, fmt.Errorf("failed to get user ID: %w", err)
	}

	if userID == "" {
		user, err := client.AuthenticatedSpotifyClient.CurrentUser(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get current user: %w", err)
		}
		userID = user.ID
	}

	limit, _ := tools.GetIntParamFromRequest(request, "Limit")

	offset, _ := tools.GetIntParamFromRequest(request, "Offset")

	if limit <= 0 || limit > 50 {
		limit = 20
	}

	opts := []spotify.RequestOption{
		spotify.Limit(limit),
		spotify.Offset(offset),
	}

	playlists, err := client.AuthenticatedSpotifyClient.GetPlaylistsForUser(ctx, userID, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to get user playlists: %w", err)
	}

	userDisplayText := userID
	if userID == "" {
		userDisplayText = "current user"
	}

	response := fmt.Sprintf("Playlists for %s (showing %d of %d total):\n\n",
		userDisplayText, len(playlists.Playlists), playlists.Total)

	for i, playlist := range playlists.Playlists {
		playlistNum := i + offset + 1

		owner := playlist.Owner.DisplayName
		if owner == "" {
			owner = playlist.Owner.ID
		}

		response += fmt.Sprintf("%d. %s\n", playlistNum, playlist.Name)
		if owner != userID && owner != "" {
			response += fmt.Sprintf("   Owner: %s\n", owner)
		}
		response += fmt.Sprintf("   Tracks: %d\n", playlist.Tracks.Total)
		response += fmt.Sprintf("   ID: %s\n", playlist.ID)

		if playlist.IsPublic {
			response += "   Public: Yes\n"
		} else {
			response += "   Public: No\n"
		}

		if playlist.Collaborative {
			response += "   Collaborative: Yes\n"
		}

		if playlist.Description != "" {
			desc := playlist.Description
			if len(desc) > 100 {
				desc = desc[:97] + "..."
			}
			response += fmt.Sprintf("   Description: %s\n", desc)
		}

		response += "\n"
	}

	if int(playlists.Total) > limit {
		response += fmt.Sprintf("\nShowing playlists %d-%d of %d. Use the Offset parameter to see more playlists.",
			offset+1, offset+len(playlists.Playlists), playlists.Total)
	}

	return mcp.NewToolResultText(response), nil
}
