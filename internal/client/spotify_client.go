package client

import (
	"context"
	"github.com/google/uuid"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"log"
	"net/http"
	"os"
	"spotify-mcp/internal/token"
	"sync"
)

const redirectURI = "http://127.0.0.1:1690/callback"

var (
	// SpotifyClient Basic client for search functionality
	SpotifyClient *spotify.Client

	// AuthenticatedSpotifyClient Client with playback permissions
	AuthenticatedSpotifyClient *spotify.Client

	playbackAuth  *spotifyauth.Authenticator
	authComplete  = make(chan struct{})
	serverRunning bool
	state         string
	authMutex     sync.Mutex
	httpServer    *http.Server
)

func InstantiateSpotifyClient(ctx context.Context) {
	accessToken, err := token.GetToken(ctx)
	if err != nil {
		log.Fatalf("couldn't get token: %v", err)
	}

	httpClient := spotifyauth.New().Client(ctx, accessToken)
	SpotifyClient = spotify.New(httpClient)

	playbackAuth = spotifyauth.New(
		spotifyauth.WithRedirectURL(redirectURI),
		spotifyauth.WithScopes(
			spotifyauth.ScopeUserReadCurrentlyPlaying,
			spotifyauth.ScopeUserReadPlaybackState,
			spotifyauth.ScopeUserModifyPlaybackState,
			spotifyauth.ScopePlaylistModifyPublic,
			spotifyauth.ScopePlaylistModifyPrivate,
			spotifyauth.ScopePlaylistReadCollaborative,
			spotifyauth.ScopePlaylistReadPrivate,
		),
		spotifyauth.WithClientID(os.Getenv("SPOTIFY_CLIENT_ID")),
		spotifyauth.WithClientSecret(os.Getenv("SPOTIFY_CLIENT_SECRET")),
	)

	state = uuid.NewString()
}

func InitiateAuth() (string, error) {
	authMutex.Lock()
	defer authMutex.Unlock()

	if AuthenticatedSpotifyClient == nil {
		authComplete = make(chan struct{})
	}

	if !serverRunning {
		startAuthServer()
	}

	return playbackAuth.AuthURL(state), nil
}

func startAuthServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/callback", completeAuth)

	httpServer = &http.Server{
		Addr:    ":1690",
		Handler: mux,
	}

	go func() {
		serverRunning = true
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP server error: %v", err)
		}
		serverRunning = false
	}()

	log.Println("Started authentication server on port 1690")
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := playbackAuth.Token(r.Context(), state, r)
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

	AuthenticatedSpotifyClient = spotify.New(playbackAuth.Client(r.Context(), tok))

	w.Header().Set("Content-Type", "text/html")
	html := `
    <html>
    <body>
        <h1>Authentication Successful</h1>
        <p>You have successfully authenticated with Spotify.</p>
        <p>You can now close this browser window and return to your MCP host.</p>
    </body>
    </html>
    `
	w.Write([]byte(html))
	close(authComplete)

	go func() {
		if err := httpServer.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server shutdown error: %v", err)
		}
	}()
}

func IsPlaybackAuthenticated() bool {
	return AuthenticatedSpotifyClient != nil
}

func WaitForAuthentication() {
	<-authComplete
}
