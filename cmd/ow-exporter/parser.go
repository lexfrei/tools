package main

import (
	"fmt"
	"log/slog"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/cockroachdb/errors"

	"github.com/PuerkitoBio/goquery"
)

// Constants for metric names.
const (
	GamesWonMetric               = "games_won"
	WinPercentageMetric          = "win_percentage"
	WeaponAccuracyMetric         = "weapon_accuracy"
	EliminationsPerLifeMetric    = "eliminations_per_life"
	KillStreakBestMetric         = "kill_streak_best"
	MultikillBestMetric          = "multikill_best"
	EliminationsPer10MinMetric   = "eliminations_per_10min"
	DeathsPer10MinMetric         = "deaths_per_10min"
	FinalBlowsPer10MinMetric     = "final_blows_per_10min"
	SoloKillsPer10MinMetric      = "solo_kills_per_10min"
	ObjectiveKillsPer10MinMetric = "objective_kills_per_10min"
	ObjectiveTimePer10MinMetric  = "objective_time_per_10min"
	HeroDamagePer10MinMetric     = "hero_damage_per_10min"
	HealingPer10MinMetric        = "healing_per_10min"
)

// FullPlayerProfile represents complete player profile with all metrics.
type FullPlayerProfile struct {
	// Basic profile information
	Username    string    `json:"username"`
	BattleTag   string    `json:"battleTag"`
	PlayerTitle string    `json:"playerTitle"`
	LastUpdate  time.Time `json:"lastUpdate"`

	// Profile-level metrics
	ProfileMetrics ProfileMetrics `json:"profileMetrics"`

	// Hero statistics matrix: platform -> gamemode -> hero -> metrics
	Platforms map[Platform]*PlatformStats `json:"platforms"`
}

// ProfileMetrics represents profile-level statistics.
type ProfileMetrics struct {
	Endorsement  EndorsementData                `json:"endorsement"`
	SkillRatings map[Platform]map[Role]RankInfo `json:"skillRatings"`
}

// EndorsementData represents endorsement level and breakdown.
type EndorsementData struct {
	Level int `json:"level"`
	// Endorsement breakdown by category (if available in HTML)
	Breakdown *EndorsementBreakdown `json:"breakdown,omitempty"`
}

// EndorsementBreakdown represents the breakdown of endorsements by category.
type EndorsementBreakdown struct {
	Sportsmanship int `json:"sportsmanship"` // Good teammate, stays positive
	Teamwork      int `json:"teamwork"`      // Team player, communicates
	ShotCaller    int `json:"shotCaller"`    // Good leadership, makes good calls
}

// Role represents player roles.
type Role string

const (
	RoleTank    Role = "tank"
	RoleDamage  Role = "damage"
	RoleSupport Role = "support"
)

// RankInfo represents competitive ranking information.
type RankInfo struct {
	Tier     string `json:"tier"`     // Bronze, Silver, Gold, Platinum, Diamond, Master, Grandmaster, Champion
	Division int    `json:"division"` // 1-5
	SR       int    `json:"sr"`       // Skill Rating (if available)
}

// PlatformStats represents statistics for a specific platform (PC/Console).
type PlatformStats struct {
	Platform  Platform                    `json:"platform"`
	GameModes map[GameMode]*GameModeStats `json:"gameModes"`
}

// GameModeStats represents statistics for a specific game mode.
type GameModeStats struct {
	GameMode       GameMode              `json:"gameMode"`
	AllHeroesStats AllHeroesStats        `json:"allHeroesStats"`
	Heroes         map[string]*HeroStats `json:"heroes"`
}

// HeroStats represents all statistics for a specific hero.
type HeroStats struct {
	HeroID   string             `json:"heroId"`
	HeroName string             `json:"heroName"`
	Metrics  map[string]float64 `json:"metrics"`
}

// Parser handles HTML parsing for Overwatch profiles.
type Parser struct {
	// Add any configuration or dependencies here
}

// NewParser creates a new parser instance.
func NewParser() *Parser {
	return &Parser{}
}

// ParseProfile parses an Overwatch profile HTML and extracts all statistics.
func (p *Parser) ParseProfile(html, username string) (*FullPlayerProfile, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse HTML")
	}

	// Extract basic profile info
	playerName := strings.TrimSpace(doc.Find(".Profile-player--name").Text())
	playerTitle := strings.TrimSpace(doc.Find(".Profile-player--title").Text())

	stats := &FullPlayerProfile{
		Username:    username,
		BattleTag:   playerName,
		PlayerTitle: playerTitle,
		LastUpdate:  time.Now(),
		Platforms:   make(map[Platform]*PlatformStats),
		ProfileMetrics: ProfileMetrics{
			SkillRatings: make(map[Platform]map[Role]RankInfo),
		},
	}

	// Parse profile-level metrics
	p.parseProfileMetrics(doc, stats)

	// Debug summary
	progressBars := doc.Find(".Profile-progressBar")
	heroElements := doc.Find("[data-hero-id]")
	slog.Info("HTML parsing debug",
		"progress_bars", progressBars.Length(),
		"hero_elements", heroElements.Length(),
		"player_name", playerName,
		"player_title", playerTitle)

	// Parse both PC and Console platforms
	for platform := range PlatformSelectors {
		platformStats, err := p.parsePlatformStats(doc, platform)
		if err != nil {
			// Skip platforms with no data
			if errors.Is(err, ErrNoPlatformData) {
				continue // Skip silently
			}

			return nil, errors.Wrapf(err, "failed to parse %s stats", platform)
		}
		if platformStats != nil {
			stats.Platforms[platform] = platformStats
		}
	}

	return stats, nil
}

