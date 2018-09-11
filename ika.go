package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gomodule/redigo/redis"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	redisAddr := fmt.Sprintf("%s:%s", os.Getenv("IKA_REDIS_HOST"), os.Getenv("IKA_REDIS_PORT"))
	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) { return redis.Dial("tcp", redisAddr) },
	}

	keyRouter, err := core.NewKeyRouter(pool)
	if err != nil {
		log.Fatal(err)
	}

	service := &core.ServiceImpl{
		KeyRouter:    keyRouter,
		QueueManager: core.NewQueueManager(core.NewStore(pool)),
	}

	server := api.BuildHTTPServer(service)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
