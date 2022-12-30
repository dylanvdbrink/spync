package applemusic

import (
	"api/internal/configuration"
	"api/internal/spotify"
	"api/internal/state"
	"context"
	applemusiclib "github.com/minchao/go-apple-music"
	spotifylib "github.com/zmb3/spotify/v2"
	"go.uber.org/zap"
	"regexp"
	"strings"
)

func SaveAuth(token string) error {
	stateObj, err := state.GetState()
	if err != nil {
		return err
	}
	stateObj.AppleMusic.AccessToken = token
	err = state.SaveState(stateObj)
	if err != nil {
		return err
	}

	return nil
}

func CreatePlaylist(name string, spotifyId string) (*applemusiclib.LibraryPlaylist, error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	playlist, _, err := client.Me.CreateLibraryPlaylist(context.TODO(), applemusiclib.CreateLibraryPlaylist{
		Attributes: applemusiclib.CreateLibraryPlaylistAttributes{
			Name:        "Spotify - " + name,
			Description: name + ". Synced with spync. (ID: " + spotifyId + ")",
		},
	}, nil)

	if err != nil {
		return nil, err
	}

	return &playlist.Data[0], nil
}

func AddTracksToPlaylist(playlistId string, trackIds []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	tracks := make([]applemusiclib.CreateLibraryPlaylistTrack, 0)
	for _, trackId := range trackIds {
		tracks = append(tracks, applemusiclib.CreateLibraryPlaylistTrack{Id: trackId, Type: "songs"})
	}

	_, err = client.Me.AddLibraryTracksToPlaylist(context.TODO(), playlistId, applemusiclib.CreateLibraryPlaylistTrackData{Data: tracks})
	if err != nil {
		return err
	}

	return nil
}

func FindTrack(item *spotifylib.FullTrack) (*applemusiclib.Song, error) {
	logger := getLogger()

	client, err := getClient()
	if err != nil {
		return nil, err
	}

	// Remove "feat." because Apple Music often does not use it in song titles
	var removeRegex = regexp.MustCompile(`(\((?:feat.|ft.|with) .+\))`)

	search, _, err := client.Catalog.Search(context.TODO(), "NL", &applemusiclib.SearchOptions{
		Offset: 0,
		Limit:  25,
		Types:  "songs",
		Term:   item.Artists[0].Name + " - " + removeRegex.ReplaceAllString(item.Name, ""),
	})
	if err != nil {
		return nil, err
	}

	if search.Results == (applemusiclib.SearchResults{}) {
		return nil, nil
	} else {
		// Check for ISRC
		for _, track := range search.Results.Songs.Data {
			if item.ExternalIDs["isrc"] == track.Attributes.ISRC {
				// If isrc matches, it is the correct match
				logger.Debug("ISRC match")
				return &track, nil
			}
		}

		// Check if joined artistnames contains the artistname
		for _, track := range search.Results.Songs.Data {
			artistsString := spotify.GetArtistNames(item)
			if strings.Contains(artistsString, track.Attributes.ArtistName) {
				logger.Debug("artistname was a rough match")
				return &track, nil
			}
		}

		logger.Debug("returning first entry")
		return &search.Results.Songs.Data[0], nil
	}
}

func GetSpotifyPlaylists() ([]applemusiclib.LibraryPlaylist, error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}
	offset := 0
	hasMore := true
	items := make([]applemusiclib.LibraryPlaylist, 0)
	for hasMore {
		playlists, _, err := client.Me.GetAllLibraryPlaylists(context.TODO(), &applemusiclib.PageOptions{
			Offset: offset,
			Limit:  100,
		})
		if err != nil {
			return nil, err
		}

		for _, playlist := range playlists.Data {
			if strings.HasPrefix(playlist.Attributes.Name, "Spotify - ") {
				items = append(items, playlist)
			}
		}

		if len(playlists.Data) < 100 {
			hasMore = false
		} else {
			offset += 100
		}
	}

	return items, nil
}

func GetSyncedPlaylist(playlistId string) (*applemusiclib.LibraryPlaylist, error) {
	playlists, err := GetSpotifyPlaylists()
	if err != nil {
		return nil, err
	}

	idRegex := regexp.MustCompile(`\(ID: (.*?)\)`)
	for _, playlist := range playlists {
		description := playlist.Attributes.Description.Standard
		matches := idRegex.FindStringSubmatch(description)
		if playlistId == matches[1] {
			return &playlist, nil
		}
	}
	return nil, nil
}

func CheckAuth() error {
	logger := getLogger()

	client, err := getClient()
	if err != nil {
		logger.Error("getClient error: " + err.Error())
		return err
	}

	logger.Debug("retrieving storefront")
	_, _, err = client.Me.GetStorefront(context.TODO(), nil)
	if err != nil {
		return err
	}

	return nil
}

func getClient() (*applemusiclib.Client, error) {
	config, err := configuration.GetConfiguration()
	if err != nil {
		return nil, err
	}
	stateObj, err := state.GetState()
	if err != nil {
		return nil, err
	}
	tp := applemusiclib.Transport{Token: config.AppleMusic.DeveloperToken, MusicUserToken: stateObj.AppleMusic.AccessToken}
	client := applemusiclib.NewClient(tp.Client())

	return client, nil
}

func getLogger() *zap.SugaredLogger {
	logger, _ := zap.NewDevelopment()
	sugar := logger.Sugar()
	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger)
	return sugar
}