// parsePlatformStats parses statistics for a specific platform.
func (p *Parser) parsePlatformStats(doc *goquery.Document, platform Platform) (*PlatformStats, error) {
	// Check if this platform has data
	platformSelector := PlatformSelectors[platform]
	platformView := doc.Find(platformSelector)

	// Debug: print what we found
	slog.Info("Platform parsing",
		"platform", platform,
		"selector", platformSelector,
		"elements_found", platformView.Length())

	if platformView.Length() == 0 {
		return nil, ErrNoPlatformData
	}

	stats := &PlatformStats{
		Platform:  platform,
		GameModes: make(map[GameMode]*GameModeStats),
	}

	// Parse both Quick Play and Competitive modes
	for gameMode := range GameModeSelectors {
		gameModeStats, err := p.parseGameModeStats(platformView, gameMode)
		if err != nil {
			// Skip game modes with no data
			if errors.Is(err, ErrNoGameModeData) {
				continue // Skip silently
			}

			return nil, errors.Wrapf(err, "failed to parse %s stats", gameMode)
		}
		if gameModeStats != nil {
			stats.GameModes[gameMode] = gameModeStats
		}
	}

	return stats, nil
}

// parseGameModeStats parses statistics for a specific game mode within a platform.
func (p *Parser) parseGameModeStats(platformView *goquery.Selection, gameMode GameMode) (*GameModeStats, error) {
	// Find the game mode view within the platform
	gameModeSelector := GameModeSelectors[gameMode]
	gameModeView := platformView.Find(gameModeSelector)

	// Debug game mode detection
	slog.Info("GameMode parsing",
		"gamemode", gameMode,
		"selector", gameModeSelector,
		"elements_found", gameModeView.Length())

	if gameModeView.Length() == 0 {
		return nil, ErrNoGameModeData
	}

	stats := &GameModeStats{
		GameMode: gameMode,
		Heroes:   make(map[string]*HeroStats),
	}

	// Parse All Heroes aggregated statistics first
	allHeroesStats := p.parseAllHeroesStats(gameModeView)
	if allHeroesStats != nil {
		stats.AllHeroesStats = *allHeroesStats
	}

	// Parse hero statistics - look for progress bars within this game mode view
	heroContainers := gameModeView.Find(HeroSelectors.Container)
	slog.Info("Hero containers found",
		"gamemode", gameMode,
		"hero_count", heroContainers.Length())

	var parseErrors []error
	heroContainers.Each(func(_ int, heroEl *goquery.Selection) {
		heroStats, err := p.parseHeroStats(heroEl)
		if err != nil {
			// Skip heroes with missing data (not real errors)
			if errors.Is(err, ErrNoHeroID) || errors.Is(err, ErrNoHeroName) {
				return // Skip silently
			}
			parseErrors = append(parseErrors, err)

			return
		}
		if heroStats != nil {
			stats.Heroes[heroStats.HeroID] = heroStats
			slog.Debug("Hero parsed",
				"hero_name", heroStats.HeroName,
				"hero_id", heroStats.HeroID,
				"metrics_count", len(heroStats.Metrics))
		}
	})

	// Check for validation errors
	if len(parseErrors) > 0 {
		for _, err := range parseErrors {
			slog.Error("Hero parsing validation error", "error", err.Error())
		}

		return nil, parseErrors[0] // Return first error to fail fast
	}

	return stats, nil
}

// parseHeroStats parses statistics for a single hero.
func (p *Parser) parseHeroStats(heroEl *goquery.Selection) (*HeroStats, error) {
	// Extract and validate hero identification
	heroID, heroName, err := p.extractHeroIdentification(heroEl)
	if err != nil {
		return nil, err
	}

	// Validate hero exists in registry
	err = p.validateHeroInRegistry(heroID, heroName)
	if err != nil {
		return nil, err
	}

	stats := &HeroStats{
		HeroID:   heroID,
		HeroName: heroName,
		Metrics:  make(map[string]float64),
	}

	// Extract all metrics for this hero
	p.extractAllHeroMetrics(heroEl, heroID, heroName, stats)

	return stats, nil
}

