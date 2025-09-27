package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/cockroachdb/errors"
)

// HeadlessParser provides JavaScript-enabled HTML parsing
type HeadlessParser struct {
	client *http.Client
}

// NewHeadlessParser creates a new headless parser
func NewHeadlessParser() *HeadlessParser {
	return &HeadlessParser{
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// FetchWithJavaScript fetches a page and executes JavaScript to get full DOM
func (h *HeadlessParser) FetchWithJavaScript(ctx context.Context, profileURL string) (*goquery.Document, error) {
	slog.Info("🌐 Fetching page with JavaScript execution", "url", profileURL)

	// Try using Node.js with Puppeteer if available
	doc, err := h.tryPuppeteer(ctx, profileURL)
	if err == nil {
		return doc, nil
	}

	slog.Debug("Puppeteer not available, trying alternative approach", "error", err)

	// Fallback to enhanced HTTP fetching with better headers and delays
	return h.fetchWithEnhancedHTTP(ctx, profileURL)
}

// tryPuppeteer attempts to use Node.js with Puppeteer for JavaScript execution
func (h *HeadlessParser) tryPuppeteer(ctx context.Context, profileURL string) (*goquery.Document, error) {
	// Create a simple Node.js script for Puppeteer
	script := fmt.Sprintf(`
const puppeteer = require('puppeteer');

(async () => {
  try {
    const browser = await puppeteer.launch({
      headless: true,
      args: ['--no-sandbox', '--disable-setuid-sandbox', '--disable-dev-shm-usage']
    });
    const page = await browser.newPage();

    await page.setUserAgent('Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36');

    console.error('Navigating to URL: %s');
    await page.goto('%s', {waitUntil: 'networkidle2', timeout: 30000});

    console.error('Waiting for JavaScript to execute...');
    // Wait longer for stats to load - using setTimeout instead of waitForTimeout
    await new Promise(resolve => setTimeout(resolve, 10000));

    // Try to wait for specific elements
    try {
      await page.waitForSelector('.mouseKeyboard-view', {timeout: 15000});
      console.error('Found mouseKeyboard-view element');
    } catch (e) {
      console.error('mouseKeyboard-view not found:', e.message);
    }

    // Check for the specific selectors
    const selectors = [
      '.stats-container',
      'blz-section',
      '.option-15',
      '[data-stat]'
    ];

    for (const selector of selectors) {
      try {
        const elements = await page.$$(selector);
        console.error('Selector ' + selector + ' found ' + elements.length + ' elements');
      } catch (e) {
        console.error('Error checking selector ' + selector + ':', e.message);
      }
    }

    const html = await page.content();
    console.log(html);

    await browser.close();
  } catch (error) {
    console.error('Puppeteer error:', error.message);
    process.exit(1);
  }
})();
`, profileURL, profileURL)

	// Try to execute the Node.js script
	cmd := exec.CommandContext(ctx, "node", "-e", script)

	// Capture both stdout and stderr to see what's happening
	output, err := cmd.CombinedOutput()
	if err != nil {
		slog.Error("Puppeteer script failed", "error", err, "output", string(output))
		return nil, errors.Wrap(err, "failed to execute Puppeteer script")
	}

	// Parse the HTML output
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(output)))
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse Puppeteer HTML output")
	}

	slog.Info("✅ Successfully fetched page with Puppeteer")
	return doc, nil
}

// fetchWithEnhancedHTTP uses enhanced HTTP with better timing and headers
func (h *HeadlessParser) fetchWithEnhancedHTTP(ctx context.Context, profileURL string) (*goquery.Document, error) {
	slog.Info("🔄 Using enhanced HTTP fetching with delays")

	req, err := http.NewRequestWithContext(ctx, "GET", profileURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}

	// Set comprehensive browser-like headers
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Cache-Control", "max-age=0")

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "HTTP request failed")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse HTML")
	}

	slog.Info("📄 Enhanced HTTP fetch completed")
	return doc, nil
}

// AnalyzeJSLoadedStructure specifically looks for the user's mentioned selectors
func (h *HeadlessParser) AnalyzeJSLoadedStructure(doc *goquery.Document) {
	slog.Info("🔍 Analyzing DOM for JS-loaded detailed hero stats...")

	// Check for the specific selector structure the user mentioned
	selectors := map[string]string{
		"mouseKeyboard-view active":     ".mouseKeyboard-view.Profile-view.is-active",
		"stats quickPlay section":       "blz-section.stats.quickPlay-view.is-active",
		"stats container option-15":     "span.stats-container.option-15.is-active",
		"full selector path":            "body > main > div.mouseKeyboard-view.Profile-view.is-active > blz-section.stats.quickPlay-view.is-active > span.stats-container.option-15.is-active",
		"any stats container":           ".stats-container",
		"any blz-section":               "blz-section",
		"data-stat attributes":          "[data-stat]",
		"cassidy hero container":        "[data-hero-id='cassidy']",
	}

	for desc, selector := range selectors {
		elements := doc.Find(selector)
		count := elements.Length()

		if count > 0 {
			slog.Info("🎯 Found target elements!", "description", desc, "selector", selector, "count", count)

			// Extract sample data from found elements
			elements.Each(func(i int, s *goquery.Selection) {
				// Get all attributes
				attrs := make(map[string]string)
				if s.Get(0) != nil {
					for _, attr := range s.Get(0).Attr {
						attrs[attr.Key] = attr.Val
					}
				}

				text := strings.TrimSpace(s.Text())
				if len(text) > 100 {
					text = text[:100] + "..."
				}

				slog.Info("Found element details",
					"index", i,
					"text", text,
					"attributes", attrs)
			})
		} else {
			slog.Debug("Selector not found", "description", desc, "selector", selector)
		}
	}

	// Look for detailed hero statistics
	h.extractDetailedHeroStats(doc)
}

// extractDetailedHeroStats attempts to extract the detailed stats that should be available
func (h *HeadlessParser) extractDetailedHeroStats(doc *goquery.Document) {
	slog.Info("🎮 Looking for detailed hero statistics...")

	// Look for Cassidy/McCree specific stats as mentioned by user
	cassidyStats := map[string]string{
		"resurrects":              "[data-stat='resurrects']",
		"damage_amplified":        "[data-stat='damage_amplified']",
		"rocket_hammer_kills":     "[data-stat='rocket_hammer_kills']",
		"pulse_bomb_kills":        "[data-stat='pulse_bomb_kills']",
		"earthshatter_kills":      "[data-stat='earthshatter_kills']",
		"deadeye_kills":           "[data-stat='deadeye_kills']",
		"flashbang_enemies":       "[data-stat='flashbang_enemies']",
		"combat_roll_kills":       "[data-stat='combat_roll_kills']",
	}

	foundStats := 0
	for statName, selector := range cassidyStats {
		elements := doc.Find(selector)
		if elements.Length() > 0 {
			foundStats++
			elements.Each(func(i int, s *goquery.Selection) {
				value := strings.TrimSpace(s.Text())
				slog.Info("🎯 Found detailed stat!", "stat", statName, "value", value, "selector", selector)
			})
		}
	}

	if foundStats > 0 {
		slog.Info("✅ Successfully found detailed hero statistics!", "count", foundStats)
	} else {
		slog.Warn("❌ No detailed hero statistics found in current DOM")
	}

	// Also check for option containers that might hold the stats
	for i := 1; i <= 20; i++ {
		optionSelector := fmt.Sprintf(".option-%d", i)
		elements := doc.Find(optionSelector)
		if elements.Length() > 0 {
			slog.Info("Found option container", "selector", optionSelector, "count", elements.Length())
		}
	}
}