package client

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"log"
	"net/http"
	"os"
	"sync"
)

const redirectURI = "http://127.0.0.1:1690/callback"

var (
	SpotifyClient  *spotify.Client
	playerState    *spotify.PlayerState
	auth           *spotifyauth.Authenticator
	clientInitOnce sync.Once
	authComplete   = make(chan struct{})
	state          = "abc123"
)

func InstantiateSpotifyClient() {
	clientInitOnce.Do(func() {
		err := godotenv.Load()
		if err != nil {
			log.Printf("Error loading .env file: %v", err)
		}

		auth = spotifyauth.New(
			spotifyauth.WithRedirectURL(redirectURI),
			spotifyauth.WithScopes(
				spotifyauth.ScopeUserReadCurrentlyPlaying,
				spotifyauth.ScopeUserReadPlaybackState,
				spotifyauth.ScopeUserModifyPlaybackState,
			),
			spotifyauth.WithClientID(os.Getenv("SPOTIFY_CLIENT_ID")),
			spotifyauth.WithClientSecret(os.Getenv("SPOTIFY_CLIENT_SECRET")),
		)

		http.HandleFunc("/callback", completeAuth)

		go func() {
			log.Println("Starting HTTP server for Spotify authentication")
			if err := http.ListenAndServe(":1690", nil); err != nil {
				log.Fatalf("HTTP server error: %v", err)
			}
		}()

		url := auth.AuthURL(state)
		fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)

		<-authComplete

		playerState, err = SpotifyClient.PlayerState(context.Background())
		if err != nil {
			log.Printf("Warning: Could not get player state: %v", err)
		} else {
			log.Printf("Connected to Spotify device: %s (%s)", playerState.Device.Type, playerState.Device.Name)
		}
	})
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(r.Context(), state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Printf("Authentication error: %v", err)
		return
	}

	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Printf("State mismatch: %s != %s", st, state)
		return
	}

	SpotifyClient = spotify.New(auth.Client(r.Context(), tok))

	fmt.Fprintf(w, "Login Completed! You can now close this window and return to the application.")

	close(authComplete)
}

func GetCurrentPlayerState(ctx context.Context) (*spotify.PlayerState, error) {
	state, err := SpotifyClient.PlayerState(ctx)
	if err != nil {
		return nil, err
	}
	playerState = state
	return state, nil
}
