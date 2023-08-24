package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
)

func main() {
	nc, _ := nats.Connect(nats.DefaultURL)
	fmt.Println("nc", nc)
	log.Println("NATS server started")
}
