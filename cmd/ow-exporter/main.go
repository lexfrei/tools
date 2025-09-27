package main

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"time"
	"unicode"

	"github.com/cockroachdb/errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	defaultPort = "9420"

	// MinArgsWithParam command argument counts.
	MinArgsWithParam = 3 // program name + command + parameter

	// DefaultServerPort sets the default Prometheus server port.
	DefaultServerPort = "9090"
	// DefaultHTTPTimeout sets the default HTTP client timeout.

	// maxSampleMetricsCount sets maximum number of sample metrics to show.
	maxSampleMetricsCount = 3
	DefaultHTTPTimeout    = 30 * time.Second
	// MaxHeroesDisplayed sets the maximum number of heroes displayed in logs.
	MaxHeroesDisplayed = 3
	// DebugFilePermissions sets file permissions for debug files.
	DebugFilePermissions = 0o600
)

// toTitle converts the first character of each word to uppercase (replacement for deprecated strings.Title).
func toTitle(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return r
		}

		return unicode.ToTitle(r)
	}, str)
}

// handleCommand processes CLI commands and returns true if a command was handled.
// commandHandler defines a function type for handling commands.
type commandHandler func([]string)

// getCommandHandlers returns a map of command names to their handlers.
func getCommandHandlers() map[string]commandHandler {
	return map[string]commandHandler{
		"parse-profile":   handleParseProfile,
		"parse-battletag": handleParseBattleTag,
		"config":          handleConfigCommand,
		"metrics":         handleMetricsCommand,
		"server":          handleServerCommand,
		"test-api":        handleTestAPI,
		"test-headless":   handleTestHeadless,
		"test-parser":     handleTestParser,
	}
}

func handleCommand(args []string) bool {
	if len(args) <= 1 {
		return false
	}

	handlers := getCommandHandlers()
	if handler, exists := handlers[args[1]]; exists {
		handler(args)

		return true
	}

	printMainUsage()
	os.Exit(1)

	return true
}

// handleParseProfile handles the parse-profile command.
func handleParseProfile(_ []string) {
	runPoC()
}

// handleTestAPI handles the test-api command.
func handleTestAPI(_ []string) {
	runAPITests()
}

// handleTestHeadless handles the test-headless command.
func handleTestHeadless(_ []string) {
	runHeadlessTests()
}

// handleTestParser handles the test-parser command.
func handleTestParser(_ []string) {
	runParserTests()
}

// handleParseBattleTag handles the parse-battletag command.
func handleParseBattleTag(args []string) {
	if len(args) < MinArgsWithParam {
		slog.Error("Usage error", "command", "parse-battletag", "message", "BattleTag argument required")
		os.Exit(1)
	}
	runParseBattleTag(args[2])
}

// handleConfigCommand handles the config command.
func handleConfigCommand(args []string) {
	if len(args) < MinArgsWithParam {
		printConfigUsage()
		os.Exit(1)
	}
	runConfigCommand(args[2:])
}

// handleMetricsCommand handles the metrics command.
func handleMetricsCommand(args []string) {
	if len(args) < MinArgsWithParam {
		printMetricsUsage()
		os.Exit(1)
	}
	runMetricsCommand(args[2:])
}

// handleServerCommand handles the server command.
func handleServerCommand(args []string) {
	port := DefaultServerPort
	if len(args) > 2 {
		port = args[2]
	}
	startPrometheusServer(port)
}

