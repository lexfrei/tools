package main

import "time"

// MetricDef defines a single metric with its parsing information.
type MetricDef struct {
	PrometheusName string `json:"prometheusName"` // "ow_hero_time_played_seconds"
	Help           string `json:"help"`           // Help text for Prometheus
	Unit           string `json:"unit"`           // "seconds", "percent", "count", "ratio"
	Selector       string `json:"selector"`       // CSS selector or data attribute
	HexID          string `json:"hexId"`          // Blizzard's hex ID for the metric
	ValueType      string `json:"valueType"`      // "duration", "number", "percentage"
}

// Platform represents PC or Console.
type Platform string

const (
	PlatformPC      Platform = "pc"
	PlatformConsole Platform = "console"
)

// GameMode represents Quick Play or Competitive.
type GameMode string

const (
	GameModeQuickPlay   GameMode = "quickplay"
	GameModeCompetitive GameMode = "competitive"
)

// MetricLabels for Prometheus metrics.
type MetricLabels struct {
	Username string   `json:"username"`
	Hero     string   `json:"hero"`
	Platform Platform `json:"platform"`
	GameMode GameMode `json:"gamemode"`
}

// GetCommonMetrics returns common metrics using the new embedded struct approach.
func GetCommonMetrics() map[string]MetricDef {
	return GetCommonMetricsForPlatform(PlatformPC, GameModeQuickPlay)
}

// GetCommonMetricsForPlatform returns common metrics for specific platform and game mode.
func GetCommonMetricsForPlatform(platform Platform, gameMode GameMode) map[string]MetricDef {
	commonStruct := CommonMetrics{}

	return GenerateMetricDefsWithContext(commonStruct, platform, gameMode)
}

// Hero-specific metrics can be added here for special abilities.

// PlatformSelectors for platform and game mode detection.
var PlatformSelectors = map[Platform]string{
	PlatformPC:      ".mouseKeyboard-view",
	PlatformConsole: ".controller-view",
}

var GameModeSelectors = map[GameMode]string{
	GameModeQuickPlay:   ".quickPlay-view",
	GameModeCompetitive: ".competitive-view",
}

// HeroSelectors for hero identification and detailed stats parsing.
var HeroSelectors = struct {
	Container         string
	Name              string
	ID                string
	TimePlayed        string
	StatsContainer    string
	StatItem          string
	StatName          string
	StatValue         string
	CategoryHeader    string
	BlzStatsSection   string
}{
	Container:         ".Profile-progressBar",
	Name:              ".Profile-progressBar-title",
	ID:                ".Profile-progressBar--bar[data-hero-id]", // data-hero-id is on the bar element
	TimePlayed:        ".Profile-progressBar-description",
	StatsContainer:    "span.stats-container",              // OverFast API style
	StatItem:          ".stat-item",                        // Individual stat items
	StatName:          "p.name",                            // Stat name within stat-item
	StatValue:         "p.value",                           // Stat value within stat-item
	CategoryHeader:    ".category .content .header p",     // Category headers
	BlzStatsSection:   "blz-section.stats",                 // Main stats sections
}

// PlatformFilters for switching views.
var PlatformFilters = map[Platform]string{
	PlatformPC:      "#mouseKeyboardFilter",
	PlatformConsole: "#controllerFilter",
}

// PrometheusMetricName generates Prometheus metric name.
func (m *MetricDef) PrometheusMetricName(_ MetricLabels) string {
	return m.PrometheusName
}

// GetSelector returns CSS selector for this metric.
func (m *MetricDef) GetSelector() string {
	return m.Selector
}

// New structures for the enhanced metrics system.

