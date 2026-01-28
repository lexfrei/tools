package main

import (
	_ "embed"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
)

// http port.
const port = "8080"

// utcPlus4 is the timezone for UTC+4.
const utcPlus4 = 4 * 60 * 60

// timeouts for HTTP server.
const (
	readHeaderTimeout = 3
	readTimeout       = 10
	writeTimeout      = 10
)

// timeZoneUTCPlus4 is the timezone of the city of Tbilisi.
var timeZoneUTCPlus4 = time.FixedZone("UTC+4", utcPlus4)

// site is the HTML template for the website.
//
//go:embed index.html
var site string

// favicon is the favicon.png.
//
//go:embed favicon.png
var favicon string

// robots is the robots.txt.
//
//go:embed robots.txt
var robots string

// logLevel is the log level.
var logLevel = slog.LevelInfo

func main() {
	programLevel := new(slog.LevelVar) // Info by default
	programLevel.Set(logLevel)
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: programLevel})))

	// Create a minifier.
	minifier := minify.New()

	// Minify the HTML.
	minifier.AddFunc("text/html", html.Minify)
	// Minify the CSS.
	minifier.AddFunc("text/css", css.Minify)
	// Minify the template.
	site, _ := minifier.String("text/html", site)

	// set birth date
	birthDate, err := time.ParseInLocation("02.01.2006", "04.08.1993", timeZoneUTCPlus4)
	if err != nil {
		slog.Error("Failed to parse birth date", "error", err)
	}

	// Render the template
	siteTemplate, err := template.New("webpage").Parse(site)
	if err != nil {
		slog.Error("Failed to parse the template", "error", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(writer http.ResponseWriter, _ *http.Request) {
		writer.Header().Set("Content-Type", "text/html; charset=utf-8")

		err = siteTemplate.Execute(writer, countFullYearsSinceBirth(birthDate, timeZoneUTCPlus4))
		if err != nil {
			slog.Error("Template execution failed", "error", err)
			http.Error(writer, "Internal Server Error", http.StatusInternalServerError)

			return
		}
	})

	// Serve the favicon
	mux.HandleFunc("GET /favicon.png", faviconHandler)

	// Serve the robots.txt
	mux.HandleFunc("GET /robots.txt", robotsHandler)

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: readHeaderTimeout * time.Second,
		ReadTimeout:       readTimeout * time.Second,
		WriteTimeout:      writeTimeout * time.Second,
	}

	slog.Info("Starting the server", "port", port)

	slog.Error("Server failed", "error", server.ListenAndServe())
}

// faviconHandler returns the favicon.png.
func faviconHandler(writer http.ResponseWriter, _ *http.Request) {
	writer.Header().Set("Content-Type", "image/png")
	fmt.Fprint(writer, favicon)
}

// robotsHandler returns the robots.txt.
func robotsHandler(writer http.ResponseWriter, _ *http.Request) {
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(writer, robots)
}

// countFullYearsSinceBirth returns the number of full years since the birth date.
func countFullYearsSinceBirth(birthDate time.Time, tz *time.Location) int {
	now := time.Now().In(tz)

	if now.Month() < birthDate.Month() || (birthDate.Month() == now.Month() && now.Day() < birthDate.Day()) {
		return now.Year() - birthDate.Year() - 1
	}

	return now.Year() - birthDate.Year()
}
