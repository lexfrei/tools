package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/cockroachdb/errors"
)

func main() {
	key1 := os.Args[1]
	hey2 := os.Args[2]

	if len(key1) == 0 || len(hey2) == 0 {
		log.Fatal("key1 and key2 must be provided")
	}

	encodedKey1, err := b64ToHex(key1)
	if err != nil {
		log.Fatal("cannot decode key1: ", err)
	}

	encodedKey2, err := b64ToHex(hey2)
	if err != nil {
		log.Fatal("cannot decode key2: ", err)
	}

	//nolint:forbidigo // This is a command line tool
	fmt.Printf("%s:%s\n", encodedKey2, encodedKey1)
}

func b64ToHex(b64 string) (string, error) {
	for len(b64)%4 != 0 {
		b64 += "="
	}

	data, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return "", errors.Wrap(err, "cannot decode base64")
	}

	var result string

	for _, r := range data {
		result += strconv.FormatUint(uint64(r), 16)
	}

	return result, nil
}