func main() {
	// Setup structured logging
	programLevel := new(slog.LevelVar)
	programLevel.Set(slog.LevelDebug) // Set to Debug level for troubleshooting
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: programLevel})))

	// Initialize systems
	err := initConfig()
	if err != nil {
		slog.Error("Failed to initialize config", "error", err)
		os.Exit(1)
	}

	initRuntimeMetrics()

	// Initialize Prometheus metrics
	initPrometheusMetrics()

	// Handle CLI commands
	if handleCommand(os.Args) {
		return
	}

	// Start automatic profile parsing
	go startPeriodicParsing()

	slog.Info("ow-exporter starting", "version", "development")

	// Create Echo instance
	server := echo.New()
	server.HideBanner = true

	// Middleware
	server.Use(middleware.Logger())
	server.Use(middleware.Recover())
	server.Use(middleware.CORS())

	// Routes
	setupRoutes(server)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Start server
	go func() {
		slog.Info("starting HTTP server", "port", port)
		err := server.Start(":" + port)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("failed to start server", "error", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	slog.Info("shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = server.Shutdown(ctx)
	if err != nil {
		slog.Error("server forced to shutdown", "error", err)
	}

	slog.Info("server shutdown complete")
}

func setupRoutes(server *echo.Echo) {
	// Health check
	server.GET("/health", healthHandler)

	// API routes
	api := server.Group("/api")
	api.GET("/users", listUsersHandler)
	api.POST("/users", createUserHandler)
	api.GET("/users/:username", getUserHandler)
	api.PUT("/users/:username", updateUserHandler)
	api.DELETE("/users/:username", deleteUserHandler)

	// Prometheus metrics
	server.GET("/metrics", metricsHandler)

	// Parse endpoints
	api.POST("/parse", parseAllPlayersHandler)
	api.POST("/parse/:username", parsePlayerHandler)

	// Development info
	server.GET("/", indexHandler)
}

// HTTP handlers.
func healthHandler(c echo.Context) error {
	err := c.JSON(http.StatusOK, map[string]string{
		"status":  "ok",
		"service": "ow-exporter",
		"version": "development",
	})
	if err != nil {
		return errors.Wrap(err, "failed to write health response")
	}

	return nil
}

func indexHandler(c echo.Context) error {
	err := c.JSON(http.StatusOK, map[string]interface{}{
		"service": "ow-exporter",
		"version": "development",
		"status":  "in development",
		"endpoints": map[string]string{
			"health":  "/health",
			"metrics": "/metrics",
			"users":   "/api/users",
		},
		"documentation": "https://github.com/lexfrei/tools/issues/439",
	})
	if err != nil {
		return errors.Wrap(err, "failed to write index response")
	}

	return nil
}

func listUsersHandler(ctx echo.Context) error {
	players := getAllPlayers()
	users := make([]map[string]interface{}, len(players))

	for i, player := range players {
		users[i] = map[string]interface{}{
			"battleTag":    player.BattleTag,
			"resolvedUrl":  player.ResolvedURL,
			"lastResolved": player.LastResolved,
		}
	}

	err := ctx.JSON(http.StatusOK, map[string]interface{}{
		"users": users,
		"total": len(players),
	})
	if err != nil {
		return errors.Wrap(err, "failed to write users response")
	}

	return nil
}

// CreateUserRequest represents the request body for creating a user.
type CreateUserRequest struct {
	BattleTag   string `json:"battleTag"             validate:"required"`
	ResolvedURL string `json:"resolvedUrl,omitempty"`
}

func createUserHandler(ctx echo.Context) error {
	var req CreateUserRequest

	err := ctx.Bind(&req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	if req.BattleTag == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "battleTag is required")
	}

	// If no resolved URL provided, try to resolve the BattleTag
	if req.ResolvedURL == "" {
		resolvedURL, resolveErr := getOrResolveURL(req.BattleTag)
		if resolveErr != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "failed to resolve BattleTag")
		}
		req.ResolvedURL = resolvedURL
	}

	err = addPlayerToConfig(req.BattleTag, req.ResolvedURL)
	if err != nil {
		return errors.Wrap(err, "failed to add player to config")
	}

	err = ctx.JSON(http.StatusCreated, map[string]interface{}{
		"message":     "User created successfully",
		"battleTag":   req.BattleTag,
		"resolvedUrl": req.ResolvedURL,
	})
	if err != nil {
		return errors.Wrap(err, "failed to send create user response")
	}

	return nil
}

func getUserHandler(ctx echo.Context) error {
	battleTag := ctx.Param("username")

	player := findPlayerByBattleTag(battleTag)
	if player == nil {
		return echo.NewHTTPError(http.StatusNotFound, "player not found")
	}

	err := ctx.JSON(http.StatusOK, map[string]interface{}{
		"battleTag":    player.BattleTag,
		"resolvedUrl":  player.ResolvedURL,
		"lastResolved": player.LastResolved,
	})
	if err != nil {
		return errors.Wrap(err, "failed to send get user response")
	}

	return nil
}

// UpdateUserRequest represents the request body for updating a user.
type UpdateUserRequest struct {
	ResolvedURL string `json:"resolvedUrl,omitempty"`
}

