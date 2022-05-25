package main

import (
	"context"
	"crypto/sha256"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

func main() {
	rds := newRedisStorage("localhost:6379", "")

	for i := 79000000000; i < 80000000000; i++ {
		h := sha256.New()
		numStr := strconv.Itoa(i)

		_, err := rds.Set(context.TODO(), numStr, string(h.Sum([]byte(numStr))), 0).Result()
		if err != nil {
			log.Println(err)
		}
	}
}

func newRedisStorage(addr, pass string) *redis.Client {
	cli := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       0,
	})

	pong := ""

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	for pong != "PONG" {
		pong, _ = cli.Ping(ctx).Result()

		time.Sleep(time.Second)
	}

	return cli
}
