package vk2tg

import (
	"context"
	"log"
	"strconv"

	"github.com/go-redis/redis/v8"
)

type storage interface {
	GetLastPost() int
	SetLastPost(int)
}

type redisStorage struct {
	serviceName string
	cli         *redis.Client
}

func (redisStorage *redisStorage) GetLastPost() int {
	res, err := redisStorage.cli.Get(context.TODO(), "LastPost").Result()
	if err != nil {
		return 0
	}

	postID, err := strconv.Atoi(res)
	if err != nil {
		return 0
	}

	return postID
}

func (redisStorage *redisStorage) SetLastPost(postID int) {
	_, err := redisStorage.cli.Set(context.TODO(), "LastPost", postID, 0).Result()
	if err != nil {
		log.Println(err)

		return
	}
}

func newRedisStorage(serviceName, addr, pass string) *redisStorage {
	cli := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       0,
	})

	rs := &redisStorage{
		serviceName: serviceName,
		cli:         cli,
	}

	return rs
}

func (vtCli *VTClinent) WithRedis(serviceName, redisAddr, redisPassword string) *VTClinent {
	vtCli.config.serviceName = serviceName
	vtCli.config.StorageEnabled = true

	vtCli.storage = newRedisStorage(serviceName, redisAddr, redisPassword)

	vtCli.config.LastPostID = vtCli.storage.GetLastPost()

	vtCli.logger.Printf("Connected to Redis at %s", redisAddr)

	return vtCli
}
