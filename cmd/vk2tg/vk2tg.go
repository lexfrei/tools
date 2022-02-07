package main

import (
	"log"
	"os"
	"strconv"
	"time"

	vt "github.com/lexfrei/tools/internal/pkg/vk2tg"
)

const (
	period      = 10 * time.Second
	serviceName = "VK2TG"
)

func main() {
	logger := log.New(os.Stdout, serviceName+": ", log.Ldate|log.Ltime|log.Lshortfile)

	user, err := strconv.ParseInt(os.Getenv("V2T_TG_USER"), 10, 64)
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

	if os.Getenv("V2T_REDIS_ADDR") != "" {
		vtClient.WithRedis(
			serviceName,
			os.Getenv("V2T_REDIS_ADDR"),
			os.Getenv("V2T_REDIS_PASS"),
		)
	}

	if err = vtClient.Start(); err != nil {
		logger.Fatalln(err)
	}

	vtClient.Wait()
}
