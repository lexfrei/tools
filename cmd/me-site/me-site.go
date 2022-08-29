package main

import (
	_ "embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
)

// http port.
const port = "8080"

// utc3seconds is the timezone for UTC+3.
const utc3seconds = 3 * 60 * 60

// utc3 is the timezone of the city of moscow.
var utc3 = time.FixedZone("UTC+3", utc3seconds)

// site is the HTML template for the website.
//
//go:embed index.html
var site string

// favicon is the favicon.png.
//
//go:embed favicon.png
var favicon string

func main() {
	// Create a minifier.
	minifier := minify.New()

	// Minify the HTML.
	minifier.AddFunc("text/html", html.Minify)
	// Minify the CSS.
	minifier.AddFunc("text/css", css.Minify)
	// Minify the template.
	site, _ := minifier.String("text/html", site)

	// set birth date
	birthDate, err := time.ParseInLocation("02.01.2006", "04.08.1993", utc3)
	if err != nil {
		log.Fatal(err)
	}

	// Render the template
	siteTemplate, err := template.New("webpage").Parse(site)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	// Serve the website
	mux.HandleFunc("/", func(responseWriter http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(responseWriter, "Method is not supported", http.StatusNotFound)

			return
		}

		err = siteTemplate.Execute(responseWriter, countFullYearsSinceBirth(birthDate, utc3))
		if err != nil {
			log.Panicln(err)
		}
	})

	// Serve the favicon
	mux.HandleFunc("/favicon.png", faviconHandler)

	log.Println("Listening on port 8080")

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  10 * time.Minute,
		WriteTimeout: 10 * time.Minute,
	}
	log.Fatal(srv.ListenAndServe())
}

// faviconHandler returns the favicon.png.
func faviconHandler(responseWriter http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(responseWriter, "Method is not supported", http.StatusNotFound)

		return
	}

	fmt.Fprint(responseWriter, favicon)
}

// countFullYearsSinceBirth returns the number of full years since the birth date.
func countFullYearsSinceBirth(birthDate time.Time, tz *time.Location) int {
	now := time.Now().In(tz)
	if now.Month() < birthDate.Month() || (birthDate.Month() == now.Month() && now.Day() < birthDate.Day()) {
		return now.Year() - birthDate.Year() - 1
	}

	return now.Year() - birthDate.Year()
}
