package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// ServerReadTimeout defines the server read timeout duration.
const ServerReadTimeout = 15 * time.Second

// ServerWriteTimeout defines the server write timeout duration.
const ServerWriteTimeout = 15 * time.Second

// ServerIdleTimeout defines the server idle timeout duration.
const ServerIdleTimeout = 60 * time.Second

// TimePlayedSecondsMetric defines the metric name for time played in seconds.
const TimePlayedSecondsMetric = "time_played_seconds"

// PrometheusMetrics contains all Prometheus metrics.
type PrometheusMetrics struct {
	// Profile-level metrics
	playerEndorsementLevel *prometheus.GaugeVec
	playerSkillRating      *prometheus.GaugeVec

	// All Heroes aggregated metrics
	playerTimePlayed             *prometheus.GaugeVec
	playerGamesWon               *prometheus.GaugeVec
	playerWinPercentage          *prometheus.GaugeVec
	playerWeaponAccuracy         *prometheus.GaugeVec
	playerEliminationsPerLife    *prometheus.GaugeVec
	playerKillStreakBest         *prometheus.GaugeVec
	playerMultikillBest          *prometheus.GaugeVec
	playerEliminationsPer10Min   *prometheus.GaugeVec
	playerDeathsPer10Min         *prometheus.GaugeVec
	playerFinalBlowsPer10Min     *prometheus.GaugeVec
	playerSoloKillsPer10Min      *prometheus.GaugeVec
	playerObjectiveKillsPer10Min *prometheus.GaugeVec
	playerObjectiveTimePer10Min  *prometheus.GaugeVec
	playerHeroDamagePer10Min     *prometheus.GaugeVec
	playerHealingPer10Min        *prometheus.GaugeVec

	// Hero-specific metrics
	heroTimePlayed             *prometheus.GaugeVec
	heroGamesWon               *prometheus.GaugeVec
	heroWinPercentage          *prometheus.GaugeVec
	heroWeaponAccuracy         *prometheus.GaugeVec
	heroEliminationsPerLife    *prometheus.GaugeVec
	heroKillStreakBest         *prometheus.GaugeVec
	heroMultikillBest          *prometheus.GaugeVec
	heroEliminationsPer10Min   *prometheus.GaugeVec
	heroDeathsPer10Min         *prometheus.GaugeVec
	heroFinalBlowsPer10Min     *prometheus.GaugeVec
	heroSoloKillsPer10Min      *prometheus.GaugeVec
	heroObjectiveKillsPer10Min *prometheus.GaugeVec
	heroObjectiveTimePer10Min  *prometheus.GaugeVec
	heroHeroDamagePer10Min     *prometheus.GaugeVec
	heroHealingPer10Min        *prometheus.GaugeVec

	// Dynamic detailed hero metrics - stores all hero-specific metrics by their Prometheus name
	detailedHeroMetrics map[string]*prometheus.GaugeVec
}

var prometheusMetrics *PrometheusMetrics

// initPrometheusMetrics initializes all Prometheus metrics.
func initPrometheusMetrics() {
	labels := createMetricLabels()
	prometheusMetrics = &PrometheusMetrics{}

	initProfileMetrics(labels)
	initAllHeroesMetrics(labels)
	initHeroSpecificMetrics(labels)

	registerAllMetrics()

	slog.Info("Prometheus metrics initialized")
}

// createMetricLabels defines label sets for different metric types.
func createMetricLabels() map[string][]string {
	return map[string][]string{
		"profile":     {"battletag", "platform", "gamemode"},
		"allHeroes":   {"battletag", "platform", "gamemode"},
		"hero":        {"battletag", "hero", "platform", "gamemode"},
		"skillRating": {"battletag", "platform", "gamemode", "rank_tier"},
	}
}

// initProfileMetrics initializes profile-level metrics.
func initProfileMetrics(labels map[string][]string) {
	prometheusMetrics.playerEndorsementLevel = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ow_player_endorsement_level",
			Help: "Player endorsement level",
		}, labels["profile"])

	prometheusMetrics.playerSkillRating = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ow_player_skill_rating",
			Help: "Player skill rating (SR) for competitive modes",
		}, labels["skillRating"])
}

// MetricDefinition defines a metric field and its configuration.
type MetricDefinition struct {
	Field **prometheus.GaugeVec
	Name  string
	Help  string
}

// GameplayMetricType represents different types of gameplay metrics.
type GameplayMetricType struct {
	PlayerField **prometheus.GaugeVec
	HeroField   **prometheus.GaugeVec
	BaseName    string
	Description string
}

