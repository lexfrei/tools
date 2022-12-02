package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/object"
	"github.com/cockroachdb/errors"
)

func main() {
	// Get token from environment variable.
	token := os.Getenv("VK_TOKEN")
	if token == "" {
		log.Fatal("VK_TOKEN is not set")
	}

	// Get directory from environment variable.
	directory := os.Getenv("VK_DIRECTORY")
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		log.Fatal("VK_DIRECTORY is not exists")
	}

	// Create VK client.
	vkClient := api.NewVK(token)

	// Get user ID of current user.
	userID, err := getMyVKUserID(vkClient)
	if err != nil {
		log.Fatal(err)
	}

	// Get all albums from current user.
	albums, err := getAllAlbumsFromVKUser(userID, vkClient)
	if err != nil {
		log.Fatal(err)
	}

	// Get all photos from each album.
	for albumID := range albums {
		photos, err := getAllPhotosFronVKAlbum(&albums[albumID], vkClient)
		if err != nil {
			log.Fatal(err)
		}

		for photoID := range photos {
			url, err := getPhotoURL(&photos[photoID])
			if err != nil {
				log.Println(err)

				continue
			}

			// Download photo.
			err = downloadAndSave(
				context.Background(),
				url,
				path.Join(
					directory,
					strconv.Itoa(albums[albumID].ID),
				),
				strconv.Itoa(photos[photoID].ID)+".jpg",
			)
			if err != nil {
				log.Println(err)

				continue
			}
		}
	}
}

// getPhotoURL returns url of photo.
func getPhotoURL(photo *object.PhotosPhoto) (string, error) {
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

// downloadAndSave downloads file from url and saves it to file.
func downloadAndSave(ctx context.Context, url, directoryPath, filePath string) error {
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

	err = os.MkdirAll(directoryPath, os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "cant create directory")
	}

	err = os.WriteFile(filepath.Join(directoryPath, filePath), body, os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "cant write file")
	}

	return nil
}

// getAllAlbumsFromVKUser returns all albums from user.
func getAllAlbumsFromVKUser(userID int, vkClient *api.VK) ([]object.PhotosPhotoAlbumFull, error) {
	albums, err := vkClient.PhotosGetAlbums(api.Params{
		"user_ids": userID,
	})
	if err != nil {
		return nil, errors.Wrap(err, "can't get albums")
	}

	return albums.Items, nil
}

// getAllPhotosFronVKAlbum returns all photos from album.
func getAllPhotosFronVKAlbum(album *object.PhotosPhotoAlbumFull, vkClient *api.VK) ([]object.PhotosPhoto, error) {
	photosRemaining := album.Size

	var offset int

	var photosResult []object.PhotosPhoto

	for photosRemaining > 0 {
		photos, err := vkClient.PhotosGet(api.Params{
			"album_id":    album.ID,
			"offset":      offset,
			"count":       100,
			"photo_sizes": "1",
		})
		if err != nil {
			return nil, errors.Wrap(err, "can't get photos")
		}

		photosRemaining -= len(photos.Items)
		offset += len(photos.Items)
		photosResult = append(photosResult, photos.Items...)
	}

	return photosResult, nil
}

// getMyVKUserID returns user ID of current user.
func getMyVKUserID(vkClient *api.VK) (int, error) {
	me, err := vkClient.UsersGet(api.Params{})
	if err != nil {
		return 0, errors.Wrap(err, "can't get me")
	}

	return me[0].ID, nil
}