// extractHeroIdentification extracts hero ID and name from the hero element.
func (p *Parser) extractHeroIdentification(heroEl *goquery.Selection) (heroID, heroName string, err error) {
	// Extract hero ID from the bar element
	barEl := heroEl.Find(".Profile-progressBar--bar[data-hero-id]")
	if barEl.Length() == 0 {
		return "", "", ErrNoHeroID
	}

	heroID, exists := barEl.Attr("data-hero-id")
	if !exists {
		return "", "", ErrNoHeroID
	}

	// Extract hero name
	heroName = strings.TrimSpace(heroEl.Find(HeroSelectors.Name).Text())
	if heroName == "" {
		slog.Warn("Hero missing name", "hero_id", heroID, "selector", HeroSelectors.Name)

		return "", "", ErrNoHeroName
	}

	return heroID, heroName, nil
}

// validateHeroInRegistry ensures the hero exists in our metrics registry.
func (p *Parser) validateHeroInRegistry(heroID, heroName string) error {
	if _, exists := HeroMetricsRegistry[heroID]; !exists {
		return errors.Wrapf(ErrUnknownHero,
			"hero_id='%s', hero_name='%s' is not in HeroMetricsRegistry - 100%% coverage validation failed",
			heroID, heroName)
	}

	return nil
}

// extractAllHeroMetrics extracts all detailed metrics for the hero.
// Now updated to work with the actual HTML structure from JavaScript-loaded content.
func (p *Parser) extractAllHeroMetrics(heroEl *goquery.Selection, heroID, heroName string, stats *HeroStats) {
	// Extract common metrics from progress bars
	p.extractCommonMetricsFromProgressBar(heroEl, heroID, heroName, stats)

	// Extract detailed hero-specific metrics from stats-container elements
	p.extractDetailedHeroMetrics(heroEl, heroID, heroName, stats)

	slog.Debug("Hero parsing completed",
		"hero_id", heroID,
		"hero_name", heroName,
		"extracted_metrics", len(stats.Metrics))
}

// extractCommonMetricsFromProgressBar extracts the common metrics from the progress bar description.
// This reads the value from the .Profile-progressBar-description element.
func (p *Parser) extractCommonMetricsFromProgressBar(heroEl *goquery.Selection, heroID, heroName string, stats *HeroStats) {
	// Get the progress bar description value (this shows the current metric value)
	descriptionEl := heroEl.Find(".Profile-progressBar-description")
	if descriptionEl.Length() == 0 {
		slog.Debug("No progress bar description found",
			"hero_id", heroID,
			"selector", ".Profile-progressBar-description")
		return
	}

	valueText := strings.TrimSpace(descriptionEl.Text())
	if valueText == "" {
		return
	}

	// Determine the current metric being displayed by finding the active category
	// Look for the active data-category-id in the parent container
	parentDoc := heroEl.Parents().Last()
	activeCategory := parentDoc.Find(".Profile-progressBars.is-active[data-category-id]")
	if activeCategory.Length() == 0 {
		slog.Debug("No active category found for hero",
			"hero_id", heroID,
			"value_text", valueText)
		return
	}

	categoryID, exists := activeCategory.Attr("data-category-id")
	if !exists {
		return
	}

	// Map category ID to metric name
	metricName := p.mapCategoryIDToMetricName(categoryID)
	if metricName == "" {
		slog.Debug("Unknown category ID",
			"hero_id", heroID,
			"category_id", categoryID)
		return
	}

	// Parse the value based on metric type
	value := p.parseMetricValue(valueText, metricName)
	if value > 0 {
		stats.Metrics[metricName] = value
		slog.Debug("Extracted progress bar metric",
			"hero_id", heroID,
			"metric", metricName,
			"value", value,
			"value_text", valueText,
			"category_id", categoryID)
	}
}

// mapCategoryIDToMetricName maps hex category IDs to metric names.
func (p *Parser) mapCategoryIDToMetricName(categoryID string) string {
	// Reverse lookup from AllHeroesHexIDs
	for metricName, hexID := range AllHeroesHexIDs {
		if hexID == categoryID {
			return metricName
		}
	}
	return ""
}

// extractMetricFromDocument extracts and parses a single metric value using its definition.
// NOTE: This function is currently not used as detailed metrics are not available in the HTML.
func (p *Parser) extractMetricFromDocument(
	heroEl *goquery.Selection,
	metricDef *MetricDef,
	metricKey, heroID string,
) float64 {
	// Find the element using the selector
	selection := p.findMetricElement(heroEl, metricDef.Selector, heroID, metricKey)
	if selection == nil {
		return -1
	}

	// Extract the value text from the element
	valueText := p.extractValueFromElement(selection)
	if valueText == "" {
		return -1
	}

	// Parse the value based on its type
	return p.parseValueByType(valueText, metricDef.ValueType)
}

// findMetricElement finds the HTML element using the provided selector.
func (p *Parser) findMetricElement(heroEl *goquery.Selection, selector, heroID, metricKey string) *goquery.Selection {
	// Try hero-specific search first
	selection := heroEl.Find(selector)

	if selection.Length() == 0 {
		// Try broader document search if hero-specific search fails
		selection = heroEl.Parent().Find(selector)
	}

	if selection.Length() == 0 {
		slog.Debug("Metric selector not found",
			"hero_id", heroID,
			"metric", metricKey,
			"selector", selector)

		return nil
	}

	return selection
}

