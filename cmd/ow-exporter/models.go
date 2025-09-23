// metrics_mapping.go - Hardcoded metrics mapping for OW2 exporter
package main

// MetricDef defines a single metric with its parsing information
type MetricDef struct {
	PrometheusName string `json:"prometheus_name"` // "ow_hero_time_played_seconds"
	Help           string `json:"help"`            // Help text for Prometheus
	Unit           string `json:"unit"`            // "seconds", "percent", "count", "ratio"
	Selector       string `json:"selector"`        // CSS selector or data attribute
	HexID          string `json:"hex_id"`          // Blizzard's hex ID for the metric
	ValueType      string `json:"value_type"`      // "duration", "number", "percentage"
}

// Platform represents PC or Console
type Platform string

const (
	PlatformPC      Platform = "pc"
	PlatformConsole Platform = "console"
)

// GameMode represents Quick Play or Competitive
type GameMode string

const (
	GameModeQuickPlay    GameMode = "quickplay"
	GameModeCompetitive  GameMode = "competitive"
)

// MetricLabels for Prometheus metrics
type MetricLabels struct {
	Username string   `json:"username"`
	Hero     string   `json:"hero"`
	Platform Platform `json:"platform"`
	GameMode GameMode `json:"gamemode"`
}

// Common metrics available for all heroes
var CommonMetrics = map[string]MetricDef{
	"time_played": {
		PrometheusName: "ow_hero_time_played_seconds",
		Help:           "Total time played on hero in seconds",
		Unit:           "seconds",
		Selector:       ".Profile-progressBar-description",
		HexID:          "0x0860000000000021",
		ValueType:      "duration",
	},
	"games_won": {
		PrometheusName: "ow_hero_games_won_total",
		Help:           "Total number of games won with hero",
		Unit:           "count",
		Selector:       "[data-category-id='0x0860000000000039'] .Profile-progressBar-description",
		HexID:          "0x0860000000000039",
		ValueType:      "number",
	},
	"win_percentage": {
		PrometheusName: "ow_hero_win_percentage",
		Help:           "Win percentage with hero",
		Unit:           "percent",
		Selector:       "[data-category-id='0x08600000000003D1'] .Profile-progressBar-description",
		HexID:          "0x08600000000003D1",
		ValueType:      "percentage",
	},
	"weapon_accuracy": {
		PrometheusName: "ow_hero_weapon_accuracy_percent",
		Help:           "Best weapon accuracy percentage with hero",
		Unit:           "percent",
		Selector:       "[data-category-id='0x08600000000001BB'] .Profile-progressBar-description",
		HexID:          "0x08600000000001BB",
		ValueType:      "percentage",
	},
	"eliminations_per_life": {
		PrometheusName: "ow_hero_eliminations_per_life",
		Help:           "Average eliminations per life with hero",
		Unit:           "ratio",
		Selector:       "[data-category-id='0x08600000000003D2'] .Profile-progressBar-description",
		HexID:          "0x08600000000003D2",
		ValueType:      "number",
	},
	"kill_streak_best": {
		PrometheusName: "ow_hero_kill_streak_best",
		Help:           "Best kill streak achieved with hero",
		Unit:           "count",
		Selector:       "[data-category-id='0x0860000000000223'] .Profile-progressBar-description",
		HexID:          "0x0860000000000223",
		ValueType:      "number",
	},
	"multikill_best": {
		PrometheusName: "ow_hero_multikill_best",
		Help:           "Best multikill achieved with hero",
		Unit:           "count",
		Selector:       "[data-category-id='0x0860000000000346'] .Profile-progressBar-description",
		HexID:          "0x0860000000000346",
		ValueType:      "number",
	},
	"eliminations_per_10min": {
		PrometheusName: "ow_hero_eliminations_per_10min_avg",
		Help:           "Average eliminations per 10 minutes with hero",
		Unit:           "rate",
		Selector:       "[data-category-id='0x08600000000004D4'] .Profile-progressBar-description",
		HexID:          "0x08600000000004D4",
		ValueType:      "number",
	},
	"deaths_per_10min": {
		PrometheusName: "ow_hero_deaths_per_10min_avg",
		Help:           "Average deaths per 10 minutes with hero",
		Unit:           "rate",
		Selector:       "[data-category-id='0x08600000000004D3'] .Profile-progressBar-description",
		HexID:          "0x08600000000004D3",
		ValueType:      "number",
	},
	"final_blows_per_10min": {
		PrometheusName: "ow_hero_final_blows_per_10min_avg",
		Help:           "Average final blows per 10 minutes with hero",
		Unit:           "rate",
		Selector:       "[data-category-id='0x08600000000004D5'] .Profile-progressBar-description",
		HexID:          "0x08600000000004D5",
		ValueType:      "number",
	},
	"solo_kills_per_10min": {
		PrometheusName: "ow_hero_solo_kills_per_10min_avg",
		Help:           "Average solo kills per 10 minutes with hero",
		Unit:           "rate",
		Selector:       "[data-category-id='0x08600000000004DA'] .Profile-progressBar-description",
		HexID:          "0x08600000000004DA",
		ValueType:      "number",
	},
	"objective_kills_per_10min": {
		PrometheusName: "ow_hero_objective_kills_per_10min_avg",
		Help:           "Average objective kills per 10 minutes with hero",
		Unit:           "rate",
		Selector:       "[data-category-id='0x08600000000004D8'] .Profile-progressBar-description",
		HexID:          "0x08600000000004D8",
		ValueType:      "number",
	},
	"objective_time_per_10min": {
		PrometheusName: "ow_hero_objective_time_per_10min_avg",
		Help:           "Average objective time per 10 minutes with hero",
		Unit:           "seconds",
		Selector:       "[data-category-id='0x08600000000004D9'] .Profile-progressBar-description",
		HexID:          "0x08600000000004D9",
		ValueType:      "duration",
	},
	"hero_damage_per_10min": {
		PrometheusName: "ow_hero_damage_per_10min_avg",
		Help:           "Average hero damage per 10 minutes",
		Unit:           "damage",
		Selector:       "[data-category-id='0x08600000000004BD'] .Profile-progressBar-description",
		HexID:          "0x08600000000004BD",
		ValueType:      "number",
	},
	"healing_per_10min": {
		PrometheusName: "ow_hero_healing_per_10min_avg",
		Help:           "Average healing done per 10 minutes",
		Unit:           "healing",
		Selector:       "[data-category-id='0x08600000000004D6'] .Profile-progressBar-description",
		HexID:          "0x08600000000004D6",
		ValueType:      "number",
	},
}