func updateUserHandler(ctx echo.Context) error {
	battleTag := ctx.Param("username")
	var req UpdateUserRequest

	err := ctx.Bind(&req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	// Check if player exists
	player := findPlayerByBattleTag(battleTag)
	if player == nil {
		return echo.NewHTTPError(http.StatusNotFound, "player not found")
	}

	// If no resolved URL provided, try to resolve the BattleTag
	if req.ResolvedURL == "" {
		resolvedURL, resolveErr := getOrResolveURL(battleTag)
		if resolveErr != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "failed to resolve BattleTag")
		}
		req.ResolvedURL = resolvedURL
	}

	err = updatePlayerURL(battleTag, req.ResolvedURL)
	if err != nil {
		return errors.Wrap(err, "failed to update player")
	}

	err = ctx.JSON(http.StatusOK, map[string]interface{}{
		"message":     "User updated successfully",
		"battleTag":   battleTag,
		"resolvedUrl": req.ResolvedURL,
	})
	if err != nil {
		return errors.Wrap(err, "failed to send update user response")
	}

	return nil
}

func deleteUserHandler(ctx echo.Context) error {
	battleTag := ctx.Param("username")

	err := removePlayerFromConfig(battleTag)
	if err != nil {
		if errors.Is(err, ErrPlayerNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "player not found")
		}

		return errors.Wrap(err, "failed to remove player")
	}

	err = ctx.JSON(http.StatusOK, map[string]interface{}{
		"message":   "User deleted successfully",
		"battleTag": battleTag,
	})
	if err != nil {
		return errors.Wrap(err, "failed to send delete user response")
	}

	return nil
}

func metricsHandler(ctx echo.Context) error {
	// Update Prometheus metrics with fresh data
	updatePrometheusMetrics()

	// Serve Prometheus metrics directly using Echo
	handler := promhttp.Handler()
	handler.ServeHTTP(ctx.Response(), ctx.Request())

	return nil
}

// parseAllPlayersHandler triggers parsing of all configured players.
func parseAllPlayersHandler(ctx echo.Context) error {
	go func() {
		slog.Info("🔄 Manual parse triggered via API")
		parseAllPlayers()
	}()

	err := ctx.JSON(http.StatusAccepted, map[string]interface{}{
		"message": "Parsing started for all configured players",
		"status":  "in_progress",
	})
	if err != nil {
		return errors.Wrap(err, "failed to send parse response")
	}

	return nil
}

// parsePlayerHandler triggers parsing of a specific player.
func parsePlayerHandler(ctx echo.Context) error {
	battleTag := ctx.Param("username")

	go func() {
		slog.Info("🔄 Manual parse triggered for player", "battletag", battleTag)
		success := parsePlayerSafely(battleTag)
		if success {
			slog.Info("✅ Manual parse completed", "battletag", battleTag)
		} else {
			slog.Error("❌ Manual parse failed", "battletag", battleTag)
		}
	}()

	err := ctx.JSON(http.StatusAccepted, map[string]interface{}{
		"message":   "Parsing started for player",
		"battleTag": battleTag,
		"status":    "in_progress",
	})
	if err != nil {
		return errors.Wrap(err, "failed to send parse response")
	}

	return nil
}

// startPeriodicParsing starts automatic parsing of all configured players.
func startPeriodicParsing() {
	const parseInterval = 30 * time.Minute

	slog.Info("🔄 Starting periodic profile parsing", "interval", parseInterval.String())

	// Initial parse on startup
	parseAllPlayers()

	// Set up periodic parsing
	ticker := time.NewTicker(parseInterval)
	defer ticker.Stop()

	for range ticker.C {
		parseAllPlayers()
	}
}

// parseAllPlayers parses all players from config.
func parseAllPlayers() {
	players := getAllPlayers()
	if len(players) == 0 {
		slog.Info("No players configured for parsing")

		return
	}

	slog.Info("🎯 Starting batch profile parsing", "players", len(players))

	successCount := 0
	for _, player := range players {
		if parsePlayerSafely(player.BattleTag) {
			successCount++
		}
		// Small delay between requests to be respectful to the server
		time.Sleep(2 * time.Second)
	}

	slog.Info("✅ Batch parsing completed",
		"success", successCount,
		"total", len(players),
		"failed", len(players)-successCount)
}

