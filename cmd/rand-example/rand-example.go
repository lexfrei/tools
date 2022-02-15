package main

import (
	"log"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

//nolint:gomnd // not interesting
func main() {
	log.Println(getRandomDate(365))
}

func getRandomDate(days uint) string {
	//nolint:gosec,gomnd // not interesting
	return time.Unix(time.Now().Unix()-rand.Int63n(int64(days*86400)), 0).Format("20060102")
}
