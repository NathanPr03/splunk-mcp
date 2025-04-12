package client

import (
	"context"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"log"
	"spotify-mcp/internal/token"
)

var SpotifyClient *spotify.Client

func InstantiateSpotifyClient(ctx context.Context) {
	accessToken, err := token.GetToken(ctx)
	if err != nil {
		log.Fatalf("couldn't get token: %v", err)
	}

	httpClient := spotifyauth.New().Client(ctx, accessToken)
	SpotifyClient = spotify.New(httpClient)
}