// parsePlayerSafely parses a single player with error handling.
func parsePlayerSafely(battleTag string) bool {
	slog.Info("🎯 Parsing profile", "battletag", battleTag)

	// Get or resolve URL
	profileURL, err := getOrResolveURL(battleTag)
	if err != nil {
		slog.Error("❌ Failed to resolve BattleTag", "battletag", battleTag, "error", err)

		return false
	}

	// Parse the profile
	parser := NewParser()
	profile, err := parser.ParseProfile(fetchProfileHTML(profileURL), battleTag)
	if err != nil {
		slog.Error("❌ Failed to parse profile", "battletag", battleTag, "error", err)

		return false
	}

	// Create/update runtime metrics
	existingMetrics, exists := getPlayerMetrics(battleTag)
	if exists {
		updatePlayerFromProfile(existingMetrics, profile)
		slog.Info("📊 Updated player metrics", "battletag", battleTag)
	} else {
		newMetrics := createPlayerMetrics(battleTag, profile)
		setPlayerMetrics(battleTag, newMetrics)
		slog.Info("📊 Created new player metrics", "battletag", battleTag)
	}

	// Log key metrics
	if profile != nil {
		slog.Info("🎖️ Endorsement", "battletag", battleTag, "level", profile.ProfileMetrics.Endorsement.Level)

		// Count heroes
		totalHeroes := 0
		for _, platform := range profile.Platforms {
			for _, gameMode := range platform.GameModes {
				totalHeroes += len(gameMode.Heroes)
			}
		}
		slog.Info("👤 Heroes parsed", "battletag", battleTag, "count", totalHeroes)
	}

	return true
}

// runPoC runs the Proof of Concept profile parser.
func runPoC() {
	profileURL := "https://overwatch.blizzard.com/en-us/career/" +
		"de5bb4aca17492e0bba120a1d1%7Ca92a11ef8d304356fccfff8df12e1dc6/"

	if len(os.Args) > 2 {
		profileURL = os.Args[2]
	}

	slog.Info("🎯 PoC: Parsing Overwatch Profile")
	slog.Info("📋 Profile URL", "url", profileURL)

	// Fetch profile HTML
	html, err := fetchProfile(profileURL)
	if err != nil {
		slog.Error("❌ Error fetching profile", "error", err)
		os.Exit(1)
	}

	slog.Info("✅ Successfully fetched profile HTML", "bytes", len(html))

	// Debug: Save HTML to file for inspection
	err = os.WriteFile("debug_profile.html", []byte(html), DebugFilePermissions)
	if err != nil {
		slog.Warn("⚠️  Warning: could not save debug HTML", "error", err)
	} else {
		slog.Info("💾 Saved HTML to debug_profile.html for inspection")
	}

	// Parse profile using our parser
	parser := NewParser()
	username := extractUsernameFromURL(profileURL)

	slog.Info("👤 Detected username", "username", username)
	slog.Info("🔍 Parsing profile data...")

	stats, err := parser.ParseProfile(html, username)
	if err != nil {
		slog.Error("❌ Error parsing profile", "error", err)
		os.Exit(1)
	}

	// Pretty print results
	printPrettyResults(stats)
}

// fetchProfile fetches the profile HTML with proper headers.
func fetchProfile(profileURL string) (string, error) {
	client := &http.Client{
		Timeout: DefaultHTTPTimeout,
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, profileURL, http.NoBody)
	if err != nil {
		return "", errors.Wrap(err, "failed to create HTTP request")
	}

	// Add browser-like headers
	req.Header.Set("User-Agent",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	// Don't request compression to avoid parsing issues
	req.Header.Set("DNT", "1")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	resp, err := client.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "failed to perform HTTP request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.Wrapf(ErrHTTPError, "%d %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "failed to read response body")
	}

	return string(body), nil
}

