package main

import (
	"fmt"
	"testing"
)

func TestHeroMetricsGeneration(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		heroID      string
		expectedMin int // Minimum expected metrics (15 common + hero-specific)
		description string
	}{
		{"soldier-76", 19, "Soldier:76 with 4 specific metrics"},
		{"mercy", 20, "Mercy with 5 specific metrics"},
		{"reinhardt", 20, "Reinhardt with 5 specific metrics"},
		{"widowmaker", 19, "Widowmaker with 4 specific metrics"},
		{"illari", 19, "Illari with 4 specific metrics"},
		{"hazard", 19, "Hazard with 4 specific metrics"},
		{"unknown-hero", 15, "Unknown hero fallback to common metrics"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.heroID, func(t *testing.T) {
			t.Parallel()
			metrics := GetHeroMetrics(testCase.heroID)

			if len(metrics) < testCase.expectedMin {
				t.Errorf("Hero %s: expected at least %d metrics, got %d",
					testCase.heroID, testCase.expectedMin, len(metrics))
			}

			// Verify common metrics are present
			commonMetrics := []string{
				"time_played", "games_won", "win_percentage",
				"weapon_accuracy", "eliminations_per_life",
			}

			for _, metricName := range commonMetrics {
				if _, exists := metrics[metricName]; !exists {
					t.Errorf("Hero %s: missing common metric %s", testCase.heroID, metricName)
				}
			}

			// Log metrics count for debugging
			t.Logf("Hero %s: found %d metrics (%s)", testCase.heroID, len(metrics), testCase.description)
		})
	}
}

func TestMetricDefGeneration(t *testing.T) {
	t.Parallel()
	// Test with Soldier:76 specifically
	soldier76 := Soldier76Metrics{}
	metrics := GenerateMetricDefs(soldier76)

	expectedMetrics := map[string]string{
		"helix_rocket_kills":      "ow_hero_helix_rocket_kills_total",
		"helix_rocket_kills_best": "ow_hero_helix_rocket_kills_best",
		"biotic_field_healing":    "ow_hero_biotic_field_healing_total",
		"tactical_visor_kills":    "ow_hero_tactical_visor_kills_total",
	}

	for metricName, expectedPrometheusName := range expectedMetrics {
		metricDef, exists := metrics[metricName]
		if !exists {
			t.Errorf("Missing Soldier:76 metric: %s", metricName)

			continue
		}

		if metricDef.PrometheusName != expectedPrometheusName {
			t.Errorf("Metric %s: expected prometheus name %s, got %s",
				metricName, expectedPrometheusName, metricDef.PrometheusName)
		}

		if metricDef.Help == "" {
			t.Errorf("Metric %s: missing help text", metricName)
		}

		if metricDef.Selector == "" {
			t.Errorf("Metric %s: missing selector", metricName)
		}
	}

	t.Logf("Soldier:76 specific metrics: %d", len(metrics))
}

func TestHeroMetricsRegistry(t *testing.T) {
	t.Parallel()
	expectedHeroCount := 29 // Current number of implemented heroes

	if len(HeroMetricsRegistry) != expectedHeroCount {
		t.Errorf("Expected %d heroes in registry, got %d", expectedHeroCount, len(HeroMetricsRegistry))
	}

	// Test that all heroes in registry can generate metrics
	for heroID, factory := range HeroMetricsRegistry {
		heroStruct := factory()
		metrics := GenerateMetricDefs(heroStruct)

		if len(metrics) == 0 {
			t.Errorf("Hero %s: factory returned struct with no metrics", heroID)
		}

		t.Logf("Hero %s: %d specific metrics", heroID, len(metrics))
	}
}

