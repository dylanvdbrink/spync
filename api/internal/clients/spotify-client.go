package clients

import (
	"api/internal/configuration"
	"context"
	"fmt"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
	"math/rand"
	"net/http"
)

func getAuthenticator() *spotifyauth.Authenticator {
	config, _ := configuration.GetConfiguration()
	return spotifyauth.New(
		spotifyauth.WithRedirectURL(config.Spotify.RedirectUrl),
		spotifyauth.WithScopes(spotifyauth.ScopePlaylistReadPrivate, spotifyauth.ScopePlaylistReadCollaborative),
		spotifyauth.WithClientID(config.Spotify.ClientId),
		spotifyauth.WithClientSecret(config.Spotify.ClientSecret),
	)
}

func GetAuthUrl() string {
	return getAuthenticator().AuthURL(fmt.Sprint("spync-", rand.Int()))
}

func SaveAuth(r *http.Request, state string) error {
	token, err := getAuthenticator().Token(r.Context(), state, r)
	if err != nil {
		return err
	}

	config, _ := configuration.GetConfiguration()
	config.Spotify.AccessToken = token.AccessToken
	config.Spotify.RefreshToken = token.RefreshToken
	config.Spotify.Expiry = token.Expiry
	config.Spotify.TokenType = token.TokenType
	_ = configuration.SaveConfiguration(config)
	return nil
}

func GetUserPlaylists() ([]spotify.SimplePlaylist, error) {
	total := 20
	offset := 0
	items := make([]spotify.SimplePlaylist, 0)
	for total != len(items) {
		playlistPage, err := getSpotifyClient().CurrentUsersPlaylists(context.TODO(),
			spotify.Limit(50), spotify.Offset(offset))
		if err != nil {
			return nil, err
		}

		items = append(items, playlistPage.Playlists...)
		total = playlistPage.Total
		offset += 50
	}

	return items, nil
}

func GetMe() (*spotify.PrivateUser, error) {
	profile, err := getSpotifyClient().CurrentUser(context.TODO())
	if err != nil {
		return nil, err
	}

	return profile, nil
}

func getSpotifyClient() *spotify.Client {
	config, _ := configuration.GetConfiguration()
	token := oauth2.Token{
		AccessToken:  config.Spotify.AccessToken,
		RefreshToken: config.Spotify.RefreshToken,
		TokenType:    config.Spotify.TokenType,
		Expiry:       config.Spotify.Expiry,
	}
	httpClient := getAuthenticator().Client(context.TODO(), &token)
	return spotify.New(httpClient, spotify.WithRetry(true))
}
