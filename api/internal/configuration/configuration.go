package configuration

import (
	"api/internal/storage"
	"encoding/json"
	"go.uber.org/zap"
	"io"
	"os"
)

type SpotifyConfig struct {
	PlaylistIds []string `json:"playlist-ids"`

	ClientId     string `json:"client-id"`
	ClientSecret string `json:"client-secret"`
	RedirectUrl  string `json:"redirect-url"`
}

type AppleMusicConfig struct {
	DeveloperToken string `json:"developer-token"`
}

type Config struct {
	SyncInterval int              `json:"sync-interval"`
	AppleMusic   AppleMusicConfig `json:"apple-music"`
	Spotify      SpotifyConfig    `json:"spotify"`
}

const configFilePath = "configuration.json"

func GetConfiguration() (Config, error) {
	logger := getLogger()
	file, _, fileErr := storage.GetOrCreateFile(configFilePath)
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	if fileErr != nil {
		logger.Error("error getting file: " + fileErr.Error())
		return Config{}, fileErr
	}

	configContentBytes, _ := io.ReadAll(file)
	var config Config
	_ = json.Unmarshal(configContentBytes, &config)

	return config, nil
}

func SaveConfiguration(config Config, replaceProtected bool) error {
	logger := getLogger()
	if replaceProtected {
		// Replace properties that should not be touched
		replaceErr := replaceProtectedProperties(&config)
		if replaceErr != nil {
			logger.Error("error replacing properties: " + replaceErr.Error())
			return replaceErr
		}
	}

	file, _, fileErr := storage.GetOrCreateFile(configFilePath)
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

	err := writeConfiguration(file, config)
	if err != nil {
		logger.Error("error while writing config: " + fileErr.Error())
		return err
	}

	return nil
}

func replaceProtectedProperties(config *Config) error {
	logger := getLogger()
	existingConfig, getErr := GetConfiguration()
	if getErr != nil {
		logger.Error("error getting config: " + getErr.Error())
		return getErr
	}

	config.Spotify.PlaylistIds = existingConfig.Spotify.PlaylistIds

	return nil
}

func writeConfiguration(file *os.File, config Config) error {
	configBytes, _ := json.MarshalIndent(config, "", "  ")
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