func TestPlatformSpecificMetrics(t *testing.T) {
	t.Parallel()
	heroID := "soldier-76"

	testCases := []struct {
		platform    Platform
		gameMode    GameMode
		description string
	}{
		{PlatformPC, GameModeQuickPlay, "PC QuickPlay"},
		{PlatformPC, GameModeCompetitive, "PC Competitive"},
		{PlatformConsole, GameModeQuickPlay, "Console QuickPlay"},
		{PlatformConsole, GameModeCompetitive, "Console Competitive"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			t.Parallel()
			validatePlatformSpecificMetrics(t, heroID, testCase.platform, testCase.gameMode, testCase.description)
		})
	}
}

// validatePlatformSpecificMetrics validates metrics for a specific platform and game mode.
func validatePlatformSpecificMetrics(
	t *testing.T, heroID string, platform Platform, gameMode GameMode, description string,
) {
	t.Helper()
	metrics := GetHeroMetricsForPlatform(heroID, platform, gameMode)

	if len(metrics) == 0 {
		t.Errorf("No metrics found for %s", description)

		return
	}

	t.Logf("%s: found %d metrics", description, len(metrics))
	logFirstFewSelectors(t, metrics, description)
	validateAllMetricSelectors(t, metrics, platform, gameMode)
}

// logFirstFewSelectors logs the first few metric selectors for debugging.
func logFirstFewSelectors(t *testing.T, metrics map[string]MetricDef, _ string) {
	t.Helper()
	count := 0
	for metricName, metricDef := range metrics {
		if count >= 2 {
			break
		}
		if metricDef.Selector != "" {
			t.Logf("  %s selector: %s", metricName, metricDef.Selector)
			count++
		}
	}
}

// validateAllMetricSelectors validates all metric selectors contain platform/gamemode context.
func validateAllMetricSelectors(t *testing.T, metrics map[string]MetricDef, platform Platform, gameMode GameMode) {
	t.Helper()
	for metricName, metricDef := range metrics {
		if metricDef.Selector != "" {
			validateMetricSelector(t, metricName, metricDef.Selector, platform, gameMode)
		}
	}
}

// validateMetricSelector validates that a metric selector contains the required platform and gamemode wrappers.
func validateMetricSelector(t *testing.T, metricName, selector string, platform Platform, gameMode GameMode) {
	t.Helper()
	platformWrapper := getPlatformWrapper(platform)
	if !contains(selector, platformWrapper) {
		t.Errorf("Metric %s selector missing platform wrapper: %s", metricName, selector)
	}

	gameModeWrapper := getGameModeWrapper(gameMode)
	if !contains(selector, gameModeWrapper) {
		t.Errorf("Metric %s selector missing gamemode wrapper: %s", metricName, selector)
	}
}

// getPlatformWrapper returns the CSS selector wrapper for the given platform.
func getPlatformWrapper(platform Platform) string {
	switch platform {
	case PlatformPC:
		return MouseKeyboardViewActiveSelector
	case PlatformConsole:
		return ControllerViewActiveSelector
	default:
		return ""
	}
}

// getGameModeWrapper returns the CSS selector wrapper for the given game mode.
func getGameModeWrapper(gameMode GameMode) string {
	switch gameMode {
	case GameModeQuickPlay:
		return QuickPlayViewActiveSelector
	case GameModeCompetitive:
		return CompetitiveViewActiveSelector
	default:
		return ""
	}
}

func contains(s, substr string) bool {
	return substr != "" && len(s) >= len(substr) &&
		(s == substr || s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr ||
			findInString(s, substr))
}

func findInString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}

	return false
}

// Example usage function for documentation.
func ExampleGetHeroMetrics() {
	// Get metrics for Soldier:76
	metrics := GetHeroMetrics("soldier-76")

	fmt.Printf("Soldier:76 has %d total metrics\n", len(metrics))

	// Access a specific metric
	if helixMetric, exists := metrics["helix_rocket_kills"]; exists {
		fmt.Printf("Helix Rocket Kills: %s\n", helixMetric.PrometheusName)
		fmt.Printf("Help: %s\n", helixMetric.Help)
	}

	// Output:
	// Soldier:76 has 19 total metrics
	// Helix Rocket Kills: ow_hero_helix_rocket_kills_total
	// Help: Total eliminations with helix rockets
}