// getGameplayMetricTypes returns all gameplay metrics with their fields and descriptions.
func getGameplayMetricTypes() []GameplayMetricType {
	basicMetrics := getBasicGameplayMetrics()
	combatMetrics := getCombatGameplayMetrics()
	per10MinMetrics := getPer10MinGameplayMetrics()

	result := make([]GameplayMetricType, 0, len(basicMetrics)+len(combatMetrics)+len(per10MinMetrics))
	result = append(result, basicMetrics...)
	result = append(result, combatMetrics...)
	result = append(result, per10MinMetrics...)

	return result
}

// getBasicGameplayMetrics returns basic gameplay metrics.
func getBasicGameplayMetrics() []GameplayMetricType {
	return []GameplayMetricType{
		{
			&prometheusMetrics.playerTimePlayed, &prometheusMetrics.heroTimePlayed,
			TimePlayedSecondsMetric, "Time played",
		},
		{
			&prometheusMetrics.playerGamesWon, &prometheusMetrics.heroGamesWon,
			"games_won", "Games won",
		},
		{
			&prometheusMetrics.playerWinPercentage, &prometheusMetrics.heroWinPercentage,
			"win_percentage", "Win percentage",
		},
		{
			&prometheusMetrics.playerWeaponAccuracy, &prometheusMetrics.heroWeaponAccuracy,
			"weapon_accuracy_percent", "Weapon accuracy percentage",
		},
	}
}

// getCombatGameplayMetrics returns combat-related gameplay metrics.
func getCombatGameplayMetrics() []GameplayMetricType {
	return []GameplayMetricType{
		{
			&prometheusMetrics.playerEliminationsPerLife, &prometheusMetrics.heroEliminationsPerLife,
			"eliminations_per_life", "Eliminations per life",
		},
		{
			&prometheusMetrics.playerKillStreakBest, &prometheusMetrics.heroKillStreakBest,
			"kill_streak_best", "Best kill streak achieved",
		},
		{
			&prometheusMetrics.playerMultikillBest, &prometheusMetrics.heroMultikillBest,
			"multikill_best", "Best multikill achieved",
		},
	}
}

// getPer10MinGameplayMetrics returns per-10-minute gameplay metrics.
func getPer10MinGameplayMetrics() []GameplayMetricType {
	return []GameplayMetricType{
		{
			&prometheusMetrics.playerEliminationsPer10Min, &prometheusMetrics.heroEliminationsPer10Min,
			"eliminations_per_10min", "Eliminations per 10 minutes",
		},
		{
			&prometheusMetrics.playerDeathsPer10Min, &prometheusMetrics.heroDeathsPer10Min,
			"deaths_per_10min", "Deaths per 10 minutes",
		},
		{
			&prometheusMetrics.playerFinalBlowsPer10Min, &prometheusMetrics.heroFinalBlowsPer10Min,
			"final_blows_per_10min", "Final blows per 10 minutes",
		},
		{
			&prometheusMetrics.playerSoloKillsPer10Min, &prometheusMetrics.heroSoloKillsPer10Min,
			"solo_kills_per_10min", "Solo kills per 10 minutes",
		},
		{
			&prometheusMetrics.playerObjectiveKillsPer10Min, &prometheusMetrics.heroObjectiveKillsPer10Min,
			"objective_kills_per_10min", "Objective kills per 10 minutes",
		},
		{
			&prometheusMetrics.playerObjectiveTimePer10Min, &prometheusMetrics.heroObjectiveTimePer10Min,
			"objective_time_per_10min", "Objective time per 10 minutes",
		},
		{
			&prometheusMetrics.playerHeroDamagePer10Min, &prometheusMetrics.heroHeroDamagePer10Min,
			"hero_damage_per_10min", "Hero damage per 10 minutes",
		},
		{
			&prometheusMetrics.playerHealingPer10Min, &prometheusMetrics.heroHealingPer10Min,
			"healing_per_10min", "Healing done per 10 minutes",
		},
	}
}

// initAllHeroesMetrics initializes aggregated all-heroes metrics.
func initAllHeroesMetrics(labels map[string][]string) {
	for _, metric := range getGameplayMetricTypes() {
		*metric.PlayerField = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "ow_player_" + getPlayerMetricName(metric.BaseName),
				Help: getPlayerHelpText(metric.Description),
			}, labels["allHeroes"])
	}
}