// extractValueFromElement extracts the text value from an HTML element.
func (p *Parser) extractValueFromElement(selection *goquery.Selection) string {
	// Try different extraction methods
	valueText := strings.TrimSpace(selection.Text())
	if valueText == "" {
		valueText, _ = selection.Attr("data-value")
	}
	if valueText == "" {
		valueText, _ = selection.Attr("value")
	}

	return valueText
}

// parseValueByType parses a value string based on the specified type.
func (p *Parser) parseValueByType(valueText, valueType string) float64 {
	switch valueType {
	case DurationMetricType:
		return p.parseTimeToSeconds(valueText)

	case PercentageMetricType:
		return p.parsePercentage(valueText)

	case NumberMetricType:
		return p.parseNumericValue(valueText)

	default:
		// Default to number parsing
		return p.parseNumericValue(valueText)
	}
}

// parseNumericValue tries to parse a value as integer first, then float.
func (p *Parser) parseNumericValue(valueText string) float64 {
	// Try integer first
	if intValue := p.parseNumber(valueText); intValue >= 0 {
		return float64(intValue)
	}
	// Try float

	return p.parseFloat(valueText)
}

// parseFloat parses float values from text, removing non-numeric characters.
func (p *Parser) parseFloat(text string) float64 {
	// Remove all non-numeric characters except decimal point and minus
	cleanStr := regexp.MustCompile(`[^\d.-]`).ReplaceAllString(text, "")
	if cleanStr == "" {
		return -1
	}

	value, err := strconv.ParseFloat(cleanStr, 64)
	if err != nil {
		return -1
	}

	return value
}

// parseTimeToSeconds converts time strings like "44:28:48" to seconds.
func (p *Parser) parseTimeToSeconds(timeStr string) float64 {
	// Handle formats: "HH:MM:SS" or "MM:SS" or just numbers
	re := regexp.MustCompile(`^(?:(\d+):)?(\d+):(\d+)$`)
	matches := re.FindStringSubmatch(timeStr)

	const expectedRegexMatches = 4
	if len(matches) != expectedRegexMatches {
		return 0
	}

	if matches[1] != "" {
		return p.parseHHMMSSFormat(matches)
	}

	return p.parseMMSSFormat(matches)
}

// parseHHMMSSFormat parses time in HH:MM:SS format.
func (p *Parser) parseHHMMSSFormat(matches []string) float64 {
	hours, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0
	}

	minutes, err := strconv.Atoi(matches[2])
	if err != nil {
		return 0
	}

	seconds, err := strconv.Atoi(matches[3])
	if err != nil {
		return 0
	}

	return float64(hours*3600 + minutes*60 + seconds)
}

// parseMMSSFormat parses time in MM:SS format.
func (p *Parser) parseMMSSFormat(matches []string) float64 {
	minutes, err := strconv.Atoi(matches[2])
	if err != nil {
		return 0
	}

	seconds, err := strconv.Atoi(matches[3])
	if err != nil {
		return 0
	}

	return float64(minutes*60 + seconds)
}

// parsePercentage converts percentage strings like "74%" to float.
func (p *Parser) parsePercentage(percentStr string) float64 {
	cleaned := strings.TrimSuffix(strings.TrimSpace(percentStr), "%")
	value, err := strconv.ParseFloat(cleaned, 64)
	if err == nil {
		return value
	}

	return 0
}

// parseNumber converts number strings to float.
func (p *Parser) parseNumber(numStr string) float64 {
	// Remove commas and other formatting
	cleaned := strings.ReplaceAll(strings.TrimSpace(numStr), ",", "")
	value, err := strconv.ParseFloat(cleaned, 64)
	if err == nil {
		return value
	}

	return 0
}

// parseProfileMetrics extracts profile-level statistics (endorsement, skill ratings).
func (p *Parser) parseProfileMetrics(doc *goquery.Document, stats *FullPlayerProfile) {
	p.parseEndorsementLevel(doc, stats)
	p.parseSkillRatings(doc, stats)
}

// parseEndorsementLevel extracts and sets the endorsement level and breakdown.
func (p *Parser) parseEndorsementLevel(doc *goquery.Document, stats *FullPlayerProfile) {
	endorsementImg := doc.Find(".Profile-playerSummary--endorsement")
	if endorsementImg.Length() == 0 {
		return
	}

	src, exists := endorsementImg.Attr("src")
	if !exists {
		return
	}

	level := extractEndorsementLevel(src)
	if level > 0 {
		stats.ProfileMetrics.Endorsement.Level = level
		slog.Debug("Parsed endorsement", "level", level)
	}

	// Try to parse endorsement breakdown if available
	breakdown := p.parseEndorsementBreakdown(doc)
	if breakdown != nil {
		stats.ProfileMetrics.Endorsement.Breakdown = breakdown
		slog.Debug("Parsed endorsement breakdown",
			"sportsmanship", breakdown.Sportsmanship,
			"teamwork", breakdown.Teamwork,
			"shotCaller", breakdown.ShotCaller)
	}
}

