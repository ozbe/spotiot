// Seeded from https://github.com/zmb3/spotify/blob/master/examples/authenticate/authcode/authenticate.go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"

	_ "github.com/joho/godotenv/autoload"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"

	"github.com/google/uuid"
)

func main() {
	const redirectURI = "http://localhost:8080/callback"
	auth := spotifyauth.New(
		spotifyauth.WithRedirectURL(redirectURI),
		spotifyauth.WithScopes(spotifyauth.ScopeUserReadCurrentlyPlaying, spotifyauth.ScopeUserReadPlaybackState, spotifyauth.ScopeUserModifyPlaybackState),
	)
	ch := make(chan *oauth2.Token)
	state := uuid.NewString()

	http.Handle("/callback", completeAuth(auth, state, ch))

	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatal(err)
		}
	}()

	url := auth.AuthURL(state)
	err := exec.Command("open", url).Start()
	if err != nil {
		log.Fatal(err)
	}

	tok := <-ch
	fmt.Printf("SPOTIFY_ACCESS_TOKEN=%s\n", tok.AccessToken)
	fmt.Printf("SPOTIFY_REFRESH_TOKEN=%s\n", tok.RefreshToken)

	assertUserID := os.Getenv("ASSERT_USER_ID")
	if assertUserID != "" {
		ctx := context.Background()
		client := spotify.New(auth.Client(ctx, tok))
		user, err := client.CurrentUser(ctx)
		if err != nil {
			log.Fatal(err)
		}
		if assertUserID != user.ID {
			log.Fatalf("Unexpected user ID: %s", user.ID)
		}
	}
}

func completeAuth(auth *spotifyauth.Authenticator, state string, ch chan<- *oauth2.Token) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tok, err := auth.Token(r.Context(), state, r)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			log.Fatal("Couldn't exchange token", err)
		}
		if st := r.FormValue("state"); st != state {
			http.NotFound(w, r)
			log.Fatalf("State mismatch: %s != %s\n", st, state)
		}
		fmt.Fprint(w, "Login complete")
		ch <- tok
	})
}
