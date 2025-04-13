# Spotify MCP Server

A Model Context Protocol (MCP) server for Spotify integration, allowing AI assistants like Claude to control and interact with your Spotify account.

## Overview

This project provides a set of MCP tools that enable AI assistants to:

- Search for playlists
- Retrieve playlist details and tracks
- Control playback (play, pause, skip, previous)
- Create and modify playlists
- Toggle shuffle mode
- View currently playing tracks
- Add tracks to the queue

Built on top of the [Model Context Protocol](https://modelcontextprotocol.io/) and [zmb3/spotify](https://github.com/zmb3/spotify) Golang SDK.

## What is MCP?

The Model Context Protocol (MCP) is an open protocol that standardizes how applications provide context to Large Language Models (LLMs). It allows LLMs like Claude to interact with external systems and data sources in a secure and standardized way.

MCP works like a USB-C port for AI applications - providing a standardized way to connect AI models to different data sources and tools. This project implements a Spotify server that follows the MCP specification, enabling AI assistants to control and interact with your Spotify account.

## Setup

### Prerequisites

1. Go 1.18 or higher
2. A Spotify account (you will need Spotify Premium for playback control)
3. Spotify Developer credentials

### Getting Spotify API Credentials

1. Visit the [Spotify Developer Dashboard](https://developer.spotify.com/dashboard/)
2. Log in with your Spotify account
3. Click "Create an App"
4. Fill in the application name and description
5. Once created, you'll see your **Client ID** and you can view your **Client Secret**
6. Set the redirect URI to `http://127.0.0.1:1690/callback`

### Using with Claude

To use this server with Claude for Desktop:

1. Open Claude for Desktop
2. Create or update your MCP configuration at:
    - macOS: `~/Library/Application Support/Claude/claude_desktop_config.json`
    - Windows: `%APPDATA%\Claude\claude_desktop_config.json`

3. Add the following configuration:

```json
{
  "mcpServers": {
    "spotify": {
      "command": "/path/to/spotify-mcp-binary",
      "env": {
        "SPOTIFY_CLIENT_ID": "your_client_id_here",
        "SPOTIFY_CLIENT_SECRET": "your_client_secret_here"
      }
    }
  }
}
```

4. Restart Claude for Desktop
5. When first using Spotify tools, you'll need to authenticate using the `spotify_login` tool

## Available Tools

### Playback
- `spotify_login` - Start Spotify authentication process for playback control
- `play` - Start or resume playback on your Spotify account
- `pause` - Pause playback on your Spotify account
- `next_track` - Skip to the next track in your Spotify queue
- `previous_track` - Skip to the previous track in your Spotify queue
- `shuffle` - Toggle shuffle mode on your Spotify account
- `current_track` - Get information about the currently playing track 
- `get_queue` - Get the current playback queue
- `add_tracks_to_queue` - Add tracks to the current playback queue

### Playlist
- `get_playlist` - Get detailed information about a specific playlist
- `get_playlist_tracks` - Get the tracks in a playlist
- `create_playlist` - Create a new Spotify playlist
- `add_tracks_to_playlist` - Add tracks to a playlist
- `remove_tracks_from_playlist` - Remove tracks from a playlist
- `get_user_playlists` - Get playlists for a Spotify user

### Search
- `simple_playlist_and_album_search` - Search for a playlist or album by name
- `simple_song_search` - Search for a song by name

## License

[MIT License](https://mit-license.org/https://mit-license.org/)

## Acknowledgments

- Built with the [zmb3/spotify](https://github.com/zmb3/spotify) Golang SDK
- Implements the [Model Context Protocol](https://modelcontextprotocol.io/)
- Inspired by the growing ecosystem of MCP servers

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.