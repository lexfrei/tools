package main

import (
	_ "embed"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
)

//go:embed index.html
var site string

//go:embed favicon.png
var favicon string

func main() {
	m := minify.New()
	m.AddFunc("text/html", html.Minify)
	m.AddFunc("text/css", css.Minify)
	site, _ := m.String("text/html", site)

	bd, err := time.Parse("02.01.2006", "04.08.1993")
	if err != nil {
		log.Fatal(err)
	}

	t, err := template.New("webpage").Parse(site)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method is not supported", http.StatusNotFound)
			return
		}

		y, _ := math.Modf(time.Until(bd).Seconds() / -31207680)
		err = t.Execute(w, y)
		if err != nil {
			log.Panicln(err)
		}
	})

	http.HandleFunc("/favicon.png", faviconHandler)

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method is not supported", http.StatusNotFound)
		return
	}

	fmt.Fprint(w, favicon)
}
