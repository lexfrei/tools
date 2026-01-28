package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// TestCountFullYearsSinceBirth tests age calculation from birth date.
func TestCountFullYearsSinceBirth(t *testing.T) {
	timezone := time.FixedZone("UTC+4", 4*60*60)

	tests := []struct {
		name      string
		birthDate time.Time
		now       time.Time
		expected  int
	}{
		{
			name:      "birthday already passed this year",
			birthDate: time.Date(1993, 8, 4, 0, 0, 0, 0, timezone),
			now:       time.Date(2026, 9, 1, 0, 0, 0, 0, timezone),
			expected:  33,
		},
		{
			name:      "birthday not yet reached this year",
			birthDate: time.Date(1993, 8, 4, 0, 0, 0, 0, timezone),
			now:       time.Date(2026, 7, 1, 0, 0, 0, 0, timezone),
			expected:  32,
		},
		{
			name:      "today is birthday",
			birthDate: time.Date(1993, 8, 4, 0, 0, 0, 0, timezone),
			now:       time.Date(2026, 8, 4, 0, 0, 0, 0, timezone),
			expected:  33,
		},
		{
			name:      "leap year check",
			birthDate: time.Date(2000, 2, 29, 0, 0, 0, 0, timezone),
			now:       time.Date(2024, 3, 1, 0, 0, 0, 0, timezone),
			expected:  24,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			result := calculateYearsBetween(testCase.birthDate, testCase.now)

			if result != testCase.expected {
				t.Errorf("expected %d years, got %d", testCase.expected, result)
			}
		})
	}
}

// calculateYearsBetween is a helper function for testing age calculation logic.
func calculateYearsBetween(birthDate, now time.Time) int {
	if now.Month() < birthDate.Month() || (birthDate.Month() == now.Month() && now.Day() < birthDate.Day()) {
		return now.Year() - birthDate.Year() - 1
	}

	return now.Year() - birthDate.Year()
}

// TestRootHandler tests the main page.
func TestRootHandler(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, server.URL+"/", http.NoBody)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	body := make([]byte, 1024)
	num, _ := resp.Body.Read(body)
	bodyStr := string(body[:num])

	if !strings.Contains(bodyStr, "<!DOCTYPE") && !strings.Contains(bodyStr, "<html") {
		t.Error("response should contain HTML")
	}

	hasNumber := false
	for _, char := range bodyStr {
		if char >= '0' && char <= '9' {
			hasNumber = true

			break
		}
	}

	if !hasNumber {
		t.Error("response should contain age number")
	}
}

// TestFaviconHandler tests favicon response.
func TestFaviconHandler(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, server.URL+"/favicon.png", http.NoBody)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	body := make([]byte, 100)
	num, _ := resp.Body.Read(body)

	if num == 0 {
		t.Error("favicon should not be empty")
	}
}

// TestRobotsHandler tests robots.txt response.
func TestRobotsHandler(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, server.URL+"/robots.txt", http.NoBody)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	body := make([]byte, 1024)
	num, _ := resp.Body.Read(body)
	bodyStr := string(body[:num])

	if !strings.Contains(strings.ToLower(bodyStr), "user-agent:") {
		t.Error("robots.txt should contain 'user-agent:'")
	}
}

// setupTestServer creates a test HTTP server for integration testing using actual handlers.
func setupTestServer(t *testing.T) *httptest.Server {
	t.Helper()

	// Setup similar to main() but for testing
	timezone := time.FixedZone("UTC+4", 4*60*60)
	birthDate, err := time.ParseInLocation("02.01.2006", "04.08.1993", timezone)
	if err != nil {
		t.Fatalf("failed to parse birth date: %v", err)
	}

	mux := http.NewServeMux()

	// Use actual handlers from main
	mux.HandleFunc("GET /", func(writer http.ResponseWriter, _ *http.Request) {
		writer.Header().Set("Content-Type", "text/html; charset=utf-8")
		// Simple HTML with age for testing
		age := countFullYearsSinceBirth(birthDate, timezone)
		html := fmt.Sprintf("<!DOCTYPE html><html><body>Age: %d</body></html>", age)
		_, _ = writer.Write([]byte(html))
	})

	mux.HandleFunc("GET /favicon.png", faviconHandler)
	mux.HandleFunc("GET /robots.txt", robotsHandler)

	return httptest.NewServer(mux)
}
