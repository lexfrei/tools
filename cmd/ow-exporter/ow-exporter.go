package main

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
	"unicode"

	owp "github.com/lexfrei/tools/internal/pkg/owparser"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var (
	playerRank = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rank",
			Help: "player rank representation",
		},
		[]string{
			"user",
			"platform",
			"role",
		},
	)
	playerEndorsment = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "endorsment",
			Help: "player endorsments representation",
		},
		[]string{
			"user",
			"platform",
			"type",
		},
	)
	stats = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "stat",
			Help: "player stats representation",
		},
		[]string{
			"user",
			"platform",
			"type",
			"stat",
			"hero",
		},
	)
)

func main() {
	registry := prometheus.NewRegistry()
	registry.MustRegister(playerRank)
	registry.MustRegister(playerEndorsment)
	registry.MustRegister(stats)

	playerURL, err := url.Parse(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		getStats(playerURL)

		for range time.Tick(1 * time.Minute) {
			getStats(playerURL)
		}
	}()

	handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	http.Handle("/metrics", handler)
	log.Fatal(http.ListenAndServe(":9420", nil))
}

func getStats(u *url.URL) {
	player := owp.NewPlayerByLink(u)

	err := player.Gather(context.TODO())
	if err != nil {
		log.Println(err)

		return
	}

	playerRank.WithLabelValues(player.Name, player.Platform, "tank").Set(float64(player.Rank.Tank))
	playerRank.WithLabelValues(player.Name, player.Platform, "heal").Set(float64(player.Rank.Heal))
	playerRank.WithLabelValues(player.Name, player.Platform, "dd").Set(float64(player.Rank.DD))
	playerEndorsment.WithLabelValues(player.Name, player.Platform, "level").Set(float64(player.Endorsment.Level))
	playerEndorsment.WithLabelValues(player.Name, player.Platform, "sportsmanship").Set(player.Endorsment.Sportsmanship)
	playerEndorsment.WithLabelValues(player.Name, player.Platform, "shotcaller").Set(player.Endorsment.Shotcaller)
	playerEndorsment.WithLabelValues(player.Name, player.Platform, "teammate").Set(player.Endorsment.Teammate)
}

//nolint:deadcode // for the future use
func normalize(str string) string {
	result, _, _ := transform.String(
		transform.Chain(
			norm.NFD,
			runes.Remove(
				runes.In(
					unicode.Mn,
				),
			),
			norm.NFC,
		),
		str,
	)

	return strings.ReplaceAll(result, " ", "_")
}