// Hero-specific metrics can be added here for special abilities
// For now, we'll use the common metrics for all heroes

// CSS Selectors for platform and game mode detection
var PlatformSelectors = map[Platform]string{
	PlatformPC:      ".mouseKeyboard-view.is-active",
	PlatformConsole: ".controller-view.is-active",
}

var GameModeSelectors = map[GameMode]string{
	GameModeQuickPlay:   ".quickPlay-view.is-active",
	GameModeCompetitive: ".competitive-view.is-active",
}

// Hero identification selectors
var HeroSelectors = struct {
	Container string
	Name      string
	ID        string
	TimePlayed string
}{
	Container:  ".Profile-progressBar",
	Name:       ".Profile-progressBar-title",
	ID:         "[data-hero-id]",
	TimePlayed: ".Profile-progressBar-description",
}

// Platform filter selectors for switching views
var PlatformFilters = map[Platform]string{
	PlatformPC:      "#mouseKeyboardFilter",
	PlatformConsole: "#controllerFilter",
}

// Helper function to generate Prometheus metric name
func (m MetricDef) PrometheusMetricName(labels MetricLabels) string {
	return m.PrometheusName
}

// Helper function to get CSS selector for a metric in a specific context
func (m MetricDef) GetSelector(platform Platform, gameMode GameMode) string {
	// For now, use the base selector
	// In the future, we might need platform/gamemode specific selectors
	return m.Selector
}