// parseEndorsementBreakdown attempts to extract endorsement breakdown data from various possible selectors.
func (p *Parser) parseEndorsementBreakdown(doc *goquery.Document) *EndorsementBreakdown {
	// Try multiple possible selectors for endorsement breakdown data
	selectors := []string{
		".Profile-endorsement--breakdown",
		".endorsement-breakdown",
		".Profile-playerSummary--endorsementBreakdown",
		"[data-endorsement-breakdown]",
		".endorsement-stats",
		".Profile-endorsement .endorsement-categories",
	}

	for _, selector := range selectors {
		if breakdown := p.tryParseEndorsementSelector(doc, selector); breakdown != nil {
			return breakdown
		}
	}

	// Log that breakdown data wasn't found (for debugging purposes)
	slog.Debug("Endorsement breakdown data not found in HTML structure")

	return nil
}

// tryParseEndorsementSelector attempts to parse endorsement breakdown from a specific selector.
func (p *Parser) tryParseEndorsementSelector(doc *goquery.Document, selector string) *EndorsementBreakdown {
	element := doc.Find(selector)
	if element.Length() == 0 {
		return nil
	}

	breakdown := &EndorsementBreakdown{}
	found := false

	// Try to find individual endorsement category elements
	categorySelectors := map[string]*int{
		"sportsmanship": &breakdown.Sportsmanship,
		"teamwork":      &breakdown.Teamwork,
		"shot-caller":   &breakdown.ShotCaller,
		"shotcaller":    &breakdown.ShotCaller,
		"leadership":    &breakdown.ShotCaller, // Alternative name
	}

	element.Find("*").Each(func(_ int, selection *goquery.Selection) {
		// Check data attributes
		for category, field := range categorySelectors {
			if attr, exists := selection.Attr("data-" + category); exists {
				if value := parseInt(attr); value > 0 {
					*field = value
					found = true
				}
			}
		}

		// Check class names
		classes := selection.AttrOr("class", "")
		text := strings.TrimSpace(selection.Text())

		for category, field := range categorySelectors {
			if strings.Contains(classes, category) && text != "" {
				if value := parseInt(text); value > 0 {
					*field = value
					found = true
				}
			}
		}
	})

	if found {
		return breakdown
	}

	return nil
}

// parseInt safely converts a string to an integer, returning 0 if parsing fails.
func parseInt(s string) int {
	// Remove any non-numeric characters except digits
	cleanStr := regexp.MustCompile(`\D`).ReplaceAllString(s, "")
	if cleanStr == "" {
		return 0
	}

	value, err := strconv.Atoi(cleanStr)
	if err != nil {
		return 0
	}

	return value
}

// parseSkillRatings extracts skill ratings for all platforms.
func (p *Parser) parseSkillRatings(doc *goquery.Document, stats *FullPlayerProfile) {
	for platform := range PlatformSelectors {
		platformRanks := p.parsePlatformRanks(doc, platform)
		if len(platformRanks) > 0 {
			stats.ProfileMetrics.SkillRatings[platform] = platformRanks
		}
	}
}

// parsePlatformRanks extracts rank information for a specific platform.
func (p *Parser) parsePlatformRanks(doc *goquery.Document, platform Platform) map[Role]RankInfo {
	platformRanks := make(map[Role]RankInfo)

	platformSelector := p.getPlatformSelector(platform)
	rankWrapper := doc.Find(platformSelector)
	slog.Debug("Rank wrapper search", "platform", platform, "selector", platformSelector, "found", rankWrapper.Length())

	if rankWrapper.Length() == 0 {
		return platformRanks
	}

	roleWrappers := rankWrapper.Find(".Profile-playerSummary--roleWrapper")
	roleWrappers.Each(func(_ int, roleEl *goquery.Selection) {
		role, rankInfo := p.parseRoleRank(roleEl, platform)
		if role != "" && rankInfo.Tier != "" {
			platformRanks[Role(role)] = rankInfo
		}
	})

	return platformRanks
}

// getPlatformSelector returns the CSS selector for a platform.
func (p *Parser) getPlatformSelector(platform Platform) string {
	if platform == PlatformPC {
		return ".Profile-playerSummary--rankWrapper.is-active"
	}

	return ".controller-view.Profile-playerSummary--rankWrapper"
}

// parseRoleRank extracts role and rank information from a role element.
func (p *Parser) parseRoleRank(roleEl *goquery.Selection, platform Platform) (string, RankInfo) {
	roleImg := roleEl.Find(".Profile-playerSummary--role img")
	if roleImg.Length() == 0 {
		return "", RankInfo{}
	}

	roleIconSrc, exists := roleImg.Attr("src")
	if !exists {
		return "", RankInfo{}
	}

	role := extractRoleFromIcon(roleIconSrc)
	if role == "" {
		return "", RankInfo{}
	}

	rankInfo := extractRankInfo(roleEl)
	if rankInfo.Tier != "" {
		slog.Debug("Parsed rank",
			"platform", platform,
			"role", role,
			"tier", rankInfo.Tier,
			"division", rankInfo.Division)
	}

	return role, rankInfo
}

