package db


import (
	"gopkg.in/redis.v4"
	"github.com/GetSimpl/go-simpl/logger"
	"fmt"
)

var redisClient *redis.Client

func Init() *redis.Client{
	options := &redis.Options{
		Addr: "localhost:6379",
		Password: "",
		DB: 0,
		PoolSize: 10,

	}
	redisClient = redis.NewClient(options)
	pong, err := redisClient.Ping().Result()
	if err != nil {
		logger.E("Redis Connection Failed!")
		panic(err)
	}

	fmt.Println(pong)

	return redisClient
}

func Get() *redis.Client {
	return redisClient
}

func CloseConnection() (error, bool) {
	err := redisClient.Close()

	if err != nil {
		logger.E("Redis Connection could not be closed")
		panic(err)
	}

	return nil, true
}