// extractUsernameFromURL extracts username from Overwatch profile URL.
func extractUsernameFromURL(profileURL string) string {
	// Extract from URL like: /en-us/career/pc/Username/ or
	// /en-us/career/de5bb4aca17492e0bba120a1d1%7Ca92a11ef8d304356fccfff8df12e1dc6/
	parts := strings.Split(profileURL, "/")
	if len(parts) >= 2 {
		username := parts[len(parts)-2] // Get second to last part
		if strings.Contains(username, "%7C") {
			// Handle encoded URLs - decode BattleTag with encoded |
			decoded, err := url.QueryUnescape(username)
			if err != nil {
				return "BattleTag-User"
			}
			// Replace | with # for BattleTag format
			return strings.ReplaceAll(decoded, "|", "#")
		}

		return username
	}

	return "Unknown"
}

// printPrettyResults prints the parsed results in a nice format.
func printPrettyResults(stats *FullPlayerProfile) {
	slog.Info("🎮 ===========================================")
	slog.Info("📊 OVERWATCH PROFILE PARSING RESULTS")
	slog.Info("===========================================")

	slog.Info("👤 Username", "username", stats.Username)
	slog.Info("🕐 Last Update", "timestamp", stats.LastUpdate.Format(time.RFC3339))
	slog.Info("🎯 Platforms Found", "count", len(stats.Platforms))

	for platform, platformStats := range stats.Platforms {
		slog.Info("🖥️  Platform", "platform", strings.ToUpper(string(platform)))
		slog.Info("📈 Game Modes", "count", len(platformStats.GameModes))

		for gameMode, gameModeStats := range platformStats.GameModes {
			slog.Info("  🎮 Game Mode", "mode", toTitle(string(gameMode)), "heroes", len(gameModeStats.Heroes))

			// Show first few heroes as examples
			count := 0
			for heroID, heroStats := range gameModeStats.Heroes {
				if count >= MaxHeroesDisplayed {
					break
				}
				slog.Info("    • Hero", "name", heroStats.HeroName, "id", heroID, "metrics", len(heroStats.Metrics))

				// Show a few metrics
				metricCount := 0
				for metricName, value := range heroStats.Metrics {
					if metricCount >= 2 {
						break
					}
					slog.Info("      - Metric", "name", metricName, "value", value)
					metricCount++
				}
				count++
			}

			if len(gameModeStats.Heroes) > MaxHeroesDisplayed {
				slog.Info("    ... and more heroes", "additional", len(gameModeStats.Heroes)-MaxHeroesDisplayed)
			}
		}
	}

	// Pretty print full JSON
	slog.Info("📋 FULL JSON OUTPUT:")
	slog.Info("====================")

	jsonData, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		slog.Error("❌ Error marshaling JSON", "error", err)

		return
	}

	slog.Info(string(jsonData))
	slog.Info("✅ PoC Complete! Profile successfully parsed and displayed.")
}

// CLI command implementations.

// runParseBattleTag parses a profile by BattleTag.
func runParseBattleTag(battleTag string) {
	slog.Info("🎯 Parsing Overwatch profile by BattleTag", "battletag", battleTag)

	// Get or resolve URL
	profileURL, err := getOrResolveURL(battleTag)
	if err != nil {
		slog.Error("❌ Failed to resolve BattleTag", "battletag", battleTag, "error", err)
		os.Exit(1)
	}

	slog.Info("✅ Resolved profile URL", "battletag", battleTag, "url", profileURL)

	// Parse the profile
	parser := NewParser()
	profile, err := parser.ParseProfile(fetchProfileHTML(profileURL), battleTag)
	if err != nil {
		slog.Error("❌ Failed to parse profile", "error", err)
		os.Exit(1)
	}

	// Create/update runtime metrics
	existingMetrics, exists := getPlayerMetrics(battleTag)
	if exists {
		updatePlayerFromProfile(existingMetrics, profile)
		slog.Info("📊 Updated existing player metrics", "battletag", battleTag)
	} else {
		newMetrics := createPlayerMetrics(battleTag, profile)
		setPlayerMetrics(battleTag, newMetrics)
		slog.Info("📊 Created new player metrics", "battletag", battleTag)
	}

	// Display results
	printParsedProfile(profile)
	slog.Info("✅ Parse complete! Metrics stored in runtime.")
}

// fetchProfileHTML fetches HTML content from URL.
func fetchProfileHTML(profileURL string) string {
	html, err := fetchProfile(profileURL)
	if err != nil {
		slog.Error("Failed to fetch profile HTML", "error", err)
		os.Exit(1)
	}

	return html
}

