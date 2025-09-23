package main

import (
	"context"
	"crypto/sha256"
	"log"
	"strconv"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/redis/go-redis/v9"
)

func main() {
	rds, err := newRedisStorage("localhost:6379", "")
	if err != nil {
		log.Println(err)

		return
	}

	for i := 79000000000; i < 80000000000; i++ {
		h := sha256.New()
		numStr := strconv.Itoa(i)

		_, err := rds.Set(context.TODO(), numStr, string(h.Sum([]byte(numStr))), 0).Result()
		if err != nil {
			log.Println(err)
		}
	}
}

func newRedisStorage(addr, pass string) (*redis.Client, error) {
	cli := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	_, err := cli.Ping(ctx).Result()
	if err != nil {
		return nil, errors.Wrap(err, "ping redis")
	}

	return cli, nil
}