// extractEndorsementLevel extracts endorsement level from icon URL.
func extractEndorsementLevel(iconURL string) int {
	// URL format: "...endorsement/2-8b9f0faa25.svg" -> level 2
	re := regexp.MustCompile(`endorsement/(\d+)-`)
	matches := re.FindStringSubmatch(iconURL)
	if len(matches) >= 2 {
		level, err := strconv.Atoi(matches[1])
		if err == nil {
			return level
		}
	}

	return 0
}

// extractRoleFromIcon determines role from icon URL.
func extractRoleFromIcon(iconURL string) string {
	switch {
	case strings.Contains(iconURL, "offense"):
		return "damage"
	case strings.Contains(iconURL, "support"):
		return "support"
	case strings.Contains(iconURL, "tank"):
		return "tank"
	default:
		return ""
	}
}

// extractRankInfo extracts tier and division from rank elements.
func extractRankInfo(roleEl *goquery.Selection) RankInfo {
	var rankInfo RankInfo

	// Find rank images
	rankImages := roleEl.Find(".Profile-playerSummary--rank")

	rankImages.Each(func(_ int, img *goquery.Selection) {
		if src, exists := img.Attr("src"); exists {
			// Extract tier (Bronze, Silver, Gold, etc.)
			if tier := extractTierFromURL(src); tier != "" && rankInfo.Tier == "" {
				rankInfo.Tier = tier
			}
			// Extract division (1-5)
			if division := extractDivisionFromURL(src); division > 0 && rankInfo.Division == 0 {
				rankInfo.Division = division
			}
		}
	})

	return rankInfo
}

// extractTierFromURL extracts tier name from rank icon URL.
func extractTierFromURL(url string) string {
	// URL format: "...Rank_DiamondTier-d775ca9c43.png"
	re := regexp.MustCompile(`Rank_(\w+)Tier`)
	matches := re.FindStringSubmatch(url)
	if len(matches) >= 2 {
		return matches[1]
	}

	return ""
}

// extractDivisionFromURL extracts division number from division icon URL.
func extractDivisionFromURL(url string) int {
	// URL format: "...TierDivision_2-4376c89b41.png"
	re := regexp.MustCompile(`TierDivision_(\d+)`)
	matches := re.FindStringSubmatch(url)
	if len(matches) >= 2 {
		div, err := strconv.Atoi(matches[1])
		if err == nil {
			return div
		}
	}

	return 0
}

// parseAllHeroesStats parses aggregated statistics across all heroes.
func (p *Parser) parseAllHeroesStats(gameModeView *goquery.Selection) *AllHeroesStats {
	stats := &AllHeroesStats{}

	// Parse each metric using the hex IDs from AllHeroesHexIDs
	for metricKey, hexID := range AllHeroesHexIDs {
		value := p.parseHeroMetricFromGameMode(gameModeView, metricKey, hexID)
		if value != 0 {
			p.assignValueToAllHeroesStats(stats, metricKey, value)
		}
	}

	p.logAllHeroesStatsSummary(stats)

	return stats
}

// parseHeroMetricFromGameMode extracts and parses a single metric value.
func (p *Parser) parseHeroMetricFromGameMode(gameModeView *goquery.Selection, metricKey, hexID string) float64 {
	selector := fmt.Sprintf("[data-category-id='%s'] .Profile-progressBar-description", hexID)
	element := gameModeView.Find(selector)

	if element.Length() == 0 {
		return 0
	}

	valueText := element.First().Text()
	value := p.parseMetricValue(valueText, metricKey)

	slog.Debug("All Heroes metric parsed",
		"metric", metricKey,
		"hex_id", hexID,
		"value_text", valueText,
		"parsed_value", value)

	return value
}

// assignValueToAllHeroesStats assigns the parsed value to the appropriate field.
// allHeroesStatAssigner defines a function type for assigning values to AllHeroesStats.
type allHeroesStatAssigner func(*AllHeroesStats, float64)

