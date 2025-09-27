package main

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/cockroachdb/errors"
)

// BrowserLikeParser attempts to handle JavaScript-loaded content
type BrowserLikeParser struct {
	client *http.Client
}

// NewBrowserLikeParser creates a new parser that tries to handle JS content
func NewBrowserLikeParser() *BrowserLikeParser {
	return &BrowserLikeParser{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// FetchWithJSSupport attempts to fetch content including JS-loaded sections
func (p *BrowserLikeParser) FetchWithJSSupport(ctx context.Context, profileURL string) (*goquery.Document, error) {
	slog.Debug("Attempting to fetch profile with JS-like behavior", "url", profileURL)

	// Create request with browser-like headers
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, profileURL, http.NoBody)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}

	// Add browser-like headers to appear more like a real browser
	req.Header.Set("User-Agent",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")

	// Perform request
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to perform HTTP request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse the HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse HTML")
	}

	slog.Debug("Successfully fetched HTML", "url", profileURL)

	return doc, nil
}

// AnalyzeJSLoadedContent looks for signs of JavaScript-loaded content
func (p *BrowserLikeParser) AnalyzeJSLoadedContent(doc *goquery.Document) {
	slog.Info("🔍 Analyzing HTML structure for JS-loaded content...")

	p.analyzeScriptTags(doc)
	p.analyzeDataAttributes(doc)
	p.analyzeSelectors(doc)
	p.analyzeHeroContainers(doc)
}

// analyzeScriptTags checks for script tags that might load additional content.
func (p *BrowserLikeParser) analyzeScriptTags(doc *goquery.Document) {
	scriptCount := 0
	doc.Find("script").Each(func(_ int, s *goquery.Selection) {
		scriptCount++
		src, exists := s.Attr("src")
		if exists && strings.Contains(src, "career") {
			slog.Debug("Found career-related script", "src", src)
		}
	})
	slog.Info("Found script tags", "count", scriptCount)
}

// analyzeDataAttributes looks for data attributes that might be populated by JS.
func (p *BrowserLikeParser) analyzeDataAttributes(doc *goquery.Document) {
	dataAttrs := make(map[string]int)
	doc.Find("*").Each(func(_ int, s *goquery.Selection) {
		for _, attr := range []string{"data-stat", "data-category-id", "data-hero-id"} {
			if val, exists := s.Attr(attr); exists {
				dataAttrs[attr]++
				if attr == "data-stat" {
					slog.Debug("Found data-stat attribute", "value", val)
				}
			}
		}
	})

	for attr, count := range dataAttrs {
		slog.Info("Found data attributes", "attribute", attr, "count", count)
	}
}

// analyzeSelectors checks for specific selectors the user mentioned.
func (p *BrowserLikeParser) analyzeSelectors(doc *goquery.Document) {
	selectors := []string{
		"blz-section.stats",
		"span.stats-container",
		".option-15",
		".mouseKeyboard-view.Profile-view.is-active",
		".quickPlay-view.is-active",
	}

	for _, selector := range selectors {
		elements := doc.Find(selector)
		count := elements.Length()
		slog.Info("Checking selector", "selector", selector, "found", count)

		if count > 0 {
			elements.Each(func(_ int, s *goquery.Selection) {
				classes, _ := s.Attr("class")
				id, _ := s.Attr("id")
				slog.Debug("Found element", "selector", selector, "classes", classes, "id", id)
			})
		}
	}
}

// analyzeHeroContainers looks for containers that might hold hero stats.
func (p *BrowserLikeParser) analyzeHeroContainers(doc *goquery.Document) {
	heroContainers := doc.Find("[data-hero-id]")
	slog.Info("Found hero containers", "count", heroContainers.Length())

	heroContainers.Each(func(_ int, s *goquery.Selection) {
		heroID, _ := s.Attr("data-hero-id")
		classes, _ := s.Attr("class")
		text := strings.TrimSpace(s.Text())
		if len(text) > 100 {
			text = text[:100] + "..."
		}
		slog.Debug("Hero container", "hero", heroID, "classes", classes, "text", text)
	})
}

// InspectFullStructure provides detailed analysis of the page structure
func (p *BrowserLikeParser) InspectFullStructure(doc *goquery.Document) {
	slog.Info("📋 Full structure analysis...")

	// Look for any elements with "stats" in class name
	statsElements := doc.Find("*[class*='stats']")
	slog.Info("Elements with 'stats' in class", "count", statsElements.Length())

	statsElements.Each(func(i int, s *goquery.Selection) {
		classes, _ := s.Attr("class")
		tagName := goquery.NodeName(s)
		slog.Debug("Stats element", "tag", tagName, "classes", classes)
	})

	// Look for view containers
	viewElements := doc.Find("*[class*='view']")
	slog.Info("Elements with 'view' in class", "count", viewElements.Length())

	viewCount := make(map[string]int)
	viewElements.Each(func(i int, s *goquery.Selection) {
		classes, _ := s.Attr("class")
		for _, class := range strings.Split(classes, " ") {
			if strings.Contains(class, "view") {
				viewCount[class]++
			}
		}
	})

	for viewClass, count := range viewCount {
		slog.Debug("View class frequency", "class", viewClass, "count", count)
	}
}