// AllHeroesStats represents aggregated statistics across all heroes.
type AllHeroesStats struct {
	TotalTimePlayed        int64   `json:"totalTimePlayedSeconds" prometheus:"ow_player_total_time_played_seconds"`
	TotalGamesWon          int     `json:"totalGamesWon"          prometheus:"ow_player_total_games_won"`
	OverallWinPercentage   float64 `json:"overallWinPercentage"   prometheus:"ow_player_overall_win_percentage"`
	WeaponAccuracy         float64 `json:"weaponAccuracyPercent"  prometheus:"ow_player_weapon_accuracy_percent"`
	EliminationsPerLife    float64 `json:"eliminationsPerLife"    prometheus:"ow_player_eliminations_per_life"`
	KillStreakBest         int     `json:"killStreakBest"         prometheus:"ow_player_kill_streak_best"`
	MultikillBest          int     `json:"multikillBest"          prometheus:"ow_player_multikill_best"`
	EliminationsPer10Min   float64 `json:"eliminationsPer10min"   prometheus:"ow_player_eliminations_per_10min"`
	DeathsPer10Min         float64 `json:"deathsPer10min"         prometheus:"ow_player_deaths_per_10min"`
	FinalBlowsPer10Min     float64 `json:"finalBlowsPer10min"     prometheus:"ow_player_final_blows_per_10min"`
	SoloKillsPer10Min      float64 `json:"soloKillsPer10min"      prometheus:"ow_player_solo_kills_per_10min"`
	ObjectiveKillsPer10Min float64 `json:"objectiveKillsPer10min" prometheus:"ow_player_objective_kills_per_10min"`
	ObjectiveTimePer10Min  float64 `json:"objectiveTimePer10min"  prometheus:"ow_player_objective_time_per_10min"`
	HeroDamagePer10Min     float64 `json:"heroDamagePer10min"     prometheus:"ow_player_hero_damage_per_10min"`
	HealingPer10Min        float64 `json:"healingPer10min"        prometheus:"ow_player_healing_per_10min"`
}

// HeroMetrics represents metrics for a specific hero.
type HeroMetrics map[string]interface{}

// RuntimeMetrics contains all runtime metrics data.
type RuntimeMetrics struct {
	Players map[string]*PlayerMetrics `json:"players"` // battletag -> metrics
}

// PlayerMetrics contains all metrics for a single player.
type PlayerMetrics struct {
	BattleTag   string    `json:"battletag"`
	DisplayName string    `json:"displayName"` // "Joe" from HTML
	PlayerTitle string    `json:"playerTitle"` // "Peasant" from HTML
	LastUpdated time.Time `json:"lastUpdated"`

	// Level 1: Profile-level metrics (SR, endorsement)
	ProfileMetrics ProfileMetrics `json:"profileMetrics"`

	// Level 2: All Heroes aggregated metrics by platform/gamemode
	AllHeroesMetrics map[Platform]map[GameMode]AllHeroesStats `json:"allHeroesMetrics"`

	// Level 3: Individual hero metrics by platform/gamemode/hero
	HeroMetrics map[Platform]map[GameMode]map[string]HeroMetrics `json:"heroMetrics"`
}

// EnhancedMetricLabels with BattleTag support.
type EnhancedMetricLabels struct {
	BattleTag  string   `json:"battletag"`  // LexFrei#21715
	PlayerName string   `json:"playerName"` // Joe
	Hero       string   `json:"hero"`       // soldier-76, widowmaker, etc.
	Platform   Platform `json:"platform"`   // pc, console
	GameMode   GameMode `json:"gamemode"`   // quickplay, competitive
	MetricType string   `json:"metricType"` // profile, all_heroes, hero
}

// AllHeroesHexIDs for parsing.
var AllHeroesHexIDs = map[string]string{
	"time_played":               "0x0860000000000021",
	"games_won":                 "0x0860000000000039",
	"win_percentage":            "0x08600000000003D1",
	"weapon_accuracy":           "0x08600000000001BB",
	"eliminations_per_life":     "0x08600000000003D2",
	"kill_streak_best":          "0x0860000000000223",
	"multikill_best":            "0x0860000000000346",
	"eliminations_per_10min":    "0x08600000000004D4",
	"deaths_per_10min":          "0x08600000000004D3",
	"final_blows_per_10min":     "0x08600000000004D5",
	"solo_kills_per_10min":      "0x08600000000004DA",
	"objective_kills_per_10min": "0x08600000000004D8",
	"objective_time_per_10min":  "0x08600000000004D9",
	"hero_damage_per_10min":     "0x08600000000004BD",
	"healing_per_10min":         "0x08600000000004D6",
}