// initHeroSpecificMetrics initializes hero-specific metrics.
func initHeroSpecificMetrics(labels map[string][]string) {
	for _, metric := range getGameplayMetricTypes() {
		*metric.HeroField = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "ow_hero_" + getHeroMetricName(metric.BaseName),
				Help: getHeroHelpText(metric.Description),
			}, labels["hero"])
	}

	// Initialize detailed hero metrics map
	prometheusMetrics.detailedHeroMetrics = make(map[string]*prometheus.GaugeVec)

	// Create dynamic metrics for all heroes based on their metric definitions
	createDetailedHeroMetrics(labels["hero"])
}

// createDetailedHeroMetrics creates all hero-specific metrics from hero metric definitions.
func createDetailedHeroMetrics(heroLabels []string) {
	processedMetrics := make(map[string]bool) // Track already created metrics

	// Iterate through all heroes in registry
	for heroID, heroFactory := range HeroMetricsRegistry {
		heroStruct := heroFactory()
		heroMetrics := GenerateMetricDefs(heroStruct)

		slog.Debug("Creating detailed metrics for hero", "hero_id", heroID, "metrics_count", len(heroMetrics))

		for _, metricDef := range heroMetrics {
			prometheusName := metricDef.PrometheusName
			if prometheusName == "" {
				continue
			}

			// Skip if metric already processed (metrics can be shared across heroes)
			if processedMetrics[prometheusName] {
				continue
			}

			// Create the Prometheus metric
			prometheusMetrics.detailedHeroMetrics[prometheusName] = prometheus.NewGaugeVec(
				prometheus.GaugeOpts{
					Name: prometheusName,
					Help: metricDef.Help,
				}, heroLabels)

			processedMetrics[prometheusName] = true

			slog.Debug("Created detailed hero metric",
				"metric_name", prometheusName,
				"help", metricDef.Help)
		}
	}

	slog.Info("Created detailed hero metrics", "total_metrics", len(prometheusMetrics.detailedHeroMetrics))
}

// getPlayerMetricName returns the metric name for player-level metrics.
func getPlayerMetricName(baseName string) string {
	if baseName == TimePlayedSecondsMetric {
		return "total_time_played_seconds"
	}
	if baseName == "games_won" {
		return "total_games_won"
	}
	if baseName == "win_percentage" {
		return "overall_win_percentage"
	}

	return baseName
}

// getHeroMetricName returns the metric name for hero-level metrics.
func getHeroMetricName(baseName string) string {
	if baseName == TimePlayedSecondsMetric {
		return baseName
	}

	if strings.Contains(baseName, "_per_10min") {
		return baseName + "_avg"
	}

	return baseName
}

// getPlayerHelpText returns help text for player metrics.
func getPlayerHelpText(description string) string {
	return description + " across all heroes"
}

// getHeroHelpText returns help text for hero metrics.
func getHeroHelpText(description string) string {
	return description + " with specific hero"
}

// registerAllMetrics registers all metrics with Prometheus.
func registerAllMetrics() {
	prometheus.MustRegister(
		prometheusMetrics.playerEndorsementLevel,
		prometheusMetrics.playerSkillRating,
		prometheusMetrics.playerTimePlayed,
		prometheusMetrics.playerGamesWon,
		prometheusMetrics.playerWinPercentage,
		prometheusMetrics.playerWeaponAccuracy,
		prometheusMetrics.playerEliminationsPerLife,
		prometheusMetrics.playerKillStreakBest,
		prometheusMetrics.playerMultikillBest,
		prometheusMetrics.playerEliminationsPer10Min,
		prometheusMetrics.playerDeathsPer10Min,
		prometheusMetrics.playerFinalBlowsPer10Min,
		prometheusMetrics.playerSoloKillsPer10Min,
		prometheusMetrics.playerObjectiveKillsPer10Min,
		prometheusMetrics.playerObjectiveTimePer10Min,
		prometheusMetrics.playerHeroDamagePer10Min,
		prometheusMetrics.playerHealingPer10Min,
		prometheusMetrics.heroTimePlayed,
		prometheusMetrics.heroGamesWon,
		prometheusMetrics.heroWinPercentage,
		prometheusMetrics.heroWeaponAccuracy,
		prometheusMetrics.heroEliminationsPerLife,
		prometheusMetrics.heroKillStreakBest,
		prometheusMetrics.heroMultikillBest,
		prometheusMetrics.heroEliminationsPer10Min,
		prometheusMetrics.heroDeathsPer10Min,
		prometheusMetrics.heroFinalBlowsPer10Min,
		prometheusMetrics.heroSoloKillsPer10Min,
		prometheusMetrics.heroObjectiveKillsPer10Min,
		prometheusMetrics.heroObjectiveTimePer10Min,
		prometheusMetrics.heroHeroDamagePer10Min,
		prometheusMetrics.heroHealingPer10Min,
	)

	// Register all dynamic detailed hero metrics
	for _, metric := range prometheusMetrics.detailedHeroMetrics {
		prometheus.MustRegister(metric)
	}

	slog.Info("Registered all Prometheus metrics", "detailed_hero_metrics", len(prometheusMetrics.detailedHeroMetrics))
}

