package main

import (
	"testing"

	scryfall "github.com/BlueMonday/go-scryfall"
)

func TestImageURL(t *testing.T) {
	t.Parallel()

	uris := &scryfall.ImageURIs{
		Grid: "https://example.com/grid.webp",
		JPEG: "https://example.com/card.jpg",
	}

	got := imageURL(uris)
	if got != uris.JPEG {
		t.Errorf("imageURL() = %q, want the JPEG url %q", got, uris.JPEG)
	}
}
