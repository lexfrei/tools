package main

import (
	"log/slog"
	"sync"
	"time"
)

var (
	// Global runtime metrics store.
	runtimeMetrics *RuntimeMetrics
	runtimeMutex   sync.RWMutex
)

// initRuntimeMetrics initializes the runtime metrics system.
func initRuntimeMetrics() {
	runtimeMutex.Lock()
	defer runtimeMutex.Unlock()

	runtimeMetrics = &RuntimeMetrics{
		Players: make(map[string]*PlayerMetrics),
	}

	slog.Info("Runtime metrics system initialized")
}

// getPlayerMetrics returns metrics for a specific player.
func getPlayerMetrics(battleTag string) (*PlayerMetrics, bool) {
	runtimeMutex.RLock()
	defer runtimeMutex.RUnlock()

	if runtimeMetrics == nil {
		return nil, false
	}

	player, exists := runtimeMetrics.Players[battleTag]

	return player, exists
}

// setPlayerMetrics stores metrics for a specific player.
func setPlayerMetrics(battleTag string, metrics *PlayerMetrics) {
	runtimeMutex.Lock()
	defer runtimeMutex.Unlock()

	if runtimeMetrics == nil {
		runtimeMetrics = &RuntimeMetrics{
			Players: make(map[string]*PlayerMetrics),
		}
	}

	runtimeMetrics.Players[battleTag] = metrics
	slog.Info("Updated runtime metrics for player", "battletag", battleTag)
}

// getAllPlayerMetrics returns all player metrics.
func getAllPlayerMetrics() map[string]*PlayerMetrics {
	runtimeMutex.RLock()
	defer runtimeMutex.RUnlock()

	if runtimeMetrics == nil {
		return make(map[string]*PlayerMetrics)
	}

	// Create a copy to avoid race conditions
	result := make(map[string]*PlayerMetrics)
	for k, v := range runtimeMetrics.Players {
		result[k] = v
	}

	return result
}

// removePlayerMetrics removes metrics for a specific player.
func removePlayerMetrics(battleTag string) bool {
	runtimeMutex.Lock()
	defer runtimeMutex.Unlock()

	if runtimeMetrics == nil {
		return false
	}

	if _, exists := runtimeMetrics.Players[battleTag]; exists {
		delete(runtimeMetrics.Players, battleTag)
		slog.Info("Removed runtime metrics for player", "battletag", battleTag)

		return true
	}

	return false
}

// clearAllMetrics clears all runtime metrics.
func clearAllMetrics() {
	runtimeMutex.Lock()
	defer runtimeMutex.Unlock()

	if runtimeMetrics != nil {
		count := len(runtimeMetrics.Players)
		runtimeMetrics.Players = make(map[string]*PlayerMetrics)
		slog.Info("Cleared all runtime metrics", "cleared_count", count)
	}
}

// getMetricsStats returns statistics about the runtime metrics store.
func getMetricsStats() map[string]interface{} {
	runtimeMutex.RLock()
	defer runtimeMutex.RUnlock()

	stats := map[string]interface{}{
		"total_players": 0,
		"last_updated":  nil,
	}

	if runtimeMetrics == nil {
		return stats
	}

	stats["total_players"] = len(runtimeMetrics.Players)

	// Find most recent update time
	var mostRecent time.Time
	for _, player := range runtimeMetrics.Players {
		if player.LastUpdated.After(mostRecent) {
			mostRecent = player.LastUpdated
		}
	}

	if !mostRecent.IsZero() {
		stats["last_updated"] = mostRecent
	}

	return stats
}

// createPlayerMetrics creates a new PlayerMetrics from a FullPlayerProfile.
func createPlayerMetrics(battleTag string, profile *FullPlayerProfile) *PlayerMetrics {
	playerMetrics := &PlayerMetrics{
		BattleTag:        battleTag,
		DisplayName:      profile.BattleTag,
		PlayerTitle:      profile.PlayerTitle,
		LastUpdated:      time.Now(),
		ProfileMetrics:   profile.ProfileMetrics,
		AllHeroesMetrics: make(map[Platform]map[GameMode]AllHeroesStats),
		HeroMetrics:      make(map[Platform]map[GameMode]map[string]HeroMetrics),
	}

	// Convert platform data to runtime format
	for platformKey, platformStats := range profile.Platforms {
		// Initialize platform maps
		playerMetrics.AllHeroesMetrics[platformKey] = make(map[GameMode]AllHeroesStats)
		playerMetrics.HeroMetrics[platformKey] = make(map[GameMode]map[string]HeroMetrics)

		for gameModeKey, gameModeStats := range platformStats.GameModes {
			// Initialize gamemode maps
			playerMetrics.HeroMetrics[platformKey][gameModeKey] = make(map[string]HeroMetrics)

			// Convert hero metrics
			for heroID, heroStats := range gameModeStats.Heroes {
				heroMetrics := make(HeroMetrics)
				for key, value := range heroStats.Metrics {
					heroMetrics[key] = value
				}
				playerMetrics.HeroMetrics[platformKey][gameModeKey][heroID] = heroMetrics
			}

			// Store AllHeroesStats from parsed data
			playerMetrics.AllHeroesMetrics[platformKey][gameModeKey] = gameModeStats.AllHeroesStats
		}
	}

	return playerMetrics
}

// updatePlayerFromProfile updates existing PlayerMetrics with new profile data.
func updatePlayerFromProfile(existing *PlayerMetrics, profile *FullPlayerProfile) {
	existing.DisplayName = profile.BattleTag
	existing.PlayerTitle = profile.PlayerTitle
	existing.LastUpdated = time.Now()
	existing.ProfileMetrics = profile.ProfileMetrics

	// Clear existing hero metrics
	existing.HeroMetrics = make(map[Platform]map[GameMode]map[string]HeroMetrics)
	existing.AllHeroesMetrics = make(map[Platform]map[GameMode]AllHeroesStats)

	// Convert new platform data
	for platformKey, platformStats := range profile.Platforms {
		existing.AllHeroesMetrics[platformKey] = make(map[GameMode]AllHeroesStats)
		existing.HeroMetrics[platformKey] = make(map[GameMode]map[string]HeroMetrics)

		for gameModeKey, gameModeStats := range platformStats.GameModes {
			existing.HeroMetrics[platformKey][gameModeKey] = make(map[string]HeroMetrics)

			for heroID, heroStats := range gameModeStats.Heroes {
				heroMetrics := make(HeroMetrics)
				for key, value := range heroStats.Metrics {
					heroMetrics[key] = value
				}
				existing.HeroMetrics[platformKey][gameModeKey][heroID] = heroMetrics
			}

			// Store AllHeroesStats from parsed data
			existing.AllHeroesMetrics[platformKey][gameModeKey] = gameModeStats.AllHeroesStats
		}
	}
}

// listPlayerBattleTags returns a list of all player BattleTags in the runtime store.
func listPlayerBattleTags() []string {
	runtimeMutex.RLock()
	defer runtimeMutex.RUnlock()

	if runtimeMetrics == nil {
		return []string{}
	}

	battleTags := make([]string, 0, len(runtimeMetrics.Players))
	for battleTag := range runtimeMetrics.Players {
		battleTags = append(battleTags, battleTag)
	}

	return battleTags
}