// updatePrometheusMetrics updates all Prometheus metrics from runtime data.
func updatePrometheusMetrics() {
	if prometheusMetrics == nil {
		slog.Error("Prometheus metrics not initialized")

		return
	}

	resetAllMetrics()
	allPlayers := getAllPlayerMetrics()
	updateAllPlayerMetrics(allPlayers)

	// Count total heroes for debug
	totalHeroes := 0
	for _, player := range allPlayers {
		for _, platforms := range player.HeroMetrics {
			for _, heroes := range platforms {
				totalHeroes += len(heroes)
			}
		}
	}

	slog.Info("Prometheus metrics updated",
		"players", len(allPlayers),
		"total_heroes", totalHeroes)
}

// resetAllMetrics clears all Prometheus metrics.
func resetAllMetrics() {
	resetPlayerMetrics()
	resetHeroMetrics()
}

// resetPlayerMetrics resets all player-level metrics.
func resetPlayerMetrics() {
	prometheusMetrics.playerEndorsementLevel.Reset()
	prometheusMetrics.playerSkillRating.Reset()
	prometheusMetrics.playerTimePlayed.Reset()
	prometheusMetrics.playerGamesWon.Reset()
	prometheusMetrics.playerWinPercentage.Reset()
	prometheusMetrics.playerWeaponAccuracy.Reset()
	prometheusMetrics.playerEliminationsPerLife.Reset()
	prometheusMetrics.playerKillStreakBest.Reset()
	prometheusMetrics.playerMultikillBest.Reset()
	prometheusMetrics.playerEliminationsPer10Min.Reset()
	prometheusMetrics.playerDeathsPer10Min.Reset()
	prometheusMetrics.playerFinalBlowsPer10Min.Reset()
	prometheusMetrics.playerSoloKillsPer10Min.Reset()
	prometheusMetrics.playerObjectiveKillsPer10Min.Reset()
	prometheusMetrics.playerObjectiveTimePer10Min.Reset()
	prometheusMetrics.playerHeroDamagePer10Min.Reset()
	prometheusMetrics.playerHealingPer10Min.Reset()
}

// resetHeroMetrics resets all hero-level metrics.
func resetHeroMetrics() {
	prometheusMetrics.heroTimePlayed.Reset()
	prometheusMetrics.heroGamesWon.Reset()
	prometheusMetrics.heroWinPercentage.Reset()
	prometheusMetrics.heroWeaponAccuracy.Reset()
	prometheusMetrics.heroEliminationsPerLife.Reset()
	prometheusMetrics.heroKillStreakBest.Reset()
	prometheusMetrics.heroMultikillBest.Reset()
	prometheusMetrics.heroEliminationsPer10Min.Reset()
	prometheusMetrics.heroDeathsPer10Min.Reset()
	prometheusMetrics.heroFinalBlowsPer10Min.Reset()
	prometheusMetrics.heroSoloKillsPer10Min.Reset()
	prometheusMetrics.heroObjectiveKillsPer10Min.Reset()
	prometheusMetrics.heroObjectiveTimePer10Min.Reset()
	prometheusMetrics.heroHeroDamagePer10Min.Reset()
	prometheusMetrics.heroHealingPer10Min.Reset()
}

// updateAllPlayerMetrics updates metrics for all players.
func updateAllPlayerMetrics(allPlayers map[string]*PlayerMetrics) {
	for battleTag, playerMetrics := range allPlayers {
		updateSinglePlayerMetrics(battleTag, playerMetrics)
	}
}

// updateSinglePlayerMetrics updates metrics for a single player.
func updateSinglePlayerMetrics(battleTag string, playerMetrics *PlayerMetrics) {
	updateProfileMetrics(battleTag, playerMetrics)
	updateAllHeroesMetrics(battleTag, playerMetrics)
	updateHeroSpecificMetrics(battleTag, playerMetrics)
}


