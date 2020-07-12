package tokenDB

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var Client *redis.Client
var ctx = context.Background()
func init() {
	e := godotenv.Load()
	if e != nil {
		fmt.Print(e)
	}
	redisDatabaseAddress := os.Getenv("TOKEN_HOST")
	Client = redis.NewClient(&redis.Options{
		Addr:     redisDatabaseAddress,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := Client.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	log.Println("Auth's Database is successfully connected to Application")
}
