package search

import "spotify-mcp/internal/server/tools"

func GetSearchTools() []tools.ToolEntry {
	allTools := append(PlayListSearchTools(), SongSearchTools()...)
	return allTools
}
