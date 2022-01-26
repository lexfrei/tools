package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"

	scryfall "github.com/BlueMonday/go-scryfall"
	"github.com/pkg/errors"
)

func main() {
	cards, err := getAllCards("all_cards")
	if err != nil {
		log.Fatal(err)
	}

	result := make(map[string][]string)

	for index := range cards {
		IDString := cards[index].Name + " " +
			cards[index].Set + " " +
			string(cards[index].Lang) + " " +
			cards[index].CollectorNumber

		log.Printf("%.2f%% %s", float64(index)*float64(100)/float64(len(cards)), IDString)
		if *cards[index].ImageStatus == scryfall.ImageStatusHighres || *cards[index].ImageStatus == scryfall.ImageStatusLowres {

			if cards[index].ImageURIs != nil {
				imagePath := "./images" + "/" + cards[index].Set + "/" + string(cards[index].Lang) + "/" + cards[index].ID + ".jpg"
				err = downloadAndSave(cards[index].ImageURIs.Normal, imagePath)
				if err != nil {
					log.Println(err)
					continue
				}
				result[IDString] = []string{imagePath}
				continue
			}

			if cards[index].CardFaces != nil {
				for i := range cards[index].CardFaces {
					if cards[index].CardFaces[i].ImageURIs.Normal == "" {
						continue
					}

					imagePath := "./images" + "/" + cards[index].Set + "/" + string(cards[index].Lang) +
						"/" + cards[index].ID + strconv.Itoa(i) + ".jpg"

					err = downloadAndSave(cards[index].CardFaces[i].ImageURIs.Normal, imagePath)
					if err != nil {
						log.Println(err)
						continue
					}

					result[IDString] = append(result[IDString], imagePath)
					continue
				}
			}
		}
	}

	jsondata, err := json.Marshal(result)
	if err != nil {
		log.Println(err)
	}

	err = os.WriteFile("./cards.json", jsondata, os.ModePerm)
	if err != nil {
		log.Println(err)
	}
}

func downloadAndSave(url, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return errors.Wrap(err, "cant dowbload file")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "cant read body")
	}

	err = os.MkdirAll(path.Dir(filepath), os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "cant create directory")
	}

	err = os.WriteFile(filepath, body, os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "cant write file")
	}
	return nil
}

func getAllCards(datatype string) (cards []scryfall.Card, err error) {
	ctx := context.Background()
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
		resp, err := http.Get(lbd[index].DownloadURI)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, errors.Wrap(err, "cant read body")
		}
		err = json.Unmarshal(body, &cards)
		if err != nil {
			return nil, errors.Wrap(err, "cant unmarshal json")
		}
	}
	return cards, nil
}
