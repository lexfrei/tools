package main

import (
	"log"
	"os"
	"strconv"
	"time"

	vt "github.com/lexfrei/tools/internal/pkg/vk2tg"
)

const period = 10 * time.Second

func main() {
	logger := log.New(os.Stdout, "VK2TG: ", log.Ldate|log.Ltime|log.Lshortfile)

	user, err := strconv.Atoi(os.Getenv("V2T_TG_USER"))
	if err != nil {
		logger.Fatalf("Invalid TG user ID: %s\n", os.Getenv("V2T_TG_USER"))
	}

	vtClient := vt.NewVTClient(
		os.Getenv("V2T_TG_TOKEN"),
		os.Getenv("V2T_VK_TOKEN"),
		user,
		period,
	).WithLogger(
		logger,
	)

	if err = vtClient.Start(); err != nil {
		logger.Fatalln(err)
	}

	vtClient.Wait()
}
