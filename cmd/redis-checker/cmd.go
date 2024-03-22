package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/cockroachdb/errors"
	"github.com/redis/go-redis/v9"
)

const (
	firstPort = 6378
	lastPort  = 6460
	firstHost = 1
	lastHost  = 80
)

type PortHosts struct {
	Port  int   `json:"port"`
	Hosts []int `json:"hosts"`
}

func main() {
	ctx := context.Background()

	portHostsList := getEmptyPortHosts(ctx)

	jsonData, err := json.MarshalIndent(portHostsList, "", "  ")
	if err != nil {
		log.Fatalf("Failed to generate JSON: %v", err)
	}

	//nolint:gomnd // 0o600 is octal a UNIX permission
	err = os.WriteFile("empty_ports_hosts.json", jsonData, 0o600)
	if err != nil {
		log.Fatalf("Failed to write JSON to file: %v", err)
	}

	log.Println("JSON file generated successfully")
}

func getEmptyPortHosts(ctx context.Context) []PortHosts {
	var portHostsList []PortHosts

	wg := sync.WaitGroup{}

	for port := firstPort; port <= lastPort; port++ {
		wg.Add(1)

		go func(port int) {
			defer wg.Done()

			emptyHosts := getEmptyHostsForPort(ctx, port)
			if len(emptyHosts) > 0 {
				portHostsList = append(portHostsList, PortHosts{Port: port, Hosts: emptyHosts})
			}
		}(port)
	}

	wg.Wait()

	return portHostsList
}

func getEmptyHostsForPort(ctx context.Context, port int) []int {
	var emptyHosts []int

	for host := firstHost; host <= lastHost; host++ {
		isAvailable, err := isRedisAvailable(ctx, host, port)
		if err != nil || !isAvailable {
			continue
		}

		if isRedisEmpty(ctx, host, port) {
			emptyHosts = append(emptyHosts, host)
		}
	}

	return emptyHosts
}

func isRedisAvailable(ctx context.Context, host, port int) (bool, error) {
	err := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("172.21.%d.248:%d", host, port),
		Password: "", // no password set
		DB:       0,  // use the default database
	}).Ping(ctx).Err()
	if err != nil {
		return false, errors.Wrap(err, "failed to ping redis")
	}

	return true, nil
}

func getDBCount(ctx context.Context, host, port int) int {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("172.21.%d.248:%d", host, port),
		Password: "", // no password set
		DB:       0,  // use the default database
	})

	result, err := client.ConfigGet(ctx, "databases").Result()
	if err != nil {
		return 0
	}

	count, err := strconv.Atoi(result["databases"])
	if err != nil {
		return 0
	}

	return count
}

func isRedisEmpty(ctx context.Context, host, port int) bool {
	dbcount := getDBCount(ctx, host, port)

	for db := range dbcount {
		client := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("172.21.%d.248:%d", host, port),
			Password: "", // no password set
			DB:       db, // use the specified database
		})

		// check if database is empty
		dbSizeCmd := client.DBSize(ctx)
		if dbSizeCmd.Err() != nil {
			return false
		}

		if dbSizeCmd.Val() > 0 {
			return false
		}
	}

	return true
}
