package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"sync"

	scryfall "github.com/BlueMonday/go-scryfall"
	"github.com/cockroachdb/errors"
	"github.com/lexfrei/tools/cmd/mtgdsgenerator/cmd"
)

// config holds the configuration for the server.
// It has a waitgroup for waiting for the server to shutdown,
// and a mutex for synchronizing access to the configuration.
type config struct {
	wg sync.WaitGroup
	mu sync.Mutex
}

func main() {
	cmd.Execute()

	var conf config

	ctx := context.Background()

	cards, err := getAllCards(ctx)
	if err != nil {
		log.Fatal(err)
	}

	result := make(map[string][]string)

	log.Printf("Spawning %d workers...\n", cmd.Parallel)

	workersCount := int(cmd.Parallel)
	jobIndex := make(chan int, workersCount)

	for w := 1; w <= workersCount; w++ {
		go worker(ctx, &conf, cards, jobIndex, result)
	}

	for index := range cards {
		jobIndex <- index
	}

	close(jobIndex)
	conf.wg.Wait()

	err = generateReport(result)
	if err != nil {
		log.Fatalf("error generating report: %s", err)
	}
}

// GenerateReport is a function that generates a report when passed a map of card names to slices of card types

func generateReport(result map[string][]string) error {
	log.Println("Generating report...")

	jsondata, err := json.Marshal(result)
	if err != nil {
		return errors.Wrap(err, "cant marshal json")
	}

	err = os.WriteFile("./cards.json", jsondata, os.ModePerm)

	if err != nil {
		return errors.Wrap(err, "cant write file")
	}

	return nil
}

func downloadAndSave(ctx context.Context, imageurl, filepath string) error {
	_, err := url.ParseRequestURI(imageurl)
	if err != nil {
		return errors.Wrap(err, "invalid url")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, imageurl, http.NoBody)
	if err != nil {
		return errors.Wrap(err, "cannot create request")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "cannot download file")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "cannot read body")
	}

	if len(body) == 0 {
		return errors.New("empty response body")
	}

	err = os.MkdirAll(path.Dir(filepath), os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "cannot create directory")
	}

	err = os.WriteFile(filepath, body, os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "cannot write file")
	}

	return nil
}

func getAllCards(ctx context.Context) ([]scryfall.Card, error) {
	client, err := scryfall.NewClient()
	if err != nil {
		return nil, errors.Wrap(err, "can't create scryfall client")
	}

	lbd, err := client.ListBulkData(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "can't get bulk data form scryfall")
	}

	cards := []scryfall.Card{}

	for index := range lbd {
		if lbd[index].Type != cmd.DataType {
			continue
		}

		log.Printf("Downloading bulk data for %s...\n", cmd.DataType)

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, lbd[index].DownloadURI, http.NoBody)
		if err != nil {
			return nil, errors.Wrap(err, "cant create request")
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, errors.Wrap(err, "cant do request")
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()

			return nil, errors.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
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
	conf *config,
	cards []scryfall.Card,
	indexChan <-chan int,
	result map[string][]string,
) {
	conf.wg.Add(1)
	defer conf.wg.Done()

	for index := range indexChan {
		IDString := cards[index].Name + " " +
			cards[index].Set + " " +
			string(cards[index].Lang) + " " +
			cards[index].CollectorNumber

		log.Printf("%.2f%% %s", float64(index)*float64(100)/float64(len(cards)), IDString)

		//nolint:lll,goconst // This line can't be shorter
		if *cards[index].ImageStatus == scryfall.ImageStatusHighres || *cards[index].ImageStatus == scryfall.ImageStatusLowres {
			if cards[index].ImageURIs != nil {
				imagePath := "./images" + "/" + cards[index].Set + "/" + string(cards[index].Lang) + "/" + cards[index].ID + ".jpg"

				err := downloadAndSave(ctx, cards[index].ImageURIs.Normal, imagePath)
				if err != nil {
					log.Printf("error downloading %s: %s", IDString, err)

					continue
				}

				conf.mu.Lock()
				result[IDString] = append(result[IDString], imagePath)
				conf.mu.Unlock()

				continue
			}

			for faceIndex := range cards[index].CardFaces {
				imagePath := "./images" + "/" +
					cards[index].Set + "/" +
					string(cards[index].Lang) + "/" +
					cards[index].ID + strconv.Itoa(faceIndex) + ".jpg"

				err := downloadAndSave(ctx, cards[index].CardFaces[faceIndex].ImageURIs.Normal, imagePath)
				if err != nil {
					log.Printf("error downloading %s: %s", IDString, err)

					continue
				}

				conf.mu.Lock()
				result[IDString] = append(result[IDString], imagePath)
				conf.mu.Unlock()

				continue
			}
		}
	}
}
