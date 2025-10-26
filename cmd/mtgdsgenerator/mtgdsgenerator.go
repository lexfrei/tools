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
	"time"

	scryfall "github.com/BlueMonday/go-scryfall"
	"github.com/cockroachdb/errors"
	"github.com/lexfrei/tools/cmd/mtgdsgenerator/cmd"
)

const (
	httpClientTimeout      = 30 * time.Second
	bulkDataTimeout        = 5 * time.Minute
	filePermission         = 0o600
	directoryPermission    = 0o755
	progressUpdateInterval = 500 * time.Millisecond
	bytesPerMegabyte       = 1024 * 1024
)

// config holds the configuration for the server.
// It has a waitgroup for waiting for the server to shutdown.
type config struct {
	wg         sync.WaitGroup
	httpClient *http.Client
}

// progressReader wraps an io.Reader to display download progress.
type progressReader struct {
	reader    io.Reader
	total     int64
	current   int64
	lastPrint time.Time
}

// Read implements io.Reader and prints progress periodically.
//
//nolint:wrapcheck // error passthrough required for io.Reader interface
func (pr *progressReader) Read(p []byte) (int, error) {
	bytesRead, err := pr.reader.Read(p)
	pr.current += int64(bytesRead)

	if time.Since(pr.lastPrint) > progressUpdateInterval || err == io.EOF {
		percentage := float64(0)
		if pr.total > 0 {
			percentage = float64(pr.current) * 100 / float64(pr.total)
		}

		log.Printf("\rDownloaded: %.2f MB / %.2f MB (%.1f%%)",
			float64(pr.current)/bytesPerMegabyte,
			float64(pr.total)/bytesPerMegabyte,
			percentage)
		pr.lastPrint = time.Now()

		if err == io.EOF {
			log.Println()
		}
	}

	return bytesRead, err
}

func main() {
	cmd.Execute()

	conf := config{
		httpClient: &http.Client{
			Timeout: httpClientTimeout,
		},
	}

	ctx := context.Background()

	cards, err := getAllCards(ctx)
	if err != nil {
		log.Fatal(err)
	}

	var result sync.Map

	log.Printf("Spawning %d workers...\n", cmd.Parallel)

	//nolint:gosec // not interesting
	workersCount := int(cmd.Parallel)
	jobIndex := make(chan int, workersCount)

	for workerID := 1; workerID <= workersCount; workerID++ {
		go worker(ctx, &conf, cards, jobIndex, &result)
	}

	for index := range cards {
		jobIndex <- index
	}

	close(jobIndex)
	conf.wg.Wait()

	err = generateReport(&result)
	if err != nil {
		log.Fatalf("error generating report: %s", err)
	}
}

// GenerateReport is a function that generates a report when passed a sync.Map of card names to slices of card types.
func generateReport(result *sync.Map) error {
	log.Println("Generating report...")

	resultMap := make(map[string][]string)
	result.Range(func(key, value any) bool {
		keyStr, keyOK := key.(string)
		mapValue, valueOK := value.(*syncMapValue)
		if keyOK && valueOK {
			mapValue.mu.Lock()
			resultMap[keyStr] = make([]string, len(mapValue.paths))
			copy(resultMap[keyStr], mapValue.paths)
			mapValue.mu.Unlock()
		}

		return true
	})

	jsondata, err := json.Marshal(resultMap)
	if err != nil {
		return errors.Wrap(err, "cant marshal json")
	}

	err = os.WriteFile("./cards.json", jsondata, filePermission)
	if err != nil {
		return errors.Wrap(err, "cant write file")
	}

	return nil
}

func downloadAndSave(ctx context.Context, client *http.Client, imageurl, filepath string) error {
	_, err := url.ParseRequestURI(imageurl)
	if err != nil {
		return errors.Wrap(err, "invalid url")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, imageurl, http.NoBody)
	if err != nil {
		return errors.Wrap(err, "cannot create request")
	}

	resp, err := client.Do(req)
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

	err = os.MkdirAll(path.Dir(filepath), directoryPermission)
	if err != nil {
		return errors.Wrap(err, "cannot create directory")
	}

	err = os.WriteFile(filepath, body, filePermission)
	if err != nil {
		return errors.Wrap(err, "cannot write file")
	}

	return nil
}

//nolint:funlen // Function handles complex bulk data download with progress
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

	bulkClient := &http.Client{
		Timeout: bulkDataTimeout,
	}

	for index := range lbd {
		if lbd[index].Type != cmd.DataType {
			continue
		}

		log.Printf("Downloading bulk data for %s...\n", cmd.DataType)

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, lbd[index].DownloadURI, http.NoBody)
		if err != nil {
			return nil, errors.Wrap(err, "cant create request")
		}

		resp, err := bulkClient.Do(req)
		if err != nil {
			return nil, errors.Wrap(err, "cant do request")
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()

			return nil, errors.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
		}

		progress := &progressReader{
			reader:    resp.Body,
			total:     resp.ContentLength,
			lastPrint: time.Now(),
		}

		body, err := io.ReadAll(progress)
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

//nolint:funlen,gocyclo,cyclop // This function handles complex card processing logic
func worker(
	ctx context.Context,
	conf *config,
	cards []scryfall.Card,
	indexChan <-chan int,
	result *sync.Map,
) {
	conf.wg.Add(1)
	defer conf.wg.Done()

	for index := range indexChan {
		select {
		case <-ctx.Done():
			return
		default:
		}

		card := cards[index]

		// Filter by language if specified.
		if cmd.Language != "" && string(card.Lang) != cmd.Language {
			continue
		}

		IDString := card.Name + " " +
			card.Set + " " +
			string(card.Lang) + " " +
			card.CollectorNumber

		log.Printf("%.2f%% %s", float64(index)*float64(100)/float64(len(cards)), IDString)

		if *card.ImageStatus != scryfall.ImageStatusHighres && *card.ImageStatus != scryfall.ImageStatusLowres {
			continue
		}

		if card.ImageURIs != nil {
			imagePath := "./images" + "/" + card.Set + "/" + string(card.Lang) + "/" + card.ID + ".jpg"

			err := downloadAndSave(ctx, conf.httpClient, card.ImageURIs.Normal, imagePath)
			if err != nil {
				log.Printf("error downloading %s: %s", IDString, err)

				continue
			}

			appendToSyncMap(result, IDString, imagePath)

			continue
		}

		for faceIndex := range card.CardFaces {
			imagePath := "./images" + "/" + card.Set + "/" + string(card.Lang) + "/" +
				card.ID + strconv.Itoa(faceIndex) + ".jpg"

			err := downloadAndSave(ctx, conf.httpClient, card.CardFaces[faceIndex].ImageURIs.Normal, imagePath)
			if err != nil {
				log.Printf("error downloading %s: %s", IDString, err)

				continue
			}

			appendToSyncMap(result, IDString, imagePath)
		}
	}
}

type syncMapValue struct {
	mu    sync.Mutex
	paths []string
}

func appendToSyncMap(syncMap *sync.Map, key, value string) {
	actual, _ := syncMap.LoadOrStore(key, &syncMapValue{paths: []string{}})
	mapValue, valueOK := actual.(*syncMapValue)
	if !valueOK {
		return
	}

	mapValue.mu.Lock()
	defer mapValue.mu.Unlock()

	mapValue.paths = append(mapValue.paths, value)
}
