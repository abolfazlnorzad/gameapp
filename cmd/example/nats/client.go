package main

import (
	"fmt"
	"gameapp/entity"
	"github.com/nats-io/nats.go"
	"os"
	"time"
)

func main() {
	// Use the env variable if running in the container, otherwise use the default.
	url := os.Getenv("NATS_URL")
	if url == "" {
		url = nats.DefaultURL
	}

	// Create an unauthenticated connection to NATS.
	nc, e := nats.Connect(url)
	if e != nil {
		fmt.Printf("eee %+v \n", e)
	}

	for {

		sub, err := nc.SubscribeSync(string(entity.MatchingUsersMatchedEvent))
		if err != nil {
			fmt.Println("err client", err)
		}
		// For a synchronous subscription, we need to fetch the next message.
		// However.. since the publish occured before the subscription was
		// established, this is going to timeout.
		msg, _ := sub.NextMsg(1 * time.Second)
		fmt.Println("message received", msg)

	}

}
