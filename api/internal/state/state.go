package state

import (
	"api/internal/storage"
	"encoding/json"
	"go.uber.org/zap"
	"io"
	"os"
	"time"
)

type SpotifyState struct {
	AccessToken  string    `json:"access-token"`
	RefreshToken string    `json:"refresh-token"`
	TokenType    string    `json:"token-type"`
	Expiry       time.Time `json:"expiry"`
}

type AppleMusicState struct {
	AccessToken string `json:"access-token"`
}

type PlaylistState struct {
	LastSyncDate time.Time `json:"last-sync-date"`
}

type State struct {
	Syncing    bool                     `json:"syncing"`
	AppleMusic AppleMusicState          `json:"apple-music"`
	Spotify    SpotifyState             `json:"spotify"`
	Playlists  map[string]PlaylistState `json:"playlists"`
}

const stateFilePath = "state.json"

func GetState() (State, error) {
	logger := getLogger()
	file, _, fileErr := storage.GetOrCreateFile(stateFilePath)
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	if fileErr != nil {
		logger.Error("error getting file: " + fileErr.Error())
		return State{}, fileErr
	}

	configContentBytes, _ := io.ReadAll(file)
	var state State
	_ = json.Unmarshal(configContentBytes, &state)

	return state, nil
}

func SaveState(state State) error {
	logger := getLogger()
	file, _, fileErr := storage.GetOrCreateFile(stateFilePath)
	defer func(file *os.File) {
		closeErr := file.Close()
		if closeErr != nil {
			logger.Error("error closing file: " + closeErr.Error())
		}
	}(file)
	if fileErr != nil {
		logger.Error("could not get or create configfile: " + fileErr.Error())
		return fileErr
	}

	err := writeState(file, state)
	if err != nil {
		logger.Error("error while writing config: " + fileErr.Error())
		return err
	}

	return nil
}

func writeState(file *os.File, state State) error {
	configBytes, _ := json.MarshalIndent(state, "", "  ")
	_ = file.Truncate(0)
	_, _ = file.Write(configBytes)

	return nil
}

func getLogger() *zap.SugaredLogger {
	logger, _ := zap.NewDevelopment()
	sugar := logger.Sugar()
	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger)
	return sugar
}
