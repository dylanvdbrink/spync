package configuration

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type SpotifyConfig struct {
	ClientId     string `json:"client-id"`
	ClientSecret string `json:"client-secret"`
}

type AppleMusicConfig struct {
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
	AppleMusic:   AppleMusicConfig{ClientId: "", ClientSecret: ""},
	Spotify:      SpotifyConfig{ClientId: "", ClientSecret: ""},
}

const configFilePath = "configuration.json"

func GetConfiguration() (Config, error) {
	var config Config

	file, fileErr := getOrCreateConfigurationFile()
	if fileErr != nil {
		fmt.Println("error getting file: " + fileErr.Error())
		return Config{}, fileErr
	}

	decoder := json.NewDecoder(file)
	decodeErr := decoder.Decode(&config)

	if decodeErr != nil {
		return Config{}, decodeErr
	}

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
			return nil, err
		}

		return newFile, nil
	} else {
		file, openErr := os.OpenFile(configFilePath, os.O_RDWR, 0644)
		if openErr != nil {
			fmt.Println("error opening file: " + openErr.Error())
			return nil, openErr
		}
		return file, nil
	}
}

func SaveConfiguration(config Config) error {
	file, fileErr := getOrCreateConfigurationFile()
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	if fileErr != nil {
		fmt.Println("could not get or create configfile: " + fileErr.Error())
		return fileErr
	}

	// Replace properties that should not be touched
	existingConfig, getErr := GetConfiguration()
	if getErr != nil {
		return getErr
	}
	config.Syncing = existingConfig.Syncing

	err := writeConfiguration(file, config)
	if err != nil {
		fmt.Println("error while writing config: " + fileErr.Error())
		return err
	}

	return nil
}

func writeConfiguration(file *os.File, config Config) error {
	encoder := json.NewEncoder(file)
	err := encoder.Encode(config)
	if err != nil {
		fmt.Println("encoding error: " + err.Error())
		return err
	}
	return nil
}
