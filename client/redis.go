package client

import (
	"log"
	"github.com/go-redis/redis"
	"github.com/pranav1698/go-proxy/config"
)

type RedisClient struct {
	client *redis.Client
}

func Connect(config config.Config)  (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr: config.CacheHost + ":" + config.CachePort,
		Password: config.Password,
		DB: config.DB,

	})

	_, err := client.Ping().Result()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	rdbClient := &RedisClient {
		client: client,
	}

	return rdbClient, nil
}

func (rdbClient *RedisClient) AddData(url string, data string) (error){
	err := rdbClient.client.Set(url, data, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func (rdbClient *RedisClient) GetData(url string) (string, error){
	val, err := rdbClient.client.Get(url).Result()
	if err != nil {
		log.Println(err)
		return "", err
	}

	return val, nil
}