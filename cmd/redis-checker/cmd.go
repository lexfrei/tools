package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	firstPort := 6378
	lastPort := 6460

	firstHost := 1
	lastHost := 80

	RedisPortsAvailibePerHost := make(map[int][]int)

	var AvailibleHosts int

	for host := firstHost; host <= lastHost; host++ {
		if pingDB(context.TODO(), host, firstPort) {
			AvailibleHosts++

			for port := firstPort; port <= lastPort; port++ {
				dbcount := getDBCount(context.TODO(), host, port)

				for db := 0; db < dbcount; db++ {
					if checkRedisIsEmpty(context.TODO(), host, port, db) {
						RedisPortsAvailibePerHost[port] = append(RedisPortsAvailibePerHost[port], host)

						break
					}
				}
			}
		}
	}

	for port, hosts := range RedisPortsAvailibePerHost {
		if len(hosts) == AvailibleHosts {
			log.Printf("Port %d is available on all hosts\n", port)
		}
	}

	// dump to json file
	jsonBytes, err := json.Marshal(RedisPortsAvailibePerHost)
	if err != nil {
		log.Println(err)

		return
	}

	// write to file
	//nolint:gomnd // coz this is just a unix permission
	err = os.WriteFile("redis-ports.json", jsonBytes, 0o600)
	if err != nil {
		log.Println(err)

		return
	}
}

func pingDB(ctx context.Context, host, port int) bool {
	client := redis.NewClient(&redis.Options{
		Addr:        fmt.Sprintf("172.21.%d.248:%d", host, port),
		Password:    "", // no password set
		DB:          0,  // use the default database
		DialTimeout: 1 * time.Second,
	})

	// check if Redis is running
	pingCmd := client.Ping(ctx)

	return pingCmd.Err() != nil
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

	i, err := strconv.Atoi(result["databases"])
	if err != nil {
		return 0
	}

	return i
}

func checkRedisIsEmpty(ctx context.Context, host, port, dbcount int) bool {
	for db := 0; db < dbcount; db++ {
		client := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("172.21.%d.248:%d", host, port),
			Password: "", // no password set
			DB:       db, // use the default database
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