// updateProfileMetrics updates profile-level metrics.
func updateProfileMetrics(battleTag string, playerMetrics *PlayerMetrics) {
	prometheusMetrics.playerEndorsementLevel.WithLabelValues(
		battleTag, "all", "all").Set(float64(playerMetrics.ProfileMetrics.Endorsement.Level))

	updateSkillRatings(battleTag, playerMetrics)
}

// updateSkillRatings updates skill rating metrics.
func updateSkillRatings(battleTag string, playerMetrics *PlayerMetrics) {
	for platform, ranks := range playerMetrics.ProfileMetrics.SkillRatings {
		for tier, sr := range ranks {
			prometheusMetrics.playerSkillRating.WithLabelValues(
				battleTag, string(platform), "competitive", string(tier)).Set(float64(sr.SR))
		}
	}
}

// updateAllHeroesMetrics updates aggregated all-heroes metrics.
func updateAllHeroesMetrics(battleTag string, playerMetrics *PlayerMetrics) {
	for platform, gameModes := range playerMetrics.AllHeroesMetrics {
		for gameMode, allHeroesStats := range gameModes {
			baseLabels := createBaseLabels(battleTag, string(platform), string(gameMode))
			setAllHeroesMetricValues(baseLabels, &allHeroesStats)
		}
	}
}

// updateHeroSpecificMetrics updates hero-specific metrics.
func updateHeroSpecificMetrics(battleTag string, playerMetrics *PlayerMetrics) {
	totalHeroes := 0
	for platform, gameModes := range playerMetrics.HeroMetrics {
		for gameMode, heroes := range gameModes {
			slog.Debug("Processing heroes for platform/gamemode",
				"battletag", battleTag,
				"platform", platform,
				"gamemode", gameMode,
				"hero_count", len(heroes))

			for heroID, heroMetrics := range heroes {
				baseLabels := createHeroLabels(battleTag, heroID, string(platform), string(gameMode))
				setHeroMetricValues(baseLabels, heroMetrics)
				totalHeroes++
			}
		}
	}

	if totalHeroes > 0 {
		slog.Info("Updated hero-specific metrics",
			"battletag", battleTag,
			"heroes_processed", totalHeroes)
	} else {
		slog.Warn("No hero-specific metrics to update",
			"battletag", battleTag,
			"hero_metrics_empty", len(playerMetrics.HeroMetrics) == 0)
	}
}

// createBaseLabels creates base labels for metrics.
func createBaseLabels(battleTag, platform, gameMode string) prometheus.Labels {
	return prometheus.Labels{
		"battletag": battleTag,
		"platform":  platform,
		"gamemode":  gameMode,
	}
}

// createHeroLabels creates labels for hero-specific metrics.
func createHeroLabels(battleTag, heroID, platform, gameMode string) prometheus.Labels {
	return prometheus.Labels{
		"battletag": battleTag,
		"hero":      heroID,
		"platform":  platform,
		"gamemode":  gameMode,
	}
}

// setAllHeroesMetricValues sets all heroes metric values.
func setAllHeroesMetricValues(labels prometheus.Labels, stats *AllHeroesStats) {
	prometheusMetrics.playerTimePlayed.With(labels).Set(float64(stats.TotalTimePlayed))
	prometheusMetrics.playerGamesWon.With(labels).Set(float64(stats.TotalGamesWon))
	prometheusMetrics.playerWinPercentage.With(labels).Set(stats.OverallWinPercentage)
	prometheusMetrics.playerWeaponAccuracy.With(labels).Set(stats.WeaponAccuracy)
	prometheusMetrics.playerEliminationsPerLife.With(labels).Set(stats.EliminationsPerLife)
	prometheusMetrics.playerKillStreakBest.With(labels).Set(float64(stats.KillStreakBest))
	prometheusMetrics.playerMultikillBest.With(labels).Set(float64(stats.MultikillBest))
	prometheusMetrics.playerEliminationsPer10Min.With(labels).Set(stats.EliminationsPer10Min)
	prometheusMetrics.playerDeathsPer10Min.With(labels).Set(stats.DeathsPer10Min)
	prometheusMetrics.playerFinalBlowsPer10Min.With(labels).Set(stats.FinalBlowsPer10Min)
	prometheusMetrics.playerSoloKillsPer10Min.With(labels).Set(stats.SoloKillsPer10Min)
	prometheusMetrics.playerObjectiveKillsPer10Min.With(labels).Set(stats.ObjectiveKillsPer10Min)
	prometheusMetrics.playerObjectiveTimePer10Min.With(labels).Set(stats.ObjectiveTimePer10Min)
	prometheusMetrics.playerHeroDamagePer10Min.With(labels).Set(stats.HeroDamagePer10Min)
	prometheusMetrics.playerHealingPer10Min.With(labels).Set(stats.HealingPer10Min)
}

