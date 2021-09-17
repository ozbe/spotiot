package main

import (
	"context"
	"fmt"
	"os"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var client *spotify.Client

func init() {
	ctx := context.Background()
	tok := &oauth2.Token{
		AccessToken:  os.Getenv("SPOTIFY_ACCESS_TOKEN"),
		TokenType:    "Bearer",
		RefreshToken: os.Getenv("SPOTIFY_REFRESH_TOKEN"),
	}
	auth := spotifyauth.New(
		spotifyauth.WithScopes(spotifyauth.ScopeUserReadCurrentlyPlaying, spotifyauth.ScopeUserReadPlaybackState, spotifyauth.ScopeUserModifyPlaybackState),
	)
	client = spotify.New(auth.Client(ctx, tok))
}

func handleRequest(ctx context.Context, event events.IoTButtonEvent) error {
	playerState, err := client.PlayerState(context.Background())
	if err != nil {
		return fmt.Errorf("player state: %w", err)
	}

	switch event.ClickType {
	case "LONG":
		return client.Shuffle(ctx, !playerState.ShuffleState)
	case "DOUBLE":
		return client.Next(ctx)
	case "SINGLE":
		if playerState.CurrentlyPlaying.Playing {
			return client.Pause(ctx)
		}
		return client.Play(ctx)
	}

	return fmt.Errorf("unknown click type: %s", event.ClickType)
}

func main() {
	lambda.Start(handleRequest)
}
