package syncer

import (
	"api/internal/applemusic"
	"api/internal/configuration"
	"api/internal/spotify"
	"api/internal/state"
	"go.uber.org/zap"
	"time"
)

func SyncPlaylist(playlistId string) error {
	logger := getLogger()
	logger.Debug("going to sync: " + playlistId)

	logger.Debug("set syncing status")
	stateObj, err := state.GetState()
	if err != nil {
		return err
	}

	stateObj.Syncing = true
	err = state.SaveState(stateObj)
	if err != nil {
		return err
	}

	SendToAll(StatusMessage{Syncing: true})

	defer func(finalState state.State) {
		finalState.Syncing = false
		_ = state.SaveState(finalState)
	}(stateObj)

	spotifyPlaylist, err := spotify.GetPlaylist(playlistId)
	if err != nil {
		return err
	}
	logger.Debug("found spotify playlist: " + spotifyPlaylist.Name)

	logger.Debug("going to retrieve tracks")
	tracks, err := spotify.GetPlaylistTracks(playlistId)
	if err != nil {
		return err
	}
	logger.Debug("got tracks. amount: ", len(tracks))

	applemusicPlaylist, err := applemusic.GetSyncedPlaylist(playlistId)
	if err != nil {
		return err
	}

	if applemusicPlaylist == nil {
		// Create the AM playlist
		applemusicPlaylist, err = applemusic.CreatePlaylist(spotifyPlaylist.Name, playlistId)
		if err != nil {
			return err
		}
		logger.Debug("created Apple Music playlist with name: " + applemusicPlaylist.Attributes.Name)
	} else {
		logger.Debug("playlist already exists")
	}

	// Find each Spotify track on AM
	playlistSyncState := stateObj.Playlists[playlistId]
	lastPlaylistSyncDate := playlistSyncState.LastSyncDate

	applemusicTrackIds := make([]string, 0)
	logger.Debug("going to search for tracks on Apple Music")
	for _, spotifyTrack := range tracks {
		addedAt, _ := time.Parse(time.RFC3339, spotifyTrack.AddedAt)
		track := spotifyTrack.Track.Track
		if addedAt.After(lastPlaylistSyncDate) {
			logger.Debug("searching for ", spotify.GetArtistNames(track), " - "+track.Name)
			amTrack, err := applemusic.FindTrack(track)
			if err != nil {
				return err
			} else if amTrack != nil {
				logger.Debug("found track on Apple Music: ", amTrack.Attributes.ArtistName, " - ", amTrack.Attributes.Name)
				applemusicTrackIds = append(applemusicTrackIds, amTrack.Id)
			} else if amTrack == nil {
				logger.Warn("could not find track: " + track.Name)
			}
			logger.Debug("----------------------------------------------------------------------")
		} else {
			logger.Debug("skipped ", spotify.GetArtistNames(track), " - "+track.Name+" because of addedAt")
		}
	}
	logger.Debug("got Apple Music tracks. Amount: ", len(applemusicTrackIds))

	if len(applemusicTrackIds) > 0 {
		err = applemusic.AddTracksToPlaylist(applemusicPlaylist.Id, applemusicTrackIds)
		if err != nil {
			return err
		}
		logger.Debug("added tracks to playlist. amount: ", len(applemusicTrackIds))
	}

	if stateObj.Playlists == nil {
		stateObj.Playlists = make(map[string]state.PlaylistState)
	}

	stateObj.Playlists[playlistId] = state.PlaylistState{LastSyncDate: time.Now()}
	err = state.SaveState(stateObj)
	if err != nil {
		return err
	}
	logger.Debug("setting last sync date to: ", time.Now())

	SendToAll(StatusMessage{Syncing: false})

	return nil
}

func SyncAllPlaylists() {
	logger := getLogger()
	config, _ := configuration.GetConfiguration()
	playlistIds := config.Spotify.PlaylistIds
	for _, playlistId := range playlistIds {
		err := SyncPlaylist(playlistId)
		if err != nil {
			logger.Error("error while syncing playlistId '", playlistId, "' : ", err.Error())
		}
	}
}

func getLogger() *zap.SugaredLogger {
	logger, _ := zap.NewDevelopment()
	sugar := logger.Sugar()
	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger)
	return sugar
}
