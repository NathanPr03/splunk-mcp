package search

import "spotify-mcp/internal/server/tools"

func SearchTools() []tools.ToolEntry {
	allTools := append(PlayListSearchTools(), SongSearchTools()...)
	return allTools
}
