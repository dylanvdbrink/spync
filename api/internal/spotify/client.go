package spotify

import (
	"api/internal/configuration"
	"api/internal/state"
	"context"
	"fmt"
	spotifylib "github.com/zmb3/spotify/v2"
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

func SaveAuth(r *http.Request, stateParam string) error {
	token, err := getAuthenticator().Token(r.Context(), stateParam, r)
	if err != nil {
		return err
	}

	stateObj, _ := state.GetState()
	stateObj.Spotify.AccessToken = token.AccessToken
	stateObj.Spotify.RefreshToken = token.RefreshToken
	stateObj.Spotify.Expiry = token.Expiry
	stateObj.Spotify.TokenType = token.TokenType
	_ = state.SaveState(stateObj)
	return nil
}

func GetPlaylist(playlistId string) (*spotifylib.FullPlaylist, error) {
	playlist, err := getSpotifyClient().GetPlaylist(context.TODO(),
		spotifylib.ID(playlistId),
		spotifylib.Fields("images,id,name,owner"))
	if err != nil {
		return nil, err
	}

	return playlist, nil
}

func GetUserPlaylists() ([]spotifylib.SimplePlaylist, error) {
	client := getSpotifyClient()
	total := 20
	offset := 0
	items := make([]spotifylib.SimplePlaylist, 0)
	for total != len(items) {
		playlistPage, err := client.CurrentUsersPlaylists(context.TODO(),
			spotifylib.Limit(50), spotifylib.Offset(offset))
		if err != nil {
			return nil, err
		}

		items = append(items, playlistPage.Playlists...)
		total = playlistPage.Total
		offset += 50
	}

	return items, nil
}

func GetPlaylistTracks(playlistId string) ([]spotifylib.PlaylistItem, error) {
	client := getSpotifyClient()
	total := 1
	offset := 0
	items := make([]spotifylib.PlaylistItem, 0)
	for total != len(items) {
		tracksPage, err := client.GetPlaylistItems(context.TODO(),
			spotifylib.ID(playlistId),
			spotifylib.Limit(100),
			spotifylib.Offset(offset),
			spotifylib.Fields("items(added_at,track(album(artists(id,name),id,name,release_date,total_tracks),"+
				"artists(id,name),duration_ms,external_ids,id,name,type)),total"))
		if err != nil {
			return nil, err
		}

		for _, playlistitem := range tracksPage.Items {
			items = append(items, playlistitem)
		}

		total = tracksPage.Total
		offset += 100
	}
	return items, nil
}

func GetTrack(trackId string) (*spotifylib.FullTrack, error) {
	track, err := getSpotifyClient().GetTrack(context.TODO(), spotifylib.ID(trackId))
	if err != nil {
		return nil, err
	}

	return track, nil
}

func GetMe() (*spotifylib.PrivateUser, error) {
	profile, err := getSpotifyClient().CurrentUser(context.TODO())
	if err != nil {
		return nil, err
	}

	return profile, nil
}

func getSpotifyClient() *spotifylib.Client {
	stateObj, _ := state.GetState()
	token := oauth2.Token{
		AccessToken:  stateObj.Spotify.AccessToken,
		RefreshToken: stateObj.Spotify.RefreshToken,
		TokenType:    stateObj.Spotify.TokenType,
		Expiry:       stateObj.Spotify.Expiry,
	}
	httpClient := getAuthenticator().Client(context.TODO(), &token)
	return spotifylib.New(httpClient, spotifylib.WithRetry(true))
}
