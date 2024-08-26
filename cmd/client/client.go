package main

import (
	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
	"io"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	// connect nats
	nc, err := nats.Connect(nats.DefaultURL, nats.Name("producer"))
	if err != nil {
		log.Fatal("Error when nats.Connect:", err)
	}
	defer nc.Close()

	open, err := os.Open("model.json")
	if err != nil {
		return
	}
	defer open.Close()

	data, _ := io.ReadAll(open)
	if err != nil {
		return
	}

	err = nc.Publish("ordersSubject", data)
	if err != nil {
		log.Fatal(err)
	}

}
