package main

import (
	"context"
	"log"
	"net/url"
	"os"

	owp "github.com/lexfrei/tools/internal/pkg/owparser"
)

func main() {
	playerURL, err := url.Parse(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	player := owp.NewPlayerByLink(playerURL)

	err = player.Gather(context.TODO())
	if err != nil {
		log.Println(err)
	}

	// pretty.Println(player)
}