// printParsedProfile prints profile information.
func printParsedProfile(profile *FullPlayerProfile) {
	slog.Info("📋 PROFILE INFORMATION")
	slog.Info("👤 Player", "battletag", profile.BattleTag, "title", profile.PlayerTitle)
	slog.Info("🎖️ Endorsement", "level", profile.ProfileMetrics.Endorsement.Level)

	// Print skill ratings
	for platform, roles := range profile.ProfileMetrics.SkillRatings {
		for role, rank := range roles {
			slog.Info("🏆 Skill Rating", "platform", platform, "role", role, "tier", rank.Tier, "division", rank.Division)
		}
	}

	// Print platform/gamemode summary
	for platform, platformStats := range profile.Platforms {
		for gameMode, gameModeStats := range platformStats.GameModes {
			slog.Info("🎮 Heroes", "platform", platform, "gamemode", gameMode, "count", len(gameModeStats.Heroes))
		}
	}
}

// runConfigCommand handles config subcommands.
func runConfigCommand(args []string) {
	if len(args) == 0 {
		printConfigUsage()

		return
	}

	switch args[0] {
	case "list-players":
		listPlayersCommand()
	case "add-player":
		if len(args) < 2 {
			slog.Error("Usage error", "command", "config add-player", "message", "BattleTag argument required")
			os.Exit(1)
		}
		addPlayerCommand(args[1])
	case "remove-player":
		if len(args) < 2 {
			slog.Error("Usage error", "command", "config remove-player", "message", "BattleTag argument required")
			os.Exit(1)
		}
		removePlayerCommand(args[1])
	case "resolve-all":
		forceResolve := len(args) > 1 && args[1] == "--force"
		resolveAllCommand(forceResolve)
	default:
		slog.Error("Unknown config command", "command", args[0])
		printConfigUsage()
		os.Exit(1)
	}
}

// runMetricsCommand handles metrics subcommands.
func runMetricsCommand(args []string) {
	if len(args) == 0 {
		printMetricsUsage()

		return
	}

	switch args[0] {
	case "show":
		if len(args) < 2 {
			slog.Error("Usage error", "command", "metrics show", "message", "BattleTag argument required")
			os.Exit(1)
		}
		showMetricsCommand(args[1])
	case "list":
		listMetricsCommand()
	case "clear":
		clearMetricsCommand()
	case "stats":
		statsCommand()
	default:
		slog.Error("Unknown metrics command", "command", args[0])
		printMetricsUsage()
		os.Exit(1)
	}
}

// Config command implementations.
func listPlayersCommand() {
	players := getAllPlayers()
	if len(players) == 0 {
		slog.Info("No players configured.")

		return
	}

	slog.Info("Configured players", "count", len(players))
	for _, player := range players {
		status := "not_resolved"
		lastResolved := ""
		if player.ResolvedURL != "" {
			status = "resolved"
			if player.LastResolved != nil {
				lastResolved = player.LastResolved.Format("2006-01-02 15:04")
			}
		}
		slog.Info("Player", "battletag", player.BattleTag, "status", status, "last_resolved", lastResolved)
	}
}

func addPlayerCommand(battleTag string) {
	slog.Info("Adding player to config", "battletag", battleTag)

	// Resolve URL immediately
	result, err := resolveBattleTagToURL(battleTag)
	if err != nil {
		slog.Error("Failed to resolve BattleTag", "battletag", battleTag, "error", err)
		os.Exit(1)
	}

	// Add to config
	err = addPlayerToConfig(battleTag, result.ResolvedURL)
	if err != nil {
		slog.Error("Failed to add player to config", "error", err)
		os.Exit(1)
	}

	slog.Info("✅ Player added to config", "battletag", battleTag, "platform", result.Platform)
}

func removePlayerCommand(battleTag string) {
	err := removePlayerFromConfig(battleTag)
	if err != nil {
		slog.Error("Failed to remove player from config", "battletag", battleTag, "error", err)
		os.Exit(1)
	}

	// Also remove from runtime metrics
	removePlayerMetrics(battleTag)
	slog.Info("✅ Player removed from config and runtime", "battletag", battleTag)
}

