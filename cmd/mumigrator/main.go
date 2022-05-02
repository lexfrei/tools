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
	"github.com/SevereCloud/vksdk/v2/object"
	"github.com/pkg/errors"
)

const photosCount = 1000

func main() {
	token := ""
	userID := 7154600
	vkClient := api.NewVK(token)

	albums, err := vkClient.PhotosGetAlbums(api.Params{
		"user_ids": userID,
	})
	if err != nil {
		log.Fatal(err)
	}

	err = albumsProcessing(albums, userID, vkClient)
	if err != nil {
		log.Fatal(err)
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

func albumsProcessing(albums api.PhotosGetAlbumsResponse, userID int, vkClient *api.VK) error {
	albums, err := vkClient.PhotosGetAlbums(api.Params{
		"user_ids": userID,
	})
	if err != nil {
		log.Fatal(err)
	}

	for albumID := range albums.Items {
		esc := regexp.MustCompile(`(?m)\~|\*|\s|\\|/|\||\?|\(|\)|\.`)

		photos, err := vkClient.PhotosGetExtended(api.Params{
			"album_id": albums.Items[albumID].ID,
			"count":    photosCount,
		})
		if err != nil {
			log.Fatal(err)
		}

		for photoID := range photos.Items {
			pwd, err := os.Getwd()
			if err != nil {
				return errors.Wrap(err, "can't get pwd")
			}

			url, err := getPhotoURL(&photos.Items[photoID])
			if err != nil {
				return errors.Wrap(err, "can't get url")
			}

			err = downloadAndSave(
				context.TODO(),
				url,
				filepath.Clean(filepath.Join(pwd, esc.ReplaceAllString(albums.Items[albumID].Title, "_"))),
				strconv.Itoa(photoID)+".jpg",
			)
			if err != nil {
				return errors.Wrap(err, "can't save photo")
			}
		}
	}

	return nil
}

func getPhotoURL(photo *object.PhotosPhotoFull) (string, error) {
	if len(photo.OrigPhoto.URL) > 0 {
		return photo.OrigPhoto.URL, nil
	}

	if len(photo.MaxSize().URL) > 0 {
		return photo.MaxSize().URL, nil
	}

	for sizeID := range photo.Sizes {
		if photo.Sizes[sizeID].BaseImage.Type == "x" {
			return photo.Sizes[sizeID].BaseImage.URL, nil
		}
	}

	return "", errors.New("photo url with fine size not found")
}
