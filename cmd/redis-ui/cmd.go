package main

import (
	_ "embed"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"sort"
	"time"
)

//go:embed template.html
var templateString string

type PortHosts struct {
	Port  int   `json:"port"`
	Hosts []int `json:"hosts"`
}

func main() {
	http.HandleFunc("/", displayTable)

	//nolint:mnd // 3 seconds is a reasonable timeout
	server := &http.Server{
		Addr:              ":8080",
		ReadHeaderTimeout: 3 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}

func displayTable(writer http.ResponseWriter, _ *http.Request) {
	tmpl, err := template.New("webpage").Parse(templateString)
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(writer, "Internal server error", http.StatusInternalServerError)

		return
	}

	data, err := os.ReadFile("empty_ports_hosts.json")
	if err != nil {
		log.Printf("Error reading JSON file: %v", err)
		http.Error(writer, "Internal server error", http.StatusInternalServerError)

		return
	}

	var portHostsList []PortHosts

	err = json.Unmarshal(data, &portHostsList)
	if err != nil {
		log.Printf("Error unmarshaling JSON: %v", err)
		http.Error(writer, "Internal server error", http.StatusInternalServerError)

		return
	}

	// Sort the portHostsList by the number of hosts
	sort.Slice(portHostsList, func(i, j int) bool {
		return len(portHostsList[i].Hosts) > len(portHostsList[j].Hosts)
	})

	err = tmpl.Execute(writer, portHostsList)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(writer, "Internal server error", http.StatusInternalServerError)

		return
	}
}
