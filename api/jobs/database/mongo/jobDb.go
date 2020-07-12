package jobDb

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

var Client *mongo.Client

func init() {
	e := godotenv.Load()
	if e != nil {
		fmt.Print(e)
	}
	host := os.Getenv("MONGO_DB_HOST")
	port := os.Getenv("MONGO_DB_PORT")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var err error
	connectionString := fmt.Sprintf("mongodb://%s:%s", host, port)
	Client, err = mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Jobs' Database is successfully connected to Application")
}
