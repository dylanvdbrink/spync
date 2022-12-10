package configuration

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

type SpotifyConfig struct {
	AccessToken  string    `json:"access-token"`
	RefreshToken string    `json:"refresh-token"`
	TokenType    string    `json:"token-type"`
	Expiry       time.Time `json:"expiry"`

	PlaylistIds []string `json:"playlist-ids"`

	ClientId     string `json:"client-id"`
	ClientSecret string `json:"client-secret"`
	RedirectUrl  string `json:"redirect-url"`
}

type AppleMusicConfig struct {
	AccessToken  string `json:"access-token"`
	RefreshToken string `json:"refresh-token"`
	ClientId     string `json:"client-id"`
	ClientSecret string `json:"client-secret"`
}

type Config struct {
	Syncing      bool             `json:"syncing"`
	SyncInterval int              `json:"sync-interval"`
	AppleMusic   AppleMusicConfig `json:"apple-music"`
	Spotify      SpotifyConfig    `json:"spotify"`
}

var defaultConfig = Config{
	Syncing:      false,
	SyncInterval: 4320, // 3 days
	AppleMusic:   AppleMusicConfig{},
	Spotify:      SpotifyConfig{},
}

const configFilePath = "configuration.json"

func GetConfiguration() (Config, error) {
	file, fileErr := getOrCreateConfigurationFile()
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	if fileErr != nil {
		fmt.Println("error getting file: " + fileErr.Error())
		return Config{}, fileErr
	}

	configContentBytes, _ := io.ReadAll(file)
	var config Config
	_ = json.Unmarshal(configContentBytes, &config)

	return config, nil
}

func getOrCreateConfigurationFile() (*os.File, error) {
	if _, existsErr := os.Stat(configFilePath); errors.Is(existsErr, os.ErrNotExist) {
		newFile, createErr := os.Create(configFilePath)
		if createErr != nil {
			fmt.Println("file creation error: " + createErr.Error())
			return nil, createErr
		}

		err := writeConfiguration(newFile, defaultConfig)
		if err != nil {
			fmt.Println("writeConfiguration error: " + err.Error())
			return nil, err
		}

		return newFile, nil
	} else {
		file, openErr := os.OpenFile(configFilePath, os.O_RDWR, 0755)
		if openErr != nil {
			fmt.Println("error opening file: " + openErr.Error())
			return nil, openErr
		}
		return file, nil
	}
}

func SaveConfiguration(config Config) error {
	fmt.Println("> SaveConfiguration")
	// Replace properties that should not be touched
	replaceErr := replaceProtectedProperties(config)
	if replaceErr != nil {
		fmt.Println("error replacing properties: " + replaceErr.Error())
		return replaceErr
	}

	file, fileErr := getOrCreateConfigurationFile()
	defer func(file *os.File) {
		closeErr := file.Close()
		if closeErr != nil {
			fmt.Println("error closing file: " + closeErr.Error())
		}
	}(file)
	if fileErr != nil {
		fmt.Println("could not get or create configfile: " + fileErr.Error())
		return fileErr
	}

	err := writeConfiguration(file, config)
	if err != nil {
		fmt.Println("error while writing config: " + fileErr.Error())
		return err
	}

	fmt.Println("< SaveConfiguration")

	return nil
}

func replaceProtectedProperties(config Config) error {
	existingConfig, getErr := GetConfiguration()
	if getErr != nil {
		fmt.Println("error getting config: " + getErr.Error())
		return getErr
	}
	config.Syncing = existingConfig.Syncing
	config.Spotify.AccessToken = existingConfig.Spotify.AccessToken
	config.Spotify.RefreshToken = existingConfig.Spotify.RefreshToken
	config.Spotify.PlaylistIds = existingConfig.Spotify.PlaylistIds

	return nil
}

func writeConfiguration(file *os.File, config Config) error {
	configBytes, _ := json.MarshalIndent(config, "", "  ")
	_ = file.Truncate(0)
	_, _ = file.Write(configBytes)

	return nil
}
