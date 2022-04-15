package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/pkg/errors"
)

func main() {
	esc := regexp.MustCompile(`(?m)\~|\*|\s|\\|/|\||\?|\(|\)|\.`)
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	token := ""
	userID := 7154600
	vkClient := api.NewVK(token)

	albums, err := vkClient.PhotosGetAlbums(api.Params{
		"user_ids": userID,
	})
	if err != nil {
		log.Fatal(err)
	}

	for i := range albums.Items {
		photos, err := vkClient.PhotosGetExtended(api.Params{
			"album_id": albums.Items[i].ID,
			"count":    1000,
		})
		if err != nil {
			log.Fatal(err)
		}

		for ii := range photos.Items {
			var url string
			if len(photos.Items[ii].OrigPhoto.URL) > 0 {
				url = photos.Items[ii].OrigPhoto.URL
			} else {
				if len(photos.Items[ii].MaxSize().URL) > 0 {
					url = photos.Items[ii].MaxSize().URL
				} else {
					for sizeID := range photos.Items[ii].Sizes {
						if photos.Items[ii].Sizes[sizeID].BaseImage.Type == "x" {
							url = photos.Items[ii].Sizes[sizeID].BaseImage.URL
						}
					}
				}
			}

			err := downloadAndSave(
				context.TODO(),
				url,
				filepath.Clean(filepath.Join(pwd, esc.ReplaceAllString(albums.Items[i].Title, "_"))),
				strconv.Itoa(ii)+".jpg",
			)
			if err != nil {
				log.Fatalln(err)
			}
		}
	}
}

func downloadAndSave(ctx context.Context, url, dpath, fpath string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return errors.Wrap(err, "cant create request")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "cant download file")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "cant read body")
	}

	err = os.MkdirAll(dpath, os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "cant create directory")
	}

	err = os.WriteFile(filepath.Join(dpath, fpath), body, os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "cant write file")
	}

	return nil
}
