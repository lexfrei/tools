package main

import "github.com/cockroachdb/errors"

// Static errors for err113 linter compliance.
var (
	// ErrPlayerNotFound occurs when a player is not found in config.
	ErrPlayerNotFound = errors.New("player not found in config")
	ErrNoConfigLoaded = errors.New("no config loaded")

	// ErrProfileNotFound occurs when a profile is not found (404).
	ErrProfileNotFound  = errors.New("profile not found (404)")
	ErrTooManyRedirects = errors.New("too many redirects")

	// ErrMetricsNil occurs when metrics cannot be nil.
	ErrMetricsNil      = errors.New("metrics cannot be nil")
	ErrEmptyBattleTag  = errors.New("battle tag cannot be empty")
	ErrZeroLastUpdated = errors.New("last updated time cannot be zero")

	// ErrBattleTagNotResolved occurs when failed to resolve BattleTag on any platform.
	ErrBattleTagNotResolved = errors.New("failed to resolve BattleTag on any platform")
	ErrUnexpectedStatusCode = errors.New("unexpected status code")
	ErrHTTPError            = errors.New("HTTP error")

	// ErrUnknownHero occurs when unknown hero detected - not in registry.
	ErrUnknownHero = errors.New("unknown hero detected - not in registry")

	// ErrNoPlatformData occurs when no data found for platform.
	ErrNoPlatformData = errors.New("no data found for platform")
	// ErrNoGameModeData occurs when no data found for game mode.
	ErrNoGameModeData = errors.New("no data found for game mode")
	// ErrNoHeroID occurs when no hero ID found in hero element.
	ErrNoHeroID = errors.New("no hero ID found in hero element")
	// ErrNoHeroName occurs when no hero name found in hero element.
	ErrNoHeroName = errors.New("no hero name found in hero element")
)
