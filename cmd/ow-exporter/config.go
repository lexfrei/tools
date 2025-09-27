package main

import (
	"os"
	"path/filepath"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

const (
	// ConfigDirPermissions sets directory permissions for config directory (rwxr-xr-x).
	ConfigDirPermissions = 0o755
	// ConfigFilePermissions sets file permissions for config file (rw-------).
	ConfigFilePermissions = 0o600
)

// PlayersConfig represents the structure of players.yaml.
type PlayersConfig struct {
	Players []PlayerEntry `yaml:"players"`
}

// PlayerEntry represents a single player configuration entry.
type PlayerEntry struct {
	BattleTag    string     `yaml:"battletag"`
	ResolvedURL  string     `yaml:"resolvedUrl"`
	LastResolved *time.Time `yaml:"lastResolved"`
}

var (
	playersConfig *PlayersConfig
	configPath    string
)

// initConfig initializes the configuration system.
func initConfig() error {
	// Set config file path
	configPath = getConfigPath()

	// Create config directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	err := os.MkdirAll(configDir, ConfigDirPermissions)
	if err != nil {
		return errors.Wrap(err, "failed to create config directory")
	}

	// Initialize viper
	viper.SetConfigName("players")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configDir)

	// Try to read existing config
	err = viper.ReadInConfig()
	if err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			// Config file not found, create default
			return createDefaultConfig()
		}

		return errors.Wrap(err, "failed to read config")
	}

	// Load config into struct
	return loadConfig()
}

// getConfigPath returns the path to the players.yaml config file.
func getConfigPath() string {
	if configFile := viper.GetString("config"); configFile != "" {
		return configFile
	}

	return "config/players.yaml"
}

// createDefaultConfig creates a default players.yaml file.
func createDefaultConfig() error {
	defaultConfig := &PlayersConfig{
		Players: []PlayerEntry{},
	}

	return saveConfig(defaultConfig)
}

// loadConfig loads the configuration from file into memory.
func loadConfig() error {
	playersConfig = &PlayersConfig{}

	return errors.Wrap(viper.Unmarshal(playersConfig), "failed to unmarshal config")
}

// saveConfig saves the configuration to file.
func saveConfig(config *PlayersConfig) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return errors.Wrap(err, "failed to marshal config")
	}

	err = os.WriteFile(configPath, data, ConfigFilePermissions)
	if err != nil {
		return errors.Wrap(err, "failed to write config file")
	}

	// Update in-memory config
	playersConfig = config

	return nil
}

// findPlayerByBattleTag finds a player entry by BattleTag.
func findPlayerByBattleTag(battleTag string) *PlayerEntry {
	if playersConfig == nil {
		return nil
	}

	for i := range playersConfig.Players {
		if playersConfig.Players[i].BattleTag == battleTag {
			return &playersConfig.Players[i]
		}
	}

	return nil
}

// addPlayerToConfig adds a new player to the configuration.
func addPlayerToConfig(battleTag, resolvedURL string) error {
	if playersConfig == nil {
		playersConfig = &PlayersConfig{}
	}

	// Check if player already exists
	if existingPlayer := findPlayerByBattleTag(battleTag); existingPlayer != nil {
		// Update existing player
		existingPlayer.ResolvedURL = resolvedURL
		now := time.Now()
		existingPlayer.LastResolved = &now
	} else {
		// Add new player
		now := time.Now()
		newPlayer := PlayerEntry{
			BattleTag:    battleTag,
			ResolvedURL:  resolvedURL,
			LastResolved: &now,
		}
		playersConfig.Players = append(playersConfig.Players, newPlayer)
	}

	return saveConfig(playersConfig)
}

// updatePlayerURL updates the resolved URL for an existing player.
func updatePlayerURL(battleTag, resolvedURL string) error {
	player := findPlayerByBattleTag(battleTag)
	if player == nil {
		return errors.Wrapf(ErrPlayerNotFound, "%s", battleTag)
	}

	player.ResolvedURL = resolvedURL
	now := time.Now()
	player.LastResolved = &now

	return saveConfig(playersConfig)
}

// getAllPlayers returns all configured players.
func getAllPlayers() []PlayerEntry {
	if playersConfig == nil {
		return []PlayerEntry{}
	}

	return playersConfig.Players
}

// removePlayerFromConfig removes a player from the configuration.
func removePlayerFromConfig(battleTag string) error {
	if playersConfig == nil {
		return ErrNoConfigLoaded
	}

	for i, player := range playersConfig.Players {
		if player.BattleTag == battleTag {
			// Remove player from slice
			playersConfig.Players = append(playersConfig.Players[:i], playersConfig.Players[i+1:]...)

			return saveConfig(playersConfig)
		}
	}

	return errors.Wrapf(ErrPlayerNotFound, "%s", battleTag)
}