func resolveAllCommand(forceResolve bool) {
	err := resolveAllPlayers(forceResolve)
	if err != nil {
		slog.Error("Failed to resolve all players", "error", err)
		os.Exit(1)
	}
	slog.Info("✅ Resolved all player URLs")
}

// Metrics command implementations.
func showMetricsCommand(battleTag string) {
	metrics, exists := getPlayerMetrics(battleTag)
	if !exists {
		slog.Warn("No metrics found", "battletag", battleTag, "suggestion", "run parse-battletag first")

		return
	}

	slog.Info("Player metrics", "battletag", battleTag)
	slog.Info("Display name", "name", metrics.DisplayName)
	slog.Info("Player title", "title", metrics.PlayerTitle)
	slog.Info("Last updated", "timestamp", metrics.LastUpdated.Format("2006-01-02 15:04:05"))
	slog.Info("Endorsement level", "level", metrics.ProfileMetrics.Endorsement.Level)

	// Count total metrics
	totalHeroes := 0
	for _, platforms := range metrics.HeroMetrics {
		for _, heroes := range platforms {
			totalHeroes += len(heroes)
		}
	}
	slog.Info("Total heroes", "count", totalHeroes)
}

func listMetricsCommand() {
	battleTags := listPlayerBattleTags()
	if len(battleTags) == 0 {
		slog.Info("No metrics in runtime store")

		return
	}

	slog.Info("Players with metrics", "count", len(battleTags))
	for _, battleTag := range battleTags {
		if metrics, exists := getPlayerMetrics(battleTag); exists {
			slog.Info("Player metrics",
				"battletag", battleTag,
				"display_name", metrics.DisplayName,
				"last_updated", metrics.LastUpdated.Format("2006-01-02 15:04"))
		}
	}
}

func clearMetricsCommand() {
	clearAllMetrics()
	slog.Info("Cleared all runtime metrics")
}

func statsCommand() {
	stats := getMetricsStats()
	slog.Info("Runtime metrics statistics")
	slog.Info("Total players", "count", stats["total_players"])
	if stats["last_updated"] != nil {
		if lastUpdated, ok := stats["last_updated"].(time.Time); ok {
			slog.Info("Last updated", "timestamp", lastUpdated.Format("2006-01-02 15:04:05"))
		}
	}
}

// Usage printing functions.
func printConfigUsage() {
	slog.Info("Config command usage", "program", os.Args[0])
	slog.Info("Available subcommands:")
	slog.Info("  list-players              List all configured players")
	slog.Info("  add-player <BattleTag>    Add a player to config")
	slog.Info("  remove-player <BattleTag> Remove a player from config")
	slog.Info("  resolve-all [--force]     Resolve URLs for all players")
}

func printMetricsUsage() {
	slog.Info("Metrics command usage", "program", os.Args[0])
	slog.Info("Available subcommands:")
	slog.Info("  show <BattleTag>     Show metrics for a player")
	slog.Info("  list                 List all players with metrics")
	slog.Info("  clear                Clear all runtime metrics")
	slog.Info("  stats                Show runtime metrics statistics")
}

func printMainUsage() {
	slog.Info("Usage", "program", os.Args[0], "format", "<command> [arguments]")
	slog.Info("Available commands:")
	slog.Info("  parse-battletag <BattleTag>  Parse Overwatch profile by BattleTag")
	slog.Info("  config <subcommand>          Manage player configuration")
	slog.Info("  metrics <subcommand>         View runtime metrics")
	slog.Info("  server [port]                Start Prometheus metrics server (default: 9090)")
	slog.Info("  parse-profile                Parse single profile (development mode)")
	slog.Info("  test-api                     Test API endpoints for detailed stats")
	slog.Info("  test-headless                Test headless browser parsing for JS content")
	slog.Info("  test-parser                  Test parser with page.html from browser")
	slog.Info("Help", "message", "Run command without arguments for command-specific help", "program", os.Args[0])
}

// runAPITests runs API endpoint discovery.
func runAPITests() {
	slog.Info("🔍 Starting API endpoint discovery...")

	ctx := context.Background()
	inspector := NewAPIInspector()

	// Get player config
	players := getAllPlayers()
	if len(players) == 0 {
		slog.Error("No players configured")

		return
	}

	// Test with first player
	player := &players[0]
	profileURL := player.ResolvedURL

	slog.Info("Testing API endpoints", "battletag", player.BattleTag, "url", profileURL)

	err := inspector.InspectPotentialAPIEndpoints(ctx, profileURL)
	if err != nil {
		slog.Error("API inspection failed", "error", err)

		return
	}

	slog.Info("✅ API endpoint discovery completed")
}