// getAllHeroesStatAssigners returns a map of metric keys to their assignment functions.
func (p *Parser) getAllHeroesStatAssigners() map[string]allHeroesStatAssigner {
	return map[string]allHeroesStatAssigner{
		"time_played": func(stats *AllHeroesStats, value float64) {
			stats.TotalTimePlayed = int64(value)
		},
		GamesWonMetric: func(stats *AllHeroesStats, value float64) {
			stats.TotalGamesWon = int(value)
		},
		WinPercentageMetric: func(stats *AllHeroesStats, value float64) {
			stats.OverallWinPercentage = value
		},
		WeaponAccuracyMetric: func(stats *AllHeroesStats, value float64) {
			stats.WeaponAccuracy = value
		},
		EliminationsPerLifeMetric: func(stats *AllHeroesStats, value float64) {
			stats.EliminationsPerLife = value
		},
		KillStreakBestMetric: func(stats *AllHeroesStats, value float64) {
			stats.KillStreakBest = int(value)
		},
		MultikillBestMetric: func(stats *AllHeroesStats, value float64) {
			stats.MultikillBest = int(value)
		},
		EliminationsPer10MinMetric: func(stats *AllHeroesStats, value float64) {
			stats.EliminationsPer10Min = value
		},
		DeathsPer10MinMetric: func(stats *AllHeroesStats, value float64) {
			stats.DeathsPer10Min = value
		},
		FinalBlowsPer10MinMetric: func(stats *AllHeroesStats, value float64) {
			stats.FinalBlowsPer10Min = value
		},
		SoloKillsPer10MinMetric: func(stats *AllHeroesStats, value float64) {
			stats.SoloKillsPer10Min = value
		},
		ObjectiveKillsPer10MinMetric: func(stats *AllHeroesStats, value float64) {
			stats.ObjectiveKillsPer10Min = value
		},
		ObjectiveTimePer10MinMetric: func(stats *AllHeroesStats, value float64) {
			stats.ObjectiveTimePer10Min = value
		},
		HeroDamagePer10MinMetric: func(stats *AllHeroesStats, value float64) {
			stats.HeroDamagePer10Min = value
		},
		HealingPer10MinMetric: func(stats *AllHeroesStats, value float64) {
			stats.HealingPer10Min = value
		},
	}
}

func (p *Parser) assignValueToAllHeroesStats(stats *AllHeroesStats, metricKey string, value float64) {
	assigners := p.getAllHeroesStatAssigners()
	if assigner, exists := assigners[metricKey]; exists {
		assigner(stats, value)
	}
}

// logAllHeroesStatsSummary logs a summary of parsed all heroes statistics.
func (p *Parser) logAllHeroesStatsSummary(stats *AllHeroesStats) {
	slog.Info("All Heroes stats parsed",
		"time_played", stats.TotalTimePlayed,
		"games_won", stats.TotalGamesWon,
		"win_percentage", stats.OverallWinPercentage)
}

// parseMetricValue parses a metric value based on its type.
func (p *Parser) parseMetricValue(valueText, metricType string) float64 {
	valueText = strings.TrimSpace(valueText)

	switch {
	case strings.Contains(metricType, "time_played"):
		return p.parseTimeToSeconds(valueText)
	case strings.Contains(metricType, "percentage") || strings.Contains(metricType, "accuracy"):
		return p.parsePercentage(valueText)
	default:
		return p.parseNumber(valueText)
	}
}

// extractDetailedHeroMetrics extracts hero-specific detailed statistics using OverFast API logic.
// Searches blz-section.stats.{gamemode}-view sections for span.stats-container.option-{N} elements.
func (p *Parser) extractDetailedHeroMetrics(heroEl *goquery.Selection, heroID, heroName string, stats *HeroStats) {
	slog.Debug("Extracting detailed hero metrics using OverFast API logic",
		"hero_id", heroID, "hero_name", heroName)

	// Get document root to search for stats sections
	doc := heroEl.Closest("html")

	// Look for blz-section.stats sections with gamemode views
	gamemodeViews := []string{".quickPlay-view", ".competitive-view"}

	extractedCount := 0
	for _, gamemodeView := range gamemodeViews {
		count := p.processGameModeStats(doc, gamemodeView, heroID, stats)
		extractedCount += count
	}

	slog.Info("Hero metrics extraction completed",
		"hero_id", heroID,
		"hero_name", heroName,
		"extracted_metrics_count", extractedCount)
}

// processGameModeStats processes stats for a specific gamemode view.
func (p *Parser) processGameModeStats(doc *goquery.Selection, gamemodeView, heroID string, stats *HeroStats) int {
	statsSection := doc.Find(HeroSelectors.BlzStatsSection + gamemodeView)
	if statsSection.Length() == 0 {
		slog.Debug("No stats section found for gamemode view",
			"gamemode_view", gamemodeView, "hero_id", heroID)
		return 0
	}

	slog.Debug("Found stats section for gamemode",
		"gamemode_view", gamemodeView, "hero_id", heroID)

	// Find all span.stats-container elements (OverFast API style)
	statsContainers := p.findStatsContainers(statsSection)
	if statsContainers.Length() == 0 {
		slog.Debug("No stats containers found in section",
			"gamemode_view", gamemodeView, "hero_id", heroID)
		return 0
	}

	slog.Debug("Found stats containers in section",
		"gamemode_view", gamemodeView,
		"containers_count", statsContainers.Length(),
		"hero_id", heroID)

	extractedCount := 0
	statsContainers.Each(func(containerIdx int, container *goquery.Selection) {
		count := p.processStatsContainer(container, containerIdx, heroID, stats)
		extractedCount += count
	})

	return extractedCount
}

