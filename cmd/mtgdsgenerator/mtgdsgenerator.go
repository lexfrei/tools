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
	"github.com/cockroachdb/errors"
	"github.com/lexfrei/tools/cmd/mtgdsgenerator/cmd"
)

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
	defer generateReport(result)

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
}

func generateReport(result map[string][]string) {
	log.Println("Generating report...")

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
	conf *config,
	iURIs *scryfall.ImageURIs,
	filepath, cardid *string,
	result map[string][]string,
) error {
	if iURIs.Normal == "" {
		return errors.Errorf("normal image url not found for %s", *cardid)
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

	conf.mu.Lock()
	result[*cardid] = append(result[*cardid], *filepath)
	conf.mu.Unlock()

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

		//nolint:lll // This line can't be shorter
		if *cards[index].ImageStatus == scryfall.ImageStatusHighres || *cards[index].ImageStatus == scryfall.ImageStatusLowres {
			if cards[index].ImageURIs != nil {
				imagePath := "./images" + "/" + cards[index].Set + "/" + string(cards[index].Lang) + "/" + cards[index].ID + ".jpg"
				if err := downloadAndSave(ctx, conf, cards[index].ImageURIs, &imagePath, &IDString, result); err != nil {
					log.Println(err)

					continue
				}

				continue
			}

			for faceIndex := range cards[index].CardFaces {
				imagePath := "./images" + "/" +
					cards[index].Set + "/" +
					string(cards[index].Lang) + "/" +
					cards[index].ID + strconv.Itoa(faceIndex) + ".jpg"

				if err := downloadAndSave(ctx, conf, &cards[index].CardFaces[faceIndex].ImageURIs, &imagePath, &IDString, result); err != nil {
					log.Println(err)

					continue
				}
			}
		}
	}
}
