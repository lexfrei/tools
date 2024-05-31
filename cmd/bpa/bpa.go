package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	line := os.Args[1]

	for len(line)%4 != 0 {
		line += "="
	}

	data, err := base64.StdEncoding.DecodeString(line)
	if err != nil {
		log.Fatal("Cannot decode base64: ", err)

		return
	}

	var result string

	for _, r := range data {
		result += strconv.FormatUint(uint64(r), 16)
	}

	//nolint:forbidigo // This is a command line tool
	fmt.Println(result)
}
