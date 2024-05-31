package main

import (
	"crypto/rand"
	"log"
	"math/big"
	"time"
)

//nolint:mnd // not interesting
func main() {
	log.Println(getRandomDate(365))
}

func getRandomDate(days uint) string {
	//nolint:mnd // not interesting
	val, _ := rand.Int(rand.Reader, big.NewInt(int64(days*86400)))

	return time.Unix(time.Now().Unix()-val.Int64(), 0).Format("20060102")
}
