// parser_example.go - Example parser with PC/Console support
package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// PlayerStats represents all statistics for a player
type PlayerStats struct {
	Username  string                          `json:"username"`
	LastUpdate time.Time                      `json:"last_update"`
	Platforms map[Platform]*PlatformStats    `json:"platforms"`
}

// PlatformStats represents statistics for a specific platform (PC/Console)
type PlatformStats struct {
	Platform  Platform                      `json:"platform"`
	GameModes map[GameMode]*GameModeStats   `json:"game_modes"`
}

// GameModeStats represents statistics for a specific game mode
type GameModeStats struct {
	GameMode GameMode                `json:"game_mode"`
	Heroes   map[string]*HeroStats   `json:"heroes"`
}

// HeroStats represents all statistics for a specific hero
type HeroStats struct {
	HeroID   string             `json:"hero_id"`
	HeroName string             `json:"hero_name"`
	Metrics  map[string]float64 `json:"metrics"`
}

// Parser handles HTML parsing for Overwatch profiles
type Parser struct {
	// Add any configuration or dependencies here
}

// NewParser creates a new parser instance
func NewParser() *Parser {
	return &Parser{}
}

// ParseProfile parses an Overwatch profile HTML and extracts all statistics
func (p *Parser) ParseProfile(html string, username string) (*PlayerStats, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	stats := &PlayerStats{
		Username:   username,
		LastUpdate: time.Now(),
		Platforms:  make(map[Platform]*PlatformStats),
	}

	// Parse both PC and Console platforms
	for platform := range PlatformSelectors {
		platformStats, err := p.parsePlatformStats(doc, platform)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %s stats: %w", platform, err)
		}
		if platformStats != nil {
			stats.Platforms[platform] = platformStats
		}
	}

	return stats, nil
}

// parsePlatformStats parses statistics for a specific platform
func (p *Parser) parsePlatformStats(doc *goquery.Document, platform Platform) (*PlatformStats, error) {
	// Check if this platform has data
	platformSelector := PlatformSelectors[platform]
	platformView := doc.Find(platformSelector)
	if platformView.Length() == 0 {
		return nil, nil // No data for this platform
	}

	stats := &PlatformStats{
		Platform:  platform,
		GameModes: make(map[GameMode]*GameModeStats),
	}

	// Parse both Quick Play and Competitive modes
	for gameMode := range GameModeSelectors {
		gameModeStats, err := p.parseGameModeStats(platformView, gameMode)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %s stats: %w", gameMode, err)
		}
		if gameModeStats != nil {
			stats.GameModes[gameMode] = gameModeStats
		}
	}

	return stats, nil
}

// parseGameModeStats parses statistics for a specific game mode within a platform
func (p *Parser) parseGameModeStats(platformView *goquery.Selection, gameMode GameMode) (*GameModeStats, error) {
	// Find the game mode view within the platform
	gameModeSelector := GameModeSelectors[gameMode]
	gameModeView := platformView.Find(gameModeSelector)
	if gameModeView.Length() == 0 {
		return nil, nil // No data for this game mode
	}

	stats := &GameModeStats{
		GameMode: gameMode,
		Heroes:   make(map[string]*HeroStats),
	}

	// Parse hero statistics
	heroContainers := gameModeView.Find(HeroSelectors.Container)
	heroContainers.Each(func(i int, heroEl *goquery.Selection) {
		heroStats := p.parseHeroStats(heroEl)
		if heroStats != nil {
			stats.Heroes[heroStats.HeroID] = heroStats
		}
	})

	return stats, nil
}

// parseHeroStats parses statistics for a single hero
func (p *Parser) parseHeroStats(heroEl *goquery.Selection) *HeroStats {
	// Extract hero ID
	heroID, exists := heroEl.Attr("data-hero-id")
	if !exists {
		return nil
	}

	// Extract hero name
	heroName := strings.TrimSpace(heroEl.Find(HeroSelectors.Name).Text())
	if heroName == "" {
		return nil
	}

	stats := &HeroStats{
		HeroID:   heroID,
		HeroName: heroName,
		Metrics:  make(map[string]float64),
	}

	// Parse time played (always visible)
	timePlayedText := strings.TrimSpace(heroEl.Find(HeroSelectors.TimePlayed).Text())
	if timePlayedText != "" {
		if timePlayedSeconds := p.parseTimeToSeconds(timePlayedText); timePlayedSeconds > 0 {
			stats.Metrics[CommonMetrics["time_played"].PrometheusName] = timePlayedSeconds
		}
	}

	// For other metrics, we would need to iterate through different views
	// or trigger JavaScript to change the dropdown selection
	// For now, we'll just get the time played metric which is always visible

	return stats
}

// parseTimeToSeconds converts time strings like "44:28:48" to seconds
func (p *Parser) parseTimeToSeconds(timeStr string) float64 {
	// Handle formats: "HH:MM:SS" or "MM:SS" or just numbers
	re := regexp.MustCompile(`^(?:(\d+):)?(\d+):(\d+)$`)
	matches := re.FindStringSubmatch(timeStr)

	if len(matches) == 4 {
		var hours, minutes, seconds int
		var err error

		if matches[1] != "" {
			// HH:MM:SS format
			hours, err = strconv.Atoi(matches[1])
			if err != nil {
				return 0
			}
			minutes, err = strconv.Atoi(matches[2])
			if err != nil {
				return 0
			}
			seconds, err = strconv.Atoi(matches[3])
			if err != nil {
				return 0
			}
		} else {
			// MM:SS format
			minutes, err = strconv.Atoi(matches[2])
			if err != nil {
				return 0
			}
			seconds, err = strconv.Atoi(matches[3])
			if err != nil {
				return 0
			}
		}

		return float64(hours*3600 + minutes*60 + seconds)
	}

	return 0
}

// parsePercentage converts percentage strings like "74%" to float
func (p *Parser) parsePercentage(percentStr string) float64 {
	cleaned := strings.TrimSuffix(strings.TrimSpace(percentStr), "%")
	if value, err := strconv.ParseFloat(cleaned, 64); err == nil {
		return value
	}
	return 0
}

// parseNumber converts number strings to float
func (p *Parser) parseNumber(numStr string) float64 {
	// Remove commas and other formatting
	cleaned := strings.ReplaceAll(strings.TrimSpace(numStr), ",", "")
	if value, err := strconv.ParseFloat(cleaned, 64); err == nil {
		return value
	}
	return 0
}

// Example usage function for testing
func ExampleParser() {
	// This would be called from the main ow-exporter application
	_ = NewParser()

	// Read one of our saved profiles
	// html := readFile("/Users/lex/git/github.com/lexfrei/tools/tmp/profile_de5bb4aca17492e0.html")
	// stats, err := parser.ParseProfile(html, "LexFrei")
	// if err != nil {
	//     log.Fatal(err)
	// }

	fmt.Println("Parser ready for integration into ow-exporter")
}