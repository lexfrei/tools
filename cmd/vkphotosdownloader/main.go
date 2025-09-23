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
	"sync"
	"time"

	"github.com/SevereCloud/vksdk/v3/api"
	"github.com/SevereCloud/vksdk/v3/object"
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

	_, err := os.Stat(directory)
	if os.IsNotExist(err) {
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

	var wg sync.WaitGroup

	ctx := context.Background()

	// Get all photos from each album.
	for albumID := range albums {
		wg.Add(1)

		go downloadAlbum(ctx, &wg, &albums[albumID], vkClient, directory)
	}

	wg.Wait()
}

// getPhotoURL returns url of photo.
func getPhotoURL(photo *object.PhotosPhoto) (string, error) {
	if photo == nil {
		return "", errors.New("photo argument is nil")
	}

	if photo.MaxSize().URL != "" {
		return photo.MaxSize().URL, nil
	}

	for sizeID := range photo.Sizes {
		if photo.Sizes[sizeID].Type == "x" {
			return photo.Sizes[sizeID].URL, nil
		}
	}

	return "", errors.New("photo url with fine size not found")
}

// downloadAndSave downloads file from url and saves it to file.
func downloadAndSave(ctx context.Context, url, directoryPath, fileName string) error {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return errors.Wrap(err, "cant create request")
	}

	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "cant download file")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Check if the directory already exists
	_, err = os.Stat(directoryPath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(directoryPath, os.ModePerm)
		if err != nil {
			return errors.Wrap(err, "cant create directory")
		}
	}

	filePath := filepath.Join(directoryPath, fileName)

	// Use io.Copy to stream data
	file, err := os.Create(filePath)
	if err != nil {
		return errors.Wrap(err, "cant create file")
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return errors.Wrap(err, "cant save file")
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

// downloadAlbum downloads all photos from album.
func downloadAlbum(
	ctx context.Context,
	wg *sync.WaitGroup,
	album *object.PhotosPhotoAlbumFull,
	vkClient *api.VK,
	directory string,
) {
	defer wg.Done()

	photos, err := getAllPhotosFronVKAlbum(album, vkClient)
	if err != nil {
		log.Println(err)

		return
	}

	for photoID := range photos {
		url, err := getPhotoURL(&photos[photoID])
		if err != nil {
			log.Printf("can't get url of photo %d: %s", photos[photoID].ID, err)

			continue
		}

		// Download photo with context
		err = downloadAndSave(
			ctx,
			url,
			path.Join(
				directory,
				strconv.Itoa(album.ID),
			),
			strconv.Itoa(photos[photoID].ID)+".jpg",
		)
		if err != nil {
			log.Printf("can't download photo %d: %s", photos[photoID].ID, err)

			continue
		}
	}
}