// findStatsContainers finds stats-container elements using multiple selectors.
func (p *Parser) findStatsContainers(statsSection *goquery.Selection) *goquery.Selection {
	// Look for pattern: span.stats-container.option-{N}
	statsContainers := statsSection.Find(HeroSelectors.StatsContainer)

	// Also try broader search if no containers found with specific selector
	if statsContainers.Length() == 0 {
		statsContainers = statsSection.Find("span.stats-container")
	}

	return statsContainers
}

// processStatsContainer processes a single stats container for hero-specific metrics.
func (p *Parser) processStatsContainer(container *goquery.Selection, containerIdx int, heroID string, stats *HeroStats) int {
	// Check if this container has a select element (OverFast API pattern)
	selectEl := container.Find("select")
	if selectEl.Length() == 0 {
		slog.Debug("Container has no select element, skipping",
			"container_idx", containerIdx, "hero_id", heroID)
		return 0
	}

	// Check if this is hero-specific data by looking for our hero in the select options
	if !p.isHeroInContainer(selectEl, heroID) {
		slog.Debug("Hero not found in container options",
			"container_idx", containerIdx, "hero_id", heroID)
		return 0
	}

	slog.Debug("Found hero-specific container",
		"container_idx", containerIdx, "hero_id", heroID)

	// Extract stat items using OverFast API selectors
	statItems := container.Find(HeroSelectors.StatItem)
	if statItems.Length() == 0 {
		slog.Debug("No stat items found in container",
			"container_idx", containerIdx, "hero_id", heroID)
		return 0
	}

	slog.Debug("Processing stat items",
		"container_idx", containerIdx,
		"stat_items_count", statItems.Length(),
		"hero_id", heroID)

	extractedCount := 0
	statItems.Each(func(itemIdx int, item *goquery.Selection) {
		if p.processStatItem(item, itemIdx, heroID, stats) {
			extractedCount++
		}
	})

	return extractedCount
}

// isHeroInContainer checks if the specified hero is available in the container's select options.
func (p *Parser) isHeroInContainer(selectEl *goquery.Selection, heroID string) bool {
	heroFound := false
	selectEl.Find("option").Each(func(_ int, option *goquery.Selection) {
		optionValue, exists := option.Attr("value")
		if exists && optionValue == heroID {
			heroFound = true
		}
	})

	return heroFound
}

// processStatItem processes a single stat item and stores it in the stats structure.
func (p *Parser) processStatItem(item *goquery.Selection, itemIdx int, heroID string, stats *HeroStats) bool {
	// Extract stat name and value using OverFast API selectors
	statName := strings.TrimSpace(item.Find(HeroSelectors.StatName).Text())
	statValue := strings.TrimSpace(item.Find(HeroSelectors.StatValue).Text())

	if statName == "" || statValue == "" {
		slog.Debug("Empty stat name or value, skipping",
			"item_idx", itemIdx, "stat_name", statName, "stat_value", statValue)
		return false
	}

	// Convert stat name to metric key using StringToSnakeCase from value_parser.go
	metricKey := StringToSnakeCase(statName)

	// Parse the value using ParseValue from value_parser.go
	parsedValue := ParseValue(statValue)

	// Convert parsed value to float64 for storage
	floatValue, ok := p.convertToFloat(parsedValue, statName)
	if !ok {
		return false
	}

	// Initialize metrics map if needed
	if stats.Metrics == nil {
		stats.Metrics = make(map[string]float64)
	}

	// Store the metric
	stats.Metrics[metricKey] = floatValue

	slog.Debug("Extracted hero metric",
		"hero_id", heroID,
		"stat_name", statName,
		"metric_key", metricKey,
		"raw_value", statValue,
		"parsed_value", parsedValue,
		"stored_value", floatValue)

	return true
}

// convertToFloat converts various parsed value types to float64.
func (p *Parser) convertToFloat(parsedValue interface{}, statName string) (float64, bool) {
	switch value := parsedValue.(type) {
	case float64:
		return value, true
	case int:
		return float64(value), true
	case string:
		// Handle string values that might be numeric
		val, err := strconv.ParseFloat(value, 64)
		if err == nil {
			return val, true
		}

		slog.Debug("Could not convert string value to float",
			"stat_name", statName, "value", value)
		return 0, false
	default:
		slog.Debug("Unsupported value type for metric storage",
			"stat_name", statName, "value", parsedValue, "type", fmt.Sprintf("%T", parsedValue))
		return 0, false
	}
}

// ExampleParser provides example usage function for testing.
func ExampleParser() {
	// This would be called from the main ow-exporter application
	_ = NewParser()

	// Read one of our saved profiles
	// html := readFile("/Users/lex/git/github.com/lexfrei/tools/tmp/profile_de5bb4aca17492e0.html")
	// stats, err := parser.ParseProfile(html, "LexFrei")
	// if err != nil {
	//     log.Fatal(err)
	// }

	slog.Info("Parser ready for integration into ow-exporter")
}
