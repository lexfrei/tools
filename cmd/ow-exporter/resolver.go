package main

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
)

// PlatformType represents the gaming platform.
type PlatformType string

const (
	PlatformTypePC  PlatformType = "pc"
	PlatformTypePSN PlatformType = "psn"
	PlatformTypeXBL PlatformType = "xbl"
)

// ResolveResult contains the result of BattleTag resolution.
type ResolveResult struct {
	BattleTag   string
	ResolvedURL string
	Platform    PlatformType
	Success     bool
	Error       error
}

// resolveBattleTagToURL resolves a BattleTag to its actual profile URL.
func resolveBattleTagToURL(battleTag string) (*ResolveResult, error) {
	slog.Info("Resolving BattleTag to URL", "battletag", battleTag)

	// Convert BattleTag format: "LexFrei#21715" -> "LexFrei-21715"
	urlTag := strings.Replace(battleTag, "#", "-", 1)

	// Try to resolve the BattleTag
	simpleURL := "https://overwatch.blizzard.com/en-us/career/" + urlTag
	slog.Debug("Trying URL", "url", simpleURL)

	resolvedURL, err := followRedirects(simpleURL)
	if err != nil {
		slog.Debug("Resolution failed", "error", err)
	} else if resolvedURL != "" && strings.Contains(resolvedURL, "/career/") && resolvedURL != simpleURL {
		// Successfully resolved to a profile-specific URL
		slog.Info("Successfully resolved BattleTag",
			"battletag", battleTag,
			"resolved_url", resolvedURL)

		return &ResolveResult{
			BattleTag:   battleTag,
			ResolvedURL: resolvedURL,
			Platform:    PlatformTypePC, // Default to PC
			Success:     true,
		}, nil
	}

	err = errors.Wrapf(ErrBattleTagNotResolved, "%s", battleTag)

	return &ResolveResult{
		BattleTag: battleTag,
		Success:   false,
		Error:     err,
	}, err
}

// followRedirects follows HTTP redirects and returns the final URL.
func followRedirects(url string) (string, error) {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(_ *http.Request, via []*http.Request) error {
			// Allow up to 10 redirects
			if len(via) >= 10 {
				return ErrTooManyRedirects
			}

			return nil
		},
	}

	// Create request with proper headers
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, http.NoBody)
	if err != nil {
		return "", errors.Wrap(err, "failed to create request")
	}

	// Set headers to mimic a real browser
	req.Header.Set("User-Agent",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 "+
			"(KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "failed to perform request")
	}
	defer resp.Body.Close()

	// Check for successful response
	if resp.StatusCode == http.StatusOK {
		finalURL := resp.Request.URL.String()
		slog.Debug("Redirect successful", "original_url", url, "final_url", finalURL)

		return finalURL, nil
	}

	// Check for profile not found
	if resp.StatusCode == http.StatusNotFound {
		slog.Debug("Profile not found", "url", url, "status", resp.StatusCode)

		return "", ErrProfileNotFound
	}

	return "", errors.Wrapf(ErrUnexpectedStatusCode, "%d", resp.StatusCode)
}

// getOrResolveURL gets the resolved URL for a BattleTag, resolving it if necessary.
func getOrResolveURL(battleTag string) (string, error) {
	// First, try to find in config
	if player := findPlayerByBattleTag(battleTag); player != nil && player.ResolvedURL != "" {
		slog.Debug("Using cached URL from config", "battletag", battleTag, "url", player.ResolvedURL)

		return player.ResolvedURL, nil
	}

	// Not in config or URL is empty, resolve it
	slog.Info("Resolving BattleTag (not in cache)", "battletag", battleTag)
	result, err := resolveBattleTagToURL(battleTag)
	if err != nil {
		return "", errors.Wrapf(err, "failed to resolve BattleTag %s", battleTag)
	}

	// Save to config for future use
	err = addPlayerToConfig(battleTag, result.ResolvedURL)
	if err != nil {
		slog.Warn("Failed to save resolved URL to config", "error", err)
		// Don't fail the whole operation just because we couldn't save to config
	} else {
		slog.Info("Saved resolved URL to config", "battletag", battleTag)
	}

	return result.ResolvedURL, nil
}

// resolveAllPlayers resolves URLs for all players in config.
func resolveAllPlayers(forceResolve bool) error {
	players := getAllPlayers()
	if len(players) == 0 {
		slog.Info("No players found in config")

		return nil
	}

	slog.Info("Resolving URLs for all players", "count", len(players), "force", forceResolve)

	for _, player := range players {
		// Skip if already resolved (unless forcing)
		if !forceResolve && player.ResolvedURL != "" {
			slog.Debug("Skipping already resolved player", "battletag", player.BattleTag)

			continue
		}

		slog.Info("Resolving player", "battletag", player.BattleTag)
		result, err := resolveBattleTagToURL(player.BattleTag)
		if err != nil {
			slog.Error("Failed to resolve player", "battletag", player.BattleTag, "error", err)

			continue
		}

		// Update config with resolved URL
		err = updatePlayerURL(player.BattleTag, result.ResolvedURL)
		if err != nil {
			slog.Error("Failed to update player URL in config", "battletag", player.BattleTag, "error", err)

			return errors.Wrap(err, "failed to update player URL in config")
		}
		slog.Info("Updated player URL in config", "battletag", player.BattleTag)
	}

	return nil
}
