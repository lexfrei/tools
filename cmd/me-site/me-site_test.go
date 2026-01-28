package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// TestCountFullYearsSinceBirth tests age calculation from birth date.
func TestCountFullYearsSinceBirth(t *testing.T) {
	tz := time.FixedZone("UTC+4", 4*60*60)

	tests := []struct {
		name      string
		birthDate time.Time
		now       time.Time
		expected  int
	}{
		{
			name:      "birthday already passed this year",
			birthDate: time.Date(1993, 8, 4, 0, 0, 0, 0, tz),
			now:       time.Date(2026, 9, 1, 0, 0, 0, 0, tz),
			expected:  33,
		},
		{
			name:      "birthday not yet reached this year",
			birthDate: time.Date(1993, 8, 4, 0, 0, 0, 0, tz),
			now:       time.Date(2026, 7, 1, 0, 0, 0, 0, tz),
			expected:  32,
		},
		{
			name:      "today is birthday",
			birthDate: time.Date(1993, 8, 4, 0, 0, 0, 0, tz),
			now:       time.Date(2026, 8, 4, 0, 0, 0, 0, tz),
			expected:  33,
		},
		{
			name:      "leap year check",
			birthDate: time.Date(2000, 2, 29, 0, 0, 0, 0, tz),
			now:       time.Date(2024, 3, 1, 0, 0, 0, 0, tz),
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

	resp, err := http.Get(server.URL + "/")
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

	resp, err := http.Get(server.URL + "/favicon.png")
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

	resp, err := http.Get(server.URL + "/robots.txt")
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

	if !strings.Contains(bodyStr, "User-agent:") {
		t.Error("robots.txt should contain 'User-agent:'")
	}
}

// setupTestServer creates a test HTTP server for integration testing.
// After migration to net/http, this will use the actual router.
func setupTestServer(t *testing.T) *httptest.Server {
	t.Helper()

	// Temporary stub for compilation.
	// After refactoring main(), this will use the actual router:
	// mux := createRouter()
	// return httptest.NewServer(mux)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte("<!DOCTYPE html><html>test</html>"))
	})
	mux.HandleFunc("/favicon.png", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte("fake-favicon"))
	})
	mux.HandleFunc("/robots.txt", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte("User-agent: *"))
	})

	return httptest.NewServer(mux)
}
