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

// utc3seconds is the timezone for UTC+3.
const utc3seconds = 3 * 60 * 60

// utc3 is the timezone of the city of moscow.
var utc3 = time.FixedZone("UTC+3", utc3seconds)

// site is the HTML template for the website.
//go:embed index.html
var site string

// favicon is the favicon.png.
//go:embed favicon.png
var favicon string

func main() {
	m := minify.New()
	m.AddFunc("text/html", html.Minify)
	m.AddFunc("text/css", css.Minify)
	site, _ := m.String("text/html", site)

	birthDate, err := time.ParseInLocation("02.01.2006", "04.08.1993", utc3)
	if err != nil {
		log.Fatal(err)
	}

	siteTemplate, err := template.New("webpage").Parse(site)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(responseWriter http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(responseWriter, "Method is not supported", http.StatusNotFound)

			return
		}

		err = siteTemplate.Execute(responseWriter, countFullYearsSinceBirth(birthDate, utc3))
		if err != nil {
			log.Panicln(err)
		}
	})

	http.HandleFunc("/favicon.png", faviconHandler)

	log.Printf("Starting server at port 8080\n")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
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
