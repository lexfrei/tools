package main

import (
	_ "embed"
	"fmt"
	"html/template"
	"log/slog"
	"os"
	"time"

	echo "github.com/labstack/echo/v4"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
)

// http port.
const port = "8080"

// utcPlus4 is the timezone for UTC+4.
const utcPlus4 = 4 * 60 * 60

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

	server := echo.New()
	server.HideBanner = true
	server.HidePort = true

	server.GET("/", func(context echo.Context) error {
		err = siteTemplate.Execute(context.Response().Writer, countFullYearsSinceBirth(birthDate, timeZoneUTCPlus4))
		if err != nil {
			slog.Error("Template execution failed", "error", err)
		}

		return nil
	})

	// Serve the favicon
	server.GET("/favicon.png", faviconHandler)

	// Serve the robots.txt
	server.GET("/robots.txt", robotsHandler)

	slog.Info("Starting the server", "port", port)

	slog.Error("Server failed", "error", server.Start(":"+port))
}

// faviconHandler returns the favicon.png.
func faviconHandler(context echo.Context) error {
	fmt.Fprint(context.Response().Writer, favicon)

	return nil
}

// robotsHandler returns the robots.txt.
func robotsHandler(context echo.Context) error {
	fmt.Fprint(context.Response().Writer, robots)

	return nil
}

// countFullYearsSinceBirth returns the number of full years since the birth date.
func countFullYearsSinceBirth(birthDate time.Time, tz *time.Location) int {
	now := time.Now().In(tz)

	if now.Month() < birthDate.Month() || (birthDate.Month() == now.Month() && now.Day() < birthDate.Day()) {
		return now.Year() - birthDate.Year() - 1
	}

	return now.Year() - birthDate.Year()
}
