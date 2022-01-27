package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"sync"

	scryfall "github.com/BlueMonday/go-scryfall"
	"github.com/pkg/errors"
)

func main() {
	ctx := context.Background()

	cards, err := getAllCards(ctx, "all_cards")
	if err != nil {
		log.Fatal(err)
	}

	result := make(map[string][]string)

	const workersCount = 5
	jobIndex := make(chan int, workersCount)

	var wg *sync.WaitGroup

	var mu *sync.Mutex

	for w := 1; w <= workersCount; w++ {
		go worker(ctx, wg, mu, cards, jobIndex, result)
		wg.Add(1)
	}

	for index := range cards {
		jobIndex <- index
	}

	close(jobIndex)
	wg.Wait()

	jsondata, err := json.Marshal(result)
	if err != nil {
		log.Println(err)
	}

	err = os.WriteFile("./cards.json", jsondata, os.ModePerm)
	if err != nil {
		log.Println(err)
	}
}

func downloadAndSave(
	ctx context.Context,
	mu *sync.Mutex,
	iURIs *scryfall.ImageURIs,
	filepath, cardid *string,
	result map[string][]string,
) error {
	if iURIs.Normal == "" {
		return errors.New("normal image url not found")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, iURIs.Normal, http.NoBody)
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

	err = os.MkdirAll(path.Dir(*filepath), os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "cant create directory")
	}

	err = os.WriteFile(*filepath, body, os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "cant write file")
	}

	mu.Lock()
	result[*cardid] = append(result[*cardid], *filepath)
	mu.Unlock()

	return nil
}

func getAllCards(ctx context.Context, datatype string) (cards []scryfall.Card, err error) {
	client, err := scryfall.NewClient()
	if err != nil {
		return
	}

	lbd, err := client.ListBulkData(ctx)
	if err != nil {
		return
	}

	for index := range lbd {
		if lbd[index].Type != datatype {
			continue
		}

		log.Println("Downloading bulk data...")

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, lbd[index].DownloadURI, http.NoBody)
		if err != nil {
			return nil, errors.Wrap(err, "cant create request")
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, errors.Wrap(err, "cant do request")
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			resp.Body.Close()

			return nil, errors.Wrap(err, "cant read body")
		}

		err = json.Unmarshal(body, &cards)
		if err != nil {
			resp.Body.Close()

			return nil, errors.Wrap(err, "cant unmarshal json")
		}

		resp.Body.Close()
	}

	return cards, nil
}

func worker(
	ctx context.Context,
	wg *sync.WaitGroup,
	mu *sync.Mutex,
	cards []scryfall.Card,
	indexChan <-chan int,
	result map[string][]string,
) {
	defer wg.Done()

	for index := range indexChan {
		IDString := cards[index].Name + " " +
			cards[index].Set + " " +
			string(cards[index].Lang) + " " +
			cards[index].CollectorNumber

		log.Printf("%.2f%% %s", float64(index)*float64(100)/float64(len(cards)), IDString)

		//nolint:lll // This line can't be shorter
		if *cards[index].ImageStatus == scryfall.ImageStatusHighres || *cards[index].ImageStatus == scryfall.ImageStatusLowres {
			if cards[index].ImageURIs != nil {
				imagePath := "./images" + "/" + cards[index].Set + "/" + string(cards[index].Lang) + "/" + cards[index].ID + ".jpg"
				if err := downloadAndSave(ctx, mu, cards[index].ImageURIs, &imagePath, &IDString, result); err != nil {
					log.Println(err)

					continue
				}
			}

			for faceIndex := range cards[index].CardFaces {
				imagePath := "./images" + "/" +
					cards[index].Set + "/" +
					string(cards[index].Lang) + "/" +
					cards[index].ID + strconv.Itoa(faceIndex) + ".jpg"

				if err := downloadAndSave(ctx, mu, &cards[index].CardFaces[faceIndex].ImageURIs, &imagePath, &IDString, result); err != nil {
					log.Println(err)

					continue
				}
			}
		}
	}
}