// runHeadlessTests runs headless browser testing.
func runHeadlessTests() {
	slog.Info("🌐 Starting headless browser testing...")

	ctx := context.Background()
	parser := NewHeadlessParser()

	// Get player config
	players := getAllPlayers()
	if len(players) == 0 {
		slog.Error("No players configured")

		return
	}

	// Test with first player
	player := &players[0]
	profileURL := player.ResolvedURL

	slog.Info("Testing headless parsing", "battletag", player.BattleTag, "url", profileURL)

	// Fetch page with JavaScript execution
	_, err := parser.FetchWithJavaScript(ctx, profileURL)
	if err != nil {
		slog.Error("Headless parsing failed", "error", err)

		return
	}

	// Note: DOM analysis was removed to simplify the function

	slog.Info("✅ Headless browser testing completed")
}

// runParserTests tests the updated parser with real HTML from page.html.
func runParserTests() {
	slog.Info("🧪 Testing parser with real HTML from page.html...")

	htmlContent, err := loadHTMLFile()
	if err != nil {
		return
	}

	profile, err := parseHTMLContent(htmlContent)
	if err != nil {
		return
	}

	analyzeParsingResults(profile)
}

// loadHTMLFile reads the HTML file from disk.
func loadHTMLFile() ([]byte, error) {
	htmlContent, err := os.ReadFile("page.html")
	if err != nil {
		slog.Error("Failed to read page.html", "error", err)
		slog.Info("Make sure page.html exists in the current directory")

		return nil, errors.Wrap(err, "failed to read page.html")
	}

	slog.Info("Loaded HTML file", "size_bytes", len(htmlContent))

	return htmlContent, nil
}

// parseHTMLContent creates a parser and parses the HTML content.
func parseHTMLContent(htmlContent []byte) (*FullPlayerProfile, error) {
	parser := NewParser()
	profile, err := parser.ParseProfile(string(htmlContent), "LexFrei#21715")
	if err != nil {
		slog.Error("Failed to parse profile", "error", err)

		return nil, err
	}

	return profile, nil
}

// analyzeParsingResults analyzes and reports on the parsing results.
func analyzeParsingResults(profile *FullPlayerProfile) {
	totalMetrics := 0
	detailedMetrics := 0

	for platform, platformStats := range profile.Platforms {
		for gameMode, gameModeStats := range platformStats.GameModes {
			slog.Info("Game mode stats",
				"platform", platform,
				"game_mode", gameMode,
				"heroes_count", len(gameModeStats.Heroes))

			for heroID, heroStats := range gameModeStats.Heroes {
				heroMetricsCount := len(heroStats.Metrics)
				totalMetrics += heroMetricsCount

				if heroMetricsCount > 1 {
					detailedMetrics += heroMetricsCount
					logHeroMetrics(heroID, heroStats, heroMetricsCount)
				}
			}
		}
	}

	logFinalResults(totalMetrics, detailedMetrics, len(profile.Platforms))
}

// logHeroMetrics logs detailed information about hero metrics.
func logHeroMetrics(heroID string, heroStats *HeroStats, count int) {
	slog.Info("Hero with detailed metrics",
		"hero_id", heroID,
		"metrics_count", count)

	// Show first few metrics for verification
	metricCount := 0
	for metricKey, value := range heroStats.Metrics {
		if metricCount < maxSampleMetricsCount {
			slog.Info("Sample metric",
				"hero_id", heroID,
				"metric", metricKey,
				"value", value)
		}
		metricCount++
	}
}

// logFinalResults logs the final parsing results.
func logFinalResults(totalMetrics, detailedMetrics, platformCount int) {
	slog.Info("🎯 Parser test results",
		"total_metrics", totalMetrics,
		"detailed_metrics", detailedMetrics,
		"platforms", platformCount)

	if detailedMetrics > 0 {
		slog.Info("✅ SUCCESS: Found detailed hero metrics!")
	} else {
		slog.Warn("❌ No detailed hero metrics found")
	}
}
