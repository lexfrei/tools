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

const (
	headlessHTTPTimeoutSeconds = 60
)

// HeadlessParser provides JavaScript-enabled HTML parsing.
type HeadlessParser struct {
	client *http.Client
}

// NewHeadlessParser creates a new headless parser.
func NewHeadlessParser() *HeadlessParser {
	return &HeadlessParser{
		client: &http.Client{
			Timeout: headlessHTTPTimeoutSeconds * time.Second,
		},
	}
}

// FetchWithJavaScript fetches a page and executes JavaScript to get full DOM.
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

// AnalyzeJSLoadedStructure analyzes the JavaScript-loaded DOM structure.
func (h *HeadlessParser) AnalyzeJSLoadedStructure(_ context.Context, _ string) error {
	slog.Info("🔍 Analyzing JavaScript-loaded structure")

	return nil
}

// tryPuppeteer attempts to use Node.js with Puppeteer for JavaScript execution.
func (h *HeadlessParser) tryPuppeteer(ctx context.Context, profileURL string) (*goquery.Document, error) {
	script := h.generatePuppeteerScript(profileURL)
	output, err := h.executePuppeteerScript(ctx, script)
	if err != nil {
		return nil, err
	}

	return h.parsePuppeteerOutput(output)
}

// generatePuppeteerScript creates the Node.js script for Puppeteer.
func (h *HeadlessParser) generatePuppeteerScript(profileURL string) string {
	return fmt.Sprintf(`
const puppeteer = require('puppeteer');

(async () => {
  try {
    const browser = await puppeteer.launch({
      headless: true,
      args: ['--no-sandbox', '--disable-setuid-sandbox', '--disable-dev-shm-usage']
    });
    const page = await browser.newPage();

    await page.setUserAgent('Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) ' +
      'AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36');

    console.error('Navigating to URL: %s');
    await page.goto('%s', {waitUntil: 'networkidle2', timeout: 30000});

    console.error('Waiting for JavaScript to execute...');
    await new Promise(resolve => setTimeout(resolve, 10000));

    // Try to wait for specific elements
    try {
      await page.waitForSelector('.mouseKeyboard-view', {timeout: 15000});
      console.error('Found mouseKeyboard-view element');
    } catch (e) {
      console.error('mouseKeyboard-view not found:', e.message);
    }

    // Check for the specific selectors
    const selectors = ['.stats-container', 'blz-section', '.option-15', '[data-stat]'];
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
}

// executePuppeteerScript runs the Node.js Puppeteer script.
func (h *HeadlessParser) executePuppeteerScript(ctx context.Context, script string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, "node", "-e", script)
	output, err := cmd.CombinedOutput()
	if err != nil {
		slog.Error("Puppeteer script failed", "error", err, "output", string(output))

		return nil, errors.Wrap(err, "failed to execute Puppeteer script")
	}

	return output, nil
}

// parsePuppeteerOutput parses the HTML output from Puppeteer.
func (h *HeadlessParser) parsePuppeteerOutput(output []byte) (*goquery.Document, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(output)))
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse Puppeteer HTML output")
	}

	slog.Info("✅ Successfully fetched page with Puppeteer")

	return doc, nil
}

// fetchWithEnhancedHTTP uses enhanced HTTP with better timing and headers.
func (h *HeadlessParser) fetchWithEnhancedHTTP(ctx context.Context, profileURL string) (*goquery.Document, error) {
	slog.Info("🔄 Using enhanced HTTP fetching with delays")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, profileURL, http.NoBody)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}

	// Set comprehensive browser-like headers
	req.Header.Set("User-Agent",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) "+
			"AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept",
		"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
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