// setHeroMetricValues sets hero-specific metric values.
func setHeroMetricValues(labels prometheus.Labels, heroMetrics map[string]interface{}) {
	heroID := labels["hero"]
	battleTag := labels["battletag"]
	metricsCount := 0

	for metricKey, metricValue := range heroMetrics {
		switch value := metricValue.(type) {
		case float64:
			updateHeroMetric(metricKey, value, labels)
			metricsCount++
		case int:
			updateHeroMetric(metricKey, float64(value), labels)
			metricsCount++
		case int64:
			updateHeroMetric(metricKey, float64(value), labels)
			metricsCount++
		default:
			slog.Debug("Skipped unsupported metric value type",
				"battletag", battleTag,
				"hero", heroID,
				"metric", metricKey,
				"type", fmt.Sprintf("%T", value))
		}
	}

	if metricsCount > 0 {
		slog.Debug("Set hero metrics",
			"battletag", battleTag,
			"hero", heroID,
			"metrics_count", metricsCount)
	}
}

// updateHeroMetric updates a specific hero metric based on the metric key.
func updateHeroMetric(metricKey string, value float64, labels prometheus.Labels) {
	// Try detailed hero metrics
	if tryDetailedMetricUpdate(metricKey, value, labels) {
		return
	}

	slog.Debug("No metric updater found",
		"battletag", labels["battletag"],
		"hero", labels["hero"],
		"metric", metricKey,
		"value", value)
}

// tryDetailedMetricUpdate attempts to update using detailed hero metrics.
func tryDetailedMetricUpdate(metricKey string, value float64, labels prometheus.Labels) bool {
	heroID, hasHero := labels["hero"]
	if !hasHero {
		slog.Debug("No hero in labels for metric", "metric", metricKey)

		return false
	}

	// Try direct prometheus name match first
	if updateByPrometheusName(metricKey, value, labels) {
		return true
	}

	// Try ow tag match
	return updateByOwTag(heroID, metricKey, value, labels)
}

// updateByPrometheusName tries to update by matching the prometheus metric name directly.
func updateByPrometheusName(prometheusName string, value float64, labels prometheus.Labels) bool {
	if dynamicMetric, exists := prometheusMetrics.detailedHeroMetrics[prometheusName]; exists {
		dynamicMetric.With(labels).Set(value)
		slog.Debug("Updated detailed hero metric",
			"battletag", labels["battletag"],
			"hero", labels["hero"],
			"prometheus_name", prometheusName,
			"value", value)

		return true
	}

	return false
}

// updateByOwTag tries to find and update a metric by matching the ow tag.
func updateByOwTag(heroID, metricKey string, value float64, labels prometheus.Labels) bool {
	heroMetrics := GetHeroMetrics(heroID)
	for taggedKey, taggedDef := range heroMetrics {
		if taggedKey == metricKey && taggedDef.PrometheusName != "" {
			if dynamicMetric, exists := prometheusMetrics.detailedHeroMetrics[taggedDef.PrometheusName]; exists {
				dynamicMetric.With(labels).Set(value)
				slog.Debug("Updated detailed hero metric via ow tag",
					"battletag", labels["battletag"],
					"hero", labels["hero"],
					"metric_key", metricKey,
					"ow_tag", taggedKey,
					"prometheus_name", taggedDef.PrometheusName,
					"value", value)

				return true
			}
		}
	}

	return false
}

// startPrometheusServer starts the HTTP server for Prometheus metrics.
func startPrometheusServer(port string) {

	// Create metrics handler that updates data before serving
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		updatePrometheusMetrics()
		promhttp.Handler().ServeHTTP(w, r)
	})

	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("OK"))
		if err != nil {
			slog.Error("Failed to write health check response", "error", err)
		}
	})

	slog.Info("Starting Prometheus metrics server", "port", port)
	server := &http.Server{
		Addr:         ":" + port,
		ReadTimeout:  ServerReadTimeout,
		WriteTimeout: ServerWriteTimeout,
		IdleTimeout:  ServerIdleTimeout,
	}
	err := server.ListenAndServe()
	if err != nil {
		slog.Error("Failed to start metrics server", "error", err)
	}
}
